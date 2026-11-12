// WaitGroup: wait until a group of goroutines is finished.
//
// wg.Add(2) says "wait for 2 goroutines". Each one calls wg.Done when it
// ends, and wg.Wait blocks main until the count is back to 0. The print
// order can change between runs.
//
// Real use:
//   - calling several services at once and waiting for all of them, like
//     loading a user's profile, orders and settings in parallel.
//   - graceful shutdown: wait for running jobs to finish before exiting.
//   - closing a results channel after every worker is done (see
//     channels/02-closing/03-many-senders).
//
// Run: go run ./02-concurrency/05-waitgroup/01-basic
package main

import (
	"fmt"
	"sync"
)

func main() {
	const goroutineCount = 2 // must match the number of "go" lines below

	var wg sync.WaitGroup
	wg.Add(goroutineCount)
	go helloworld(&wg)
	go goodbye(&wg)
	wg.Wait()
}

func helloworld(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Hello World!")
}

func goodbye(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Good Bye!")
}
