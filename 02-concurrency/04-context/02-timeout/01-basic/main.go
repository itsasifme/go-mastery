// Timeout: stop a slow function with context.WithTimeout.
//
// slowTask needs 2 seconds, but ctx times out after 500ms. slowTask takes
// ctx as its first parameter and selects on ctx.Done(), so it stops
// itself and returns ctx.Err(): "context deadline exceeded".
//
// Compare with select/04-timeout: there, main stopped waiting, but the
// slow work kept running. Here the work itself stops.
//
// Try: change the timeout to 3 * time.Second. slowTask finishes and
// returns nil after 2s.
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

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel() // still needed, even when the timeout fires

	start := time.Now()
	err := slowTask(ctx)
	fmt.Println("took:", time.Since(start).Round(100*time.Millisecond))
	fmt.Println("err:", err)
	fmt.Println("timed out:", errors.Is(err, context.DeadlineExceeded))
}

// slowTask takes 2 seconds, unless ctx stops it first.
func slowTask(ctx context.Context) error {
	select {
	case <-time.After(2 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
