// Send-only channel: a chan<- int can only send.
//
// Problem: this channel is made send-only, so nothing can ever receive
// from it. sendData gets stuck, and main ends without getting a value.
//
// Fix: make a normal channel, pass it in as send-only, and receive in
// main. Comment out the 2 lines marked "problem", then uncomment the
// "fix:" lines.
//
// Run: go run ./02-concurrency/02-channels/03-directional/02-send-only
package main

import "fmt"

func sendData(ch chan<- int) {
	ch <- 5 // ok: send
	// <-ch // compile error: can't receive from a send-only channel
}

func main() {
	ch := make(chan<- int) // problem: send-only, so nobody can receive from it
	go sendData(ch)        // problem: gets stuck on the send

	// ch := make(chan int)    // fix: normal channel
	// go sendData(ch)         // fix: Go passes it to sendData as chan<- int
	// fmt.Println("x:", <-ch) // fix: receive in main

	fmt.Println("main done")
}
