// Flag: turn a feature on and off while workers keep running.
//
// Workers process jobs in a loop. Before each job, they check paused with
// Load. main turns maintenance mode on and off with Store. Workers don't
// stop. They just skip jobs while paused, and pick up again right after.
//
// stop is still a channel, because it's a one-time event that every
// worker must notice. paused flips back and forth, and workers only need
// to read its current value. That's what atomic.Bool is good at.
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

	var paused atomic.Bool      // main -> workers: "skip jobs for now"
	stop := make(chan struct{}) // main -> workers: "please stop"

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(id+1, &paused, stop, &wg)
	}

	time.Sleep(phaseTime)
	fmt.Println("main: maintenance mode on")
	paused.Store(true)

	time.Sleep(phaseTime)
	fmt.Println("main: maintenance mode off")
	paused.Store(false)

	time.Sleep(phaseTime)
	close(stop)
	wg.Wait()
}

// worker processes a job every jobTime, or skips it while paused is true,
// until stop is closed.
func worker(id int, paused *atomic.Bool, stop <-chan struct{}, wg *sync.WaitGroup) {
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
		case <-stop:
			fmt.Printf("worker %d: %d processed, %d skipped\n", id, processed, skipped)
			return
		}
	}
}
