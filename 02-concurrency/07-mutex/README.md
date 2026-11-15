# Mutex

`sync.Mutex` lets only one goroutine touch shared data at a time.

## Why

**Problem:** when two goroutines change the same variable at once, updates get lost. Even `count++` is three steps (read, add, write back), and another goroutine can sneak in between them. This is a data race. The results are wrong, and they change from run to run.

**Solution:** `sync.Mutex`. Call `Lock` before you touch the data and `Unlock` after. While one goroutine holds the lock, the others wait at `Lock`.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `var mu sync.Mutex` | | 1.0 | Ready to use, starts unlocked | [01-race](01-race/) |
| `mu.Lock()` | | 1.0 | Wait until the lock is free, then take it | [01-race](01-race/), [02-safe-type](02-safe-type/) |
| `mu.Unlock()` | | 1.0 | Give the lock back | [02-safe-type](02-safe-type/) |
| `mu.TryLock()` | `bool` | 1.18 | Take the lock only if it's free right now. Rarely the right tool | [05-try-lock](05-try-lock/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-race` | `go run ./02-concurrency/07-mutex/01-race` | Problem: a page view counter loses updates without a lock |
| `02-safe-type` | `go run ./02-concurrency/07-mutex/02-safe-type` | A wallet that locks inside its own methods, so callers can't forget to lock |
| `03-check-then-act` | `go run ./02-concurrency/07-mutex/03-check-then-act` | Problem: locking each step still sells the last ticket five times |
| `04-nested-lock` | `go run ./02-concurrency/07-mutex/04-nested-lock` | Problem: a method that locks calls another that locks, and gets stuck |
| `05-try-lock` | `go run ./02-concurrency/07-mutex/05-try-lock` | Several admins click "sync" at once: one sync runs, the others skip instead of waiting |

## Things to remember

- A zero `sync.Mutex` is ready to use.
- **Put `defer mu.Unlock()` right after `mu.Lock()`.** Then an early `return` or a panic can't leave it locked.
- **Make the type safe by itself: keep the mutex inside it, right above the fields it guards, and lock in every method.** Callers never touch the mutex. Use pointer receivers (`02`).
- **Never copy a mutex.** `go vet` warns you about it.
- **A lock on each step isn't enough.** A check and the change that depends on it must happen under one lock (`03`).
- **Locking the same mutex twice in one goroutine gets it stuck forever.** Go mutexes don't know who holds them (`04`).
- **`TryLock` is for "skip it if it's already running", and it only works inside one server.** If it returns `false`, don't call `Unlock`. Across servers, use a database lock (`05`).
- Find races with `go run -race` or `go test -race`. On Windows, `-race` needs `CGO_ENABLED=1` and a C compiler like gcc.
