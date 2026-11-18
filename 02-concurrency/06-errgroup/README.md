# errgroup

`errgroup` is like a WaitGroup that also returns the first error, and can cancel the other goroutines when one fails.

## Why

**Problem:** a WaitGroup only waits. If a goroutine fails, you need your own error variable, a mutex to set it safely, and your own way to stop the others.

**Solution:** `errgroup.Group`. Start each job with `g.Go`, and `g.Wait` returns the first error. `WithContext` also gives you a `ctx` that's cancelled when a job fails, and `SetLimit` caps how many jobs run at once.

## API

| API | Kind | Returns | What it does | Example |
| --- | --- | --- | --- | --- |
| `var g errgroup.Group` | type | | Ready to use, no setup | [01-basic](01-basic/) |
| `g.Go(f)` | method | | Run `f func() error` in a new goroutine | [01-basic](01-basic/) |
| `g.Wait()` | method | `error` | Wait for all jobs, return the first error | [01-basic](01-basic/) |
| `g.SetLimit(n)` | method | | At most `n` jobs at once. `g.Go` waits for a free slot | [02-limit](02-limit/) |
| `g.TryGo(f)` | method | `bool` | Like `Go`, but returns `false` instead of waiting when the limit is full | none |
| `errgroup.WithContext(ctx)` | function | `*Group, ctx` | `ctx` is cancelled when a job fails (or when `Wait` returns) | [03-with-context](03-with-context/), [04-collect-results](04-collect-results/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/06-errgroup/01-basic` | Startup checks run at once, and `Wait` returns the one that failed |
| `02-limit` | `go run ./02-concurrency/06-errgroup/02-limit` | Resize photos two at a time with `SetLimit` |
| `03-with-context` | `go run ./02-concurrency/06-errgroup/03-with-context` | One download fails, so the others stop early |
| `04-collect-results` | `go run ./02-concurrency/06-errgroup/04-collect-results` | Fetch every cart price at once, write each to its own slice index, or fail as a whole |

## Things to remember

- `Wait` waits for all goroutines, then **returns the first error**. Later errors are dropped.
- **The context from `WithContext` is cancelled after the first error.** Your jobs must watch `ctx.Done()`, or they keep running (`03`).
- **Don't use that `ctx` after `Wait` returns.** It's always cancelled by then.
- `g.Go` takes a `func() error`. Keep the closure to one line that calls a named function.
- Each goroutine can write to its own slice index without a lock. Read the slice only after `Wait` (`04`).
- If you want every job to run even when some fail, use a worker pool that collects errors instead (see [patterns/01-worker-pool/03-errors](../15-patterns/01-worker-pool/03-errors/)).

## Package note

`errgroup` isn't in the standard library. It comes from `golang.org/x/sync`. This repo uses `v0.16.0`.
