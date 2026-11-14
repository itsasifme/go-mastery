// No retry: Once keeps a failure forever.
//
// The first time the app connects to a payment API, the network is down,
// so connect fails. A moment later the network is fine again.
//
// With sync.OnceValues, the first result is final, even an error. Every
// later call gets the same old error, and the app can't recover without a
// restart. (sync.Once is the same: if the function panics, Do still counts
// it as done.)
//
// lazyClient fixes it with a mutex. It saves the client only when connect
// works. After a failure, the next call tries again.
//
// Real use:
//   - connecting to a database, cache or API on first use, when the first
//     try can fail for a moment (DNS not ready, service still starting).
//     Use Once only when a failure should stop the program, like a bad
//     config file.
//
// Run: go run ./02-concurrency/10-once/03-no-retry
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// client stands in for a real API client.
type client struct {
	attempt int // which connect attempt made it
}

// connector fails on its first attempt, like a short network outage.
type connector struct {
	attempts int
}

// connect returns a client, or an error on the first attempt.
func (c *connector) connect() (*client, error) {
	c.attempts++
	if c.attempts == 1 {
		return nil, errors.New("network is unreachable")
	}
	return &client{attempt: c.attempts}, nil
}

// lazyClient connects on first use, and tries again after a failure.
type lazyClient struct {
	mu      sync.Mutex // guards c
	c       *client
	connect func() (*client, error)
}

// Get returns the saved client, or connects. It saves the client only
// when connect works.
func (l *lazyClient) Get() (*client, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.c != nil {
		return l.c, nil
	}
	c, err := l.connect()
	if err != nil {
		return nil, err // nothing saved, so the next call tries again
	}
	l.c = c
	return c, nil
}

const (
	callCount = 3
	callGap   = 50 * time.Millisecond // time between two payments
)

func main() {
	fmt.Println("with sync.OnceValues:")
	onceGet := sync.OnceValues((&connector{}).connect)
	done := make(chan struct{}) // pay -> main: "I'm finished"
	go pay(onceGet, done)
	<-done

	fmt.Println("with lazyClient:")
	lazy := &lazyClient{connect: (&connector{}).connect}
	done = make(chan struct{})
	go pay(lazy.Get, done)
	<-done
}

// pay makes a few payments, getting the client from getClient each time.
func pay(getClient func() (*client, error), done chan<- struct{}) {
	defer close(done)

	for i := range callCount {
		c, err := getClient()
		if err != nil {
			fmt.Printf("  payment %d: failed: %v\n", i+1, err)
		} else {
			fmt.Printf("  payment %d: ok (client from attempt %d)\n", i+1, c.attempt)
		}
		time.Sleep(callGap)
	}
}
