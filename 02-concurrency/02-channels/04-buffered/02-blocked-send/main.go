// Blocked send: a send on a buffered channel waits when the buffer is full.
//
// Problem: 2 sends fit in the buffer, but the 3rd has no room and nobody
// receives. Every goroutine is blocked, so Go stops with
// "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: receive one value in main to free a slot. Uncomment the "fix:" line.
//
// Real use:
//   - a full job queue makes the HTTP handler wait. That's backpressure:
//     sometimes what you want, sometimes a hidden slowdown (see
//     patterns/05-backpressure).
//
// Run: go run ./02-concurrency/02-channels/04-buffered/02-blocked-send
package main

import "fmt"

func main() {
	const bufferSize = 2 // worker sends one more value than this

	ch := make(chan int, bufferSize)
	done := make(chan struct{}) // worker -> main: "I'm finished"

	go worker(ch, done)

	// fmt.Println(<-ch) // fix: receive one value to free a slot

	<-done // stuck here: worker never finishes its send
}

// worker sends 3 values into a buffer that only has room for 2.
func worker(out chan<- int, done chan<- struct{}) {
	defer close(done) // tells main we're finished (see channels/02-closing)
	out <- 1
	out <- 2
	out <- 3 // blocks: buffer full
	fmt.Println("sent 3")
}
