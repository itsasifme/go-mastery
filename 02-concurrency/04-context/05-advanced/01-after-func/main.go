// AfterFunc: run a function when a context is done.
//
// Some blocking calls don't take a ctx, like reading from an old-style
// connection. conn.Read here waits forever, because nobody sends.
// context.AfterFunc runs c.Close in its own goroutine once ctx times out.
// Closing the connection unblocks Read, so the call stops on time anyway.
//
// stop cancels the hook if it hasn't run yet, like when the read finishes
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

	c := &conn{data: make(chan string)}
	stop := context.AfterFunc(ctx, c.Close) // close the connection when ctx is done
	defer stop()

	_, ok := c.Read() // blocks: nobody sends, and Read doesn't know about ctx
	fmt.Println("read returned, ok =", ok, "because:", ctx.Err())
}
