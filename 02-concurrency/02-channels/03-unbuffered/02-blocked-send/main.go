// Blocked send: a send on an unbuffered channel waits for a receiver.
//
// Problem: nobody receives from ch. worker is stuck on the send,
// and main is stuck on <-done. Every goroutine is blocked, so Go stops
// with "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: receive the value in main. Uncomment the "fix:" line.
//
// Real use:
//   - a handler returns early on an error and never reads the result. In a
//     server this doesn't deadlock, it leaks a goroutine on every request
//     (see 05-leak).
//
// Run: go run ./02-concurrency/02-channels/03-unbuffered/02-blocked-send
package main

import "fmt"

func main() {
	ch := make(chan int)
	done := make(chan struct{}) // worker -> main: "I'm finished"

	go worker(ch, done)

	// fmt.Println(<-ch) // fix: receive the value

	<-done // stuck here: worker never finishes its send
}

// worker sends one value, then tells main it's finished.
func worker(out chan<- int, done chan<- struct{}) {
	defer close(done) // tells main we're finished (see channels/02-closing)
	out <- 5          // blocks: no receiver
	fmt.Println("sent")
}
