# Directional Channels

You can limit a channel to send-only (`chan<- T`) or receive-only (`<-chan T`). Then the compiler stops you from using it the wrong way.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-bidirectional` | `go run ./02-concurrency/02-channels/03-directional/01-bidirectional` | A normal channel passed to a sender and a reader |
| `02-send-only` | `go run ./02-concurrency/02-channels/03-directional/02-send-only` | A send-only channel type |
| `03-receive-only` | `go run ./02-concurrency/02-channels/03-directional/03-receive-only` | A receive-only channel type |
| `04-conversion` | `go run ./02-concurrency/02-channels/03-directional/04-conversion` | Turn a normal channel into send-only and receive-only |

## Things to remember

- In real code, make a normal channel and pass it to functions as send-only or receive-only, like in `04`.
- `02` and `03` only show the types. Their goroutines get stuck, so they usually print nothing.
- `01` waits with `time.Sleep`. It works here, but it's not a real way to sync. Use a WaitGroup or a channel instead.
