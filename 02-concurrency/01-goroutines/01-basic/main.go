// Basic goroutine: main doesn't wait for the goroutines it starts.
//
// Problem: when main returns, the program ends. So "Hello World!" may
// never print.
//
// Fix: uncomment both "fix" lines. That's a quick hack for a demo.
// The real fix is a WaitGroup (see 05-waitgroup).
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

func helloWorld() {
	fmt.Println("Hello World!")
}

func goodbye() {
	fmt.Println("Good Bye!")
}
