// Request ID: carry request data down the call chain with context values.
//
// main starts two requests at the same time, each in its own goroutine,
// like net/http does. Each request puts its own ID in its ctx once. Every
// function below it gets ctx anyway, so logf can read the ID without a new
// parameter on every function. The log lines of the two requests mix, but
// each line still shows the right ID.
//
// The key is a private type, requestIDKey. No other package can create
// the same key, so nobody can overwrite your value by accident. Don't use
// a plain string like "id" as the key.
//
// Note: the user ID is still a normal parameter. Values are only for data
// about the request (IDs, the logged-in user, tracing), never for normal
// inputs a function needs.
//
// Try: call logf with context.Background(). It prints "no-id".
//
// Real use:
//   - a request ID or trace ID in every log line, so you can find all
//     logs for one request.
//   - middleware that checks the login and stores the user in ctx for
//     the handlers.
//
// Run: go run ./02-concurrency/04-context/04-values/01-request-id
package main

import (
	"context"
	"fmt"
	"time"
)

// requestIDKey is the context key for the request ID.
type requestIDKey struct{}

// queryTime is how long loadUser's database call takes.
const queryTime = 50 * time.Millisecond

func main() {
	const requestCount = 2 // must match the number of "go handleRequest" lines below

	done := make(chan struct{}) // requests -> main: "I'm finished"

	go handleRequest(context.WithValue(context.Background(), requestIDKey{}, "req-42"), 7, done)
	go handleRequest(context.WithValue(context.Background(), requestIDKey{}, "req-43"), 9, done)

	for range requestCount {
		<-done // wait for every request
	}
}

// handleRequest stands in for an HTTP handler.
func handleRequest(ctx context.Context, userID int, done chan<- struct{}) {
	defer func() { done <- struct{}{} }() // send, not close: others share done

	logf(ctx, "handling request")
	loadUser(ctx, userID)
	logf(ctx, "sending response")
}

// loadUser stands in for a database call.
func loadUser(ctx context.Context, userID int) {
	logf(ctx, "loading user %d", userID)
	time.Sleep(queryTime) // simulate the query, so the two requests overlap
}

// logf prints a formatted message with the request ID from ctx.
func logf(ctx context.Context, format string, args ...any) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		id = "no-id"
	}
	fmt.Printf("[%s] %s\n", id, fmt.Sprintf(format, args...))
}
