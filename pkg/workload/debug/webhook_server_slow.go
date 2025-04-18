// Copyright 2024 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package debug

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/spf13/cobra"
)

var webhookServerSlowCmd = &cobra.Command{
	Use:   "webhook-server-slow",
	Short: "webhook-server-slow opens an http server on 3000 to which cdc's webhook can emit a table with a numeric unique 'id' column",
	RunE:  webhookServerSlow,
	Args:  cobra.NoArgs,
}

const (
	WebhookServerSlowPort = 9707
)

func webhookServerSlow(cmd *cobra.Command, args []string) error {
	var (
		mu    syncutil.Mutex
		seen  = map[string]struct{}{}
		size  int64
		dupes int
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Length  int `json:"length"`
			Payload []struct {
				After struct {
					ID  int `json:"id"`
					VAL int `json:"val"`
				} `json:"after"`
				Before struct {
					ID  int `json:"id"`
					VAL int `json:"val"`
				} `json:"before"`
				Updated string `json:"updated"`
			} `json:"payload"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("AMF: error decoding body: %v", err)
			return
		}

		sleepDurations := []time.Duration{
			0 * time.Millisecond, 10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
			0 * time.Millisecond, 10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
			0 * time.Millisecond, 10 * time.Millisecond, 20 * time.Millisecond, 30 * time.Millisecond,
		}

		var before, after, d int
		func() {
			mu.Lock()
			defer mu.Unlock()
			before = len(seen)
			after = before
			// TODO(cdc): add check for ordering guarantees using resolved timestamps and event timestamps
			for _, i := range req.Payload {
				id := i.After.ID
				seenKey := fmt.Sprintf("%d-%s", id, i.Updated)
				if _, ok := seen[seenKey]; !ok {
					seen[seenKey] = struct{}{}
					after++

					if id > 100 {
						time.Sleep(sleepDurations[id/100])
					}
					if (id+i.After.VAL)%317 == 0 {
						http.Error(w, "transient sink error", 500)
						return
					}

				} else {
					log.Printf("AMF: seen %d dupe", id)
					dupes++
				}
			}
			if r.ContentLength > 0 {
				size += r.ContentLength
			}
			d = dupes
		}()
		const printEvery = 100
		if before/printEvery != after/printEvery {
			log.Printf("AMF: keys seen: %d (%d dupes); %.1f MB", after, d, float64(size)/float64(1<<20))
		}
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		func() {
			mu.Lock()
			defer mu.Unlock()
			seen = make(map[string]struct{}, len(seen))
			dupes = 0
			size = 0
		}()
		log.Printf("AMF: reset")
		log.Printf("reset")
	})
	mux.HandleFunc("/unique", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		l := len(seen)
		log.Printf("AMF: keys seen: %d", l)
		fmt.Fprintf(w, "%d", l)
	})
	mux.HandleFunc("/dupes", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		log.Printf("AMF: dupes: %d", dupes)
		fmt.Fprintf(w, "%d", dupes)
	})

	mux.HandleFunc("/exit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		go func() {
			time.Sleep(time.Millisecond * 5)
			exit.WithCode(exit.Success())
		}()
	})

	cert, err := genKeyPair()
	if err != nil {
		return err
	}
	log.Printf("starting server on port %d", WebhookServerSlowPort)
	return (&http.Server{
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler:   mux,
		Addr:      fmt.Sprintf(":%d", WebhookServerSlowPort),
	}).ListenAndServeTLS("", "")
}
