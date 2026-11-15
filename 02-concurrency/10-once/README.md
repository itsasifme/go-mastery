# Once

`sync.Once` runs a function only one time, even when many goroutines call it.

## Why

**Problem:** some setup is slow or must happen only once, like opening a database pool. If several goroutines check "is it ready yet?" at the same time, they all see "no", and each one sets it up again.

**Solution:** `sync.Once`. `once.Do(f)` runs `f` for the first caller only. The other callers wait until `f` is finished, then go on. Since Go 1.21, `OnceFunc`, `OnceValue` and `OnceValues` wrap a function so you don't need a separate `Once` variable.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `var once sync.Once` | | 1.0 | Ready to use | [01-lazy-init](01-lazy-init/) |
| `once.Do(f)` | | 1.0 | Run `f` the first time only. Others wait for it | [01-lazy-init](01-lazy-init/) |
| `sync.OnceFunc(f)` | `func()` | 1.21 | A function that runs `f` only once | [02-once-value](02-once-value/) |
| `sync.OnceValue(f)` | `func() T` | 1.21 | Runs `f` once and returns its result every time | none |
| `sync.OnceValues(f)` | `func() (T, error)` | 1.21 | Same, for a value and an error | [02-once-value](02-once-value/), [03-no-retry](03-no-retry/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-lazy-init` | `go run ./02-concurrency/10-once/01-lazy-init` | Open a database pool on first use, once, while requests arrive together |
| `02-once-value` | `go run ./02-concurrency/10-once/02-once-value` | Read a config file once with `OnceValues`, and warn once with `OnceFunc` |
| `03-no-retry` | `go run ./02-concurrency/10-once/03-no-retry` | `OnceValues` keeps a failed connect forever. A mutex version tries again |

## Things to remember

- **The first result is final, even an error.** If setup can fail for a moment, use a mutex and save the result only on success (`03`).
- **If the function panics, `once.Do` still counts it as done.** Later calls do nothing. `OnceFunc`, `OnceValue` and `OnceValues` are different: they panic again with the same value on every call.
- Calling `once.Do` inside its own function gets stuck forever.
- Prefer `OnceValues` over a `Once` plus a package variable. The value and the "done" state can't get out of sync (`02`).

## Go version notes

- `OnceFunc`, `OnceValue`, and `OnceValues` came in Go 1.21.
