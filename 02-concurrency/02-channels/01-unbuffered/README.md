# Unbuffered Channels

An unbuffered channel can't hold any value. A send waits until someone receives, and a receive waits until someone sends.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/02-channels/01-unbuffered/01-basic` | `main` waits for a value from a goroutine |
| `02-blocked-receive` | `go run ./02-concurrency/02-channels/01-unbuffered/02-blocked-receive` | A goroutine waits to receive, but nobody sends |
| `03-blocked-send` | `go run ./02-concurrency/02-channels/01-unbuffered/03-blocked-send` | A goroutine waits to send, but nobody receives |

## Things to remember

- `02` and `03` were first written as "deadlock" examples, but they don't crash.
- The stuck goroutine runs in the background, and the program simply ends when `main` returns.
- Go only shows `all goroutines are asleep - deadlock!` when every goroutine is stuck, including `main`.
