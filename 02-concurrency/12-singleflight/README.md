# singleflight

`singleflight` makes many calls for the same key run only once. All callers get the same result.

## Why

**Problem:** when a popular cache entry expires, hundreds of requests miss at the same moment, and each one asks the database for the same row. This is a cache stampede, and it can take the database down.

**Solution:** `singleflight.Group`. The first call for a key does the work. Calls for the same key that arrive while it's running wait and share its result. When it's finished, the key is free again.

## API

| API | Returns | What it does | Example |
| --- | --- | --- | --- |
| `var g singleflight.Group` | | Ready to use | [01-basic](01-basic/) |
| `g.Do(key, fn)` | `v any, err error, shared bool` | Run `fn` once for all callers of `key` at the same time | [01-basic](01-basic/), [02-cache](02-cache/) |
| `g.DoChan(key, fn)` | `<-chan Result` | Like `Do`, but returns a channel so you can stop waiting | [03-dochan-timeout](03-dochan-timeout/) |
| `g.Forget(key)` | | Let the next call for `key` start a new run, even if one is still going | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/12-singleflight/01-basic` | Five requests for one profile, and the database is called once |
| `02-cache` | `go run ./02-concurrency/12-singleflight/02-cache` | A cache stampede: 20 database queries without singleflight, 1 with it |
| `03-dochan-timeout` | `go run ./02-concurrency/12-singleflight/03-dochan-timeout` | Each caller waits with its own time limit, and the work goes on for the others |

## Things to remember

- **It joins calls that overlap in time. It isn't a cache.** Put it around the "load on miss" part of a cache (`02`).
- **Callers share the error too, not just the result.**
- **Callers also share a returned slice or map.** Don't change it, or copy it first.
- **One caller's ctx shouldn't cancel the shared work.** Use `DoChan` and select on your own `ctx.Done()` (`03`).
- **The `DoChan` channel is never closed.** It gets one `Result`, so read it once. Don't `range` over it.
- If `fn` panics, every caller waiting in `Do` panics too.
- `Forget` removes a key, so the next call runs again.

## Package note

It comes from `golang.org/x/sync`, not the standard library. This repo uses `v0.16.0`.
