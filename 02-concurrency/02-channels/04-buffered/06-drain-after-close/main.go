// Drain after close: values already in the buffer survive a close.
//
// main puts 3 values in a buffered channel and closes it. Closing only
// means "no more sends". Receivers still get every value in the buffer,
// and only then see the channel as closed: ok becomes false and the value
// is the zero value. main keeps reading until ok is false, so it doesn't
// need to know how many values are left.
//
// Try: uncomment "ch <- 4". Sending on a closed channel panics.
//
// Real use:
//   - graceful shutdown: close the job queue, workers finish the jobs
//     already in it, then exit. No queued job is lost.
//
// Run: go run ./02-concurrency/02-channels/04-buffered/06-drain-after-close
package main

import "fmt"

func main() {
	const bufferSize = 3 // room for all the values below

	ch := make(chan int, bufferSize)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	// ch <- 4 // panic: send on closed channel

	for {
		v, ok := <-ch // ok is false once the channel is closed and empty
		fmt.Println(v, ok)
		if !ok {
			break
		}
	}
}
