// Timeout cause: say why a timeout happened with WithTimeoutCause.
//
// ctx.Err() after a timeout is always "context deadline exceeded", for
// every timeout in the program. WithTimeoutCause takes an error, and
// context.Cause(ctx) returns it once the time runs out. query reads rows
// in its own goroutine. When the timeout stops it, it reports how far it
// got and which limit was hit, and checks the cause with errors.Is.
//
// Note: the cancel function it returns takes no cause. If you call
// cancel before the timeout, Cause is just context.Canceled.
//
// Try: call cancel() right after creating ctx. Cause becomes "context
// canceled", not errSlowDB.
//
// Real use:
//   - logs and metrics that say "database too slow" or "payment API too
//     slow", instead of the same "deadline exceeded" for every timeout.
//   - returning a clear error to the user, like 504 for a slow upstream.
//
// Run: go run ./02-concurrency/04-context/02-timeout/03-timeout-cause
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errSlowDB = errors.New("database took longer than 500ms")

func main() {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 500*time.Millisecond, errSlowDB)
	defer cancel()

	done := make(chan struct{}) // query -> main: "I'm finished"

	go query(ctx, done)

	<-done // wait until query stopped
}

// query reads rows one by one until it has them all, or ctx stops it.
func query(ctx context.Context, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	const totalRows = 20 // at 100ms per row, the whole query needs 2s

	rows := 0
	for rows < totalRows {
		select {
		case <-time.After(100 * time.Millisecond): // simulate reading one row
			rows++
		case <-ctx.Done(): // the timeout fired first
			fmt.Println("query: read", rows, "of", totalRows, "rows")
			fmt.Println("Err:  ", ctx.Err())          // the same for every timeout
			fmt.Println("Cause:", context.Cause(ctx)) // which timeout it was
			fmt.Println("slow DB:", errors.Is(context.Cause(ctx), errSlowDB))
			return
		}
	}

	fmt.Println("query: read all", totalRows, "rows")
}
