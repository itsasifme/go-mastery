// Receive-only channel: a <-chan int can only receive.
//
// Problem: this channel is made receive-only, so nothing can ever send
// on it. readData gets stuck, and main ends without printing a value.
//
// Fix: make a normal channel, send from a goroutine, and read in main.
// Comment out the 2 lines marked "problem", then uncomment the "fix:"
// lines.
//
// Run: go run ./02-concurrency/02-channels/01-directional/03-receive-only
package main

import "fmt"

func readData(ch <-chan int) {
	x := <-ch // ok: receive
	// ch <- 1 // compile error: can't send to a receive-only channel
	fmt.Println("x:", x)
}

func main() {
	ch := make(<-chan int) // problem: receive-only, so nobody can send to it
	go readData(ch)        // problem: gets stuck on the receive

	// ch := make(chan int)    // fix: normal channel
	// go func() { ch <- 5 }() // fix: send from a goroutine
	// readData(ch)            // fix: read in main, so main waits for it

	fmt.Println("main done")
}
