// Cancel cause: say why a context was cancelled, and act on it.
//
// ctx.Err() only ever says "context canceled". With WithCancelCause, the
// cancel function takes an error, and context.Cause(ctx) returns it.
//
// An export job can be stopped for two reasons, and each needs different
// cleanup. On a server shutdown it saves its progress, so it can resume
// after the restart. When the user cancels it, it deletes the partial file.
// ctx.Err() can't tell these apart. The cause can.
//
// Try: replace errShutdown with errUserCancelled in main.
//
// Note: WithTimeoutCause and WithDeadlineCause do the same for timeouts.
//
// Real use:
//   - a SIGTERM handler cancels running jobs with "server shutting down",
//     so they save their progress instead of throwing it away.
//   - a DELETE /exports/{id} handler finds that job's cancel function and
//     calls it with "cancelled by user".
//   - logs that say why a job stopped, instead of just "context canceled".
//
// Run: go run ./02-concurrency/04-context/01-cancellation/03-cancel-cause
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	errShutdown      = errors.New("server shutting down")
	errUserCancelled = errors.New("cancelled by user")
)

const (
	partTime = 30 * time.Millisecond  // how long export needs for one part
	runFor   = 100 * time.Millisecond // how long main lets export run
)

func main() {
	ctx, cancel := context.WithCancelCause(context.Background())
	done := make(chan struct{}) // export -> main: "I'm finished"

	go export(ctx, done)

	time.Sleep(runFor) // let export write a few parts
	cancel(errShutdown)
	// cancel(errUserCancelled) // try this instead of errShutdown to see the other cleanup

	<-done // wait until export cleaned up
}

// export writes the file part by part until ctx is cancelled, then cleans
// up depending on why.
func export(ctx context.Context, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	parts := 0
	for {
		select {
		case <-time.After(partTime): // simulate writing one part
			parts++
			fmt.Println("export: wrote part", parts)
		case <-ctx.Done():
			fmt.Println("Err:  ", ctx.Err())          // always "context canceled"
			fmt.Println("Cause:", context.Cause(ctx)) // the real reason

			switch cause := context.Cause(ctx); {
			case errors.Is(cause, errShutdown):
				fmt.Println("export: saving progress after part", parts, "to resume after restart")
			case errors.Is(cause, errUserCancelled):
				fmt.Println("export: deleting the partial file with", parts, "parts")
			}
			return
		}
	}
}
