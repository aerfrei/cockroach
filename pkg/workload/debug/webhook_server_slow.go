// Copyright 2025 The Cockroach Authors.
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
	Use: "webhook-server-slow",
	// revisit this stuff
	Short: "webhook-server opens an http server on 3000 to which cdc's webhook can emit a table with a numeric unique 'id' column",
	RunE:  webhookServerSlow,
	Args:  cobra.NoArgs,
}

func webhookServerSlow(cmd *cobra.Command, args []string) error {
	var (
		mu    syncutil.Mutex
		seen  = map[int]struct{}{}
		size  int64
		dupes int
	)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Length  int `json:"length"`
			Payload []struct {
				After struct {
					ID int `json:"id"`
				} `json:"after"`
			} `json:"payload"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("decoding body: %v", err)
			return
		}
		var before, after, d int
		func() {
			mu.Lock()
			defer mu.Unlock()
			before = len(seen)
			after = before
			// TODO(cdc): add check for ordering guarantees using resolved timestamps and event timestamps
			for _, i := range req.Payload {
				if _, ok := seen[i.After.ID]; !ok {
					seen[i.After.ID] = struct{}{}
					after++
				} else {
					dupes++
				}
			}
			if r.ContentLength > 0 {
				size += r.ContentLength
			}
			d = dupes
		}()
		const printEvery = 10000
		if before/printEvery != after/printEvery {
			log.Printf("keys seen: %d (%d dupes); %.1f MB", after, d, float64(size)/float64(1<<20))
			log.Print("I AM HERE REMOVE ME")
		}
		if (after*7+13)%10_000 == 0 {
			log.Print("throwing for no reason REMOVE ME")
			http.Error(w, "retryable error for test", 500)
		}
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		func() {
			mu.Lock()
			defer mu.Unlock()
			seen = make(map[int]struct{}, len(seen))
			dupes = 0
			size = 0
		}()
		log.Printf("reset")
	})
	mux.HandleFunc("/unique", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		l := len(seen)
		log.Printf("keys seen*: %d", l)
		fmt.Fprintf(w, "%d", l)
	})
	mux.HandleFunc("/dupes", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
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
	log.Printf("starting server on port %d", WebhookServerPort)
	return (&http.Server{
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler:   mux,
		Addr:      fmt.Sprintf(":%d", WebhookServerPort),
	}).ListenAndServeTLS("", "")
}
