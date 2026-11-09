# errgroup

> 📝 **Planned.** No code yet. This is what I plan to add here.

`errgroup` is like a WaitGroup that also returns the first error.

## Examples to add

- `01-basic`: `g.Go` and `g.Wait`
- `02-with-context`: stop the other goroutines when one fails
- `03-limit`: run only a few at a time with `SetLimit`

## Things to remember

- `Wait` waits for all goroutines, then **returns the first error**.
- The context from `WithContext` is cancelled after the first error.

## Package note

`errgroup` isn't in the standard library. It comes from `golang.org/x/sync`. Pick a version that works with Go 1.23, so `go get` doesn't bump the Go version.
