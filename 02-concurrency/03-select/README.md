# Select

`select` waits on several channel operations at once and runs the first one that's ready. It's how a goroutine listens to many things at the same time: results, timeouts, and stop signals.

## Why

**Problem:** `<-ch` waits on one channel only. A goroutine that must react to a result, a timeout and a stop signal can't wait on all three.

**Solution:** `select` waits on several channel operations at once and runs whichever is ready first.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `select { case ...: }` | Wait on several channels | [01-basic](01-basic/) |
| several cases ready | One is picked at random | [02-random-pick](02-random-pick/) |
| `default:` | Don't wait if nothing is ready | [03-default](03-default/) |
| `time.After(d)` | A channel that gets a value once, after `d` | [04-timeout](04-timeout/) |
| `for { select { ... } }` | A long-running loop | [05-for-select](05-for-select/) |
| `ch = nil` | Turn off a case | [06-nil-channel](06-nil-channel/), [07-nil-or-stop](07-nil-or-stop/) |
| `time.NewTicker(d)`, `ticker.C`, `ticker.Stop()` | A channel that gets a value every `d` | [08-ticker](08-ticker/) |
| `select` + `default`, then `select` | Check one case first: priority | [09-priority](09-priority/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/03-select/01-basic` | Wait on two channels and take whichever is ready first |
| `02-random-pick` | `go run ./02-concurrency/03-select/02-random-pick` | When several cases are ready, the pick is random |
| `03-default` | `go run ./02-concurrency/03-select/03-default` | `trySend` and `tryReceive` that never block |
| `04-timeout` | `go run ./02-concurrency/03-select/04-timeout` | Give up after a fixed time with `time.After` |
| `05-for-select` | `go run ./02-concurrency/03-select/05-for-select` | A worker loop that handles jobs until it's told to stop |
| `06-nil-channel` | `go run ./02-concurrency/03-select/06-nil-channel` | Read two channels until both close, turning off a closed one with `nil` |
| `07-nil-or-stop` | `go run ./02-concurrency/03-select/07-nil-or-stop` | Read two channels until both close, or stop early when told |
| `08-ticker` | `go run ./02-concurrency/03-select/08-ticker` | Do something on a timer until it's told to stop |
| `09-priority` | `go run ./02-concurrency/03-select/09-priority` | `select` has no priority. How to make urgent jobs go first |

## Things to remember

- If more than one case is ready, Go picks one **at random**. The order of the cases doesn't matter (`02`).
- With a `default` case, `select` never waits (`03`). **Don't use `default` in a loop that only waits**: it spins and burns CPU.
- `select {}` with no cases waits forever.
- Inside `select`, **`break` only leaves the `select`, not the surrounding `for` loop**. Use `return` or a labeled `break` (`05`).
- **A closed channel is always ready.** Set it to `nil` to stop `select` from picking it (`06`).
- A `nil` channel is never ready, so its case is simply skipped.
- One channel until closed: `for range`. Several until all closed: `nil`. Stop on a signal: `return`. Both: `07`.
- Always `Stop` a ticker when you're done with it (`08`).
- If one case must win, check it first with its own `select` + `default` (`09`).
- `09`, and the Try lines in `06` and `07`, show a problem on purpose. Each has a fix or a note to try.

## Go version notes

- Since Go 1.23, a timer from `time.After` that's no longer needed is cleaned up even before it fires. Before that, calling `time.After` in a busy loop could pile up timers. For timeouts in real code, `context.WithTimeout` is usually the better choice (see [context](../04-context/)).
