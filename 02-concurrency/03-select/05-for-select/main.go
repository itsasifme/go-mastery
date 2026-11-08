// For-select loop: keep handling several channels until told to stop.
//
// A worker loops forever. Each round, select waits for either a new job
// or the quit signal. This "for + select" loop is the most common shape
// of a long-running goroutine in Go.
//
// main sends 3 jobs, then closes quit. The worker sees the close and
// stops, then closes done so main knows it really finished.
//
// Note: inside select, break only leaves the select, not the for loop.
// Use return (or a labeled break) to stop the loop.
//
// Run: go run ./02-concurrency/03-select/05-for-select
package main

import "fmt"

func main() {
	jobs := make(chan int)
	quit := make(chan struct{})
	done := make(chan struct{})

	go worker(jobs, quit, done)

	for i := 1; i <= 3; i++ {
		jobs <- i
	}
	close(quit) // tell the worker to stop (see closing channels)
	<-done      // wait until it really stopped
}

// worker handles jobs until quit is closed.
func worker(jobs <-chan int, quit <-chan struct{}, done chan<- struct{}) {
	defer close(done)

	for {
		select {
		case j := <-jobs:
			fmt.Println("job", j)
		case <-quit:
			fmt.Println("worker: stopping")
			return // not break: break would only leave the select
		}
	}
}
