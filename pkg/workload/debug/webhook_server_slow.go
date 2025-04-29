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
	"strconv"
	"time"

	"github.com/cockroachdb/cockroach/pkg/cli/exit"
	"github.com/cockroachdb/cockroach/pkg/util/syncutil"
	"github.com/spf13/cobra"
)

var webhookServerSlowCmd = &cobra.Command{
	Use:   "webhook-server-slow [transient error freq ms] ...[range delays ms]",
	Short: "webhook-server-slow opens an http server on 3000 to which cdc's webhook can emit a table with a numeric unique 'id' column",
	RunE:  webhookServerSlow,
	Args:  cobra.MinimumNArgs(0),
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
	// Add a variable to track the last time a transient error occurred.
	lastTransientErrorTime := time.Now()

	const (
		defaultTransientErrorFrequency = 1000 // Default transient error frequency in ms
	)
	defaultRangeDelays := []time.Duration{} // Default range delays in ms (no delay)

	// Parse the arguments
	transientErrorFrequency := defaultTransientErrorFrequency
	rangeDelays := defaultRangeDelays

	if len(args) > 0 {
		if freq, err := strconv.Atoi(args[0]); err == nil {
			transientErrorFrequency = freq
		} else {
			log.Printf("Invalid transient error frequency, using default: %d ms", defaultTransientErrorFrequency)
		}
	}

	if len(args) > 1 {
		rangeDelays = make([]time.Duration, len(args[1:]))
		for i, delay := range args[1:] {
			if d, err := strconv.Atoi(delay); err == nil {
				rangeDelays[i] = time.Duration(d) * time.Millisecond
			} else {
				log.Printf("Invalid range delay at index %d, using default: 0 ms", i)
				rangeDelays[i] = 0 * time.Millisecond
			}
		}
	}

	log.Printf("AF: Parsed transient error frequency: %d ms", transientErrorFrequency)
	log.Printf("AF: Parsed range delays: %v ms", rangeDelays)

	// Log the parsed or default values
	log.Printf("Transient Error Frequency: %d ms", transientErrorFrequency)
	log.Printf("Range Delays: %v ms", rangeDelays)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Length  int `json:"length"`
			Payload []struct {
				After struct {
					CRDBRegion string `json:"crdb_region"`
					O          int    `json:"no_o_id"`
					D          int    `json:"no_d_id"`
					W          int    `json:"no_w_id"`
					VAL        int    `json:"val"`
				} `json:"after"`
				Before struct {
					ID  int `json:"id"`
					VAL int `json:"val"`
				} `json:"before"`
				Key     []any  `json:"key"`
				Updated string `json:"updated"`
			} `json:"payload"`
		}

		// byts, _ := httputil.DumpRequest(r, true)
		// log.Println(string(byts), "request bytes")

		// var obj map[string]interface{}

		err := json.NewDecoder(r.Body).Decode(&req)

		// log.Printf("request body obj %s", obj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		now := time.Now()
		if transientErrorFrequency > 0 && now.Sub(lastTransientErrorTime) >= time.Duration(transientErrorFrequency)*time.Millisecond {
			lastTransientErrorTime = now
			log.Printf("Simulating transient sink error")
			http.Error(w, "transient sink error", http.StatusInternalServerError)
			return
		}

		var before, after, dupeCount int
		func() {
			mu.Lock()
			defer mu.Unlock()
			before = len(seen)
			after = before
			// TODO(cdc): add check for ordering guarantees using resolved timestamps and event timestamps
			for _, i := range req.Payload {
				var keyParts []string
				for _, part := range i.Key {
					keyParts = append(keyParts, fmt.Sprintf("%v", part))
				}
				w := i.Key[1]
				d := i.Key[2]
				o := i.Key[3]
				// fmt.Println(strings.Join(keyParts, "-"), "key")
				seenKey := fmt.Sprintf("%d-%d-%d-%s", w, d, o, i.Updated)
				if _, ok := seen[seenKey]; !ok {
					seen[seenKey] = struct{}{}
					after++

					if wInt, ok := w.(float64); ok {
						if int(wInt)%2000 < 10 {
							fmt.Println("w is eligible", wInt)
							sleepDurationIndex := int(wInt / 2000)
							if sleepDurationIndex < len(rangeDelays) {
								timeToSleep := rangeDelays[sleepDurationIndex]
								fmt.Printf("Sleeping for %v ms for seenkey %s", timeToSleep.Milliseconds(), seenKey)
								time.Sleep(timeToSleep)
							} else {
								fmt.Println("wInt is too large for rangeDelays", wInt)
							}
						}
					} else {
						fmt.Printf("w not an int: %T", wInt)
					}
				} else {
					dupes++
				}
			}
			if r.ContentLength > 0 {
				size += r.ContentLength
			}
			dupeCount = dupes
		}()
		const printEvery = 100
		if before/printEvery != after/printEvery {
			log.Printf("keys seen: %d (%d dupes); %.1f MB", after, dupeCount, float64(size)/float64(1<<20))
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
		log.Printf("reset")
	})
	mux.HandleFunc("/unique", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		l := len(seen)
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
	log.Printf("starting server on port %d", WebhookServerSlowPort)
	return (&http.Server{
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
		Handler:   mux,
		Addr:      fmt.Sprintf(":%d", WebhookServerSlowPort),
	}).ListenAndServeTLS("", "")
}
