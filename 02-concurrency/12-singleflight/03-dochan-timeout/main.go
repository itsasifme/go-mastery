// DoChan timeout: join a shared call, but stop waiting when you run out of time.
//
// A slow report is built once for everyone who asks at the same time.
// Each caller has its own time limit. DoChan returns a channel instead of
// blocking, so each caller can select on it and on its own ctx.Done().
//
// The impatient caller gives up and shows a "try again later" message.
// The report keeps building, and the patient callers still get it. One
// caller giving up doesn't cancel the work for the others.
//
// Note: with Do (see 01-basic), every caller waits as long as the slowest
// call, with no way out.
//
// Real use:
//   - an HTTP handler that joins a shared call to another service, but
//     stops waiting at r.Context().Done(). Then a client that hangs up
//     doesn't keep a goroutine waiting.
//   - a cache refresh where callers with a short deadline fall back to
//     old data while the refresh finishes.
//
// Run: go run ./02-concurrency/12-singleflight/03-dochan-timeout
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// buildTime is how long building the fake report takes.
const buildTime = 200 * time.Millisecond

func main() {
	const (
		impatient = 50 * time.Millisecond  // shorter than buildTime, gives up
		patient   = 500 * time.Millisecond // longer than buildTime, gets the report
	)

	var g singleflight.Group
	limits := map[string]time.Duration{ // caller -> how long it will wait
		"mobile app": impatient,
		"dashboard":  patient,
		"email job":  patient,
	}

	var wg sync.WaitGroup
	wg.Add(len(limits))
	for caller, limit := range limits {
		go fetchReport(caller, limit, &g, &wg)
	}
	wg.Wait()
}

// fetchReport joins the shared report build and waits up to limit for it.
func fetchReport(caller string, limit time.Duration, g *singleflight.Group, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()

	// DoChan's channel has a buffer of 1. If this caller gives up, the result
	// can still be sent, so nothing leaks. buildReport keeps running, though.
	ch := g.DoChan("daily-report", buildReport)
	select {
	case res := <-ch:
		if res.Err != nil {
			fmt.Printf("%-10s: error: %v\n", caller, res.Err)
			return
		}
		fmt.Printf("%-10s: got %q (shared=%v)\n", caller, res.Val.(string), res.Shared)
	case <-ctx.Done():
		fmt.Printf("%-10s: gave up after %v, showing \"try again later\"\n", caller, limit)
	}
}

// buildReport pretends to build a slow report.
func buildReport() (any, error) {
	fmt.Println("building the report...")
	time.Sleep(buildTime)
	return "sales: 1,204 orders", nil
}
