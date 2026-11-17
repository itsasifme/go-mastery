# Worker Pool

A fixed number of workers take jobs from one channel and send results to another.

## Why

**Problem:** one goroutine per job is fine for ten jobs, but not for a million. You'd start too many goroutines at once and overload the CPU, memory, or the service you're calling.

**Solution:** start a set number of workers. They all read from one jobs channel, so each job goes to the next free worker. The number of workers is your limit, however many jobs there are.

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `for job := range jobs` | A worker takes jobs until the channel is closed | [01-basic](01-basic/) |
| `close(jobs)` | The sender says "no more jobs", so workers stop | [01-basic](01-basic/) |
| `wg.Wait(); close(results)` | One goroutine closes results after every worker is done | [01-basic](01-basic/), [03-errors](03-errors/) |
| `select { case job, ok := <-jobs: ... case <-ctx.Done(): ... }` | A worker that also stops on cancel | [02-with-context](02-with-context/) |
| `type result struct { ...; err error }` | Each result carries its own error | [03-errors](03-errors/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/13-patterns/01-worker-pool/01-basic` | Three workers make thumbnails from a queue of uploads |
| `02-with-context` | `go run ./02-concurrency/13-patterns/01-worker-pool/02-with-context` | A newsletter send with a time budget: the whole pool stops when it runs out |
| `03-errors` | `go run ./02-concurrency/13-patterns/01-worker-pool/03-errors` | Charge invoices, keep going when cards are declined, list the failures at the end |

## Things to remember

- Close the jobs channel after the last job, so the workers can stop.
- **Close the results channel only after all workers are done.** Use a WaitGroup for that.
- **Send jobs from their own goroutine** when main also reads results. Otherwise main blocks on a send while workers block on results.
- **Both the sender and the workers must watch `ctx.Done()`** if the pool can be stopped early (`02`).
- Decide up front: stop on the first error (errgroup) or keep going and collect errors (`03`).
