// Stop workers: one cancel call stops every goroutine that shares the ctx.
//
// This is channels/02-closing/02-broadcast with a context instead of a stop
// channel. Every worker gets the same ctx. cancel closes ctx.Done() once,
// and a closed channel wakes everyone waiting on it, so every worker stops.
//
//	broadcast                 this example
//	close(stop)               cancel()
//	<-stop                    case <-ctx.Done():        (in a select, so it can work while it waits)
//	for range workerCount {   for range workerCount {   (same: one done per worker)
//		<-done                    <-done
//	}                         }
//
// Real use:
//   - graceful shutdown: signal.NotifyContext gives you a ctx that's
//     cancelled on Ctrl+C or SIGTERM. Pass it to every worker, and they
//     all stop together.
//
// Run: go run ./02-concurrency/04-context/01-cancellation/02-stop-workers
package main

import (
	"context"
	"fmt"
	"time"
)

const (
	stepTime = 100 * time.Millisecond // how long one step of work takes
	runFor   = 250 * time.Millisecond // how long main lets the workers run
)

func main() {
	const workerCount = 3

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{}) // workers -> main: "I'm finished"

	for id := range workerCount {
		go worker(ctx, id, done)
	}

	time.Sleep(runFor)
	cancel() // one call stops every worker

	for range workerCount {
		<-done // one per worker
	}
	fmt.Println("all workers stopped")
}

// worker works every 100ms until ctx is cancelled.
func worker(ctx context.Context, id int, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	for {
		select {
		case <-ctx.Done():
			fmt.Println("worker", id, "stopped")
			return
		case <-time.After(stepTime):
			fmt.Println("worker", id, "working")
		}
	}
}
