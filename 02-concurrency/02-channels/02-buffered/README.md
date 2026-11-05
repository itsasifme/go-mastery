# Buffered Channels

A buffered channel can hold a few values. A send only waits when the buffer is full.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/02-channels/02-buffered/01-basic` | Send two values with nobody receiving yet, then read them |
| `02-full-buffer` | `go run ./02-concurrency/02-channels/02-buffered/02-full-buffer` | A goroutine tries to send to a full buffer and has to wait |

## Things to remember

- In `02`, `main` reads both values and then ends. `Sent 3 to the channel` may not print, but the program doesn't crash.
