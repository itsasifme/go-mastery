// Result or timeout: wait for goroutine results, but give up in time.
//
// main starts two fetches that share one ctx with a timeout, and selects
// on the results and on ctx.Done():
//
//	fast fetch   finishes before the timeout   main gets the data
//	slow fetch   still working at the timeout  main gives up, fetch stops too
//
// fetch watches the same ctx, so when main gives up, fetch stops too. Its
// send also selects on ctx.Done(), so it never gets stuck sending a result
// nobody reads. That works even with an unbuffered channel (see
// select/07-nil-or-stop).
//
// main waits on done for every fetch, so a stuck fetch would show up as a
// deadlock instead of a silent leak.
//
// Try: change the timeout to 3 * time.Second. Now both fetches send data.
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
	const fetchCount = 2 // must match the number of "go fetch" lines below

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result := make(chan string) // fetches -> main: the data
	done := make(chan struct{}) // fetches -> main: "I'm finished"

	go fetch(ctx, "fast fetch", 200*time.Millisecond, result, done) // result wins
	go fetch(ctx, "slow fetch", 2*time.Second, result, done)        // timeout wins

	for range fetchCount { // one answer per fetch: its data, or the timeout
		select {
		case r := <-result: // a fetch finished in time
			fmt.Println("main: got", r)
		case <-ctx.Done(): // the timeout fired before this result came
			fmt.Println("main: gave up:", ctx.Err())
		}
	}

	for range fetchCount {
		<-done // proves every fetch stopped too
	}
}

// fetch takes "needs" to get a result, unless ctx stops it first.
func fetch(ctx context.Context, name string, needs time.Duration, result chan<- string, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	select {
	case <-time.After(needs): // simulate slow work
	case <-ctx.Done(): // main gave up while we were working
		fmt.Println(name + ": cancelled, stopping")
		return
	}

	select {
	case result <- "data from " + name: // main is waiting: hand it over
	case <-ctx.Done(): // main already gave up: don't wait for it forever
	}
}
