// Result or timeout: wait for a goroutine's result, but give up in time.
//
// main selects on the result and on ctx.Done(). fetch watches the same
// ctx, so when main gives up, fetch stops too. Its send also selects on
// ctx.Done(), so it never gets stuck sending a result nobody reads. That
// works even with an unbuffered channel (see select/07-nil-or-stop).
//
// main waits on done at the end, so a stuck fetch would show up as a
// deadlock instead of a silent leak.
//
// Try: change the timeout to 3 * time.Second. Now the result wins.
//
// Real use:
//   - a handler starts a call to another service in the background and
//     waits for it, but answers the user with an error after 1 second.
//
// Run: go run ./02-concurrency/04-context/02-timeout/02-result-or-timeout
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result := make(chan string)
	done := make(chan struct{})
	go fetch(ctx, result, done)

	select {
	case r := <-result:
		fmt.Println("main: got", r)
	case <-ctx.Done():
		fmt.Println("main: gave up:", ctx.Err())
	}

	<-done // proves fetch stopped too
}

// fetch takes 2 seconds to get a result, unless ctx stops it first.
func fetch(ctx context.Context, result chan<- string, done chan<- struct{}) {
	defer close(done)

	select {
	case <-time.After(2 * time.Second): // simulate slow work
	case <-ctx.Done():
		fmt.Println("fetch: cancelled, stopping")
		return
	}

	select {
	case result <- "data":
	case <-ctx.Done(): // main already gave up: don't wait for it forever
	}
}
