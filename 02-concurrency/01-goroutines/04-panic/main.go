// Panic in a goroutine: it crashes the whole program.
//
// Problem: recover only works in the goroutine that panics. The recover
// in main never runs, and the program stops with "panic: boom".
//
// Fix: recover inside the goroutine. Uncomment the "fix:" line.
//
// Run: go run ./02-concurrency/01-goroutines/04-panic
package main

import (
	"fmt"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("main recovered:", r) // never runs
		}
	}()

	go func() {
		// defer func() { fmt.Println("recovered:", recover()) }() // fix: recover in the same goroutine
		panic("boom")
	}()

	time.Sleep(time.Second) // hack: give the goroutine time to panic
}
