# WaitGroup

`sync.WaitGroup` lets you wait until a group of goroutines finish.

## Why

**Problem:** waiting for many goroutines with channels means one `done` signal per goroutine, and counting the receives by hand (like [channels/02-closing/02-broadcast](../02-channels/02-closing/02-broadcast/)).

**Solution:** `sync.WaitGroup`, a counter. Add to it before starting goroutines, each one marks itself done, and `Wait` blocks until the count is 0.

## API

| API | Go | What it does | Example |
| --- | --- | --- | --- |
| `var wg sync.WaitGroup` | 1.0 | Ready to use, no setup | [01-basic](01-basic/) |
| `wg.Add(n)` | 1.0 | Expect `n` more goroutines | [01-basic](01-basic/), [05-add-too-late](05-add-too-late/) |
| `wg.Done()` | 1.0 | One goroutine finished | [01-basic](01-basic/) |
| `wg.Wait()` | 1.0 | Block until the count is 0 | [01-basic](01-basic/), [03-close-results](03-close-results/) |
| `wg.Go(f)` | 1.25 | `Add(1)`, `go f()` and `Done()` in one call | [06-wg-go](06-wg-go/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/05-waitgroup/01-basic` | Wait for two goroutines before the program ends |
| `02-parallel-fetch` | `go run ./02-concurrency/05-waitgroup/02-parallel-fetch` | Load profile, orders and settings at once; total time is the slowest call |
| `03-close-results` | `go run ./02-concurrency/05-waitgroup/03-close-results` | Check many sites at once and close the results channel after the last one |
| `04-graceful-shutdown` | `go run ./02-concurrency/05-waitgroup/04-graceful-shutdown` | Let running jobs finish before the program exits |
| `05-add-too-late` | `go run ./02-concurrency/05-waitgroup/05-add-too-late` | Problem: `Add` inside the goroutine lets `Wait` return too early |
| `06-wg-go` | `go run ./02-concurrency/05-waitgroup/06-wg-go` | `wg.Go` does `Add`, `go` and `Done` in one call, so `Add` can't come too late |

## Things to remember

- **Call `Add` before you start the goroutine.** Inside it is too late (see [05-add-too-late](05-add-too-late/)). On Go 1.25+, `wg.Go` does this for you (see [06-wg-go](06-wg-go/)).
- **Call `Done` with `defer`, so it always runs.**
- **Pass the WaitGroup as a pointer** (`*sync.WaitGroup`). A copy has its own count, so `Done` on it never reaches `Wait`. `go vet` warns about this.
- **To close a shared results channel, call `wg.Wait()` then `close` in one extra goroutine** (see [03-close-results](03-close-results/)).
- Goroutines can write to different fields or slice indexes without a lock. Read them only after `Wait`.
- The print order can change between runs.
