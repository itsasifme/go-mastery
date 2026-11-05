package main

import "fmt"

func sendData(ch chan<- int) {
	ch <- 42 // Sending data to the channel
}

func main() {
	ch := make(chan int) // Declaring a bidirectional channel

	sendCh := (chan<- int)(ch) // Converting to send-only channel
	go sendData(sendCh)        // Sending data to the channel

	readCh := (<-chan int)(ch) // Converting to receive-only channel
	x := <-readCh              // Receiving data from the channel
	fmt.Println("x: ", x)
}
