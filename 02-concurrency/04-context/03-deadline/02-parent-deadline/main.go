// Parent deadline: a child can make the deadline earlier, but never later.
//
// The parent's deadline is 1 second from now. One child asks for 5
// seconds, but it still stops at the parent's deadline. Another child asks
// for 200ms and gets it, because that's earlier. The earliest deadline
// always wins.
//
// A worker does steps on each context and prints when it stopped and how
// many steps it did, so you can see the real stop times, not just the
// numbers that were asked for.
//
// Real use:
//   - a server gives each request 2 seconds. A database call inside asks
//     for 5 seconds, but still stops at the request's 2 seconds. A cache
//     lookup asks for 50ms, and gets it.
//
// Run: go run ./02-concurrency/04-context/03-deadline/02-parent-deadline
package main

import (
	"context"
	"fmt"
	"time"
)

const (
	stepTime = 50 * time.Millisecond  // how long one step of work takes
	roundTo  = 100 * time.Millisecond // round printed times to hide timer noise
)

func main() {
	const (
		workerCount = 3                      // must match the number of "go worker" lines below
		parentLimit = time.Second            // the parent's deadline, from start
		laterLimit  = 5 * time.Second        // later than the parent: capped
		soonerLimit = 200 * time.Millisecond // earlier than the parent: kept
	)

	start := time.Now() // every deadline below is counted from here

	parent, cancel := context.WithDeadline(context.Background(), start.Add(parentLimit))
	defer cancel()

	later, cancelLater := context.WithDeadline(parent, start.Add(laterLimit))
	defer cancelLater()

	sooner, cancelSooner := context.WithDeadline(parent, start.Add(soonerLimit))
	defer cancelSooner()

	done := make(chan struct{}) // workers -> main: "I'm finished"

	go worker(parent, "parent", start, done)
	go worker(later, fmt.Sprint("child asked ", laterLimit), start, done)
	go worker(sooner, fmt.Sprint("child asked ", soonerLimit), start, done)

	for range workerCount {
		<-done // wait for every worker
	}
}

// worker does one step every stepTime until its ctx's deadline passes, then
// says when that was and how many steps it did.
func worker(ctx context.Context, name string, start time.Time, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	steps := 0
	for {
		select {
		case <-time.After(stepTime): // simulate one step of work
			steps++
		case <-ctx.Done():
			took := time.Since(start).Round(roundTo)
			fmt.Printf("%-18s stopped after %v, %d steps: %v\n", name+":", took, steps, ctx.Err())
			return
		}
	}
}
