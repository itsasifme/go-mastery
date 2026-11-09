// Priority: select has no priority, but you can build one.
//
// A worker gets urgent jobs and normal jobs. Both channels already have 3
// jobs waiting.
//
// Problem: handleRandom uses one plain select. When both cases are ready,
// select picks at random (see 02-random-pick), so urgent jobs end up
// mixed in with normal ones.
//
// Fix: handleUrgentFirst first checks the urgent channel with a select
// that has a default. Only when no urgent job is waiting does it wait on
// both channels. All urgent jobs now run first. Use the "fix:" line
// instead of the line above it.
//
// Real use:
//   - control messages (shutdown, config reload) before normal work.
//   - paid users' jobs before free users' jobs.
//
// Run: go run ./02-concurrency/03-select/09-priority
package main

import "fmt"

// jobsPerQueue is how many jobs wait in each channel. The handlers read
// exactly totalJobs, so the numbers can't get out of sync.
const (
	jobsPerQueue = 3
	totalJobs    = 2 * jobsPerQueue
)

func main() {
	urgent := make(chan string, jobsPerQueue)
	normal := make(chan string, jobsPerQueue)
	for i := range jobsPerQueue {
		urgent <- fmt.Sprint("urgent ", i)
		normal <- fmt.Sprint("normal ", i)
	}

	handleRandom(urgent, normal)
	// handleUrgentFirst(urgent, normal) // fix: use this instead
}

// handleRandom takes whichever job select picks. Urgent jobs get no
// priority.
func handleRandom(urgent, normal <-chan string) {
	for range totalJobs {
		select {
		case j := <-urgent:
			fmt.Println(j)
		case j := <-normal:
			fmt.Println(j)
		}
	}
}

// handleUrgentFirst always takes waiting urgent jobs before normal ones.
func handleUrgentFirst(urgent, normal <-chan string) {
	for range totalJobs {
		select {
		case j := <-urgent: // 1. is an urgent job waiting? take it
			fmt.Println(j)
			continue
		default: // 2. no urgent job right now: go on
		}

		select { // 3. wait for whichever job comes next
		case j := <-urgent:
			fmt.Println(j)
		case j := <-normal:
			fmt.Println(j)
		}
	}
}
