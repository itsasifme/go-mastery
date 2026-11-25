# sync.Pool

`sync.Pool` keeps used objects so you can use them again. You don't have to make new ones and leave the old ones to the garbage collector.

## Why

**Problem:** a busy server makes the same short-lived objects over and over, like a buffer for every response or log line. Each one becomes garbage right after the request, and the garbage collector spends CPU cleaning them up.

**Solution:** `sync.Pool`. `Get` hands out an object from the pool, or makes one with `New` if the pool is empty. `Put` gives it back for the next caller. A few objects end up serving thousands of requests. It's safe to use from many goroutines.

## API

| API | Kind | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- | --- |
| `sync.Pool{New: f}` | type | | 1.3 | A pool. `New` makes an object when the pool is empty | [01-buffer](01-buffer/) |
| `p.New` | field | | 1.3 | A `func() any`. Without it, `Get` returns `nil` on an empty pool | [01-buffer](01-buffer/) |
| `p.Get()` | method | `any` | 1.3 | Take an object from the pool, or a new one from `New` | [01-buffer](01-buffer/), [02-reset](02-reset/) |
| `p.Put(x)` | method | | 1.3 | Give an object back for reuse | [01-buffer](01-buffer/), [03-big-buffer](03-big-buffer/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-buffer` | `go run ./02-concurrency/14-pool/01-buffer` | 1000 log lines without a pool make 1000 buffers. With a pool, only a few |
| `02-reset` | `go run ./02-concurrency/14-pool/02-reset` | Problem: without `Reset`, one user's email shows up in the next user's response |
| `03-big-buffer` | `go run ./02-concurrency/14-pool/03-big-buffer` | Problem: one big export leaves a 4 MB buffer in the pool for small requests. Fix: don't put big buffers back |

## Things to remember

- **Reset what you get.** A pooled object still holds the last user's data. Leaving it in is a data leak between requests (`02`).
- **The pool can drop objects at any time,** like during a garbage collection, without telling you. Only pool things that are cheap to lose. Never pool connections or open files: use a real connection pool for those, like `database/sql` has.
- **Don't put back objects that grew too big.** Check the size before `Put`, like `fmt` does with its 64 KB limit (`03`).
- **Store pointers, like `*bytes.Buffer`.** If you put a plain value or slice in, wrapping it in `any` allocates memory on every `Put`.
- Measure first. A pool only helps on hot paths that make a lot of garbage. Use `go test -bench` with `-benchmem` to see if it does.
- Don't copy a `sync.Pool` after first use.

## Go version notes

- Since Go 1.13, objects in the pool survive one garbage collection before they're dropped. So a collection doesn't empty the pool right away.
