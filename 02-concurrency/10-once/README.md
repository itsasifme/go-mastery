# Once

> 📝 **Planned.** No code yet. This is what I plan to add here.

`sync.Once` runs a function only one time, even when many goroutines call it.

## Examples to add

- `01-lazy-init`: set something up only on first use
- `02-once-func`: `sync.OnceFunc` and `sync.OnceValue`

## Things to remember

- **If the function panics, `Once` still counts it as done.**
- Calling `once.Do` inside its own function gets stuck forever.

## Go version notes

- `OnceFunc`, `OnceValue`, and `OnceValues` came in Go 1.21.
