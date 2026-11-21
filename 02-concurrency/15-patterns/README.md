# Concurrency Patterns

Common ways to build things with goroutines and channels.

The topics before this one each teach one tool. This one puts them together.

## Why

**Problem:** goroutines and channels are simple on their own. But real programs ask the same questions again and again: how many goroutines to start, who closes which channel, how results come back, and what happens when one side is slower than the other.

**Solution:** a few well-known patterns. Each one answers one of those questions, and they fit together: a worker pool is fan-out plus fan-in, and a pipeline step can have its own workers.

## Patterns

| Pattern | Use it when |
| --- | --- |
| [Worker pool](03-worker-pool/) | You have many jobs and want a fixed limit on how many run at once |
| [Fan-in](02-fan-in/) | Values come from several channels and you want to read them in one loop |
| Fan-out | One stream of work is too slow for one goroutine. Several goroutines read from one channel: a [worker pool](03-worker-pool/) does this, and so does [pipeline/03-parallel-step](04-pipeline/03-parallel-step/) |
| [Pipeline](04-pipeline/) | The work has several steps, and each step can work on a different item at the same time |
| [Backpressure](05-backpressure/) | A fast producer meets a slow consumer, and you must choose: wait, drop, or reject |

## Things to remember

- **Every goroutine needs a clear way to stop.** Otherwise it leaks.
- Only the sender closes a channel, and only after all sending is done.
- **Pass a `ctx` when the reader can stop early,** so every goroutine behind it stops too.
- Use errgroup when one failure should stop the whole group (see [pipeline/04-errors](04-pipeline/04-errors/)).
