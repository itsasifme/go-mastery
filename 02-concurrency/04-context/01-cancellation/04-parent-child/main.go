// Parent and child: cancelling a parent cancels its children, not the other way.
//
// child1 and child2 are made from parent. Cancelling child1 stops only
// child1. Cancelling parent stops parent and every child still running.
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
)

func main() {
	parent, cancelParent := context.WithCancel(context.Background())
	child1, cancelChild1 := context.WithCancel(parent)
	child2, cancelChild2 := context.WithCancel(parent)

	cancelChild1()
	show("after cancelling child1:", parent, child1, child2)

	cancelParent()
	show("after cancelling parent:", parent, child1, child2)

	cancelChild2() // already stopped by its parent, but always call cancel
}

// show prints whether each context is still running.
func show(title string, parent, child1, child2 context.Context) {
	fmt.Println(title)
	fmt.Println("  parent:", status(parent))
	fmt.Println("  child1:", status(child1))
	fmt.Println("  child2:", status(child2))
}

// status returns "running", or the reason the context stopped.
func status(ctx context.Context) string {
	if ctx.Err() == nil {
		return "running"
	}
	return ctx.Err().Error()
}
