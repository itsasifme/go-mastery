// Flag: turn a feature on and off while workers keep running.
//
// Workers process jobs in a loop. Before each job, they check paused with
// Load. main turns maintenance mode on and off with Store. Workers don't
// stop. They just skip jobs while paused, and pick up again right after.
//
// Stopping uses ctx, because it's a one-time event that every worker must
// notice. paused flips back and forth, and workers only need to read its
// current value. That's what atomic.Bool is good at. A ctx can't do it:
// once canceled, it stays canceled.
//
// Real use:
//   - a "maintenance mode" or "read-only mode" switch that an admin
//     endpoint flips while the server keeps running.
//   - a ready flag that a /healthz handler reads, and that startup code
//     sets once caches are warm.
//
// Run: go run ./02-concurrency/09-atomic/02-flag
// Race check: go run -race ./02-concurrency/09-atomic/02-flag
package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// jobTime is how long one fake job takes.
const jobTime = 20 * time.Millisecond

func main() {
	const (
		workerCount = 2
		phaseTime   = 200 * time.Millisecond // how long each phase (on, paused, on) lasts
	)

	var paused atomic.Bool                                  // main -> workers: "skip jobs for now"
	ctx, cancel := context.WithCancel(context.Background()) // cancel tells every worker to stop

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(ctx, id+1, &paused, &wg)
	}

	time.Sleep(phaseTime)
	fmt.Println("main: maintenance mode on")
	paused.Store(true)

	time.Sleep(phaseTime)
	fmt.Println("main: maintenance mode off")
	paused.Store(false)

	time.Sleep(phaseTime)
	cancel()
	wg.Wait()
}

// worker processes a job every jobTime, or skips it while paused is true,
// until ctx is canceled.
func worker(ctx context.Context, id int, paused *atomic.Bool, wg *sync.WaitGroup) {
	defer wg.Done()

	processed, skipped := 0, 0
	for {
		select {
		case <-time.After(jobTime):
			if paused.Load() {
				skipped++
				continue
			}
			processed++ // simulate one job
		case <-ctx.Done():
			fmt.Printf("worker %d: %d processed, %d skipped\n", id, processed, skipped)
			return
		}
	}
}
