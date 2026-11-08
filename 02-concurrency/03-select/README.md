# Select

`select` waits on several channel operations at once and runs the first one that's ready. It's how a goroutine listens to many things at the same time: results, timeouts, and stop signals.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/03-select/01-basic` | Wait on two channels and take whichever is ready first |
| `02-random-pick` | `go run ./02-concurrency/03-select/02-random-pick` | When several cases are ready, the pick is random |
| `03-default` | `go run ./02-concurrency/03-select/03-default` | `trySend` and `tryReceive` that never block |
| `04-timeout` | `go run ./02-concurrency/03-select/04-timeout` | Give up after a fixed time with `time.After` |
| `05-for-select` | `go run ./02-concurrency/03-select/05-for-select` | A worker loop that handles jobs until it's told to stop |
| `06-nil-channel` | `go run ./02-concurrency/03-select/06-nil-channel` | Merge two channels, and turn off a closed one with `nil` |
| `07-ticker` | `go run ./02-concurrency/03-select/07-ticker` | Do something on a timer until it's told to stop |
| `08-priority` | `go run ./02-concurrency/03-select/08-priority` | `select` has no priority. How to make urgent jobs go first |

## Things to remember

- If more than one case is ready, Go picks one **at random**. The order of the cases doesn't matter (`02`).
- With a `default` case, `select` never waits (`03`). Don't use `default` in a loop that only waits: it spins and burns CPU.
- `select {}` with no cases waits forever.
- Inside `select`, `break` only leaves the `select`, not the surrounding `for` loop. Use `return` or a labeled `break` (`05`).
- A closed channel is always ready. Set it to `nil` to stop `select` from picking it (`06`).
- A `nil` channel is never ready, so its case is simply skipped.
- Always `Stop` a ticker when you're done with it (`07`).
- If one case must win, check it first with its own `select` + `default` (`08`).
- `08` and `06`'s Try line show a problem on purpose. Each has a fix or a note to try.

## Go version notes

- Since Go 1.23, a timer from `time.After` that's no longer needed is cleaned up even before it fires. Before that, calling `time.After` in a busy loop could pile up timers. For timeouts in real code, `context.WithTimeout` is usually the better choice (see [context](../04-context/)).
