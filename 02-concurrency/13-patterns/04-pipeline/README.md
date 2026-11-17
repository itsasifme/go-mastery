# Pipeline

A chain of steps connected by channels. Each step takes values, changes them, and passes them on.

## Why

**Problem:** one big loop that reads, changes and saves each item does one thing at a time. While it waits on a slow save, nothing else moves.

**Solution:** split the work into steps, each in its own goroutine, connected by channels. Every step works on a different item at the same moment. A slow step can get more workers without changing the others.

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `func step(in <-chan A, out chan<- B)` | One step: read from `in`, write to `out` | [01-basic](01-basic/) |
| `defer close(out)` + `for v := range in` | A step ends when its input closes, then closes its output | [01-basic](01-basic/) |
| `select { case out <- v: case <-ctx.Done(): return }` | A step that stops when the pipeline is cancelled | [02-cancel](02-cancel/) |
| several workers on one step + `wg.Wait(); close(out)` | Speed up the slowest step | [03-parallel-step](03-parallel-step/) |
| `g.Go(func() error { return step(ctx, in, out) })` | Run each step with errgroup, so one error stops them all | [04-errors](04-errors/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/13-patterns/04-pipeline/01-basic` | An order report: generate → add tax → format → print |
| `02-cancel` | `go run ./02-concurrency/13-patterns/04-pipeline/02-cancel` | Search an endless log, stop after three errors, and every step shuts down |
| `03-parallel-step` | `go run ./02-concurrency/13-patterns/04-pipeline/03-parallel-step` | A slow geocoding step gets three workers |
| `04-errors` | `go run ./02-concurrency/13-patterns/04-pipeline/04-errors` | A CSV import hits a bad row: every step stops, and the import rolls back |

## Things to remember

- Each step closes its own output channel.
- If a later step stops early, **the earlier steps get stuck unless they also watch `ctx.Done()`** (`02`).
- **A pipeline is only as fast as its slowest step.** Give that step more workers (`03`).
- A step with several workers loses the order of items. Carry an index if order matters.
- **If one bad item should stop everything, run the steps with errgroup** and return the error. Commit the result only when `g.Wait()` returns nil (`04`).
