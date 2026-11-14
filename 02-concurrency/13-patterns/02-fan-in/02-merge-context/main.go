// Merge with context: stop every goroutine when the reader is done early.
//
// Three servers stream CPU readings that never end. main reads the merged
// stream until one reading is over alertAt, raises an alert, and calls
// cancel. Every watch and forward goroutine also watches ctx.Done(), so
// they all stop. Without ctx, they'd block forever on a send that no one
// reads. That's a goroutine leak.
//
// After cancel, main keeps reading until the merged channel is closed.
// That proves every goroutine stopped. If one leaked, the channel would
// never close, and Go would report a deadlock.
//
// Note: compare with 01-merge. There, the inputs ended on their own.
//
// Real use:
//   - watching several metric or event streams and stopping when you find
//     what you need, or when r.Context() is done.
//   - querying several replicas and taking the first good answer, then
//     cancelling the rest.
//
// Run: go run ./02-concurrency/13-patterns/02-fan-in/02-merge-context
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// reading is one CPU sample from one server.
type reading struct {
	server string
	cpu    int // percent
}

const (
	alertAt   = 90 // CPU percent that raises an alert
	readEvery = 20 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	web1 := make(chan reading)
	web2 := make(chan reading)
	db1 := make(chan reading)

	// Sample readings each server repeats. web-2 climbs past alertAt.
	go watch(ctx, "web-1", []int{40, 42, 41}, web1)
	go watch(ctx, "web-2", []int{50, 70, 93}, web2)
	go watch(ctx, "db-1", []int{30, 31, 29}, db1)

	merged := merge(ctx, web1, web2, db1)
	for r := range merged {
		fmt.Printf("%-6s cpu %d%%\n", r.server, r.cpu)
		if r.cpu > alertAt {
			fmt.Println("ALERT:", r.server, "is overloaded, stopping the watch")
			cancel()
			break
		}
	}

	extra := 0
	for range merged { // wait until merge closes it: proves nothing leaked
		extra++
	}
	fmt.Println("all watchers stopped,", extra, "readings came in after the alert")
}

// watch sends the server's readings, over and over, until ctx is done.
func watch(ctx context.Context, server string, samples []int, out chan<- reading) {
	defer close(out)
	sent := 0
	for {
		select {
		case <-time.After(readEvery):
		case <-ctx.Done():
			return
		}

		select {
		case out <- reading{server: server, cpu: samples[sent%len(samples)]}:
			sent++
		case <-ctx.Done():
			return
		}
	}
}

// merge returns one channel that gets every value from inputs. It's
// closed after every input is closed or ctx is done.
func merge(ctx context.Context, inputs ...<-chan reading) <-chan reading {
	out := make(chan reading)

	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, in := range inputs {
		go forward(ctx, in, out, &wg)
	}
	go closeWhenDone(&wg, out)

	return out
}

// forward copies values from in to out until in is closed or ctx is done.
func forward(ctx context.Context, in <-chan reading, out chan<- reading, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		select {
		case out <- v:
		case <-ctx.Done():
			return
		}
	}
}

// closeWhenDone waits for every forward, then closes out.
func closeWhenDone(wg *sync.WaitGroup, out chan<- reading) {
	wg.Wait()
	close(out)
}
