// Bigger buffer: it only hides a deadlock, it doesn't fix it.
//
// sendThenRead sends every value first and reads them afterwards. That
// only works while all values fit in the buffer. With a buffer of 10, it
// works for 10 values and deadlocks on the 11th.
//
// Problem: the send loop and the read loop are in the same goroutine, one
// after the other. The 11th send waits for a reader, but the read loop
// can't start until sending is done. Each waits for the other: deadlock.
// A bigger buffer only changes when it breaks, so the bug comes back as
// soon as there's more data.
//
// Fix: put the send loop in its own goroutine. Then sending and reading
// run at the same time: when the buffer is full, the sender waits only
// until main reads one value. Any buffer size works, even 0. Use the
// "fix:" line instead of the line above it.
//
// Try: change 11 to 10. It "works", and that's the trap.
//
// Run: go run ./02-concurrency/02-channels/02-buffered/05-bigger-buffer
package main

import "fmt"

func main() {
	ch := make(chan int, 10) // made bigger to "fix" a deadlock

	sendThenRead(ch, 11)
	// sendWhileReading(ch, 11) // fix: use this instead
}

// sendThenRead sends everything first, then reads, all in one goroutine.
func sendThenRead(ch chan int, n int) {
	for i := 1; i <= n; i++ {
		ch <- i // the 11th send waits for a reader...
	}
	close(ch)

	for v := range ch { // ...but the reader only starts after the loop above: deadlock
		fmt.Println("got", v)
	}
}

// sendWhileReading sends from its own goroutine while the caller reads.
func sendWhileReading(ch chan int, n int) {
	go func() { // sender: runs at the same time as the loop below
		for i := 1; i <= n; i++ {
			ch <- i // when the buffer is full, waits until the reader takes one
		}
		close(ch)
	}()

	for v := range ch { // reader: starts right away, doesn't wait for the sender
		fmt.Println("got", v)
	}
}
