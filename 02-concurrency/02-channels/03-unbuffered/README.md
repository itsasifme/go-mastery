# Unbuffered Channels

An unbuffered channel can't hold any value. A send waits until someone receives, and a receive waits until someone sends.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/02-channels/03-unbuffered/01-basic` | `main` waits for a value from a goroutine |
| `02-blocked-send` | `go run ./02-concurrency/02-channels/03-unbuffered/02-blocked-send` | A goroutine waits to send, but nobody receives: deadlock |
| `03-blocked-receive` | `go run ./02-concurrency/02-channels/03-unbuffered/03-blocked-receive` | A goroutine waits to receive, but nobody sends: deadlock |
| `04-ping-pong` | `go run ./02-concurrency/02-channels/03-unbuffered/04-ping-pong` | `main` and a goroutine take strict turns |
| `05-leak` | `go run ./02-concurrency/02-channels/03-unbuffered/05-leak` | Senders stuck forever with no error: a goroutine leak |

## Things to remember

- `02` and `03` deadlock on purpose. Go reports it only when every goroutine is stuck, including `main`.
- If `main` keeps running, stuck goroutines leak silently (`05`).
- `02`, `03` and `05` have a `// fix:` line to try.
