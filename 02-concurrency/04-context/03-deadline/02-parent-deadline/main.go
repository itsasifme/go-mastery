// Parent deadline: a child can make the deadline earlier, but never later.
//
// The parent has 1 second. A child that asks for 5 seconds still ends at
// the parent's deadline. A child that asks for 200ms gets 200ms, because
// that's earlier. The earliest deadline always wins.
//
// Real use:
//   - a server gives each request 2 seconds. A database call inside asks
//     for 5 seconds, but still stops at the request's 2 seconds. A cache
//     lookup asks for 50ms, and gets it.
//
// Run: go run ./02-concurrency/04-context/03-deadline/02-parent-deadline
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	parent, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	later, cancelLater := context.WithTimeout(parent, 5*time.Second)
	defer cancelLater()

	sooner, cancelSooner := context.WithTimeout(parent, 200*time.Millisecond)
	defer cancelSooner()

	fmt.Println("parent:           ", left(parent))
	fmt.Println("child asked 5s:   ", left(later)) // capped by the parent
	fmt.Println("child asked 200ms:", left(sooner))
}

// left returns how much time ctx has before its deadline.
func left(ctx context.Context) time.Duration {
	d, _ := ctx.Deadline()
	return time.Until(d).Round(10 * time.Millisecond)
}
