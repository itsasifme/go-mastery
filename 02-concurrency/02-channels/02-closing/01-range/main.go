// Closing a channel: a for range loop stops when the channel is closed.
//
// The goroutine sends 1 to 5, then closes ch. The loop in main reads
// every value and stops after the close.
//
// Try: remove close(out). The loop keeps waiting for more values, and the
// program deadlocks.
//
// Run: go run ./02-concurrency/02-channels/02-closing/01-range
package main

import "fmt"

func main() {
	ch := make(chan int)

	go produce(ch)

	for x := range ch {
		fmt.Println(x) // 1 2 3 4 5
	}
}

// produce sends 1 to 5, then closes the channel.
func produce(out chan<- int) {
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out) // tells range there are no more values
}
