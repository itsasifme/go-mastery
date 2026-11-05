package main

import "fmt"

func readData(ch <-chan int) {
	x := <-ch // Receiving data from the channel
	fmt.Println("x: ", x)
}

func main() {
	ch := make(<-chan int) // Declaring a receive-only channel

	go readData(ch) // Reading data from the channel
}
