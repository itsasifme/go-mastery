// Goroutine leak: stuck goroutines don't cause any error.
//
// Problem: main takes the first result and moves on. The other workers
// wait forever. Go doesn't complain, because main is still running.
//
// Fix: give the channel room for every worker, so they can all finish.
// Change the channel as shown in the "fix:" comment.
//
// Real use:
//   - "first response wins": ask several replicas or mirrors and use the
//     fastest. With an unbuffered channel, the slower ones leak on every
//     request, and memory grows until the server crashes.
//
// Run: go run ./02-concurrency/02-channels/03-unbuffered/05-leak
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	const workerCount = 3

	ch := make(chan int) // fix: make(chan int, workerCount) gives every worker room

	for i := range workerCount {
		go worker(ch, i)
	}

	fmt.Println("first:", <-ch)

	time.Sleep(100 * time.Millisecond)                 // give the others time to reach their send
	fmt.Println("goroutines:", runtime.NumGoroutine()) // main + (workerCount - 1) stuck. With the fix: just main
}

// worker sends one value. It's stuck forever if nobody receives it.
func worker(out chan<- int, v int) {
	out <- v
}
