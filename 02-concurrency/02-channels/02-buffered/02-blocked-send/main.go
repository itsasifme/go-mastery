// Blocked send: a send on a buffered channel waits when the buffer is full.
//
// Problem: 2 sends fit in the buffer, but the 3rd has no room and nobody
// receives. Every goroutine is blocked, so Go stops with
// "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: receive one value in main to free a slot. Uncomment the "fix:" line.
//
// Run: go run ./02-concurrency/02-channels/02-buffered/02-blocked-send
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	done := make(chan struct{})

	go func() {
		defer close(done)
		ch <- 1
		ch <- 2
		ch <- 3 // blocks: buffer full
		fmt.Println("sent 3")
	}()

	// fmt.Println(<-ch) // fix: receive one value to free a slot

	<-done
}
