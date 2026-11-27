// Reset: a buffer from the pool still holds the last user's data.
//
// A handler builds each user's JSON response in a pooled buffer. Put
// doesn't clear anything, so the next Get usually returns the same buffer,
// with the last response still in it.
//
// Problem: the handler forgets to call Reset. Bob's response starts with
// Alice's email, and Carol's has both. That's one user's private data
// sent to another user.
//
// Fix: call buf.Reset() right after Get. Uncomment the "fix:" line.
//
// Real use:
//   - every pooled bytes.Buffer, slice or struct in an HTTP handler. Reset
//     it after Get (or before Put), so nothing leaks from one request into
//     the next.
//
// Run: go run ./02-concurrency/14-pool/02-reset
package main

import (
	"bytes"
	"fmt"
	"sync"
)

// bufPool hands out *bytes.Buffer values. It calls newBuffer when it's empty.
var bufPool = sync.Pool{New: newBuffer}

func main() {
	users := []string{"alice", "bob", "carol"}

	var wg sync.WaitGroup
	wg.Go(func() { serve(users) })
	wg.Wait()
}

// newBuffer makes a fresh buffer.
func newBuffer() any {
	return new(bytes.Buffer)
}

// serve handles one request per user, one after another, like a single
// worker would.
func serve(users []string) {
	for _, user := range users {
		fmt.Printf("%-5s gets: %s\n", user, respond(user))
	}
}

// respond builds the user's JSON response in a pooled buffer.
func respond(user string) string {
	buf := bufPool.Get().(*bytes.Buffer)
	// buf.Reset() // fix: clear what the last request left in the buffer
	defer bufPool.Put(buf)

	fmt.Fprintf(buf, `{"email":%q}`, user+"@example.com")
	return buf.String() // a copy, so it stays valid after Put gives buf to someone else
}
