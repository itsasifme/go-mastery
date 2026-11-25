# Once

`sync.Once` runs a function only one time, even when many goroutines call it.

## Why

**Problem:** some setup is slow or must happen only once, like opening a database pool. If several goroutines check "is it ready yet?" at the same time, they all see "no", and each one sets it up again.

**Solution:** `sync.Once`. `once.Do(f)` runs `f` for the first caller only. The other callers wait until `f` is finished, then go on. Since Go 1.21, `OnceFunc`, `OnceValue` and `OnceValues` wrap a function so you don't need a separate `Once` variable.

## API

| API | Kind | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- | --- |
| `var once sync.Once` | type | | 1.0 | Ready to use | [01-lazy-init](01-lazy-init/) |
| `once.Do(f)` | method | | 1.0 | Run `f` the first time only. Others wait for it | [01-lazy-init](01-lazy-init/) |
| `sync.OnceFunc(f)` | function | `func()` | 1.21 | A function that runs `f` only once | [02-once-value](02-once-value/) |
| `sync.OnceValue(f)` | function | `func() T` | 1.21 | Runs `f` once and returns its result every time | none |
| `sync.OnceValues(f)` | function | `func() (T, error)` | 1.21 | Same, for a value and an error | [02-once-value](02-once-value/), [03-no-retry](03-no-retry/) |

### `once.Do` or `OnceFunc`?

Both run the function once, and other callers wait for it. `OnceFunc` uses a `sync.Once` inside. The difference is where the `Once` lives, and what happens after a panic.

| | `once.Do(f)` | `sync.OnceFunc(f)` |
| --- | --- | --- |
| What you keep | a `sync.Once` variable, plus your own variable for the result | one function. Use `OnceValue` or `OnceValues` when you need a result |
| Best for | once per struct value: a `sync.Once` field works with no setup | once for the whole program, as a package variable |
| If `f` panics | the first call panics. Later calls do nothing, so the code goes on as if setup worked | every call panics again with the same value, so a broken setup can't go unnoticed |

```go
// once per struct value
type Client struct {
	initOnce sync.Once
	conn     *conn
}

func (c *Client) get() *conn {
	c.initOnce.Do(c.connect)
	return c.conn
}

// once for the whole program
var loadConfig = sync.OnceValues(readConfig)
```

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-lazy-init` | `go run ./02-concurrency/10-once/01-lazy-init` | Open a database pool on first use, once, while requests arrive together |
| `02-once-value` | `go run ./02-concurrency/10-once/02-once-value` | Read a config file once with `OnceValues`, and warn once with `OnceFunc` |
| `03-no-retry` | `go run ./02-concurrency/10-once/03-no-retry` | `OnceValues` keeps a failed connect forever. A mutex version tries again |

## Things to remember

- **The first result is final, even an error.** If setup can fail for a moment, use a mutex and save the result only on success (`03`).
- **If the function panics, `once.Do` still counts it as done.** Later calls do nothing, and the code goes on without the setup. `OnceFunc`, `OnceValue` and `OnceValues` panic again on every call, so you keep seeing the problem.
- **Pick by where the `Once` lives.** Use `once.Do` for a `Once` field inside a struct. Use `OnceFunc`, `OnceValue` or `OnceValues` for a package variable.
- Calling `once.Do` inside its own function gets stuck forever.
- Use `OnceValues` instead of a `Once` plus a package variable. Then the value and the "done" state can't get out of sync (`02`).

## Go version notes

- `OnceFunc`, `OnceValue`, and `OnceValues` came in Go 1.21.
