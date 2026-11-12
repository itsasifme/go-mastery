// Deadline cause: say why a deadline stopped the work with WithDeadlineCause.
//
// batch runs jobs in a goroutine until a fixed time, when maintenance
// starts. When the deadline hits, ctx.Err() only says "context deadline
// exceeded". context.Cause(ctx) returns the error given to
// WithDeadlineCause, so batch can report the real reason it stopped.
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

const (
	untilMaintenance = 250 * time.Millisecond // time until the deadline
	jobTime          = 100 * time.Millisecond // how long one job takes
)

func main() {
	maintenanceAt := time.Now().Add(untilMaintenance)
	ctx, cancel := context.WithDeadlineCause(context.Background(), maintenanceAt, errMaintenance)
	defer cancel()

	done := make(chan struct{}) // batch -> main: "I'm finished"

	go batch(ctx, done)

	<-done // wait until batch stopped
}

// batch runs one job every jobTime until ctx is done.
func batch(ctx context.Context, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	for i := 1; ; i++ { // no end: only the deadline stops it
		select {
		case <-time.After(jobTime): // simulate a job
			fmt.Println("job", i, "done")
		case <-ctx.Done(): // maintenance starts
			fmt.Println("Err:  ", ctx.Err())          // the same for every deadline
			fmt.Println("Cause:", context.Cause(ctx)) // the real reason
			return
		}
	}
}
