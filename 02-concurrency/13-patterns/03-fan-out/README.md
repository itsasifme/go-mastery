# Fan-Out

Many goroutines read from the same channel to share the work.

## Why

**Problem:** one goroutine reading a channel handles one value at a time. If each value takes a while (a network call, heavy CPU work), the queue grows.

**Solution:** start several goroutines that all read from the same channel. Go gives each value to exactly one of them, whichever is ready first. No one has to split up the work.

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `for range n { go worker(ch, ...) }` | Several readers on one channel | [01-basic](01-basic/) |
| `for v := range ch` | Each reader takes values until the channel is closed | [01-basic](01-basic/) |
| `merge(outs...)` | Collect each reader's output in one place (fan-in) | [02-fan-out-fan-in](02-fan-out-fan-in/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/13-patterns/03-fan-out/01-basic` | Packers share a stream of orders. The fastest one packs the most |
| `02-fan-out-fan-in` | `go run ./02-concurrency/13-patterns/03-fan-out/02-fan-out-fan-in` | Check many stores for a price at once, merge the quotes, pick the cheapest |

## Things to remember

- **Each value goes to only one goroutine.** It shares work. It doesn't copy it to everyone. To send to everyone, close a channel (see [channels/02-closing/02-broadcast](../../02-channels/02-closing/02-broadcast/)).
- Work balances by itself: a faster reader just comes back for more sooner (`01`).
- It's often used together with fan-in to collect the results (`02`).
