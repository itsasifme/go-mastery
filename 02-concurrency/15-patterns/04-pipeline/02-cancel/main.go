// Cancel: stop every step when the end of the pipeline has had enough.
//
// A log search tool reads a log that never ends:
//
//	readLines -> parse -> keepErrors -> main
//
// main only wants the first few errors. Once it has them, it calls cancel.
// Every step sends with a select that also watches ctx.Done(). So each
// one stops, instead of waiting forever on a send nobody reads. Each step
// prints how many items it handled before it stopped.
//
// A step that's waiting to receive stops a different way: the step before
// it closes its channel, so its range loop ends. That's why parse sometimes
// prints "input closed" instead of "stopped". Both mean it stopped.
//
// main waits for every step with a WaitGroup. If one didn't stop, Go
// would report a deadlock.
//
// Note: compare with 01-basic. There, the input ended on its own, so no
// step had to be stopped.
//
// Real use:
//   - "find the first N matches" over a big file or a stream, like
//     grep -m.
//   - a pipeline started for one HTTP request. When the client leaves,
//     r.Context() is done, and every step stops.
//
// Run: go run ./02-concurrency/15-patterns/04-pipeline/02-cancel
package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// entry is one parsed log line.
type entry struct {
	level string
	msg   string
}

func main() {
	const (
		stepCount  = 3 // readLines, parse, keepErrors
		wantErrors = 3
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lines := make(chan string)  // readLines -> parse
	entries := make(chan entry) // parse -> keepErrors
	errs := make(chan entry)    // keepErrors -> main

	var wg sync.WaitGroup
	wg.Add(stepCount)
	go readLines(ctx, lines, &wg)
	go parse(ctx, lines, entries, &wg)
	go keepErrors(ctx, entries, errs, &wg)

	found := 0
	for e := range errs {
		found++
		fmt.Printf("error %d: %s\n", found, e.msg)
		if found == wantErrors {
			break
		}
	}

	cancel()
	wg.Wait() // proves every step stopped
	fmt.Println("main: done, every step stopped")
}

// readLines sends log lines forever, like tailing a file, until ctx is done.
func readLines(ctx context.Context, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	sample := []string{ // a log that repeats, so it never ends
		"INFO user logged in",
		"INFO page viewed",
		"ERROR payment timeout",
		"INFO page viewed",
		"ERROR disk almost full",
	}
	sent := 0
	for {
		select {
		case out <- sample[sent%len(sample)]:
			sent++
		case <-ctx.Done():
			fmt.Println("readLines: stopped after", sent, "lines")
			return
		}
	}
}

// parse splits each line into a level and a message.
func parse(ctx context.Context, in <-chan string, out chan<- entry, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	parsed := 0
	for line := range in {
		level, msg, _ := strings.Cut(line, " ")
		select {
		case out <- entry{level: level, msg: msg}:
			parsed++
		case <-ctx.Done():
			fmt.Println("parse: stopped after", parsed, "entries")
			return
		}
	}
	fmt.Println("parse: input closed after", parsed, "entries")
}

// keepErrors passes on only the ERROR entries.
func keepErrors(ctx context.Context, in <-chan entry, out chan<- entry, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(out)

	kept := 0
	for e := range in {
		if e.level != "ERROR" {
			continue
		}
		select {
		case out <- e:
			kept++
		case <-ctx.Done():
			fmt.Println("keepErrors: stopped after", kept, "errors")
			return
		}
	}
	fmt.Println("keepErrors: input closed after", kept, "errors")
}
