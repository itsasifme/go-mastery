// Broadcast: closing a channel wakes every goroutine waiting on it.
//
// Every worker waits on <-stop. main closes stop once, and they all wake up
// and finish. A receive from a closed channel never blocks, so every waiting
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
// Real use:
//   - graceful shutdown: on SIGTERM, close one stop channel and every
//     worker stops.
//   - context cancellation works this way: ctx.Done() is a channel that
//     gets closed.
//
// Run: go run ./02-concurrency/02-channels/02-closing/02-broadcast
package main

import "fmt"

func main() {
	const workers = 3

	stop := make(chan struct{}) // main -> workers: "please stop"
	done := make(chan struct{}) // workers -> main: "I'm finished"

	for id := range workers {
		go worker(id, stop, done)
	}

	close(stop) // one close wakes every worker

	for range workers {
		<-done // one per worker
	}
}

// worker waits for the stop signal, then reports that it stopped.
func worker(id int, stop <-chan struct{}, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	<-stop // wait for the signal
	fmt.Println("worker", id, "stopped")
}
