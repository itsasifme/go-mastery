// Timeout: stop a slow function with context.WithTimeout.
//
// query reads rows in its own goroutine. Reading every row takes 2
// seconds, but ctx times out after 500ms. query takes ctx as its first
// parameter and checks ctx.Done() between rows, so it stops itself and
// reports ctx.Err(): "context deadline exceeded".
//
// Compare with select/04-timeout: there, main stopped waiting, but the
// slow work kept running. Here the work itself stops.
//
// Try: change the timeout to 3 * time.Second. query reads every row and
// finishes after 2s with no error.
//
// Real use:
//   - every slow call takes a ctx: db.QueryContext, http.NewRequestWithContext,
//     gRPC calls. Give each one a time limit, so a slow database or API
//     can't hang your server.
//
// Run: go run ./02-concurrency/04-context/02-timeout/01-basic
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	timeout   = 500 * time.Millisecond
	rowTime   = 100 * time.Millisecond // how long one row takes to read
	totalRows = 20                     // the whole query needs totalRows * rowTime
	roundTo   = 100 * time.Millisecond // round printed times to hide timer noise
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // still needed, even when the timeout fires

	done := make(chan struct{}) // query -> main: "I'm finished"

	go query(ctx, done)

	<-done // wait until query stopped
}

// query reads rows one by one until it has them all, or ctx stops it.
func query(ctx context.Context, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	start := time.Now()
	rows := 0
	for rows < totalRows {
		select {
		case <-time.After(rowTime): // simulate reading one row
			rows++
		case <-ctx.Done(): // closed when the timeout fires
			err := ctx.Err()
			fmt.Println("query: read", rows, "of", totalRows, "rows")
			fmt.Println("took:", time.Since(start).Round(roundTo))              // round off the noise: 500.7ms -> 500ms
			fmt.Println("err:", err)                                            // context deadline exceeded
			fmt.Println("timed out:", errors.Is(err, context.DeadlineExceeded)) // the right way to check for a timeout
			return
		}
	}

	fmt.Println("query: read all", totalRows, "rows in", time.Since(start).Round(roundTo))
}
