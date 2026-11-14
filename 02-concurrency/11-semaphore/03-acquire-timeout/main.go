// Acquire timeout: wait for a slot, but not forever.
//
// A service may use only maxConns database connections. Requests arrive
// all at once. Each one waits for a connection, but only up to waitLimit.
// If none frees up in time, Acquire returns the context's error, and the
// request answers "busy, try again" (HTTP 503) instead of hanging.
//
// TryAcquire is the extreme version: it doesn't wait at all.
//
// Note: compare with 02-weighted. There, jobs waited as long as needed.
//
// Real use:
//   - an HTTP handler that calls sem.Acquire(r.Context(), 1), so the wait
//     also ends if the client hangs up.
//   - protecting a slow downstream service: fail fast instead of piling
//     up waiting requests.
//
// Run: go run ./02-concurrency/11-semaphore/03-acquire-timeout
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

const (
	queryTime = 300 * time.Millisecond
	waitLimit = 100 * time.Millisecond // shorter than queryTime, so late requests give up
)

func main() {
	const (
		maxConns     = 2
		requestCount = 5
	)

	sem := semaphore.NewWeighted(maxConns)

	var wg sync.WaitGroup
	wg.Add(requestCount)
	for id := range requestCount {
		go handle(id+1, sem, &wg)
	}
	wg.Wait()
}

// handle is one request. It waits up to waitLimit for a connection, then
// runs its query or answers busy.
func handle(id int, sem *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(context.Background(), waitLimit)
	defer cancel()

	if err := sem.Acquire(ctx, 1); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("request %d: 503 busy, try again (waited %v)\n", id, waitLimit)
			return
		}
		fmt.Printf("request %d: %v\n", id, err)
		return
	}
	defer sem.Release(1)

	time.Sleep(queryTime) // simulate the query
	fmt.Printf("request %d: 200 ok\n", id)
}
