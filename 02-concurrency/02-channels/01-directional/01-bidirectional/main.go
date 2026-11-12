// Bidirectional channel: a normal chan can both send and receive.
//
// sendData and readData take the same chan int. One sends, the other
// receives.
//
// Note: time.Sleep is a hack to wait for the goroutines. The real fix is
// a WaitGroup (see waitgroup/01-basic).
//
// Real use:
//   - inside the one function that creates a channel and wires producers
//     and consumers together. Everywhere else, pass a directional type.
//
// Run: go run ./02-concurrency/02-channels/01-directional/01-bidirectional
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int) // can send and receive

	go sendData(ch)
	go readData(ch)

	time.Sleep(time.Second) // hack: wait for the goroutines
}

// sendData sends one value on ch.
func sendData(ch chan int) {
	ch <- 42
}

// readData receives one value from ch and prints it.
func readData(ch chan int) {
	x := <-ch
	fmt.Println("x:", x)
}
