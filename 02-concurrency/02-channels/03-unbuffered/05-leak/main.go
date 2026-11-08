// Goroutine leak: stuck goroutines don't cause any error.
//
// Problem: main takes the first result and moves on. The other 2 senders
// wait forever. Go doesn't complain, because main is still running.
//
// Fix: give the channel room for every sender, so they can all finish.
// Change the channel as shown in the "fix:" comment.
//
// Run: go run ./02-concurrency/02-channels/03-unbuffered/05-leak
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	ch := make(chan int) // fix: make(chan int, 3) gives every sender room

	for i := 0; i < 3; i++ {
		go send(ch, i)
	}

	fmt.Println("first:", <-ch)

	time.Sleep(100 * time.Millisecond)                 // give the others time to reach their send
	fmt.Println("goroutines:", runtime.NumGoroutine()) // 3 = main + 2 stuck (1 with the fix)
}

// send sends one value. It's stuck forever if nobody receives it.
func send(out chan<- int, v int) {
	out <- v
}
