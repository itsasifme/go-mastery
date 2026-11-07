// Blocked receive: a receive on an unbuffered channel waits for a sender.
//
// Problem: nobody sends on ch. The goroutine is stuck on the receive,
// and main is stuck on <-done. Every goroutine is blocked, so Go stops
// with "fatal error: all goroutines are asleep - deadlock!".
//
// Fix: send a value from main. Uncomment the "fix:" line.
//
// Run: go run ./02-concurrency/02-channels/01-unbuffered/03-blocked-receive
package main

import "fmt"

func main() {
	ch := make(chan int)
	done := make(chan struct{})

	go func() {
		defer close(done)
		x := <-ch // blocks: no sender
		fmt.Println("x:", x)
	}()

	// ch <- 5 // fix: send a value

	<-done
}
