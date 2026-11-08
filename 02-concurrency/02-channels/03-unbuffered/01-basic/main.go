// Unbuffered channel: a receive waits until someone sends.
//
// main waits at <-ch until the goroutine sends, so no Sleep or WaitGroup
// is needed.
//
// Run: go run ./02-concurrency/02-channels/03-unbuffered/01-basic
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)

	go produce(ch)

	x := <-ch // blocks ~1s until the send
	fmt.Println(x)
}

// produce does some work, then sends the result.
// chan<- int means produce can only send on it (see 01-directional).
func produce(out chan<- int) {
	time.Sleep(time.Second) // simulate work
	out <- 5
}
