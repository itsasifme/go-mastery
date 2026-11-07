// Ping-pong: main and a goroutine take strict turns.
//
// Each send waits until the other side receives, so the output is always
// in the same order: ping 1, pong 2, ping 3, pong 4, ping 5, pong 6.
//
// Run: go run ./02-concurrency/02-channels/01-unbuffered/04-ping-pong
package main

import "fmt"

func main() {
	ball := make(chan int) // one channel both ways: safe only because it's unbuffered

	go func() {
		for i := 0; i < 3; i++ {
			n := <-ball // wait for ping
			fmt.Println("pong", n+1)
			ball <- n + 1 // waits until main receives
		}
	}()

	n := 1
	for i := 0; i < 3; i++ {
		fmt.Println("ping", n)
		ball <- n      // waits until the goroutine receives
		n = <-ball + 1 // wait for pong, then add 1 for the next ping
	}
}
