// Deadline: stop at a fixed time with context.WithDeadline.
//
// WithDeadline takes a point in time. WithTimeout(d) is the same as
// WithDeadline(time.Now().Add(d)).
//
// ctx.Deadline() tells you when the deadline is, so run can skip a job
// that can't finish in time, instead of starting it and failing halfway.
// The big job needs 500ms but only about 200ms are left, so it's skipped.
//
// Real use:
//   - a request with a 2 second budget skips an optional slow step, like
//     loading recommendations, instead of failing the whole page.
//   - gRPC sends the deadline to the next service, so it knows how much
//     time is left too.
//
// Run: go run ./02-concurrency/04-context/03-deadline/01-basic
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	deadline := time.Now().Add(300 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	run(ctx, "quick", 100*time.Millisecond)
	run(ctx, "big", 500*time.Millisecond)
	run(ctx, "tiny", 50*time.Millisecond)
}

// run does a job that needs some time, or skips it if there isn't enough left.
func run(ctx context.Context, name string, needs time.Duration) {
	if d, ok := ctx.Deadline(); ok && time.Until(d) < needs {
		left := time.Until(d).Round(10 * time.Millisecond)
		fmt.Printf("%s: skipped, needs %v but only %v left\n", name, needs, left)
		return
	}

	select {
	case <-time.After(needs): // simulate the work
		fmt.Println(name + ": done")
	case <-ctx.Done():
		fmt.Println(name+": stopped:", ctx.Err())
	}
}
