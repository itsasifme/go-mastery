// Cancellation: stop a goroutine with context.WithCancel.
//
// ctx.Done() is a channel that's closed when cancel is called, just like
// the stop channel in select/05-for-select. The worker selects on it, and
// ctx.Err() says why it stopped: "context canceled".
//
// Try: comment out the ctx.Done() case. The worker never stops, main waits
// on done forever, and you have to stop it with Ctrl+C. Go can't stop a
// goroutine for you: it has to check ctx.Done() itself.
//
// Real use:
//   - a user cancels a long export or upload, and the server stops the
//     work behind it.
//   - every net/http request has r.Context(), which is cancelled when the
//     client disconnects.
//
// Run: go run ./02-concurrency/04-context/01-cancellation/01-basic
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go worker(ctx, done)

	time.Sleep(350 * time.Millisecond)
	cancel() // closes ctx.Done()
	<-done
}

// worker does one step every 100ms until ctx is cancelled.
func worker(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	for i := 1; ; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("worker: stopped:", ctx.Err())
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Println("worker: step", i)
		}
	}
}
