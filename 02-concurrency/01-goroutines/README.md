# Goroutines

A goroutine is a function that runs at the same time as the rest of the program. Start one with `go f()`.

## Why

**Problem:** a program that does one thing at a time waits on every slow step, like a network call. OS threads can run things at the same time, but they're heavy: about 1 MB each and slow to start, so you can't have thousands.

**Solution:** goroutines. They start with about 2 KB and are scheduled by the Go runtime, so thousands or even millions are fine. Start one with `go f()`.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `go f(x)` | Run `f` in a new goroutine | [01-basic](01-basic/), [02-many](02-many/) |
| `go func() { ... }()` | Run an anonymous function in a new goroutine | [03-loop-variable](03-loop-variable/) |
| `runtime.NumGoroutine()` | How many goroutines exist right now | [02-many](02-many/) |
| `defer` + `recover()` | Stop a panic, only in the same goroutine | [04-panic](04-panic/), [05-recover](05-recover/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/01-goroutines/01-basic` | Start a goroutine without waiting for it |
| `02-many` | `go run ./02-concurrency/01-goroutines/02-many` | 5 goroutines print in a different order every run |
| `03-loop-variable` | `go run ./02-concurrency/01-goroutines/03-loop-variable` | Since Go 1.22, each loop round gets its own variable |
| `04-panic` | `go run ./02-concurrency/01-goroutines/04-panic` | A panic in a goroutine crashes the whole program |
| `05-recover` | `go run ./02-concurrency/01-goroutines/05-recover` | Recover inside each goroutine, so one failing job doesn't crash the rest |

## Things to remember

- **`main` doesn't wait for goroutines.** When `main` ends, the program ends.
- `time.Sleep` is a hack used here to wait. The proper way is in [WaitGroup](../05-waitgroup/).
- **Never rely on the order goroutines run in.**
- **`recover` only catches panics in its own goroutine.** Uncomment the `// fix:` line in `04` to see it.
- In `05`, each job recovers its own panic. Use it as a safety net for bugs, not instead of returning errors.
- For `defer`, `panic` and `recover` in general (not just goroutines), see [defer, panic, and recover](../../01-fundamentals/12-defer-panic-recover/).
