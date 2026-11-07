// Drain after close: values already in the buffer survive a close.
//
// main puts 3 values in a buffered channel and closes it. Closing only
// means "no more sends". Receivers still get every value in the buffer,
// and only then see the channel as closed: ok becomes false and the value
// is the zero value.
//
// Try: uncomment "ch <- 4". Sending on a closed channel panics.
//
// Run: go run ./02-concurrency/02-channels/04-closing/02-drain-after-close
package main

import "fmt"

func main() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	// ch <- 4 // panic: send on closed channel

	for i := 0; i < 4; i++ {
		v, ok := <-ch // ok is false once the channel is closed and empty
		fmt.Println(v, ok)
	}
}
