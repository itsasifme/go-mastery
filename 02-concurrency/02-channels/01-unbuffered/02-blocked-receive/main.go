package main

import "fmt"

func main() {
	ch := make(chan int) // Creating an unbuffered channel
	fmt.Println("Hey")

	go func() {
		fmt.Println("Hello")
		// No sender is sending a value to the channel
		x := <-ch // Receiving a value from the channel
		fmt.Println("x: ", x)
	}()

	// The program will deadlock here because there is no sender to send a value
}
