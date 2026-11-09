// Nil channel: turn off a select case once its channel is closed.
//
// printBoth reads from two channels until both are closed. A closed channel
// is always ready, so after one closes, select would keep picking it
// forever.
//
// The trick: when a channel closes, set it to nil. A nil channel is never
// ready, so select stops picking that case. When both are nil, the loop
// ends.
//
// Try: comment out "a = nil". The loop never ends and burns CPU, because
// the closed channel is always ready. Stop it with Ctrl+C.
//
// Real use:
//   - combining events from several sources, like two message queues or
//     file watchers, until all of them are finished.
//
// Run: go run ./02-concurrency/03-select/06-nil-channel
package main

import "fmt"

func main() {
	a := make(chan int)
	b := make(chan int)

	go produce([]int{1, 2, 3}, a)
	go produce([]int{10, 20}, b)

	printBoth(a, b)
}

// produce sends each value, then closes the channel.
func produce(values []int, out chan<- int) {
	for _, v := range values {
		out <- v
	}
	close(out)
}

// printBoth prints values from a and b until both are closed.
func printBoth(a, b <-chan int) {
	for a != nil || b != nil {
		select {
		case v, ok := <-a:
			if !ok {
				a = nil // closed: never pick this case again
				continue
			}
			fmt.Println("a:", v)
		case v, ok := <-b:
			if !ok {
				b = nil // closed: never pick this case again
				continue
			}
			fmt.Println("b:", v)
		}
	}
	fmt.Println("both closed")
}
