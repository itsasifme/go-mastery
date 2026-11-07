// Bidirectional channel: a normal chan can both send and receive.
//
// sendData and readData take the same chan int. One sends, the other
// receives.
//
// Note: time.Sleep is a hack to wait for the goroutines. The real fix is
// a WaitGroup (see 05-waitgroup).
//
// Run: go run ./02-concurrency/02-channels/03-directional/01-bidirectional
package main

import (
	"fmt"
	"time"
)

func sendData(ch chan int) {
	ch <- 42
}

func readData(ch chan int) {
	x := <-ch
	fmt.Println("x:", x)
}

func main() {
	ch := make(chan int) // can send and receive

	go sendData(ch)
	go readData(ch)

	time.Sleep(time.Second) // hack: wait for the goroutines
}
