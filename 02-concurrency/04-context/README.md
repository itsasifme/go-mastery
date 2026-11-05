# Context

A `context.Context` tells goroutines when to stop. It can carry a cancel signal, a timeout, or a deadline.

1. [Cancellation](01-cancellation/)
2. [Timeout](02-timeout/)
3. [Deadline](03-deadline/)

## Things to remember

- Pass `ctx` as the first parameter. Don't keep it inside a struct.
- Always call `cancel`, usually with `defer`. `go vet` warns you if you forget.
- Only put request data in context values, not normal function arguments.

## Go version notes

- `WithCancelCause` came in Go 1.20. `WithTimeoutCause`, `AfterFunc`, and `WithoutCancel` came in Go 1.21.
