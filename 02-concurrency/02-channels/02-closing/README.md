# Closing Channels

The sender can close a channel to say "no more values". A `for range` loop over the channel stops once it's closed and empty.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-range` | `go run ./02-concurrency/02-channels/02-closing/01-range` | Send five values, close, and read them with `range` |
| `02-broadcast` | `go run ./02-concurrency/02-channels/02-closing/02-broadcast` | One `close` wakes every waiting goroutine. A send wakes only one |
| `03-many-senders` | `go run ./02-concurrency/02-channels/02-closing/03-many-senders` | With many senders, one separate goroutine closes the channel once |

What happens to values still in a buffer after `close` is shown in [`04-buffered/06-drain-after-close`](../04-buffered/06-drain-after-close/).

## Cheat sheet

| Operation | `nil` channel | Open channel | Closed channel |
| --- | --- | --- | --- |
| Send `ch <- v` | blocks forever | sends, or waits | **panic** |
| Receive `<-ch` | blocks forever | gets a value, or waits | zero value right away, `ok == false` |
| `close(ch)` | **panic** | closes it | **panic** |

## Things to remember

- Only the sender should close a channel, and only once.
- With many senders, none of them closes it. One separate goroutine waits for all of them, then closes (`03`).
- Use `v, ok := <-ch` to tell a real zero from a closed channel: `ok` is `false` once it's closed and empty.
- `close` is a broadcast to every receiver. That's why `close(done)` is the usual way to say "I'm finished" (`02`).
- You don't have to close every channel. Close it only when a receiver needs to know there's nothing more, like a `for range` loop.
- `03` panics and `02`'s Try line deadlocks on purpose. Each has a fix or a note to try.
