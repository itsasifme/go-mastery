package main

func sendData(ch chan<- int) {
	ch <- 5 // Sending data to the channel
}

func main() {
	ch := make(chan<- int) // Declaring a send-only channel

	go sendData(ch) // Sending data to the channel
}
