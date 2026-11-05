# RWMutex

> 📝 **Planned.** No code yet. This is what I plan to add here.

`sync.RWMutex` lets many goroutines read at the same time, but only one write.

## Examples to add

- `01-cache`: a map that's read a lot and changed rarely
- `02-benchmark`: compare `Mutex` and `RWMutex` with a benchmark

## Things to remember

- It only helps when there are many more reads than writes. Measure it first.
- Don't call `RLock` twice in the same goroutine. It can get stuck if a writer is waiting.
