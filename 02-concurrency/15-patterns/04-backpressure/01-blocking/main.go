// Blocking: a small buffer makes the fast side wait for the slow side.
//
// An app writes log lines much faster than the shipper can send them to
// the log server. Between them is a buffer of bufferSize lines. At first,
// writes are instant, because the buffer has room. Once it's full, each
// write waits until the shipper takes a line out. The app is slowed down
// to the shipper's speed. That push back is called backpressure.
//
// The "waited" column shows it: zero while the buffer fills, then about
// shipTime for each line.
//
// Try: set bufferSize to logCount. The app never waits, but it needs room
// for every line in memory. With no limit at all, a slow shipper would
// let the buffer grow until the program runs out of memory.
//
// Real use:
//   - a log or metrics shipper that must never lose data, so it's fine to
//     slow the app down a little.
//   - a bounded job queue between an HTTP handler and background workers.
//
// Run: go run ./02-concurrency/15-patterns/04-backpressure/01-blocking
package main

import (
	"fmt"
	"time"
)

const (
	logCount = 8
	shipTime = 50 * time.Millisecond // the shipper is slow
	roundTo  = 10 * time.Millisecond
)

func main() {
	const bufferSize = 3

	logs := make(chan string, bufferSize) // app -> shipper
	done := make(chan struct{})           // shipper -> main: "I'm finished"

	go ship(logs, done)
	writeLogs(logs) // main plays the app
	<-done
}

// writeLogs writes logCount lines as fast as it can, and prints how long
// each write had to wait. Then it closes logs.
func writeLogs(logs chan<- string) {
	defer close(logs)
	for i := range logCount {
		start := time.Now()
		logs <- fmt.Sprintf("line %d", i+1) // blocks while the buffer is full
		fmt.Printf("app:     wrote line %d (waited %v)\n", i+1, time.Since(start).Round(roundTo))
	}
}

// ship sends each line to the fake log server, until logs is closed.
func ship(logs <-chan string, done chan<- struct{}) {
	defer close(done)
	for line := range logs {
		time.Sleep(shipTime) // simulate the network call
		fmt.Println("shipper: sent", line)
	}
}
