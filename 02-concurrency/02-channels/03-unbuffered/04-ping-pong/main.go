// Ping-pong: main and a goroutine take strict turns.
//
// Each send waits until the other side receives, so the output is always
// in the same order: ping 1, pong 2, ping 3, pong 4, ping 5, pong 6.
//
// Real use:
//   - one goroutine owns a resource (a connection, a cache) and serves
//     requests one at a time: send a request, wait for its reply.
//
// Run: go run ./02-concurrency/02-channels/03-unbuffered/04-ping-pong
package main

import "fmt"

// rounds is how many times the ball goes back and forth.
// main and pong must use the same number, or one of them waits forever.
const rounds = 3

func main() {
	ball := make(chan int) // one channel both ways: safe only because it's unbuffered

	go pong(ball)

	n := 1
	for range rounds {
		fmt.Println("ping", n)
		ball <- n      // waits until pong receives
		n = <-ball + 1 // wait for pong, then add 1 for the next ping
	}
}

// pong receives the ball, adds 1, and sends it back, rounds times.
// It takes a two-way chan int, because it both receives and sends.
func pong(ball chan int) {
	for range rounds {
		n := <-ball // wait for ping
		fmt.Println("pong", n+1)
		ball <- n + 1 // waits until main receives
	}
}
