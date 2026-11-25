// Race: many goroutines update one counter without a lock.
//
// A site counts page views. Every request runs in its own goroutine and
// does views++. That looks like one step, but it's three: read the value,
// add 1, write it back. Two goroutines can read the same value, and then
// one of the updates is lost.
//
// Problem: the final count is lower than expected, and it changes between
// runs. "go run -race" prints "WARNING: DATA RACE" with the lines
// involved.
//
// Fix: lock a mutex around views++, so only one goroutine changes it at a
// time. Uncomment the two "fix:" lines.
//
// Real use:
//   - every HTTP handler runs in its own goroutine. So a plain counter or
//     map shared between handlers needs a lock. For a counter, an atomic
//     works too (see atomic/01-counter).
//   - run your tests with "go test -race" in CI to catch this early.
//
// Run: go run ./02-concurrency/07-mutex/01-race
// Race check: go run -race ./02-concurrency/07-mutex/01-race
package main

import (
	"fmt"
	"sync"
)

// viewsPerWorker is how many views each goroutine counts.
const viewsPerWorker = 1000

func main() {
	const workerCount = 100

	views := 0
	var mu sync.Mutex // guards views

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go countViews(&views, &mu, &wg)
	}
	wg.Wait()

	fmt.Println("expected", workerCount*viewsPerWorker)
	fmt.Println("got     ", views)
}

// countViews adds viewsPerWorker to the shared counter, one at a time.
func countViews(views *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for range viewsPerWorker {
		// mu.Lock() // fix:
		*views++ // read, add, write back: two goroutines can overwrite each other
		// mu.Unlock() // fix:
	}
}
