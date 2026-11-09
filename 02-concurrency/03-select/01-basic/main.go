// Basic select: wait on two channels and take whichever is ready first.
//
// Two lookups run at the same time. One takes 100ms, the other 300ms.
// select waits on both channels and runs the case that's ready first, so
// the fast one always wins.
//
// The channels have a buffer of 1, so the slow lookup can still send its
// answer and exit after main stops listening (no goroutine leak, see
// channels/04-buffered/01-basic).
//
// Real use:
//   - query two replicas or caches and use whichever answers first.
//   - wait for a result and a shutdown signal at the same time.
//
// Run: go run ./02-concurrency/03-select/01-basic
package main

import (
	"fmt"
	"time"
)

func main() {
	fast := make(chan string, 1)
	slow := make(chan string, 1)

	go lookup("fast server", 100*time.Millisecond, fast)
	go lookup("slow server", 300*time.Millisecond, slow)

	select { // waits until one case is ready, then runs only that one
	case r := <-fast:
		fmt.Println("got:", r)
	case r := <-slow:
		fmt.Println("got:", r)
	}
}

// lookup pretends to ask a server, then sends the answer.
func lookup(name string, delay time.Duration, out chan<- string) {
	time.Sleep(delay)
	out <- name
}
