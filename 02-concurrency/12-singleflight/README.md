# singleflight

`singleflight` makes many calls for the same key run only once. All callers get the same result.

## Why

**Problem:** when a popular cache entry expires, hundreds of requests miss it at the same moment. Each one asks the database for the same row. This is a cache stampede, and it can take the database down.

**Solution:** `singleflight.Group`. The first call for a key does the work. Calls for the same key that arrive while it's running wait and share its result. When it's finished, the key is free again.

## Import

```go
import "golang.org/x/sync/singleflight"
```

Not in the standard library. It comes from the `golang.org/x/sync` module. `go.mod` says which version this repo uses.

- `golang.org/x/sync/singleflight`: the package name is `singleflight`, the last part of the path, so you write `singleflight.Group`.

## API

| API | Kind | Returns | What it does | Example |
| --- | --- | --- | --- | --- |
| `var g singleflight.Group` | type | | Ready to use | [01-basic](01-basic/) |
| `g.Do(key, fn)` | method | `v any, err error, shared bool` | Run `fn` once for all callers of `key` at the same time | [01-basic](01-basic/), [02-cache](02-cache/) |
| `g.DoChan(key, fn)` | method | `<-chan Result` | Like `Do`, but returns a channel so you can stop waiting | [03-dochan-timeout](03-dochan-timeout/) |
| `g.Forget(key)` | method | | Let the next call for `key` start a new run, even if one is still going | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/12-singleflight/01-basic` | Five requests for one profile, and the database is called once |
| `02-cache` | `go run ./02-concurrency/12-singleflight/02-cache` | A cache stampede: 20 database queries without singleflight, 1 with it |
| `03-dochan-timeout` | `go run ./02-concurrency/12-singleflight/03-dochan-timeout` | Each caller waits with its own time limit, and the work goes on for the others |

## singleflight or `Once`?

Both let one caller do the work while the others wait and get the same result. Both work inside one server only. The difference is how often the work runs, and what's kept afterwards.

| | `sync.Once` / `OnceValues` | `singleflight.Group` |
| --- | --- | --- |
| Runs the work | once, ever | once per key, while a call is running |
| After it finishes | keeps the result forever | forgets the key, so the next call runs again |
| Keys | none | yes: each key runs on its own |
| A failed call | the error is kept forever (see [once/03-no-retry](../10-once/03-no-retry/)) | the next call tries again |
| Stop waiting early | no | yes, with `DoChan` ([03](03-dochan-timeout/)) |

Five requests for `product:42` arrive together, then five more a minute later:

| | First five | A minute later |
| --- | --- | --- |
| `Once` | 1 load | 0 loads: the old result is reused |
| singleflight | 1 load | 1 load: it forgot the key |

singleflight is like `Once` while a call is running, but it remembers nothing afterwards. That's why it's usually paired with a cache ([02](02-cache/)).

## Things to remember

- **It joins calls that overlap in time. It isn't a cache.** Put it around the "load on miss" part of a cache (`02`).
- **It joins calls inside one server only.** Fine for loading a product when the cache is empty: you get one call per server. For one load in the whole app, use a Redis lock (`01`). See [One server or many?](../README.md#one-server-or-many).
- **Callers share the error too, not just the result.**
- **Callers also share a returned slice or map.** Don't change it, or copy it first.
- **One caller's ctx shouldn't cancel the shared work.** Use `DoChan` and select on your own `ctx.Done()` (`03`).
- **The `DoChan` channel is never closed.** It gets one `Result`, so read it once. Don't `range` over it.
- If `fn` panics, every caller waiting in `Do` panics too.
- `Forget` removes a key, so the next call runs again.
