// Many senders: who closes the channel?
//
// Two senders each send a value, then each closes ch.
//
// Problem: the first sender to finish closes ch while the other one may
// still use it. The program panics with "close of closed channel" or
// "send on closed channel", depending on timing.
//
// Fix: senders only send. One separate goroutine waits until all senders
// are done (with a WaitGroup, see waitgroup/01-basic), then closes ch once.
// Comment out the line marked "problem", then uncomment the "fix:" line.
//
// Real use:
//   - a worker pool where many workers send results to one channel. One
//     extra goroutine closes it after wg.Wait().
//
// Run: go run ./02-concurrency/02-channels/02-closing/03-many-senders
package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	for id := range 2 {
		wg.Add(1)
		go send(id, ch, &wg)
	}

	// go closeWhenDone(ch, &wg) // fix: one closer, after all senders are done

	for v := range ch {
		fmt.Println("got", v)
	}

	wg.Wait() // wait for the senders, so the program can't exit before a panic
}

// send sends one value, then closes the channel.
func send(id int, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	out <- id
	close(out) // problem: every sender closes, but only one close is allowed
}

// closeWhenDone waits until every sender is finished, then closes the
// channel once.
func closeWhenDone(out chan<- int, wg *sync.WaitGroup) {
	wg.Wait()
	close(out)
}
