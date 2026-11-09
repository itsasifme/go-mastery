# Context Cancellation

Stop work by hand with `context.WithCancel`.

See the [full API overview](../README.md#api).

## Why

**Problem:** Go has no way to stop a goroutine from outside. You had to build your own `stop` channel and pass it to every function (like [select/05-for-select](../../03-select/05-for-select/)).

**Solution:** `WithCancel` gives you a `ctx` and a `cancel` function. Any code that takes `ctx` can stop when `ctx.Done()` closes.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `WithCancel(parent)` | `ctx, CancelFunc` | 1.7 | Cancel by hand | [01-basic](01-basic/), [02-stop-workers](02-stop-workers/), [04-parent-child](04-parent-child/) |
| `WithCancelCause(parent)` | `ctx, CancelCauseFunc` | 1.20 | Cancel with a reason | [03-cancel-cause](03-cancel-cause/) |
| `Cause(ctx)` | `error` | 1.20 | Read the reason | [03-cancel-cause](03-cancel-cause/) |
| `ctx.Done()` | `<-chan struct{}` | 1.7 | Closed on cancel | [01-basic](01-basic/) |
| `ctx.Err()` | `error` | 1.7 | `context.Canceled` after cancel | [01-basic](01-basic/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/04-context/01-cancellation/01-basic` | A goroutine that stops when `ctx.Done()` closes |
| `02-stop-workers` | `go run ./02-concurrency/04-context/01-cancellation/02-stop-workers` | One `cancel` call stops many workers |
| `03-cancel-cause` | `go run ./02-concurrency/04-context/01-cancellation/03-cancel-cause` | Say why it was cancelled with `WithCancelCause`, and clean up differently for each reason |
| `04-parent-child` | `go run ./02-concurrency/04-context/01-cancellation/04-parent-child` | Cancelling a parent stops its children, not the other way |

## Things to remember

- **A goroutine has to check `ctx.Done()` itself.** Go won't stop it for you (`01`).
- `ctx.Err()` only says "canceled". `context.Cause(ctx)` gives the real reason (`03`).
- **Cancelling a parent also cancels all its children, but not the other way around** (`04`).
