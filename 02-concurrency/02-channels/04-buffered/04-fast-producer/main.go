// Fast producer: a buffer lets a fast sender run ahead of a slow receiver.
//
// The producer sends 5 values with no delay. The receiver takes 200ms for
// each one. With a buffer of 3, this happens:
//
//	0ms     sent 1      the receiver is already waiting, so 1 goes straight to it
//	0ms     sent 2-4    no receiver free, so 2, 3, 4 go into the buffer (now full)
//	0ms     send 5      waits: the buffer is full
//	200ms   got 1       the receiver finishes 1 and takes 2 from the buffer
//	200ms   sent 5      that frees a slot, so 5 goes in. The producer is done.
//	1000ms  got 5       the receiver works through the rest, one every 200ms
//
// The times are a few ms later in real runs, because starting goroutines
// and waking them up takes a little time.
//
// Note: the producer is done at ~200ms, but the total is still ~1s. The
// slow receiver needs 5 x 200ms no matter what. A buffer doesn't make the
// work faster. It lets the producer finish early and do something else.
//
// Try: change the buffer to 0, make(chan int). Now every send waits for
// the receiver, so "sent" and "got" take turns and the producer is only
// done at ~800ms.
//
// Run: go run ./02-concurrency/02-channels/04-buffered/04-fast-producer
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 3)
	start := time.Now()

	go produce(ch, start)

	for i := 1; i <= 5; i++ {
		v := <-ch
		time.Sleep(200 * time.Millisecond) // slow work
		fmt.Printf("%4dms got  %d\n", time.Since(start).Milliseconds(), v)
	}
}

// produce sends 5 values as fast as it can.
func produce(out chan<- int, start time.Time) {
	for i := 1; i <= 5; i++ {
		out <- i // only waits when the buffer is full
		fmt.Printf("%4dms sent %d\n", time.Since(start).Milliseconds(), i)
	}
}
