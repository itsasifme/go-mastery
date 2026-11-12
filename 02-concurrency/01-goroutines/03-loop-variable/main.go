// Loop variable: since Go 1.22, each loop round gets its own variable.
//
// The first loop declares i with ":=" in the for line. Each round, Go
// copies i into a new variable and runs i++ on that one. So the count
// keeps going, but every goroutine keeps its own i. It prints 0, 1 and 2
// in any order.
//
// The second loop declares shared outside the loop, so there's only one
// variable, and all goroutines share it. They usually print 3 3 3.
// Before Go 1.22, "for i := ..." worked like this too.
//
// Problem: the second loop is also a data race. go run -race reports it.
//
// Fix: pass the value in as an argument, so each goroutine gets a copy.
// Use the "fix:" line instead of the goroutine above it.
//
// Real use:
//   - starting one goroutine per item in a slice (URLs, files, user IDs).
//     Older code passes the item as an argument because of this bug, so
//     you'll see that style a lot.
//
// Run: go run ./02-concurrency/01-goroutines/03-loop-variable
package main

import (
	"fmt"
	"time"
)

func main() {
	count := 3 // goroutines per loop

	for i := 0; i < count; i++ {
		go func() {
			fmt.Println("new:", i) // own i: copied from the last round, then i++
		}()
	}
	time.Sleep(time.Second) // hack: wait for them

	// Declared outside the loop: one variable, reused every round.
	// This is how "for i := ..." worked before Go 1.22.
	var shared int
	for shared = 0; shared < count; shared++ {
		go func() {
			fmt.Println("old:", shared) // same variable, usually 3 by the time it prints
		}()
		// go func(n int) { fmt.Println("fixed:", n) }(shared) // fix: pass a copy
	}
	time.Sleep(time.Second)
}
