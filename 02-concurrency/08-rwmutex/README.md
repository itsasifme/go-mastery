# RWMutex

`sync.RWMutex` lets many goroutines read at the same time, but only one write.

## Why

**Problem:** with a plain `Mutex`, readers wait for each other, even though reading doesn't change anything. For data that's read on every request and changed rarely, that's a lot of waiting for nothing.

**Solution:** `sync.RWMutex` has two kinds of lock. `RLock` can be held by many readers at once. `Lock` waits until every reader is gone, then lets one writer in alone.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `var mu sync.RWMutex` | | 1.0 | Ready to use, starts unlocked | [01-cache](01-cache/) |
| `mu.RLock()` / `mu.RUnlock()` | | 1.0 | Read lock. Many goroutines can hold it at once | [01-cache](01-cache/), [02-compare](02-compare/) |
| `mu.Lock()` / `mu.Unlock()` | | 1.0 | Write lock. Only one, and no readers at the same time | [01-cache](01-cache/), [03-upgrade](03-upgrade/), [04-nested-rlock](04-nested-rlock/) |
| `mu.TryRLock()` / `mu.TryLock()` | `bool` | 1.18 | Take the lock only if it's free right now. The `sync` docs say correct uses of `TryRLock` are rare | `TryLock` works like Mutex's (see [mutex/05-try-lock](../07-mutex/05-try-lock/)) |
| `mu.RLocker()` | `sync.Locker` | 1.0 | A `Locker` whose `Lock`/`Unlock` call `RLock`/`RUnlock` | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-cache` | `go run ./02-concurrency/08-rwmutex/01-cache` | Exchange rates read by many requests and updated by one job |
| `02-compare` | `go run ./02-concurrency/08-rwmutex/02-compare` | The same read-only work timed with a `Mutex` and an `RWMutex` |
| `03-upgrade` | `go run ./02-concurrency/08-rwmutex/03-upgrade` | Problem: taking `Lock` while holding `RLock` gets stuck. Fix: drop it, lock, look again |
| `04-nested-rlock` | `go run ./02-concurrency/08-rwmutex/04-nested-rlock` | Problem: taking `RLock` twice gets stuck when a refresh calls `Lock` in between |

## Things to remember

- **It only helps when there are many more reads than writes.** Measure it first (`02`).
- **You can't upgrade `RLock` to `Lock`.** Release the read lock, take the write lock, then check again (`03`).
- **Don't call `RLock` twice in the same goroutine.** It works in tests, then gets stuck when a writer calls `Lock` in between. Use a `...Locked` helper (`04`).
- A waiting writer blocks new readers, so writers aren't starved by a steady stream of reads.
