# Context Values

Carry request data, like a request ID, down the call chain with `context.WithValue`.

See the [full API overview](../README.md#api).

## Why

**Problem:** data about the request, like a request ID or the logged-in user, is needed deep in the call chain. Adding a parameter to every function is noisy, and a global variable mixes up requests that run at the same time.

**Solution:** `WithValue` attaches the data to the ctx, which every function already gets.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `WithValue(parent, key, val)` | `ctx` | 1.7 | Attach one key and value | [01-request-id](01-request-id/) |
| `ctx.Value(key)` | `any` | 1.7 | Find the value, searching up the parents | [01-request-id](01-request-id/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-request-id` | `go run ./02-concurrency/04-context/04-values/01-request-id` | Two requests run at the same time, and every log line gets the right request ID without an extra parameter |

## Things to remember

- Use a private key type, like `type requestIDKey struct{}`. **Never a plain string.**
- **Only request data**: IDs, the logged-in user, tracing. Normal inputs stay normal parameters.
- `ctx.Value` returns `any`. Check the type with `v, ok := ...(string)`.
