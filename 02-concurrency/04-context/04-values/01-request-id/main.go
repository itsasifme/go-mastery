// Request ID: carry request data down the call chain with context values.
//
// handleRequest puts a request ID in ctx once. Every function below it
// gets ctx anyway, so logf can read the ID without a new parameter on
// every function.
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
)

// requestIDKey is the context key for the request ID.
type requestIDKey struct{}

func main() {
	ctx := context.WithValue(context.Background(), requestIDKey{}, "req-42")
	handleRequest(ctx)
}

// handleRequest stands in for an HTTP handler.
func handleRequest(ctx context.Context) {
	logf(ctx, "handling request")
	loadUser(ctx, 7)
}

// loadUser stands in for a database call.
func loadUser(ctx context.Context, userID int) {
	logf(ctx, "loading user %d", userID)
}

// logf prints a formatted message with the request ID from ctx.
func logf(ctx context.Context, format string, args ...any) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		id = "no-id"
	}
	fmt.Printf("[%s] %s\n", id, fmt.Sprintf(format, args...))
}
