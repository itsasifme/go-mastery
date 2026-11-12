// Nil or stop: read several channels until all close, or stop early when told.
//
// readUntilDone reads a and b and sets each one to nil once it's closed
// (see 06-nil-channel). It also listens to stop, so someone else can stop
// it before both channels are closed (see 05-for-select).
//
// Here a finishes fast, but b is slow. stop is closed after 250ms, so the
// loop stops before b is done.
//
// After stop, nobody reads b anymore. A plain "out <- v" in produce would
// wait forever: a goroutine leak (see channels/03-unbuffered/05-leak).
// So produce watches stop too, and gives up its send once the reader is
// gone. main waits on done for both producers, so a stuck producer shows
// up as a deadlock instead of hiding.
//
// Try: comment out "close(stop)". The loop then ends by itself once b closes.
//
// Try: replace the select in produce with a plain "out <- v". b's producer
// gets stuck, main waits on done forever, and Go reports a deadlock.
//
// Real use:
//   - a consumer reading several queues that must stop on shutdown
//     (SIGTERM or ctx.Done()), even if the queues aren't empty yet.
//
// Run: go run ./02-concurrency/03-select/07-nil-or-stop
package main

import (
	"fmt"
	"time"
)

func main() {
	const producerCount = 2 // must match the number of "go produce" lines below

	a := make(chan int)
	b := make(chan int)
	stop := make(chan struct{}) // stopAfter -> everyone: "please stop"
	done := make(chan struct{}) // producers -> main: "I'm finished"

	go produce([]int{1, 2, 3}, 0, a, stop, done)
	go produce([]int{10, 20, 30, 40, 50}, 100*time.Millisecond, b, stop, done)

	go stopAfter(250*time.Millisecond, stop)

	readUntilDone(a, b, stop)

	for range producerCount {
		<-done // one per producer: proves none is stuck
	}
	fmt.Println("producers finished")
}

// stopAfter waits for d, then closes stop. One close reaches the reader
// and both producers.
func stopAfter(d time.Duration, stop chan<- struct{}) {
	time.Sleep(d)
	close(stop)
}

// produce sends each value after a delay, then closes the channel.
// It stops early if stop is closed.
func produce(values []int, delay time.Duration, out chan<- int, stop <-chan struct{}, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done
	defer close(out)

	for _, v := range values {
		time.Sleep(delay)
		select {
		case out <- v:
		case <-stop:
			return // reader is gone: don't wait for it forever
		}
	}
}

// readUntilDone prints values from a and b until both are closed or stop is closed.
func readUntilDone(a, b <-chan int, stop <-chan struct{}) {
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // closed: never pick this case again
				fmt.Println("a closed")
				continue
			}
			fmt.Println("a:", v)
		case v, ok := <-b:
			if !ok {
				b = nil
				fmt.Println("b closed")
				continue
			}
			fmt.Println("b:", v)
		case <-stop:
			fmt.Println("told to stop")
			return // stop early, even if b isn't done
		}
	}
	fmt.Println("both closed")
}
