# sync.Map

`sync.Map` is a map that's safe to use from many goroutines without your own lock.

## Why

**Problem:** a plain Go map isn't safe to use from many goroutines at once. If one goroutine writes while another reads, the program crashes with "concurrent map read and map write". The usual fix is a map with a `Mutex` or `RWMutex`. But under heavy load, every goroutine waits on that one lock.

**Solution:** `sync.Map`. It does its own locking inside, so it's always safe. But it's only **faster** than a map with a lock in two cases:

- **Keys written once, then read many times.** For example, one client per partner host, made on first use (`02`, `04`).
- **Each goroutine uses its own keys.** For example, one entry per connected user, which only that user's goroutine changes (`01`).

For anything else, the `sync` docs say a plain map with a lock is better. It's simpler, and you get type safety.

## Import

```go
import "sync"
```

Standard library, so there is nothing to download.

- `sync`: you write `sync.Map`.

## API

Every `key`, `value`, `old` and `new` below is `any`. A blank **Returns** cell means the method returns nothing.

| API | Kind | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- | --- |
| `var m sync.Map` | type | | 1.9 | Ready to use. Keys and values are `any` | [01-registry](01-registry/) |
| `m.Store(key, value)` | method | | 1.9 | Set the value for a key | [01-registry](01-registry/), [03-lost-update](03-lost-update/), [04-compare](04-compare/) |
| `m.Load(key)` | method | `value any, ok bool` | 1.9 | Get the value. `ok` is false if the key is missing | [01-registry](01-registry/), [02-load-or-store](02-load-or-store/), [03-lost-update](03-lost-update/), [04-compare](04-compare/) |
| `m.Delete(key)` | method | | 1.9 | Remove a key | [01-registry](01-registry/) |
| `m.Range(f)` | method | | 1.9 | `f` is a `func(key, value any) bool`. `Range` calls it for every key. Return `true` to go on, `false` to stop | [01-registry](01-registry/), [02-load-or-store](02-load-or-store/) |
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
| `02-load-or-store` | `go run ./02-concurrency/13-sync-map/02-load-or-store` | Problem: every request builds a client, and all but one are thrown away. Fix: store a `sync.OnceValue` |
| `03-lost-update` | `go run ./02-concurrency/13-sync-map/03-lost-update` | Problem: `Load` then `Store` loses page views, and `-race` finds nothing. Fix: one atomic counter per page |
| `04-compare` | `go run ./02-concurrency/13-sync-map/04-compare` | The same work timed with a map + `Mutex`, a map + `RWMutex`, and `sync.Map`: `sync.Map` is much faster for reading keys stored once, and slowest for writing one key |

## `sync.Map` or map + `Mutex`?

Both are safe to use from many goroutines. The difference is speed, and how easy the code is to get right.

| | map + `Mutex` | `sync.Map` |
| --- | --- | --- |
| Types | `map[string]int`, checked by the compiler | `any`, so every `Load` needs `v.(T)` |
| Read, change, write back (`m[k]++`) | Easy: do it all under one `Lock` | `Load` then `Store` loses updates (`03`). Use `LoadOrStore`, `CompareAndSwap`, or an atomic value |
| `len(m)` | ✅ | None. Count with `Range` |
| Many goroutines read keys stored once | Readers take turns on the lock. `RWMutex` helps only a little | ✅ Much faster (`04`) |
| Each goroutine writes its own keys | Goroutines wait for each other, even on different keys | ✅ A little faster |
| Many goroutines write the same key | ✅ Usually a bit faster (`04`) | No help. It does more work per write |

Start with a map and a `Mutex` (or `RWMutex` for mostly reads). Switch to `sync.Map` only when your case is one of its two, and the numbers show the lock slows you down.

## Things to remember

- **Start with a plain map and a `Mutex` or `RWMutex`.** `sync.Map` is always safe, but only faster for keys written once and read many times, or goroutines that each use their own keys (`01`, `04`).
- **`Load` then `Store` isn't one step.** Use `LoadOrStore`, `CompareAndSwap`, or store something safe like an `*atomic.Int64` as the value (`03`).
- **Each server has its own map.** Fine for this server's own connections. For user sessions that any server may get, use Redis. See [One server or many?](../README.md#one-server-or-many).
- **`LoadOrStore` needs the value before it knows if it'll store it.** If building the value is costly, store a `sync.OnceValue` instead, and call `Load` first so a stored key costs nothing (`02`).
- **Keys and values are `any`, so wrap the `sync.Map` in your own type.** Every `Load` needs a type assertion like `v.(*session)`. Keep those inside the type's methods, so callers can't store the wrong type (`01`).
- `Range` visits keys in no fixed order, and it isn't a snapshot. Keys added or deleted during `Range` may or may not show up.
- Don't copy a `sync.Map` after first use. Pass a pointer.

## Go version notes

- `Swap`, `CompareAndSwap` and `CompareAndDelete` came in Go 1.20. `Clear` came in Go 1.23.
- Go 1.24 changed how `sync.Map` works inside. Changes to the map got faster, especially when goroutines use different keys.
