# Closing Channels

The sender can close a channel to say "no more values". A `for range` loop over the channel stops once it's closed and empty.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-range` | `go run ./02-concurrency/02-channels/04-closing/01-range` | Send five values, close, and read them with `range` |
| `02-drain-after-close` | `go run ./02-concurrency/02-channels/04-closing/02-drain-after-close` | Values in the buffer can still be read after `close` |

## Things to remember

- Only the sender should close a channel.
- Sending to a closed channel panics.
- Reading from a closed channel returns the zero value right away.
- Closing doesn't throw away buffered values. Receivers get them all first (`02`).
- Use `v, ok := <-ch` to tell a real zero from a closed channel: `ok` is `false` once it's closed and empty.
