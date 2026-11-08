// Channel conversion: turn a normal channel into send-only or receive-only.
//
// This is how directional channels are used in real code: make one normal
// channel, then give each side only the access it needs.
//
// Note: the conversions here are optional. Go converts automatically when
// you pass ch to a function that takes chan<- int or <-chan int.
//
// Run: go run ./02-concurrency/02-channels/01-directional/04-conversion
package main

import "fmt"

func sendData(ch chan<- int) {
	ch <- 42
}

func main() {
	ch := make(chan int) // normal channel

	sendCh := (chan<- int)(ch) // send-only view
	go sendData(sendCh)

	readCh := (<-chan int)(ch) // receive-only view
	x := <-readCh
	fmt.Println("x:", x)
}
