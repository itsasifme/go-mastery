package main

import "fmt"

func main() {
	ch := make(chan int, 2) // Creating a buffered channel with a capacity of 2

	ch <- 1 // Sending the value 1 to the channel
	ch <- 2 // Sending the value 2 to the channel
	fmt.Println("Hey")

	go func() {
		fmt.Println("Hello")

		ch <- 3                              // Attempting to send the value 3 to the channel (channel already full)
		fmt.Println("Sent 3 to the channel") // This line will never be reached
	}()

	fmt.Println(<-ch) // Receiving the value 1 from the channel
	fmt.Println(<-ch) // Receiving the value 2 from the channel
}
