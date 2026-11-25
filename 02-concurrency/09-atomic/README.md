# Atomic

`sync/atomic` safely updates a single value without a lock.

## Why

**Problem:** a mutex is more than you need for one counter or one flag. Every reader and writer has to lock and unlock, and they wait for each other.

**Solution:** atomic types. `count++` is really three steps: read the number, add 1, write it back. Two goroutines can mix up those steps and lose an update. An atomic method like `count.Add(1)` does all three as one step, so no other goroutine can get in between. There is no lock, so nobody waits.

## API

### Types

Pick the type that matches your value.

| Type | Go | What it holds | Example |
| --- | --- | --- | --- |
| `atomic.Int32`, `atomic.Int64` | 1.19 | A whole number, like a counter | [01-counter](01-counter/), [04-compare-and-swap](04-compare-and-swap/) |
| `atomic.Uint32`, `atomic.Uint64`, `atomic.Uintptr` | 1.19 | A whole number that can't go below zero | none |
| `atomic.Bool` | 1.19 | `true` or `false`, like an on/off switch | [02-flag](02-flag/) |
| `atomic.Pointer[T]` | 1.19 | A pointer to a `T`, like `*Config`, swapped as a whole | [03-pointer](03-pointer/) |
| `atomic.Value` | 1.4 | Any value. Older and untyped, so use `Pointer[T]` instead | none |

### Methods

Call these on a variable of one of the types above, like `count.Add(1)`.

| Method | Returns | Go | Works on | What it does | Example |
| --- | --- | --- | --- | --- | --- |
| `Load()` | the value | 1.19 | every type | Read the current value | [01-counter](01-counter/), [02-flag](02-flag/), [03-pointer](03-pointer/) |
| `Store(v)` | | 1.19 | every type | Set a new value | [02-flag](02-flag/), [03-pointer](03-pointer/) |
| `Swap(v)` | the old value | 1.19 | every type | Set a new value and get the old one back | none |
| `CompareAndSwap(old, new)` | `bool` | 1.19 | every type | Set `new` only if the value is still `old`. Returns `false` if someone changed it first | [04-compare-and-swap](04-compare-and-swap/) |
| `Add(delta)` | the new value | 1.19 | numbers | Add to the number. On `Int` types, use a negative number to subtract | [01-counter](01-counter/) |
| `And(mask)`, `Or(mask)` | the old value | 1.23 | numbers | Bit operations. Rarely needed | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-counter` | `go run ./02-concurrency/09-atomic/01-counter` | Request and byte counters updated by many handlers, with no lost updates |
| `02-flag` | `go run ./02-concurrency/09-atomic/02-flag` | A maintenance mode switch that workers check before each job |
| `03-pointer` | `go run ./02-concurrency/09-atomic/03-pointer` | Hot reload: swap a whole config while handlers keep reading it |
| `04-compare-and-swap` | `go run ./02-concurrency/09-atomic/04-compare-and-swap` | Problem: checking, then taking, hands out more coupons than exist, and `-race` finds nothing. CompareAndSwap fixes it |
| `05-late-read` | `go run ./02-concurrency/09-atomic/05-late-read` | Jobs read a shared atomic ID too late and get the same one. No race, still wrong |

## Things to remember

- **Use atomic for one value. Use a mutex when several fields must change together.**
- **Don't mix atomic and normal reads or writes on the same variable.** The typed atomics make that hard to do by mistake.
- **Never change a value after you `Store` its pointer.** Build a new one and `Store` that (`03`).
- A check followed by a separate change isn't safe. Use a `CompareAndSwap` loop (`04`).
- **Give a `CompareAndSwap` loop a retry limit.** A swap only fails when another goroutine changed the value first, so the loop does end. But one unlucky goroutine can lose again and again. At the limit, return a "try again" error, so the user isn't left waiting (`04`).
- **Atomic stops the race, but not wrong logic.** If a goroutine needs the value from its own loop round, give it a copy before `go`. Don't let it read the shared value later (`05`).
- For a one-time "stop" that many goroutines must notice, use `ctx`. `atomic.Bool` is for a value that flips back and forth (`02`).

## Go version notes

- Types like `atomic.Int64` came in Go 1.19. Use them instead of the older `atomic.AddInt64` functions.
- `And` and `Or` came in Go 1.23.
