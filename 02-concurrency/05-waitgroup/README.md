# WaitGroup

`sync.WaitGroup` lets you wait until a group of goroutines finish.

## Why

**Problem:** waiting for many goroutines with channels means one `done` signal per goroutine, and counting the receives by hand (like [channels/02-closing/02-broadcast](../02-channels/02-closing/02-broadcast/)).

**Solution:** `sync.WaitGroup`, a counter. Add to it before starting goroutines, each one marks itself done, and `Wait` blocks until the count is 0.

## API

| API | Go | What it does | Example |
| --- | --- | --- | --- |
| `var wg sync.WaitGroup` | 1.0 | Ready to use, no setup | [01-basic](01-basic/) |
| `wg.Add(n)` | 1.0 | Expect `n` more goroutines | [01-basic](01-basic/) |
| `wg.Done()` | 1.0 | One goroutine finished | [01-basic](01-basic/) |
| `wg.Wait()` | 1.0 | Block until the count is 0 | [01-basic](01-basic/) |
| `wg.Go(f)` | 1.25 | `Add(1)`, `go f()` and `Done()` in one call | none (this repo uses Go 1.23) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/05-waitgroup/01-basic` | Wait for two goroutines before the program ends |

## Things to remember

- **Call `Add` before you start the goroutine.**
- **Call `Done` with `defer`, so it always runs.**
- The print order can change between runs.
