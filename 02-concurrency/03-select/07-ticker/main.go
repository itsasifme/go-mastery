// Ticker: do something every few milliseconds until told to stop.
//
// heartbeat prints a beat every 200ms. Its select waits on the ticker's
// channel and on stop at the same time. After about 1 second, main closes
// stop, and heartbeat ends. You see 4 or 5 beats, depending on timing.
//
// Always call ticker.Stop when you're done with a ticker.
//
// Run: go run ./02-concurrency/03-select/07-ticker
package main

import (
	"fmt"
	"time"
)

func main() {
	stop := make(chan struct{})
	done := make(chan struct{})

	go heartbeat(200*time.Millisecond, stop, done)

	time.Sleep(time.Second)
	close(stop) // tell heartbeat to stop
	<-done      // wait until it did
}

// heartbeat prints a beat on every tick until stop is closed.
func heartbeat(every time.Duration, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)

	ticker := time.NewTicker(every)
	defer ticker.Stop()

	beats := 0
	for {
		select {
		case <-ticker.C: // gets a value every tick
			beats++
			fmt.Println("beat", beats)
		case <-stop:
			fmt.Println("stopped after", beats, "beats")
			return
		}
	}
}
