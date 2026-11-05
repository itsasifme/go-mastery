package main

import "fmt"

func main() {
	ch := make(chan int) // Creating an unbuffered channel
	fmt.Println("Hey")

	go func() {
		fmt.Println("Hello")

		ch <- 5 // Sending a value to the channel
	}()

	// x := <-ch // Receiving the value from the channel
	// fmt.Println("x: ", x)

	// The program will deadlock at this point because the sending operation is blocked
}
