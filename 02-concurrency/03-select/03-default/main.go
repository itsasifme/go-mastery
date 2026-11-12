// Default: try a send or receive without waiting.
//
// A select with a default case never blocks. If no other case is ready,
// default runs right away. That gives you "try" operations:
//
//	trySend     send if there's room, otherwise give up
//	tryReceive  take a value if one is waiting, otherwise move on
//
// Note: don't put a default case in a loop that only waits. The loop
// spins and burns CPU. Leave default out, so select waits properly.
//
// Real use:
//   - a metrics or log channel that drops data when full, instead of
//     slowing down requests.
//   - a quick "is there work?" check before doing something else.
//
// Run: go run ./02-concurrency/03-select/03-default
package main

import "fmt"

func main() {
	const bufferSize = 2

	ch := make(chan int, bufferSize)

	v, ok := tryReceive(ch)
	fmt.Println("receive:", v, ok) // 0 false: nothing there yet

	for i := range bufferSize + 1 { // one more than fits
		if trySend(ch, i) {
			fmt.Println("sent", i)
		} else {
			fmt.Println("dropped", i) // buffer full, nobody reading
		}
	}

	v, ok = tryReceive(ch)
	fmt.Println("receive:", v, ok) // 0 true: the first value sent
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
