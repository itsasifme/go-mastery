// Late read: atomic stops the race, but not a goroutine reading too late.
//
// This is goroutines/03-loop-variable again, with shared as an
// atomic.Int64 instead of a plain int. The loop reuses one shared variable
// for every round. Each goroutine reads it when it runs, not when the loop
// started it. By then the loop has usually finished, so every job gets the
// same ID: the value that ended the loop.
//
// Problem: in loop-variable, "go run -race" reports this. Here it finds
// nothing, because every read and write is atomic. The IDs are still
// wrong. Atomic fixed the race, not the bug: each job needs its own copy
// of the value from its round.
//
// Fix: read shared before "go" and pass that copy in. Use the "fix:" line
// instead of the goroutine above it.
//
// Real use:
//   - giving each request, job or upload a sequence number from a shared
//     counter. Take it with id := next.Add(1) before "go", and pass id in.
//
// Run: go run ./02-concurrency/09-atomic/05-late-read
// Race check: go run -race ./02-concurrency/09-atomic/05-late-read
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// jobCount is how many jobs the loop starts.
const jobCount = 3

func main() {
	// Declared outside the loop: one variable, reused every round, like
	// shared in goroutines/03-loop-variable. The C-style loop is the lesson.
	var shared atomic.Int64

	var wg sync.WaitGroup
	wg.Add(jobCount)
	for shared.Store(0); shared.Load() < jobCount; shared.Add(1) {
		go func() { startJob(shared.Load(), &wg) }() // problem: reads shared when it runs, not now
		// go startJob(shared.Load(), &wg) // fix: read it now and pass a copy
	}
	wg.Wait()
}

// startJob pretends to start the job with the given ID.
func startJob(id int64, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("job", id, "started")
}
