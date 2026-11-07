// Blocked send: a send on an unbuffered channel waits for a receiver.
//
// Problem: nobody receives from ch. The goroutine is stuck on the send,
// and main is stuck on <-done. Every goroutine is blocked, so Go stops
// with "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: receive the value in main. Uncomment the "fix:" line.
//
// Run: go run ./02-concurrency/02-channels/01-unbuffered/02-blocked-send
package main

import "fmt"

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		defer close(done)
		ch <- 5 // blocks: no receiver
		fmt.Println("sent")
	}()

	// fmt.Println(<-ch) // fix: receive the value

	<-done
}
