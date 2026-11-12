# Unbuffered Channels

An unbuffered channel can't hold any value. A send waits until someone receives, and a receive waits until someone sends.

See the [full API overview](../README.md#api).

## Why

**Problem:** two goroutines need to hand off a value, and the sender needs to know the receiver really got it.

**Solution:** an unbuffered channel. The send and the receive happen together, so a send only finishes once someone took the value.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `make(chan T)` | A channel with no buffer | [01-basic](01-basic/) |
| `ch <- v` | Waits until someone receives | [02-blocked-send](02-blocked-send/) |
| `<-ch` | Waits until someone sends | [01-basic](01-basic/), [03-blocked-receive](03-blocked-receive/) |
| send and receive meet | Strict turns between two goroutines | [04-ping-pong](04-ping-pong/) |
| `runtime.NumGoroutine()` | Count goroutines stuck on a channel | [05-leak](05-leak/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/02-channels/03-unbuffered/01-basic` | `main` waits for a value from a goroutine |
| `02-blocked-send` | `go run ./02-concurrency/02-channels/03-unbuffered/02-blocked-send` | A worker waits to send, but nobody receives: deadlock |
| `03-blocked-receive` | `go run ./02-concurrency/02-channels/03-unbuffered/03-blocked-receive` | A goroutine waits to receive, but nobody sends: deadlock |
| `04-ping-pong` | `go run ./02-concurrency/02-channels/03-unbuffered/04-ping-pong` | `main` and a goroutine take strict turns |
| `05-leak` | `go run ./02-concurrency/02-channels/03-unbuffered/05-leak` | Workers stuck forever on a send with no error: a goroutine leak |

## Things to remember

- `02` and `03` deadlock on purpose. **Go reports it only when every goroutine is stuck, including `main`.**
- If `main` keeps running, **stuck goroutines leak silently** (`05`).
- `02`, `03` and `05` have a `// fix:` line to try.
