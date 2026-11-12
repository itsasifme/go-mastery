// AfterFunc: run a function when a context is done.
//
// Some blocking calls don't take a ctx, like reading from an old-style
// connection. readMessages reads from conn in a loop in its own goroutine.
// It handles the messages that are already there, then conn.Read waits,
// because no more come. context.AfterFunc runs c.Close in its own
// goroutine once ctx times out. Closing the connection unblocks Read, so
// the loop stops on time anyway.
//
// stop cancels the hook if it hasn't run yet, like when the work finishes
// first.
//
// Real use:
//   - closing a net.Conn or setting its deadline when the request is
//     cancelled, so a blocked Read or Write returns.
//   - waking goroutines that wait on a sync.Cond, which has no ctx
//     support.
//
// Run: go run ./02-concurrency/04-context/05-advanced/01-after-func
package main

import (
	"context"
	"fmt"
	"time"
)

// conn stands in for a connection whose Read doesn't take a ctx.
type conn struct {
	data chan string
}

// Read waits for data. ok is false once the connection is closed.
func (c *conn) Read() (msg string, ok bool) {
	msg, ok = <-c.data
	return msg, ok
}

// Close closes the connection, which unblocks Read.
func (c *conn) Close() {
	fmt.Println("conn: closed")
	close(c.data)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	messages := []string{"hello", "how are you?"}
	c := &conn{data: make(chan string, len(messages))} // room for the messages already received
	for _, m := range messages {
		c.data <- m
	}

	stop := context.AfterFunc(ctx, c.Close) // close the connection when ctx is done
	defer stop()

	done := make(chan struct{}) // readMessages -> main: "I'm finished"

	go readMessages(ctx, c, done)

	<-done // wait until the read loop stopped
}

// readMessages handles each message until the connection is closed.
func readMessages(ctx context.Context, c *conn, done chan<- struct{}) {
	defer close(done) // tells main we're finished, however we return

	for {
		msg, ok := c.Read() // blocks when no message is waiting, and Read doesn't know about ctx
		if !ok {
			fmt.Println("read loop: stopped because:", ctx.Err())
			return
		}
		fmt.Println("read loop: got", msg)
	}
}
