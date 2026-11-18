// Buffer: reuse buffers with sync.Pool instead of making one per request.
//
// A server builds a log line for every request in a bytes.Buffer. Without
// a pool, every request makes a new buffer, and the garbage collector has
// to clean up each one. With a pool, a request takes a buffer with Get,
// uses it, and gives it back with Put. New only runs when the pool is
// empty, so a handful of buffers serve every request.
//
// Each request calls Reset first, because a buffer from the pool still
// holds the last request's data (see 02-reset).
//
// Note: the pool can drop its buffers at any time, like during a garbage
// collection, without telling you. So New must always be able to make a
// fresh one. Use a pool only for things that are cheap to lose, never for
// connections or open files.
//
// Real use:
//   - fmt, encoding/json and net/http use pools inside, for print state,
//     encode buffers and bufio readers and writers.
//   - a JSON API or logger that builds each response or log line in a
//     pooled bytes.Buffer.
//
// Run: go run ./02-concurrency/14-pool/01-buffer
package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// requestsPerWorker is how many requests each worker handles.
const requestsPerWorker = 250

// buffersMade counts how many buffers were made, with or without the pool.
var buffersMade atomic.Int64

// bufPool hands out *bytes.Buffer values. It calls newBuffer when it's empty.
var bufPool = sync.Pool{New: newBuffer}

func main() {
	const workerCount = 4

	run(workerCount, logLineNew)
	fmt.Println("without a pool:", buffersMade.Load(), "buffers made for", workerCount*requestsPerWorker, "requests")

	buffersMade.Store(0)
	run(workerCount, logLinePooled)
	fmt.Println("with a pool:   ", buffersMade.Load(), "buffers made for", workerCount*requestsPerWorker, "requests")
}

// newBuffer makes a fresh buffer and counts it.
func newBuffer() any {
	buffersMade.Add(1)
	return new(bytes.Buffer)
}

// run starts workerCount workers that each handle requestsPerWorker
// requests with logLine, and waits for them.
func run(workerCount int, logLine func(request int)) {
	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(func() { work(logLine) })
	}
	wg.Wait()
}

// work handles requestsPerWorker requests.
func work(logLine func(request int)) {
	for request := range requestsPerWorker {
		logLine(request)
	}
}

// logLineNew builds the log line in a new buffer every time.
func logLineNew(request int) {
	buf := newBuffer().(*bytes.Buffer)
	writeLogLine(buf, request)
}

// logLinePooled builds the log line in a buffer from the pool, then puts
// it back.
func logLinePooled(request int) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset() // clear what the last user left in it
	writeLogLine(buf, request)
	bufPool.Put(buf)
}

// writeLogLine builds one log line in buf and writes it out.
func writeLogLine(buf *bytes.Buffer, request int) {
	fmt.Fprintf(buf, `{"level":"info","request":%d,"path":"/orders"}`+"\n", request)
	buf.WriteTo(io.Discard) // stands in for writing to the log file
}
