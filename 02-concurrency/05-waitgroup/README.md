# WaitGroup

`sync.WaitGroup` lets you wait until a group of goroutines finish.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./02-concurrency/05-waitgroup/01-basic` | Wait for two goroutines before the program ends |

## Things to remember

- Call `Add` before you start the goroutine.
- Call `Done` with `defer`, so it always runs.
- The print order can change between runs.
