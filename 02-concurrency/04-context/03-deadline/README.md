# Context Deadline

Stop work at a fixed time with `context.WithDeadline`.

See the [full API overview](../README.md#api).

## Why

**Problem:** a request has a total time budget. Each step needs to know how much time is left, and a step deep inside must not stretch the budget.

**Solution:** `WithDeadline` sets a fixed end time. `ctx.Deadline()` tells every step when that is, and a child can only make it earlier.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `WithDeadline(parent, t)` | `ctx, CancelFunc` | 1.7 | Cancel at time `t` | [01-basic](01-basic/), [02-parent-deadline](02-parent-deadline/) |
| `WithDeadlineCause(parent, t, cause)` | `ctx, CancelFunc` | 1.21 | Same, with a reason for `Cause` | [03-deadline-cause](03-deadline-cause/) |
| `ctx.Deadline()` | `time.Time, bool` | 1.7 | When it ends. `false` = no deadline | [01-basic](01-basic/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/04-context/03-deadline/01-basic` | Skip a job that can't finish before the deadline |
| `02-parent-deadline` | `go run ./02-concurrency/04-context/03-deadline/02-parent-deadline` | A child can make the deadline earlier, but never later |
| `03-deadline-cause` | `go run ./02-concurrency/04-context/03-deadline/03-deadline-cause` | Say why the deadline stopped the work with `WithDeadlineCause` |

## Things to remember

- `WithTimeout(d)` is the same thing as `WithDeadline(time.Now().Add(d))`.
- `ctx.Deadline()` tells you if there's a deadline, so you can skip work that won't finish in time (`01`).
- **The earliest deadline in the chain always wins** (`02`).
- `Cause(ctx)` gives the error passed to the `...Cause` function. **The returned `cancel` takes no cause**, so a manual cancel gives `context.Canceled` (`03`).
