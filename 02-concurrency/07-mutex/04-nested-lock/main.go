// Nested lock: a method that locks calls another method that locks.
//
// stats counts cache hits and misses. HitRate locks the mutex. Report also
// locks it, then calls HitRate. In Go, a mutex doesn't know who holds it,
// so the second Lock in the same goroutine waits for the first to be
// unlocked. That never happens.
//
// Problem: report gets stuck forever. Nothing else is running, so Go
// stops the program with "fatal error: all goroutines are asleep -
// deadlock!". In a real server, other goroutines keep running, so there's
// no error. The requests just hang.
//
// Fix: move the work into a helper that doesn't lock (hitRateLocked), and
// say in its comment that the caller must hold the lock. Both methods call
// it. Comment out the line marked "problem", then uncomment the "fix:" line.
//
// Real use:
//   - any type with a mutex and several methods: Get calls a validate
//     method, Save calls Get, and so on. The "...Locked" helper is a
//     common naming habit in Go code, including the standard library.
//
// Run: go run ./02-concurrency/07-mutex/04-nested-lock
package main

import (
	"fmt"
	"sync"
)

// stats counts cache hits and misses.
type stats struct {
	mu     sync.Mutex // guards hits and misses
	hits   int
	misses int
}

// HitRate returns the share of lookups that were hits, from 0 to 1.
func (s *stats) HitRate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hitRateLocked()
}

// hitRateLocked does the work of HitRate. The caller must hold s.mu.
func (s *stats) hitRateLocked() float64 {
	total := s.hits + s.misses
	if total == 0 {
		return 0
	}
	return float64(s.hits) / float64(total)
}

// Report returns a one-line summary.
func (s *stats) Report() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	rate := s.HitRate() // problem: HitRate locks s.mu again, and we already hold it
	// rate := s.hitRateLocked() // fix: call the helper that expects the lock to be held
	return fmt.Sprintf("hits=%d misses=%d rate=%.2f", s.hits, s.misses, rate)
}

func main() {
	s := &stats{hits: 8, misses: 2}

	done := make(chan struct{}) // report -> main: "I'm finished"
	go report(s, done)
	<-done // a stuck report shows up here as a deadlock
}

// report prints the stats summary.
func report(s *stats, done chan<- struct{}) {
	defer close(done)
	fmt.Println(s.Report())
}
