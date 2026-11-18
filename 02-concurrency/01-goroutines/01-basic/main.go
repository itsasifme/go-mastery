// Basic goroutine: main doesn't wait for the goroutines it starts.
//
// Problem: when main returns, the program ends. So "Hello World!" may
// never print.
//
// Fix: uncomment both "fix" lines. That's a quick hack for a demo.
// The real fix is a WaitGroup (see waitgroup/01-basic).
//
// Real use:
//   - fire-and-forget work in a server, like sending a welcome email after
//     signup. The server keeps running, so the goroutine gets to finish.
//   - in a CLI tool or batch job, main must wait, or background work is
//     lost when the program exits.
//
// Run: go run ./02-concurrency/01-goroutines/01-basic
package main

import (
	"fmt"
	// "time" // fix
)

func main() {
	go helloWorld()

	// time.Sleep(time.Second) // fix: give the goroutine time to run

	goodbye()
}

// helloWorld prints a greeting. main starts it as a goroutine.
func helloWorld() {
	fmt.Println("Hello World!")
}

// goodbye prints a goodbye. main calls it directly.
func goodbye() {
	fmt.Println("Good Bye!")
}
