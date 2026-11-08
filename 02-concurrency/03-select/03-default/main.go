// Default: try a send or receive without waiting.
//
// A select with a default case never blocks. If no other case is ready,
// default runs right away. That gives you "try" operations:
//
//	trySend     send if there's room, otherwise give up
//	tryReceive  take a value if one is waiting, otherwise move on
//
// A common use: a log or metrics channel that drops data when it's full,
// instead of slowing the whole program down.
//
// Note: don't put a default case in a loop that only waits. The loop
// spins and burns CPU. Leave default out, so select waits properly.
//
// Run: go run ./02-concurrency/03-select/03-default
package main

import "fmt"

func main() {
	ch := make(chan int, 2)

	v, ok := tryReceive(ch)
	fmt.Println("receive:", v, ok) // 0 false: nothing there yet

	for i := 1; i <= 3; i++ {
		if trySend(ch, i) {
			fmt.Println("sent", i)
		} else {
			fmt.Println("dropped", i) // buffer full, nobody reading
		}
	}

	v, ok = tryReceive(ch)
	fmt.Println("receive:", v, ok) // 1 true
}

// trySend sends v if there's room, and reports whether it did.
func trySend(ch chan<- int, v int) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

// tryReceive returns a value if one is waiting, without blocking.
func tryReceive(ch <-chan int) (int, bool) {
	select {
	case v := <-ch:
		return v, true
	default:
		return 0, false
	}
}
