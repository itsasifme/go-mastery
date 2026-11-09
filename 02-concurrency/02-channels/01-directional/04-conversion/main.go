// Channel conversion: give each side of one channel only the access it needs.
//
// One normal channel ch is shared by a sender and a reader. sendData gets
// it as send-only (chan<- int), and readData as receive-only (<-chan int).
//
// Part 1 converts by hand. sendCh and readCh aren't copies: they print the
// same address as ch, and each value is still received only once.
//
// Part 2 writes no conversion at all. Go converts ch automatically when
// it's passed to a chan<- int or <-chan int parameter. Real code does this.
//
// Try: uncomment a "compile error" line. go build refuses it.
//
// Real use:
//   - producers and consumers declared as func(out chan<- T) and
//     func(in <-chan T), so neither side can misuse the channel.
//   - a constructor that returns only the receive side, like a file
//     watcher returning <-chan Event.
//
// Run: go run ./02-concurrency/02-channels/01-directional/04-conversion
package main

import "fmt"

// count is how many values sendData sends and readData reads.
const count = 3

func main() {
	// Part 1: convert by hand.
	ch := make(chan int)
	sendCh := (chan<- int)(ch) // send-only view of ch
	readCh := (<-chan int)(ch) // receive-only view of the same ch
	fmt.Printf("same channel: %p %p %p\n", ch, sendCh, readCh)

	go sendData(sendCh)
	readData(readCh)

	// Part 2: let Go convert for you.
	fmt.Println("--- automatic conversion")
	go sendData(ch) // chan int -> chan<- int, automatic
	readData(ch)    // chan int -> <-chan int, automatic
}

// sendData can only send on out.
func sendData(out chan<- int) {
	for i := range count {
		out <- i + 1
	}
	// <-out // compile error: can't receive from a send-only channel
}

// readData can only receive from in.
func readData(in <-chan int) {
	for range count {
		fmt.Println("got", <-in)
	}
	// in <- 1 // compile error: can't send to a receive-only channel
}
