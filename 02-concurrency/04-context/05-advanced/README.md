# Advanced Context

Two tools for code that doesn't fit the normal "check `ctx.Done()`" shape.

See the [full API overview](../README.md#api).

## Why

**Problem 1:** some blocking calls don't take a ctx, like reading from a connection. Cancelling can't stop them.

**Solution:** `AfterFunc` runs a cleanup, like closing the connection, the moment ctx is done. That unblocks the call.

**Problem 2:** background work started by a request, like writing an audit log, gets cancelled when the request ends. Using `Background()` instead loses the request's values.

**Solution:** `WithoutCancel` keeps the values but drops the cancellation.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `AfterFunc(ctx, f)` | `stop func() bool` | 1.21 | Run `f` in a new goroutine when ctx is done | [01-after-func](01-after-func/) |
| `WithoutCancel(parent)` | `ctx` | 1.21 | Keep the values, drop the cancel and deadline | [02-without-cancel](02-without-cancel/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-after-func` | `go run ./02-concurrency/04-context/05-advanced/01-after-func` | Close a connection when ctx is done, to unblock a call with no ctx |
| `02-without-cancel` | `go run ./02-concurrency/04-context/05-advanced/02-without-cancel` | Background work that keeps the request's values but outlives it |

## Things to remember

- `AfterFunc` runs its function in a new goroutine. Call the returned `stop` if you don't need it anymore.
- **`WithoutCancel` has no deadline either.** Add your own `WithTimeout`.
- Both need Go 1.21 or newer.
