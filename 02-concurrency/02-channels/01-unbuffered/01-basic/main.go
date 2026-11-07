// Unbuffered channel: a receive waits until someone sends.
//
// main waits at <-ch until the goroutine sends, so no Sleep or WaitGroup
// is needed.
//
// Run: go run ./02-concurrency/02-channels/01-unbuffered/01-basic
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go func() {
		time.Sleep(time.Second) // simulate work
		ch <- 5
	}()

	x := <-ch // blocks ~1s until the send
	fmt.Println(x)
}
