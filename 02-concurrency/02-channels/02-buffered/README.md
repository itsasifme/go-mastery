# Buffered Channels

A buffered channel can hold a few values. A send only waits when the buffer is full, and a receive only waits when it's empty.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/02-channels/02-buffered/01-basic` | Send two values with nobody receiving yet, then read them |
| `02-blocked-send` | `go run ./02-concurrency/02-channels/02-buffered/02-blocked-send` | 2 sends fit, the 3rd finds the buffer full: deadlock |
| `03-blocked-receive` | `go run ./02-concurrency/02-channels/02-buffered/03-blocked-receive` | 1 receive works, the 2nd finds the buffer empty: deadlock |
| `04-fast-producer` | `go run ./02-concurrency/02-channels/02-buffered/04-fast-producer` | A fast sender runs ahead of a slow receiver |
| `05-single-result` | `go run ./02-concurrency/02-channels/02-buffered/05-single-result` | A buffer of 1 stops a goroutine leak when nobody reads the result |
| `06-bigger-buffer` | `go run ./02-concurrency/02-channels/02-buffered/06-bigger-buffer` | A bigger buffer hides a deadlock, it doesn't fix it |

## Things to remember

- Sends block only when the buffer is full. Uncomment `ch <- 3` in `01` to see it.
- `len(ch)` = values waiting, `cap(ch)` = buffer size.
- Make the buffer as big as the number of senders, so none get stuck if the receiver stops early (no goroutine leak).
- `01` runs only in `main`, which is itself a goroutine, just to demo the buffer. In real code, channels are for passing data between goroutines.
- `02`, `03`, `05` and `06` show a problem on purpose. Each has a `// fix:` line to try.
- A buffer helps a fast and a slow side work at their own speed (`04`). It doesn't fix a design where nobody reads (`06`).
- For a single result, a buffer of 1 is enough (`05`). Uber's guide says channels should be size 0 or 1.
