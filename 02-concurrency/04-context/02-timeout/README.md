# Context Timeout

Stop work after some time with `context.WithTimeout`.

See the [full API overview](../README.md#api).

## Why

**Problem:** a slow database or API call can hang a request forever. `time.After` in a `select` only stops the waiting, not the work, which keeps running ([select/04-timeout](../../03-select/04-timeout/)).

**Solution:** `WithTimeout` cancels the ctx after a set time. Work that watches `ctx.Done()` stops itself.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `WithTimeout(parent, d)` | `ctx, CancelFunc` | 1.7 | Cancel after `d` | [01-basic](01-basic/), [02-result-or-timeout](02-result-or-timeout/) |
| `WithTimeoutCause(parent, d, cause)` | `ctx, CancelFunc` | 1.21 | Same, with a reason for `Cause` | [03-timeout-cause](03-timeout-cause/) |
| `context.DeadlineExceeded` | `error` | 1.7 | `ctx.Err()` after the timeout | [01-basic](01-basic/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/04-context/02-timeout/01-basic` | A slow function that stops itself when the timeout hits |
| `02-result-or-timeout` | `go run ./02-concurrency/04-context/02-timeout/02-result-or-timeout` | Wait for a goroutine's result or the timeout, and stop the goroutine too |
| `03-timeout-cause` | `go run ./02-concurrency/04-context/02-timeout/03-timeout-cause` | Say which timeout was hit with `WithTimeoutCause` |

## Things to remember

- **Still call `cancel`, even when the timeout already fired.**
- A timeout error is `context.DeadlineExceeded`. Check it with `errors.Is`.
- The worker must also watch `ctx.Done()`, including when it sends a result. **Otherwise it leaks** (`02`).
- `Cause(ctx)` gives the error passed to the `...Cause` function. **The returned `cancel` takes no cause**, so a manual cancel gives `context.Canceled` (`03`).
