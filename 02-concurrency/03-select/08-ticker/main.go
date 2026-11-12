// Ticker: do something every few milliseconds until told to stop.
//
// heartbeat prints a beat every 200ms. Its select waits on the ticker's
// channel and on stop at the same time. After about 1 second, main closes
// stop, and heartbeat ends. You see 4 or 5 beats, depending on timing.
//
// Always call ticker.Stop when you're done with a ticker.
//
// Real use:
//   - health checks, flushing metrics every few seconds, refreshing a
//     cache, or sending heartbeats to a server.
//
// Run: go run ./02-concurrency/03-select/08-ticker
package main

import (
	"fmt"
	"time"
)

func main() {
	const (
		every  = 200 * time.Millisecond // time between beats
		runFor = time.Second            // how long main lets heartbeat run
	)

	stop := make(chan struct{}) // main -> heartbeat: "please stop"
	done := make(chan struct{}) // heartbeat -> main: "I'm finished"

	go heartbeat(every, stop, done)

	time.Sleep(runFor)
	close(stop) // tell heartbeat to stop
	<-done      // wait until it did, so its last line gets printed
}

// heartbeat prints a beat on every tick until stop is closed.
func heartbeat(every time.Duration, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	// The ticker sends the current time on ticker.C once per "every".
	// The first tick comes after one "every", not right away.
	ticker := time.NewTicker(every)
	defer ticker.Stop() // no more ticks once heartbeat returns

	beatCount := 0
	for {
		select {
		case <-ticker.C: // a tick arrived. If we're too slow, missed ticks are dropped, not queued
			beatCount++
			fmt.Println("beat", beatCount)
		case <-stop: // main closed stop
			fmt.Println("stopped after", beatCount, "beats")
			return // not break: break would only leave the select
		}
	}
}
