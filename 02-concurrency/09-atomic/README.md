# Atomic

`sync/atomic` safely updates a single value without a lock.

## Why

**Problem:** a mutex is more than you need for one counter or one flag. Every reader and writer has to lock and unlock, and they wait for each other.

**Solution:** atomic types. Methods like `Add`, `Load`, `Store` and `CompareAndSwap` change one value in a single step that no other goroutine can split. Nobody waits.

## API

| API | Returns | Go | What it does | Example |
| --- | --- | --- | --- | --- |
| `atomic.Int32`, `Int64`, `Uint32`, `Uint64`, `Uintptr` | | 1.19 | An integer you can update safely | [01-counter](01-counter/) |
| `atomic.Bool` | | 1.19 | A true/false value | [02-flag](02-flag/) |
| `atomic.Pointer[T]` | | 1.19 | A pointer to `T`, swapped as a whole | [03-pointer](03-pointer/) |
| `atomic.Value` | | 1.4 | Older, untyped version of `Pointer[T]` | none |
| `x.Load()` | value | 1.19 | Read the current value | [01-counter](01-counter/), [02-flag](02-flag/) |
| `x.Store(v)` | | 1.19 | Set a new value | [02-flag](02-flag/), [03-pointer](03-pointer/) |
| `x.Add(delta)` | new value | 1.19 | Add to an integer (use a negative number to subtract) | [01-counter](01-counter/) |
| `x.Swap(v)` | old value | 1.19 | Set a new value and get the old one | none |
| `x.CompareAndSwap(old, new)` | `bool` | 1.19 | Set `new` only if the value is still `old` | [04-compare-and-swap](04-compare-and-swap/) |
| `x.And(mask)`, `x.Or(mask)` | old value | 1.23 | Bit operations on integers | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-counter` | `go run ./02-concurrency/09-atomic/01-counter` | Request and byte counters updated by many handlers, with no lost updates |
| `02-flag` | `go run ./02-concurrency/09-atomic/02-flag` | A maintenance mode switch that workers check before each job |
| `03-pointer` | `go run ./02-concurrency/09-atomic/03-pointer` | Hot reload: swap a whole config while handlers keep reading it |
| `04-compare-and-swap` | `go run ./02-concurrency/09-atomic/04-compare-and-swap` | Many users claim a few coupons, and the stock never goes below zero |

## Things to remember

- **Use atomic for one value. Use a mutex when several fields must change together.**
- **Don't mix atomic and normal reads or writes on the same variable.** The typed atomics make that hard to do by mistake.
- **Never change a value after you `Store` its pointer.** Build a new one and `Store` that (`03`).
- A check followed by a separate change isn't safe. Use a `CompareAndSwap` loop (`04`).
- For a one-time signal that many goroutines must notice, a closed channel is still the clearer choice. `atomic.Bool` is for a value that flips back and forth (`02`).

## Go version notes

- Types like `atomic.Int64` came in Go 1.19. Use them instead of the older `atomic.AddInt64` functions.
- `And` and `Or` came in Go 1.23.
