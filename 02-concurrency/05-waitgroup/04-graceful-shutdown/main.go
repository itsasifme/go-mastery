// Graceful shutdown: let running jobs finish before the program exits.
//
// Workers keep processing jobs one after another. When main closes stop,
// each worker finishes the job it's on, then returns. It never cuts a job
// in half. main calls wg.Wait, so it doesn't exit while a job is still
// running.
//
// Try: comment out wg.Wait(). main exits right away, and the workers never
// print "stopping".
//
// Note: in a real program, stop comes from Ctrl+C with
// signal.NotifyContext(ctx, os.Interrupt) (see context/01-cancellation/02-stop-workers).
//
// Real use:
//   - a queue consumer that finishes the current message, acks it, then
//     exits, so no message is lost or handled twice.
//   - http.Server.Shutdown(ctx), which waits for requests in progress the
//     same way before the server stops.
//
// Run: go run ./02-concurrency/05-waitgroup/04-graceful-shutdown
package main

import (
	"fmt"
	"sync"
	"time"
)

// jobTime is how long one fake job takes.
const jobTime = 150 * time.Millisecond

func main() {
	const (
		workerCount = 3
		runTime     = 400 * time.Millisecond // not a multiple of jobTime, so stop lands mid-job
		roundTo     = 10 * time.Millisecond
	)

	stop := make(chan struct{}) // main -> workers: "please stop"

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(id+1, stop, &wg)
	}

	time.Sleep(runTime) // the program runs for a while
	fmt.Println("main: shutting down, waiting for jobs in progress")
	close(stop)

	shutdownStart := time.Now()
	wg.Wait()
	fmt.Println("main: all workers stopped after", time.Since(shutdownStart).Round(roundTo))
}

// worker processes jobs until stop is closed. It checks stop only between
// jobs, so a job that started always finishes.
func worker(id int, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	jobs := 0
	for {
		select {
		case <-stop:
			fmt.Printf("worker %d: stopping after %d jobs\n", id, jobs)
			return
		default:
		}

		time.Sleep(jobTime) // simulate one job; it isn't stopped halfway
		jobs++
	}
}
