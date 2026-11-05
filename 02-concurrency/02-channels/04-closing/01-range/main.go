package main

import "fmt"

func main() {
	ch := make(chan int) // Creating an unbuffered channel

	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i // Sending values to the channel
		}
		close(ch) // Closing the channel after all values are sent
	}()

	for x := range ch {
		fmt.Println(x) // Output: 1, 2, 3, 4, 5
	}
}
