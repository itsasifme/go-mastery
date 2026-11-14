// Add too late: calling wg.Add inside the goroutine.
//
// Each worker calls wg.Add(1) itself, as its first step.
//
// Problem: main reaches wg.Wait before the workers have started. The count
// is still 0, so Wait returns right away and main exits. The workers'
// output is usually missing.
//
// Fix: call wg.Add in main, before the "go" line. Comment out the line
// marked "problem", then uncomment the "fix:" line.
//
// Real use:
//   - a handler that starts background jobs in a loop and calls Add inside
//     each job. Under load, Wait sometimes returns too early, so the bug
//     only shows up now and then. Always call Add before "go".
//
// Run: go run ./02-concurrency/05-waitgroup/05-add-too-late
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
	// wg.Add(workerCount) // fix: count every worker before any of them starts
	for id := range workerCount {
		go worker(id+1, &wg)
	}

	wg.Wait()
	fmt.Println("main: all workers done?")
}

// worker does a small job and prints when it's finished.
func worker(id int, wg *sync.WaitGroup) {
	wg.Add(1) // problem: too late, main may already be past wg.Wait
	defer wg.Done()

	time.Sleep(workTime)
	fmt.Println("worker", id, "finished")
}
