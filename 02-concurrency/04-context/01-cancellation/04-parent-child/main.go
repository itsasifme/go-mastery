// Parent and child: cancelling a parent cancels its children, not the other way.
//
// child1 and child2 are made from parent, and a worker does steps on each
// of the three contexts. Cancelling child1 stops only child1's worker.
// Cancelling parent stops parent's worker and every child's worker still
// running. The step counts show who kept working longer.
//
// Real use:
//   - a request ctx is the parent, and each database query gets a child
//     with its own timeout. If the client disconnects, every query
//     stops. If one query times out, the request keeps going.
//
// Run: go run ./02-concurrency/04-context/01-cancellation/04-parent-child
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parent, cancelParent := context.WithCancel(context.Background())
	child1, cancelChild1 := context.WithCancel(parent)
	child2, cancelChild2 := context.WithCancel(parent)
	defer cancelChild2() // stopped by its parent below, but always call cancel

	done := make(chan struct{}) // workers -> main: "I'm finished"

	go worker(parent, "parent", done)
	go worker(child1, "child1", done)
	go worker(child2, "child2", done)

	time.Sleep(200 * time.Millisecond) // let every worker do some steps
	fmt.Println("main: cancel child1")
	cancelChild1()
	<-done // only child1's worker stops. parent and child2 keep running

	time.Sleep(200 * time.Millisecond) // parent and child2 keep working
	fmt.Println("main: cancel parent")
	cancelParent()
	<-done // parent's and child2's workers, in any order
	<-done // child2 stopped only because its parent was cancelled
}

// worker does one step every 50ms until its ctx is cancelled, then says
// how many steps it did and why it stopped.
func worker(ctx context.Context, name string, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	steps := 0
	for {
		select {
		case <-time.After(50 * time.Millisecond): // simulate one step of work
			steps++
		case <-ctx.Done():
			fmt.Println(name+": stopped after", steps, "steps:", ctx.Err())
			return
		}
	}
}
