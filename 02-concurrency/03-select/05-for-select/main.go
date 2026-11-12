// For-select loop: keep handling several channels until told to stop.
//
// A worker loops forever. Each round, select waits for either a new job
// or the stop signal. This "for + select" loop is the most common shape
// of a long-running goroutine in Go.
//
// main sends a few jobs, then closes stop. The worker sees the close and
// stops, then closes done so main knows it really finished.
//
// Note: inside select, break only leaves the select, not the for loop.
// Use return (or a labeled break) to stop the loop.
//
// Real use:
//   - the main loop of a background service: take messages from a queue
//     until shutdown.
//
// Run: go run ./02-concurrency/03-select/05-for-select
package main

import "fmt"

func main() {
	const jobCount = 3

	jobs := make(chan int)      // main -> worker: the work
	stop := make(chan struct{}) // main -> worker: "please stop"
	done := make(chan struct{}) // worker -> main: "I'm finished"

	go worker(jobs, stop, done)

	for i := range jobCount {
		jobs <- i
	}
	close(stop) // tell the worker to stop (see channels/02-closing/02-broadcast)
	<-done      // wait until it really stopped
}

// worker handles jobs until stop is closed.
func worker(jobs <-chan int, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	for {
		select {
		case j := <-jobs:
			fmt.Println("job", j)
		case <-stop:
			fmt.Println("worker: stopping")
			return // not break: break would only leave the select
		}
	}
}
