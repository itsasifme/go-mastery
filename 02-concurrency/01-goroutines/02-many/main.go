// Many goroutines: they don't run in a fixed order.
//
// main starts 5 goroutines, and the print order changes on every run.
// runtime.NumGoroutine shows how many are still running.
//
// Real use:
//   - net/http runs every incoming request in its own goroutine.
//   - processing many items at once, like resizing uploaded images or
//     calling several APIs.
//
// Run: go run ./02-concurrency/01-goroutines/02-many
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	const goroutineCount = 5

	for i := range goroutineCount {
		go printID(i)
	}

	fmt.Println("running:", runtime.NumGoroutine()) // main + the ones not done yet
	time.Sleep(time.Second)                         // hack: wait for them (real fix: WaitGroup)
}

// printID prints which goroutine is running.
func printID(id int) {
	fmt.Println("goroutine", id)
}
