# Defer, Panic, and Recover

> 📝 **Planned.** No code yet. This is what I plan to add here.

`defer` runs code when a function ends. `panic` stops the program. `recover` can catch a panic.

## Examples to add

- `01-defer-order`: deferred calls run last-in, first-out
- `02-defer-return`: a deferred function changes the return value
- `03-recover`: catch a panic and turn it into an error

## Things to remember

- The arguments of a deferred call are **evaluated right away**, not when it runs.
- `recover` only works inside a deferred function.
- **A panic in another goroutine crashes the whole program.**
- A `defer` inside a loop runs when the function ends, not after each loop round.
