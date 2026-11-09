# Context

A `context.Context` tells goroutines when to stop. It can carry a cancel signal, a timeout, a deadline, and request data.

1. [Cancellation](01-cancellation/): stop work by hand, and how parents and children cancel
2. [Timeout](02-timeout/): stop work after some time
3. [Deadline](03-deadline/): stop work at a fixed time
4. [Values](04-values/): carry request data like a request ID
5. [Advanced](05-advanced/): `AfterFunc` and `WithoutCancel`

## Why

**Problem:** a server starts goroutines for every request: database queries, calls to other services. When the client disconnects or the request takes too long, they should stop. There was no standard way to tell them. Every library made its own stop channel or timeout parameter (like the `stop` channels in [select](../03-select/)), and they didn't work together. Goroutines kept working for nobody, or leaked. Request data, like a request ID, had to be passed by hand through every function.

**Solution:** one standard type for all of it. Google published it as `golang.org/x/net/context` in 2014, and it joined the standard library in Go 1.7 (2016). You pass `ctx` as the first parameter through the whole call chain, and it carries the cancel signal, the deadline, and request values. The standard library uses it everywhere: `r.Context()` in net/http, `db.QueryContext`, `exec.CommandContext`.

## API

The whole package: 3 types, 4 interface methods, 2 errors and 12 functions.

### Types (3)

| Type | What it is | Example |
| --- | --- | --- |
| `Context` | The main interface | every example |
| `CancelFunc` | `func()`. Call it to cancel | [01-cancellation/01-basic](01-cancellation/01-basic/) |
| `CancelCauseFunc` | `func(cause error)`. Cancel with a reason (Go 1.20+) | [01-cancellation/03-cancel-cause](01-cancellation/03-cancel-cause/) |

### `Context` methods (4)

| Method | Returns | Meaning | Example |
| --- | --- | --- | --- |
| `ctx.Deadline()` | `time.Time, bool` | When it expires. `false` = no deadline | [03-deadline/01-basic](03-deadline/01-basic/) |
| `ctx.Done()` | `<-chan struct{}` | Closed on cancel or timeout | [01-cancellation/01-basic](01-cancellation/01-basic/) |
| `ctx.Err()` | `error` | `nil`, then `Canceled` or `DeadlineExceeded` | [02-timeout/01-basic](02-timeout/01-basic/) |
| `ctx.Value(key)` | `any` | The value for `key`, or `nil` | [04-values/01-request-id](04-values/01-request-id/) |

There's no `Cancel()` method. Code that gets a ctx can only watch it. Only its creator holds `cancel`.

### Errors (2)

| Variable | When | Example |
| --- | --- | --- |
| `context.Canceled` | `cancel()` was called | [01-cancellation/01-basic](01-cancellation/01-basic/) |
| `context.DeadlineExceeded` | The time ran out | [02-timeout/01-basic](02-timeout/01-basic/) |

### Functions (12)

| Function | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `Background()` | `Context` | 1.7 | Empty root context | every example |
| `TODO()` | `Context` | 1.7 | Empty root, "pass a real ctx here later" | none |
| `WithCancel(parent)` | `ctx, CancelFunc` | 1.7 | Cancel by hand | [01-cancellation/01-basic](01-cancellation/01-basic/), [04-parent-child](01-cancellation/04-parent-child/) |
| `WithTimeout(parent, d)` | `ctx, CancelFunc` | 1.7 | Cancel after `d` | [02-timeout/01-basic](02-timeout/01-basic/), [03-deadline/02-parent-deadline](03-deadline/02-parent-deadline/) |
| `WithDeadline(parent, t)` | `ctx, CancelFunc` | 1.7 | Cancel at time `t` | [03-deadline/01-basic](03-deadline/01-basic/) |
| `WithValue(parent, key, val)` | `ctx` | 1.7 | Attach a value | [04-values/01-request-id](04-values/01-request-id/) |
| `WithCancelCause(parent)` | `ctx, CancelCauseFunc` | 1.20 | Cancel with a reason | [01-cancellation/03-cancel-cause](01-cancellation/03-cancel-cause/) |
| `Cause(ctx)` | `error` | 1.20 | Read the reason | [01-cancellation/03-cancel-cause](01-cancellation/03-cancel-cause/), [02-timeout/03-timeout-cause](02-timeout/03-timeout-cause/), [03-deadline/03-deadline-cause](03-deadline/03-deadline-cause/) |
| `WithTimeoutCause(parent, d, cause)` | `ctx, CancelFunc` | 1.21 | Timeout with a reason | [02-timeout/03-timeout-cause](02-timeout/03-timeout-cause/) |
| `WithDeadlineCause(parent, t, cause)` | `ctx, CancelFunc` | 1.21 | Deadline with a reason | [03-deadline/03-deadline-cause](03-deadline/03-deadline-cause/) |
| `WithoutCancel(parent)` | `ctx` | 1.21 | Keep the values, drop the cancel | [05-advanced/02-without-cancel](05-advanced/02-without-cancel/) |
| `AfterFunc(ctx, f)` | `stop func() bool` | 1.21 | Run `f` when ctx is done | [05-advanced/01-after-func](05-advanced/01-after-func/) |

### How it works

Each `With...` function wraps its parent, so contexts form a chain:

```text
Background()                 never cancelled
└── WithValue(req-42)        adds a request ID
    └── WithCancel           adds a cancel button
        └── WithTimeout(1s)  adds a timer that calls cancel
```

- A cancel context keeps a list of its children and cancels them too. Nothing points back up, so a child never cancels its parent ([04-parent-child](01-cancellation/04-parent-child/)).
- `Value` checks its own key, then asks its parent, up to the root. Deep chains make lookups slower.
- `Background()` is empty: no deadline, no values, and `Done()` returns a `nil` channel, which is never ready ([select/06-nil-channel](../03-select/06-nil-channel/)). It's the root of the chain.
- `TODO()` is the same as `Background()`. Use it as a placeholder where a real ctx isn't passed in yet.

## Things to remember

- **Pass `ctx` as the first parameter.** Don't keep it inside a struct.
- **Always call `cancel`, usually with `defer`.** `go vet` warns you if you forget.
- Only put request data in context values, not normal function arguments.
