// Reject: wait a little for room, then tell the caller "busy".
//
// An API puts each upload into a queue for a background worker. When the
// queue is full, enqueue waits up to enqueueLimit for room. If none frees
// up, it returns errBusy. Then the handler answers 503, so the client can
// retry later. The client hears about the overload right away. It doesn't
// hang, and it doesn't lose its upload without knowing.
//
// Note: this sits between 01-blocking (wait forever) and 02-drop (don't
// wait at all, and don't tell anyone). semaphore/03-acquire-timeout answers
// 503 the same way, but it limits calls running at once. This one limits
// the queue.
//
// Real use:
//   - an HTTP handler that returns 503 or 429 with a Retry-After header
//     when the job queue is full.
//   - turning away extra requests in front of a slow database or service
//     (load shedding).
//
// Run: go run ./02-concurrency/15-patterns/05-backpressure/03-reject
package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	requestCount = 10
	arriveEvery  = 20 * time.Millisecond  // requests come in fast
	processTime  = 100 * time.Millisecond // the worker is slow
	enqueueLimit = 30 * time.Millisecond  // how long a request waits for room
)

// errBusy means the queue stayed full for the whole wait.
var errBusy = errors.New("server busy, try again later")

func main() {
	const queueSize = 2

	queue := make(chan int, queueSize) // handlers -> worker: upload IDs
	done := make(chan struct{})        // worker -> main: "I'm finished"

	go worker(queue, done)

	for id := range requestCount {
		time.Sleep(arriveEvery) // the next request arrives
		if err := enqueue(queue, id+1); err != nil {
			fmt.Printf("upload %2d: 503 %v\n", id+1, err)
			continue
		}
		fmt.Printf("upload %2d: 202 accepted\n", id+1)
	}
	close(queue)
	<-done
}

// enqueue puts an upload in the queue, waiting up to enqueueLimit for room.
func enqueue(queue chan<- int, upload int) error {
	select {
	case queue <- upload:
		return nil
	case <-time.After(enqueueLimit):
		return errBusy
	}
}

// worker processes uploads until queue is closed.
func worker(queue <-chan int, done chan<- struct{}) {
	defer close(done)
	for upload := range queue {
		time.Sleep(processTime) // simulate processing the upload
		fmt.Printf("worker:    processed upload %d\n", upload)
	}
}
