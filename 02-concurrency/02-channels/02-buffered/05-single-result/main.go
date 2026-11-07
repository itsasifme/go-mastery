// Single result: a buffer of 1 lets a goroutine finish even if nobody reads.
//
// fetch starts a slow job and returns a channel for its result. main waits
// only 100ms, then gives up. The job finishes later and tries to send.
//
// Problem: the channel is unbuffered and nobody receives any more, so the
// job's goroutine is stuck on the send forever: a goroutine leak.
//
// Fix: give the channel a buffer of 1. The send always has room, so the
// goroutine finishes and exits. Change the channel as shown in the "fix:"
// comment.
//
// Note: this is why Uber's guide says a channel should be size 0 or 1.
//
// Run: go run ./02-concurrency/02-channels/02-buffered/05-single-result
package main

import (
	"fmt"
	"runtime"
	"time"
)

func fetch() <-chan string {
	ch := make(chan string) // fix: make(chan string, 1)

	go func() {
		time.Sleep(300 * time.Millisecond) // slow work
		ch <- "result"                     // stuck if nobody reads and there's no room
	}()

	return ch
}

func main() {
	// select waits for whichever is ready first (see 03-select)
	select {
	case r := <-fetch():
		fmt.Println("got:", r)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout: gave up waiting")
	}

	time.Sleep(500 * time.Millisecond)                 // give the job time to finish
	fmt.Println("goroutines:", runtime.NumGoroutine()) // 2 = main + stuck job (1 with the fix)
}
