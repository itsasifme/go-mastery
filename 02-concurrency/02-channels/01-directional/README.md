# Directional Channels

You can limit a channel to send-only (`chan<- T`) or receive-only (`<-chan T`). Then the compiler stops you from using it the wrong way.

See the [full API overview](../README.md#api).

## Why

**Problem:** with a plain `chan`, every function can send, receive and close. A consumer can close the channel by mistake, or a producer can read its own output.

**Solution:** directional channel types. The compiler stops each side from doing what it shouldn't.

## API

| Type | Can do | Example |
| --- | --- | --- |
| `chan T` | Send and receive | [01-bidirectional](01-bidirectional/) |
| `chan<- T` | Send only (and `close`) | [02-send-only](02-send-only/) |
| `<-chan T` | Receive only. Can't `close` | [03-receive-only](03-receive-only/) |
| `chan T` → `chan<- T` / `<-chan T` | Converts automatically when passed | [04-conversion](04-conversion/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-bidirectional` | `go run ./02-concurrency/02-channels/01-directional/01-bidirectional` | A normal channel passed to a sender and a reader |
| `02-send-only` | `go run ./02-concurrency/02-channels/01-directional/02-send-only` | A send-only channel type |
| `03-receive-only` | `go run ./02-concurrency/02-channels/01-directional/03-receive-only` | A receive-only channel type |
| `04-conversion` | `go run ./02-concurrency/02-channels/01-directional/04-conversion` | A sender and a reader share one channel through send-only and receive-only views, by hand and automatically |

## Things to remember

- In real code, **make a normal channel and pass it to functions as send-only or receive-only**, like in `04`.
- `02` and `03` only show the types. Their goroutines get stuck, so they only print `main done`. Each has `// fix:` lines to try.
- `01` waits with `time.Sleep`. It works here, but it's not a real way to sync. Use a WaitGroup or a channel instead.
