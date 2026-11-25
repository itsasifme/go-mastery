# Semaphore

A semaphore limits how many goroutines can do something at the same time.

## Why

**Problem:** starting one goroutine per job is easy. But a thousand jobs at once can open too many connections, use too much memory, or get you blocked by an API.

**Solution:** a semaphore with a fixed number of slots. A job takes a slot before it starts and gives it back when it ends. When all slots are taken, the next job waits. A buffered channel works for simple cases. `semaphore.Weighted` also handles jobs of different sizes, and it can wait with a context.

## API

| API | Kind | Returns | What it does | Example |
| --- | --- | --- | --- | --- |
| `make(chan struct{}, n)` | built-in | | A channel semaphore with `n` slots | [01-channel](01-channel/) |
| `sem <- struct{}{}` / `<-sem` | syntax | | Take a slot / give it back | [01-channel](01-channel/) |
| `semaphore.NewWeighted(n)` | function | `*Weighted` | A semaphore with `n` units in total | [02-weighted](02-weighted/) |
| `s.Acquire(ctx, n)` | method | `error` | Wait for `n` units, or return `ctx.Err()` | [02-weighted](02-weighted/), [03-acquire-timeout](03-acquire-timeout/) |
| `s.Release(n)` | method | | Give `n` units back | [02-weighted](02-weighted/) |
| `s.TryAcquire(n)` | method | `bool` | Take `n` units only if they're free right now | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-channel` | `go run ./02-concurrency/11-semaphore/01-channel` | Download files three at a time with a buffered channel |
| `02-weighted` | `go run ./02-concurrency/11-semaphore/02-weighted` | Videos of different sizes share a memory budget |
| `03-acquire-timeout` | `go run ./02-concurrency/11-semaphore/03-acquire-timeout` | Requests wait a short time for a slot to call a partner API, then answer 503 |

## Things to remember

- **Take a slot before you start the goroutine, not inside it.** Then only `n` goroutines exist at a time (`01`).
- **Give the slot back with `defer`**, so a panic or an early `return` doesn't lose it.
- **Wait with a deadline in servers.** A request that waits forever is worse than a quick "busy" (`03`).
- `Weighted` is fair: a big job that's waiting isn't passed by small jobs behind it.
- **`Acquire(ctx, n)` with `n` bigger than the total size can never succeed.** It waits until `ctx` is done, then returns its error. With `context.Background()` it waits forever.
- If you also need errors, `errgroup.SetLimit` is often simpler (see [errgroup/02-limit](../06-errgroup/02-limit/)).

## Package note

`semaphore` comes from `golang.org/x/sync`, not the standard library. This repo uses `v0.16.0`.
