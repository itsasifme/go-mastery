package main

import (
	"fmt"
	"time"
)

func main() {
	// Creating an unbuffered channel
	ch := make(chan int)

	go func() {
		time.Sleep(time.Second) // Simulating some work
		ch <- 5                 // Sending a value to the channel
	}()

	x := <-ch      // Receiving the value from the channel
	fmt.Println(x) // Output: 5
}
