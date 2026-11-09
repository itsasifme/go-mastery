# Mutex

> 📝 **Planned.** No code yet. This is what I plan to add here.

`sync.Mutex` lets only one goroutine touch shared data at a time.

## Examples to add

- `01-race`: a counter without a lock, so `go run -race` catches the bug
- `02-counter`: the same counter, fixed with a mutex
- `03-struct`: a struct with a mutex next to the fields it protects

## Things to remember

- A zero `sync.Mutex` is ready to use.
- **Never copy a mutex.** `go vet` warns you about it.
- Put `defer mu.Unlock()` right after `mu.Lock()`.
- **Locking the same mutex twice in one goroutine gets it stuck forever.**
