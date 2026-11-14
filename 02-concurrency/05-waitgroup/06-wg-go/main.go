// wg.Go: Add, go and Done in one call (Go 1.25+).
//
// wg.Go(f) calls Add(1), starts f in a new goroutine, and calls Done when
// f returns. You can't forget Done, and you can't call Add too late (see
// 05-add-too-late), because Go does the Add before the goroutine starts.
//
// The worker no longer needs the WaitGroup at all. It just does its job.
//
// Note: wg.Go takes a func() with no arguments, so a small closure passes
// id to the worker. The closure only carries id. The WaitGroup stays in main.
//
// Real use:
//   - any "start N goroutines and wait for all" code written for Go 1.25+.
//     Most code that used Add/go/defer Done can switch to wg.Go.
//
// Run: go run ./02-concurrency/05-waitgroup/06-wg-go
package main

import (
	"fmt"
	"sync"
	"time"
)

// workTime is how long one fake job takes.
const workTime = 50 * time.Millisecond

func main() {
	const workerCount = 3

	var wg sync.WaitGroup
	for id := range workerCount {
		wg.Go(func() { worker(id + 1) })
	}

	wg.Wait()
	fmt.Println("main: all workers done")
}

// worker does a small job and prints when it's finished. It knows nothing
// about the WaitGroup.
func worker(id int) {
	time.Sleep(workTime)
	fmt.Println("worker", id, "finished")
}
