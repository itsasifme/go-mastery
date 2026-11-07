// Blocked receive: a receive on a buffered channel waits when the buffer is empty.
//
// Problem: the 1st receive gets the value in the buffer, but the 2nd finds
// it empty and nobody sends. Every goroutine is blocked, so Go stops with
// "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: send one more value from main. Uncomment the "fix:" line.
//
// Run: go run ./02-concurrency/02-channels/02-buffered/03-blocked-receive
package main

import "fmt"

func main() {
	ch := make(chan int, 2)
	done := make(chan struct{})
	ch <- 1

	go func() {
		defer close(done)
		fmt.Println(<-ch) // 1: buffer has a value
		fmt.Println(<-ch) // blocks: buffer empty
	}()

	// ch <- 2 // fix: send another value

	<-done
}
