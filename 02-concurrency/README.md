# Concurrency

Running code at the same time with goroutines, and keeping it safe with channels and the `sync` tools.

1. [Goroutines](01-goroutines/)
2. [Channels](02-channels/)
3. [Select](03-select/)
4. [Context](04-context/)
5. [WaitGroup](05-waitgroup/)
6. [errgroup](06-errgroup/)
7. [Mutex](07-mutex/)
8. [RWMutex](08-rwmutex/)
9. [Atomic](09-atomic/)
10. [Once](10-once/)
11. [Semaphore](11-semaphore/)
12. [singleflight](12-singleflight/)
13. [sync.Map](13-sync-map/)
14. [sync.Pool](14-pool/)
15. [Patterns](15-patterns/)

## One server or many?

Every tool here works inside **one running program**. If your app runs on 3 servers (or pods), each server has its own mutex, its own `Once`, its own cache. They can't see each other.

Two short stories show when that matters.

**❌ Wrong result: 1 ticket left, 2 servers.**

```text
server A: locks its mutex → sees 1 left → sells it → 0 left
server B: locks its mutex → sees 1 left → sells it → 0 left    (B's mutex never saw A's)
result:   2 people paid for 1 seat
```

Each mutex did its job, but the stock belongs to the whole app. Let the database check and change it in one step: `UPDATE tickets SET left = left - 1 WHERE id = $1 AND left > 0`.

**✅ Only extra work: a product page, 2 servers.**

```text
server A: cache is empty → singleflight → 1 database query
server B: cache is empty → singleflight → 1 database query
result:   2 queries instead of 1, and both pages show the right product
```

Nothing is wrong. The database just did one extra query.

So ask one question: **if every server does this on its own, does someone get a wrong result?**

- **Yes** (sold twice, email sent twice, a limit broken): use something all servers share, like the database, Redis or a queue.
- **No, only extra work** or data that's a bit old for a short time: the tool here is fine.

| Tool | Fine per server | Needs a shared store |
| --- | --- | --- |
| Mutex | a counter or cache inside this server | ticket stock, account balances |
| RWMutex (cache in memory) | old data is harmless, like a banner | old data is a mistake, like a price or a balance |
| Atomic | a request counter, config in `atomic.Pointer` | coupon stock |
| `Once` | opening a database pool or a Redis client, parsing templates | "send the welcome email once" |
| Semaphore | keeping this server's memory or CPU safe | a partner API's limit for your whole account |
| singleflight | loading a product when the cache is empty | rarely needed: one load for the whole app |
| `sync.Map` | this server's own connections | user sessions that any server may get |

The examples that hit this limit say so in a `Note:`, with what to use across servers.

## How often will you use them?

"Web" means API servers and websites. "Infra" means tools that run or watch other systems: CLIs, Kubernetes operators, proxies, agents and background workers.

| Tool | Web | Infra | Where you see it |
| --- | --- | --- | --- |
| Goroutines | Every day | Every day | `net/http` already runs each request in its own goroutine. You start more for background work. |
| `context` | Every day | Every day | The first parameter of almost every handler, database call and HTTP call. |
| Channels | Sometimes | Often | Web: a background job queue. Infra: event loops, watchers and worker queues. |
| `select` | Sometimes | Often | Waiting for a result, a timeout or `ctx.Done()`, whichever comes first. |
| WaitGroup | Sometimes | Often | Waiting for a batch of goroutines, or for workers at shutdown. |
| errgroup | Often | Often | Calling several services at once in one handler, and stopping at the first error. |
| Mutex | Often | Every day | Guarding a small cache, a client or any shared struct. |
| RWMutex | Sometimes | Often | Config or a cache that's read a lot and changed rarely. |
| Atomic | Sometimes | Often | Counters and metrics (the Prometheus client uses atomics), and swapping config with `atomic.Pointer`. |
| `Once` | Sometimes | Sometimes | Making a client or parsing templates the first time they're needed. |
| Semaphore | Rarely | Sometimes | Limiting how many uploads, builds or API calls run at once. |
| singleflight | Rarely | Sometimes | Stopping a cache stampede. Go's own DNS lookup uses one inside. |
| `sync.Map` | Rarely | Rarely | Mostly inside libraries, like the type cache in `encoding/json`. |
| `sync.Pool` | Rarely | Sometimes | Mostly inside libraries: `fmt`, `encoding/json`, Gin's request `Context`, loggers. |

Two things to notice:

- **Web code shares less in memory.** Most state lives in the database or Redis, because the app runs on many servers (see [One server or many?](#one-server-or-many)). Infra tools often run as one long-lived process, so they keep more state in memory and need more locks and channels.
- **Rarely doesn't mean "skip it".** You use `sync.Map` and `sync.Pool` on every request through `fmt`, `encoding/json` and your framework. Write your own only when a profile shows a real cost.
