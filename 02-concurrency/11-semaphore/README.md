# Semaphore

A semaphore limits how many goroutines can do something at the same time.

## Why

**Problem:** starting one goroutine per job is easy. But a thousand jobs at once can open too many connections, use too much memory, or get you blocked by an API.

**Solution:** a semaphore with a fixed number of slots. A job takes a slot before it starts and gives it back when it ends. When all slots are taken, the next job waits. A buffered channel works for simple cases. `semaphore.Weighted` also handles jobs of different sizes, and it can wait with a context.

## Import

```go
import "golang.org/x/sync/semaphore"
```

Not in the standard library. It comes from the `golang.org/x/sync` module. `go.mod` says which version this repo uses.

The channel semaphore needs no import. It is plain Go.

- `golang.org/x/sync/semaphore`: the package name is `semaphore`, the last part of the path, so you write `semaphore.NewWeighted`.

## API

### Built in: a buffered channel

| API | Kind | What it does | Example |
| --- | --- | --- | --- |
| `make(chan struct{}, n)` | built-in | A channel semaphore with `n` slots | [01-channel](01-channel/) |
| `sem <- struct{}{}` / `<-sem` | syntax | Take a slot / give it back | [01-channel](01-channel/) |

### `golang.org/x/sync/semaphore`

| API | Kind | Returns | What it does | Example |
| --- | --- | --- | --- | --- |
| `semaphore.NewWeighted(n)` | function | `*Weighted` | A semaphore with `n` units in total | [02-weighted](02-weighted/) |
| `s.Acquire(ctx, n)` | method | `error` | Wait for `n` units, or return `ctx.Err()` | [02-weighted](02-weighted/), [03-acquire-timeout](03-acquire-timeout/) |
| `s.Release(n)` | method | | Give `n` units back. Giving back more than you took panics | [02-weighted](02-weighted/) |
| `s.TryAcquire(n)` | method | `bool` | Take `n` units only if they're free right now | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-channel` | `go run ./02-concurrency/11-semaphore/01-channel` | Download files three at a time with a buffered channel |
| `02-weighted` | `go run ./02-concurrency/11-semaphore/02-weighted` | Videos of different sizes share a memory budget |
| `03-acquire-timeout` | `go run ./02-concurrency/11-semaphore/03-acquire-timeout` | Requests wait a short time for a slot to call a partner API, then answer 503 |

## Things to remember

- **In your own loop, take a slot before you start the goroutine.** Then only `n` goroutines exist at a time (`01`, `02`).
- **When the goroutine already exists, like an HTTP handler, take the slot inside it,** with a deadline (`03`).
- **Give the slot back with `defer`**, so a panic or an early `return` doesn't lose it.
- **Wait with a deadline in servers.** A request that waits forever is worse than a quick "busy" (`03`).
- **The limit is per server.** Fine for keeping this server's memory or CPU safe. For a partner API's limit for your whole account, keep the count in Redis (`03`). See [One server or many?](../README.md#one-server-or-many).
- `Weighted` is fair: when many goroutines wait, they get their units in the order they asked, so a big job isn't passed by small jobs behind it. `02` doesn't show this, because only main asks there.
- **`Acquire(ctx, n)` with `n` bigger than the total size can never succeed.** It waits until `ctx` is done, then returns its error. With `context.Background()` it waits forever.
- If you also need errors, `errgroup.SetLimit` is often simpler (see [errgroup/02-limit](../06-errgroup/02-limit/)).
