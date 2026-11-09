// WithoutCancel: keep the values, but not the cancellation.
//
// A request ctx is cancelled as soon as the request ends. The handler
// starts writeAudit in the background, and then the request ends right
// away.
//
// Problem: writeAudit uses the request ctx, so it's cancelled with the
// request and the audit record is lost.
//
// Fix: give it context.WithoutCancel(ctx). It keeps the values, like the
// request ID, but is never cancelled by the parent. Use the "fix:" line
// instead of the line above it.
//
// Note: a detached ctx has no deadline either. In real code, add your own
// WithTimeout, so the background work can't run forever.
//
// Real use:
//   - after answering a request: write an audit log, send an email, or
//     update a cache. The work must finish even though the client already
//     got its response.
//
// Run: go run ./02-concurrency/04-context/05-advanced/02-without-cancel
package main

import (
	"context"
	"fmt"
	"time"
)

// requestIDKey is the context key for the request ID
// (see context/04-values/01-request-id).
type requestIDKey struct{}

func main() {
	base := context.WithValue(context.Background(), requestIDKey{}, "req-42")
	ctx, cancel := context.WithCancel(base)
	done := make(chan struct{})

	go writeAudit(ctx, done) // problem
	// go writeAudit(context.WithoutCancel(ctx), done) // fix: detach from the request

	cancel() // the request ends right away
	<-done
}

// writeAudit takes 200ms to save an audit record, unless ctx stops it.
func writeAudit(ctx context.Context, done chan<- struct{}) {
	defer close(done)

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("audit saved for", ctx.Value(requestIDKey{}))
	case <-ctx.Done():
		fmt.Println("audit lost:", ctx.Err())
	}
}
