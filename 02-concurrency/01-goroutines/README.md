# Goroutines

A goroutine is a function that runs at the same time as the rest of the program. Start one with `go f()`.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/01-goroutines/01-basic` | Start a goroutine without waiting for it |

## Things to remember

- `main` doesn't wait for goroutines. When `main` ends, the program ends.
- So `Hello World!` may or may not print. That's the point of this example.
- To wait for goroutines, see [WaitGroup](../05-waitgroup/).
