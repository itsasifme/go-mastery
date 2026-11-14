// Drop: when the buffer is full, throw the new value away.
//
// An app records a metric sample on every request. The exporter that
// sends them out is slower. Losing a few samples is fine, but slowing
// down requests isn't. So record uses select with default: if the buffer
// has room, the sample goes in. If not, it's dropped right away and
// counted.
//
// Note: compare with 01-blocking. There, the app waited and nothing was
// lost. Here, the app never waits, and some data is lost. Choose on
// purpose.
//
// Real use:
//   - metrics, traces or debug logs, where a dropped sample is better
//     than a slow request. Count the drops, so you know it's happening.
//   - a live video or game update stream, where only the newest data
//     matters.
//
// Run: go run ./02-concurrency/13-patterns/05-backpressure/02-drop
package main

import (
	"fmt"
	"time"
)

const (
	sampleCount = 20
	recordEvery = 10 * time.Millisecond // requests come in fast
	exportTime  = 50 * time.Millisecond // the exporter is slow
)

func main() {
	const bufferSize = 3

	samples := make(chan int, bufferSize) // app -> exporter
	done := make(chan struct{})           // exporter -> main: "I'm finished"

	go export(samples, done)

	dropped := 0
	for i := range sampleCount {
		time.Sleep(recordEvery) // simulate a request
		if !record(samples, i+1) {
			dropped++
		}
	}
	close(samples)
	<-done

	fmt.Printf("%d samples: kept %d, dropped %d\n", sampleCount, sampleCount-dropped, dropped)
}

// record tries to put a sample in the buffer without waiting. It reports
// whether the sample was kept.
func record(samples chan<- int, sample int) bool {
	select {
	case samples <- sample:
		return true
	default: // buffer full: drop it instead of waiting
		return false
	}
}

// export sends each sample to the fake metrics server, until samples is
// closed.
func export(samples <-chan int, done chan<- struct{}) {
	defer close(done)

	exported := 0
	for range samples {
		time.Sleep(exportTime) // simulate the network call
		exported++
	}
	fmt.Println("exporter: sent", exported, "samples")
}
