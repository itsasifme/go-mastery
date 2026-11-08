// Timeout: stop waiting after a fixed time with time.After.
//
// fetch takes 2 seconds, but main only waits 500ms. time.After returns a
// channel that gets a value once the time is up, so select can wait on
// "the result" and "the deadline" at the same time.
//
// The result channel has a buffer of 1. When main gives up, fetch can
// still send its late result and exit, so it doesn't leak.
//
// Try: change timeout to 3 * time.Second. Now the result wins.
//
// Note: real code usually uses context.WithTimeout instead, so the slow
// work can be cancelled too (see 04-context).
//
// Run: go run ./02-concurrency/03-select/04-timeout
package main

import (
	"fmt"
	"time"
)

func main() {
	timeout := 500 * time.Millisecond
	result := make(chan string, 1)

	go fetch(2*time.Second, result)

	select {
	case r := <-result:
		fmt.Println("got:", r)
	case <-time.After(timeout):
		fmt.Println("timeout: gave up after", timeout)
	}
}

// fetch pretends to do slow work, then sends the result.
func fetch(delay time.Duration, out chan<- string) {
	time.Sleep(delay)
	out <- "data"
}
