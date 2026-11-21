# Fan-Out

Many goroutines read from the same channel to share the work.

## Why

**Problem:** one goroutine reading a channel handles one value at a time. If each value takes a while (a network call, heavy CPU work), the queue grows.

**Solution:** start several goroutines that all read from the same channel. Go gives each value to exactly one of them, whichever is ready first. No one has to split up the work. When values that belong together must stay in order, give each reader its own channel and pick the channel by key.

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `for range n { go worker(ch, ...) }` | Several readers on one channel | [01-basic](01-basic/), [02-by-key](02-by-key/) |
| `for v := range ch` | Each reader takes values until the channel is closed | [01-basic](01-basic/) |
| `queues[key % len(queues)] <- v` | Send every value with the same key to the same reader, so they stay in order | [02-by-key](02-by-key/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/15-patterns/01-fan-out/01-basic` | Packers share a stream of orders. The fastest one packs the most |
| `02-by-key` | `go run ./02-concurrency/15-patterns/01-fan-out/02-by-key` | Problem: an order's status events finish out of order. Fix: send each order's events to the same worker |

## Things to remember

- **Each value goes to only one goroutine.** It shares work. It doesn't copy it to everyone. To send to everyone, close a channel (see [channels/02-closing/02-broadcast](../../02-channels/02-closing/02-broadcast/)).
- Work balances by itself: a faster reader just comes back for more sooner (`01`).
- **Plain fan-out loses order.** If values for the same key must happen in order, route by key: one channel per reader, and the key picks the channel (`02`).
- Fan-out plus fan-in, with a fixed number of readers, is a [worker pool](../03-worker-pool/).
