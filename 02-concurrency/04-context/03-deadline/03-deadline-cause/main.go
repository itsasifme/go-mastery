// Deadline cause: say why a deadline stopped the work with WithDeadlineCause.
//
// A batch runs jobs until a fixed time, when maintenance starts. When the
// deadline hits, ctx.Err() only says "context deadline exceeded".
// context.Cause(ctx) returns the error given to WithDeadlineCause, so the
// batch can report the real reason it stopped.
//
// Jobs take 100ms each and the deadline is in 250ms, so 2 jobs finish.
//
// Note: like WithTimeoutCause, the cancel function takes no cause.
//
// Real use:
//   - a nightly job that must stop before a maintenance window, and logs
//     "maintenance window starts" instead of a generic timeout.
//   - a request deadline set by a client, reported as "client deadline
//     reached" so it's not confused with your own timeouts.
//
// Run: go run ./02-concurrency/04-context/03-deadline/03-deadline-cause
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var errMaintenance = errors.New("maintenance window starts")

func main() {
	maintenanceAt := time.Now().Add(250 * time.Millisecond)
	ctx, cancel := context.WithDeadlineCause(context.Background(), maintenanceAt, errMaintenance)
	defer cancel()

	runBatch(ctx)
}

// runBatch runs one job every 100ms until ctx is done.
func runBatch(ctx context.Context) {
	for i := 1; ; i++ {
		select {
		case <-time.After(100 * time.Millisecond): // simulate a job
			fmt.Println("job", i, "done")
		case <-ctx.Done():
			fmt.Println("Err:  ", ctx.Err())
			fmt.Println("Cause:", context.Cause(ctx))
			return
		}
	}
}
