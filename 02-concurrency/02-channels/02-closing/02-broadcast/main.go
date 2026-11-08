// Broadcast: closing a channel wakes every goroutine waiting on it.
//
// 3 workers wait on <-stop. main closes stop once, and all 3 wake up and
// finish. A receive from a closed channel never blocks, so every waiting
// receiver gets through.
//
// This is why the later examples use close(done) to say "I'm finished".
//
// Here the workers do nothing but wait for stop. To keep working and
// still listen for stop, a worker needs select (see select/05-for-select).
//
// Try: replace close(stop) with stop <- struct{}{}. A send reaches only
// one receiver, so only one worker stops. The other two wait forever, and
// the program deadlocks.
//
// Run: go run ./02-concurrency/02-channels/02-closing/02-broadcast
package main

import "fmt"

func main() {
	stop := make(chan struct{}) // struct{} carries no data: it's only a signal
	done := make(chan int)

	for id := 1; id <= 3; id++ {
		go worker(id, stop, done)
	}

	close(stop) // one close wakes all 3 workers

	for i := 0; i < 3; i++ {
		<-done
	}
}

// worker waits for the stop signal, then reports that it stopped.
func worker(id int, stop <-chan struct{}, done chan<- int) {
	<-stop // wait for the signal
	fmt.Println("worker", id, "stopped")
	done <- id
}
