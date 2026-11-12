// Many senders: who closes the channel?
//
// Two workers each send a value on ch, then each closes it.
//
// Problem: the first worker to finish closes ch while the other one may
// still use it. The program panics with "close of closed channel" or
// "send on closed channel", depending on timing.
//
// Fix: workers only send, then tell a done channel they're finished. One
// separate goroutine reads one done signal per worker, then closes ch once.
// Comment out the line marked "problem", then uncomment the "fix:" line.
//
// Real use:
//   - a worker pool where many workers send results to one channel. One
//     extra goroutine closes it after every worker is finished. Later,
//     sync.WaitGroup does this counting for you (see waitgroup/01-basic).
//
// Run: go run ./02-concurrency/02-channels/02-closing/03-many-senders
package main

import (
	"fmt"
	"time"
)

func main() {
	const workerCount = 2

	ch := make(chan int)                     // workers -> main: the values
	done := make(chan struct{}, workerCount) // workers -> closeWhenDone: "I'm finished"
	// done has room for every worker, so the send in worker's defer never
	// blocks. Without it, a panicking worker would get stuck in its defer,
	// and the panic would never show.

	for id := range workerCount {
		go worker(id, ch, done)
	}

	// go closeWhenDone(ch, done, workerCount) // fix: one closer, after all workers are done

	for v := range ch {
		fmt.Println("got", v)
	}

	time.Sleep(100 * time.Millisecond) // hack: give the other worker time to panic before main exits
}

// worker sends one value, closes the channel, then says it's finished.
func worker(id int, out chan<- int, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	out <- id
	close(out) // problem: every worker closes, but only one close is allowed
}

// closeWhenDone waits for one done signal from each worker, then closes
// the channel once.
func closeWhenDone(out chan<- int, done <-chan struct{}, workerCount int) {
	for range workerCount {
		<-done
	}
	close(out)
}
