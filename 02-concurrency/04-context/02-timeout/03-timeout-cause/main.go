// Timeout cause: say why a timeout happened with WithTimeoutCause.
//
// ctx.Err() after a timeout is always "context deadline exceeded", for
// every timeout in the program. WithTimeoutCause takes an error, and
// context.Cause(ctx) returns it once the time runs out. So the caller
// knows which limit was hit, and can check it with errors.Is.
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

	query(ctx)

	fmt.Println("Err:  ", ctx.Err())          // the same for every timeout
	fmt.Println("Cause:", context.Cause(ctx)) // which timeout it was
	fmt.Println("slow DB:", errors.Is(context.Cause(ctx), errSlowDB))
}

// query stands in for a database call that takes 2 seconds.
func query(ctx context.Context) {
	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
	}
}
