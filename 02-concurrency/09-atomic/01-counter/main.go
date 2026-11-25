// Counter: count from many goroutines without a lock.
//
// A server keeps two numbers for its metrics page: how many requests it
// handled and how many bytes it sent. Every handler adds to them. With
// atomic.Int64, Add does read, add and write back as one step. No other
// goroutine can get in between. No mutex, and no lost updates.
//
// The two counters are separate. Each one is always right on its own. But
// a reader might see the new requests value next to the old bytes value.
// If several values must change together, use a mutex (see
// mutex/02-safe-type).
//
// Note: compare with mutex/01-race. The same kind of counter, without a
// lock, lost updates.
//
// Real use:
//   - request, error and byte counters behind a /metrics endpoint.
//   - counting jobs done by a worker pool to show progress.
//
// Run: go run ./02-concurrency/09-atomic/01-counter
// Race check: go run -race ./02-concurrency/09-atomic/01-counter
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// metrics holds counters that many handlers update at once.
type metrics struct {
	requests  atomic.Int64
	bytesSent atomic.Int64
}

const (
	requestsPerWorker = 1000
	responseSize      = 512 // bytes in each fake response
)

func main() {
	const workerCount = 100

	var m metrics

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go handleRequests(&m, &wg)
	}
	wg.Wait()

	fmt.Println("requests:", m.requests.Load(), "expected", workerCount*requestsPerWorker)
	fmt.Println("bytes:   ", m.bytesSent.Load(), "expected", workerCount*requestsPerWorker*responseSize)
}

// handleRequests pretends to serve requestsPerWorker requests and counts
// each one.
func handleRequests(m *metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for range requestsPerWorker {
		m.requests.Add(1)
		m.bytesSent.Add(responseSize)
	}
}
