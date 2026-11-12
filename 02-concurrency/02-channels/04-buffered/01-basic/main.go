// Buffered channel: a send doesn't wait until the buffer is full.
//
// main sends 2 values with nobody receiving, then reads them back in the
// same order (first in, first out).
//
// Try: uncomment "ch <- 3". The buffer is full and nobody receives, so
// the program deadlocks.
//
// Real use:
//   - a small queue between two parts of a program, like an HTTP handler
//     handing jobs to a background worker.
//
// Run: go run ./02-concurrency/02-channels/04-buffered/01-basic
package main

import "fmt"

func main() {
	const bufferSize = 2

	ch := make(chan int, bufferSize)

	ch <- 1 // doesn't block: the buffer has room
	ch <- 2
	// ch <- 3 // buffer full: deadlock

	fmt.Println(len(ch), cap(ch)) // 2 2

	fmt.Println(<-ch) // 1: first in, first out
	fmt.Println(<-ch) // 2
}
