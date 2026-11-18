# sync.Map

`sync.Map` is a map that's safe to use from many goroutines without your own lock.

## Why

**Problem:** a plain Go map isn't safe for concurrent use. If one goroutine writes while another reads, the program crashes with "concurrent map read and map write". The usual fix is a map with a `Mutex` or `RWMutex`, but under heavy load every goroutine waits on that one lock.

**Solution:** `sync.Map`. It does its own locking inside, and it's built for two cases: keys that are written once and read many times, and goroutines that read and write different keys. For anything else, the `sync` docs say a plain map with a lock is better: it's simpler, and you get type safety.

## API

| API | Kind | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- | --- |
| `var m sync.Map` | type | | 1.9 | Ready to use. Keys and values are `any` | [01-registry](01-registry/) |
| `m.Store(key, value)` | method | | 1.9 | Set the value for a key | [01-registry](01-registry/), [03-lost-update](03-lost-update/) |
| `m.Load(key)` | method | `value any, ok bool` | 1.9 | Get the value. `ok` is false if the key is missing | [01-registry](01-registry/), [03-lost-update](03-lost-update/) |
| `m.Delete(key)` | method | | 1.9 | Remove a key | [01-registry](01-registry/) |
| `m.Range(f)` | method | | 1.9 | Call `f(key, value)` for every key. Stops when `f` returns false | [01-registry](01-registry/) |
| `m.LoadOrStore(key, value)` | method | `actual any, loaded bool` | 1.9 | Get the value if the key exists, or store `value`. One step | [02-load-or-store](02-load-or-store/), [03-lost-update](03-lost-update/) |
| `m.LoadAndDelete(key)` | method | `value any, loaded bool` | 1.15 | Remove a key and get the value it had | none |
| `m.Swap(key, value)` | method | `previous any, loaded bool` | 1.20 | Set a new value and get the old one | none |
| `m.CompareAndSwap(key, old, new)` | method | `bool` | 1.20 | Set `new` only if the value is still `old` | none |
| `m.CompareAndDelete(key, old)` | method | `bool` | 1.20 | Remove the key only if the value is still `old` | none |
| `m.Clear()` | method | | 1.23 | Remove every key | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-registry` | `go run ./02-concurrency/13-sync-map/01-registry` | Users connect and leave in their own goroutines, and a broadcast reaches everyone online |
| `02-load-or-store` | `go run ./02-concurrency/13-sync-map/02-load-or-store` | Problem: five requests each build a client and four are thrown away. Fix: store a `sync.OnceValue` |
| `03-lost-update` | `go run ./02-concurrency/13-sync-map/03-lost-update` | Problem: `Load` then `Store` loses page views, and `-race` finds nothing. Fix: one atomic counter per page |

## Things to remember

- **Start with a plain map and a `Mutex` or `RWMutex`.** Switch to `sync.Map` only for keys written once and read many times, or goroutines that use different keys (`01`).
- **`Load` then `Store` isn't one step.** Use `LoadOrStore`, `CompareAndSwap`, or store something safe like an `*atomic.Int64` as the value (`03`).
- **`LoadOrStore` needs the value before it knows if it'll store it.** If building the value is costly, store a `sync.OnceValue` instead (`02`).
- Keys and values are `any`, so every `Load` needs a type assertion like `v.(*session)`.
- `Range` visits keys in no fixed order, and it isn't a snapshot. Keys added or deleted during `Range` may or may not show up.
- Don't copy a `sync.Map` after first use. Pass a pointer.

## Go version notes

- `Swap`, `CompareAndSwap` and `CompareAndDelete` came in Go 1.20. `Clear` came in Go 1.23.
- Go 1.24 changed how `sync.Map` works inside, which made changes to the map faster, especially when goroutines use different keys.
