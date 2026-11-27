// Acquire timeout: wait for a slot, but not forever.
//
// A shipping partner's API allows only maxCalls calls at once from us.
// Requests arrive all at once. Each one waits for a slot, but only up to
// waitLimit. If no slot frees up in time, Acquire returns the context's
// error. Then the request answers "busy, try again" (HTTP 503) instead of
// hanging.
//
// TryAcquire goes one step further: it doesn't wait at all.
//
// Note: compare with 02-weighted. There, jobs waited as long as needed.
// Database connections don't need this: database/sql already limits them
// with db.SetMaxOpenConns.
//
// Note: maxCalls counts calls from this server only. The partner's limit is
// for the whole app, so with 20 servers the partner can get 20 times
// maxCalls at once. For a limit shared by every server, keep the count in
// Redis or use a shared rate limiter.
//
// Real use:
//   - an HTTP handler that calls sem.Acquire(r.Context(), 1) before calling
//     a partner API, so the wait also ends if the client hangs up.
//   - protecting a slow service that your code calls. Answer "busy" right
//     away, instead of letting waiting requests pile up.
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
	callTime  = 300 * time.Millisecond // how long one partner API call takes
	waitLimit = 100 * time.Millisecond // shorter than callTime, so late requests give up
)

func main() {
	const (
		maxCalls     = 2 // the partner's limit
		requestCount = 5
	)

	sem := semaphore.NewWeighted(maxCalls)

	var wg sync.WaitGroup
	wg.Add(requestCount)
	for id := range requestCount {
		go handle(id+1, sem, &wg)
	}
	wg.Wait()
}

// handle is one request. It waits up to waitLimit for a slot, then calls
// the partner API or answers busy.
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

	time.Sleep(callTime) // simulate the partner API call
	fmt.Printf("request %d: 200 ok\n", id)
}
