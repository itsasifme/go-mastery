// Deadline: stop at a fixed time with context.WithDeadline.
//
// WithDeadline takes a point in time. WithTimeout(d) is the same as
// WithDeadline(time.Now().Add(d)).
//
// main starts three workers that share one ctx. ctx.Deadline() tells each
// worker when the deadline is, so it can skip a job that can't finish in
// time, instead of starting it and failing halfway. The big job needs more
// time than is left, so it's skipped. The others finish.
//
// Real use:
//   - a request with a 2 second budget skips an optional slow step, like
//     loading recommendations, instead of failing the whole page.
//   - gRPC sends the deadline to the next service, so it knows how much
//     time is left too.
//
// Run: go run ./02-concurrency/04-context/03-deadline/01-basic
package main

import (
	"context"
	"fmt"
	"time"
)

// roundTo rounds the printed time left, to hide timer noise.
const roundTo = 10 * time.Millisecond

func main() {
	const (
		workerCount = 3                      // must match the number of "go worker" lines below
		budget      = 300 * time.Millisecond // time until the deadline
		quickNeeds  = 100 * time.Millisecond
		bigNeeds    = 500 * time.Millisecond // more than budget: skipped
		tinyNeeds   = 50 * time.Millisecond
	)

	deadline := time.Now().Add(budget)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	done := make(chan struct{}) // workers -> main: "I'm finished"

	go worker(ctx, "quick", quickNeeds, done)
	go worker(ctx, "big", bigNeeds, done)
	go worker(ctx, "tiny", tinyNeeds, done)

	for range workerCount {
		<-done // wait for every worker
	}
}

// worker does a job that needs some time, or skips it if there isn't enough left.
func worker(ctx context.Context, name string, needs time.Duration, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	if d, ok := ctx.Deadline(); ok && time.Until(d) < needs { // not enough time: don't even start
		left := time.Until(d).Round(roundTo)
		fmt.Printf("%s: skipped, needs %v but only %v left\n", name, needs, left)
		return
	}

	select {
	case <-time.After(needs): // simulate the work
		fmt.Println(name + ": done")
	case <-ctx.Done(): // the deadline passed while working
		fmt.Println(name+": stopped:", ctx.Err())
	}
}
