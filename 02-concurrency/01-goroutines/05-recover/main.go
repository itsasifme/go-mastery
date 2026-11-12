// Recover in a goroutine: catch a panic so the program keeps running.
//
// How it works (like try/catch in other languages):
//   - panic stops the function right away, like "throw".
//   - defer runs a function when the current function ends, even during
//     a panic, like "finally".
//   - recover, called inside that deferred function, stops the panic and
//     returns the panic value, like "catch". With no panic, it returns nil.
//
// Each job has its own defer + recover. When job 2 panics, the recover
// runs in that same goroutine, so the other jobs and main keep going.
//
// Note: recover only works inside a deferred function, and only for a
// panic in the same goroutine (see 04-panic). Use it as a safety net for
// bugs. For normal failures, return an error instead.
//
// Real use:
//   - net/http recovers a panic per request, so one broken handler doesn't
//     kill the server.
//   - background job runners recover per job, log the stack trace, and
//     move on to the next job.
//
// Run: go run ./02-concurrency/01-goroutines/05-recover
package main

import (
	"fmt"
	"time"
)

const (
	jobCount   = 3 // how many jobs main starts
	failingJob = 2 // the id of the job that panics
)

func main() {
	for i := range jobCount {
		go job(i)
	}

	time.Sleep(time.Second) // hack: wait for them (real fix: WaitGroup)
	fmt.Println("main still running")
}

func job(id int) {
	// 1. Register the "catch". It runs when job ends, panic or not.
	defer func() {
		// 3. recover stops the panic and returns its value (nil if no panic).
		if r := recover(); r != nil {
			fmt.Println("job", id, "recovered:", r)
		}
		// 4. job now returns normally, and the program keeps running.
	}()

	if id == failingJob {
		panic("something went wrong") // 2. job stops here and the deferred function runs
	}
	fmt.Println("job", id, "done") // never runs for the failing job
}
