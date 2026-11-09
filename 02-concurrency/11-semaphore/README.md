# Semaphore

> 📝 **Planned.** No code yet. This is what I plan to add here.

A semaphore limits how many goroutines can do something at the same time.

## Examples to add

- `01-channel`: a buffered channel used as a semaphore
- `02-weighted`: `semaphore.Weighted` from `golang.org/x/sync`

## Things to remember

- **Take a slot before you start the goroutine, not inside it.**
- If you also need errors, `errgroup.SetLimit` is often simpler.
