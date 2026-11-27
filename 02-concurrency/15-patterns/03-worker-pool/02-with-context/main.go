// With context: stop the whole pool early.
//
// A newsletter goes out to many subscribers through a pool of workers.
// The job has a time budget. When ctx is done, queueEmails stops adding
// emails, and each worker stops after the email it's sending now. Every
// worker reports how far it got, so you can see what's left to send.
//
// Both sides watch ctx.Done() in a select. Without that, queueEmails
// would block forever on a send that no worker will take.
//
// When ctx is done, queueEmails also closes emails. Then both cases in a
// worker's select are ready, and Go picks one at random. So each worker
// prints "stopped" or "no more emails", and it changes between runs. Both
// mean it stopped.
//
// Note: compare with 01-basic. There, the pool always ran until jobs was
// empty.
//
// Real use:
//   - a batch job with a deadline, or one that stops on Ctrl+C with
//     signal.NotifyContext(ctx, os.Interrupt).
//   - a pool started per request that must stop when r.Context() is done.
//
// Run: go run ./02-concurrency/15-patterns/03-worker-pool/02-with-context
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// sendTime is how long sending one fake email takes.
const sendTime = 50 * time.Millisecond

func main() {
	const (
		workerCount     = 3
		subscriberCount = 100
		budget          = 300 * time.Millisecond // not enough time for every email
	)

	ctx, cancel := context.WithTimeout(context.Background(), budget)
	defer cancel()

	emails := make(chan int) // queueEmails -> workers: subscriber IDs

	go queueEmails(ctx, subscriberCount, emails)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(ctx, id+1, emails, &wg)
	}
	wg.Wait()

	fmt.Println("main: stopped:", ctx.Err())
}

// queueEmails sends subscriber IDs until they run out or ctx is done.
func queueEmails(ctx context.Context, count int, emails chan<- int) {
	defer close(emails)
	for id := range count {
		select {
		case emails <- id + 1:
		case <-ctx.Done():
			fmt.Printf("queue: stopped with %d of %d emails queued\n", id, count)
			return
		}
	}
}

// worker sends emails until emails is closed or ctx is done.
func worker(ctx context.Context, id int, emails <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	sent, last := 0, 0 // last is the last subscriber emailed
	for {
		select {
		case subscriber, ok := <-emails:
			if !ok {
				fmt.Printf("worker %d: no more emails, sent %d, last subscriber %d\n", id, sent, last)
				return
			}
			time.Sleep(sendTime) // simulate sending one email. A started email always finishes
			sent, last = sent+1, subscriber
		case <-ctx.Done():
			fmt.Printf("worker %d: stopped, sent %d, last subscriber %d\n", id, sent, last)
			return
		}
	}
}
