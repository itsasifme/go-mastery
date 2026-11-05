package main

import (
	"fmt"
	"time"
)

func sendData(ch chan int) {
	ch <- 42 // Sending data to the channel
}

func readData(ch chan int) {
	x := <-ch // Receiving data from the channel
	fmt.Println("x: ", x)
}

func main() {
	ch := make(chan int) // Declaring a bidirectional channel

	go sendData(ch) // Sending data to the channel
	go readData(ch) // Reading data from the channel

	time.Sleep(time.Second) // Wait for goroutines to complete
}
