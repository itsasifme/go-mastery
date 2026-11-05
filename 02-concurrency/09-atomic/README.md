# Atomic

> 📝 **Planned.** No code yet. This is what I plan to add here.

`sync/atomic` safely updates a single value without a lock.

## Examples to add

- `01-counter`: `atomic.Int64`
- `02-flag`: `atomic.Bool` as a stop signal
- `03-pointer`: swap a config with `atomic.Pointer[T]`

## Things to remember

- Use atomic for one value. Use a mutex when several fields must change together.
- Don't mix atomic and normal reads or writes on the same variable.

## Go version notes

- Types like `atomic.Int64` came in Go 1.19. Use them instead of the older `atomic.AddInt64` functions.
