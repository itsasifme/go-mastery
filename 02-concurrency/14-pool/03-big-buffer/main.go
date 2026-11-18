// Big buffer: don't put a huge buffer back in the pool.
//
// Most requests build a small response, but one CSV export builds 4 MB. Its
// buffer grows to fit, and Put keeps that big buffer in the pool. Later
// small requests get it, so the 4 MB stays in memory for nothing. With many
// workers and a few exports, a server can hold many of these.
//
// Problem: after the export, the next small request gets a buffer with
// about 4 MB reserved.
//
// Fix: putBufferCapped drops buffers bigger than maxKeep instead of putting
// them back. The garbage collector frees them, and the pool makes a new
// small buffer when it needs one. fmt does the same: it doesn't keep print
// buffers bigger than 64 KB. Use the "fix:" line instead of the line above
// it.
//
// Real use:
//   - pooled buffers in servers where a few requests are much bigger than
//     the rest, like exports, file uploads or large JSON responses.
//
// Run: go run ./02-concurrency/14-pool/03-big-buffer
package main

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
)

// request is one response to build: a name for the output, and its size.
type request struct {
	name string
	size int
}

// Sizes, in bytes.
const (
	smallSize  = 200     // a normal response
	exportSize = 4 << 20 // a 4 MB CSV export
	maxKeep    = 64 << 10
)

// bufPool hands out *bytes.Buffer values. It calls newBuffer when it's empty.
var bufPool = sync.Pool{New: newBuffer}

func main() {
	requests := []request{
		{"small", smallSize},
		{"export", exportSize},
		{"small", smallSize},
	}

	var wg sync.WaitGroup
	wg.Go(func() { serve(requests) })
	wg.Wait()
}

// newBuffer makes a fresh buffer.
func newBuffer() any {
	return new(bytes.Buffer)
}

// serve handles the requests one after another, like a single worker would.
func serve(requests []request) {
	for _, r := range requests {
		handle(r.name, r.size)
	}
}

// handle builds a response of size bytes in a pooled buffer, and shows how
// much room the buffer it got already had.
func handle(name string, size int) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	fmt.Printf("%-6s got a buffer with %7d bytes reserved\n", name, buf.Cap())

	buf.WriteString(strings.Repeat("x", size)) // build the response

	putBuffer(buf) // problem: keeps the buffer, however big it grew
	// putBufferCapped(buf) // fix: drop buffers bigger than maxKeep
}

// putBuffer puts buf back in the pool.
func putBuffer(buf *bytes.Buffer) {
	bufPool.Put(buf)
}

// putBufferCapped puts buf back only if it's at most maxKeep bytes.
func putBufferCapped(buf *bytes.Buffer) {
	if buf.Cap() > maxKeep {
		return // too big to keep: let the garbage collector free it
	}
	bufPool.Put(buf)
}
