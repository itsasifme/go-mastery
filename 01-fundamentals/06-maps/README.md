# Maps

> 📝 **Planned.** No code yet. This is what I plan to add here.

A map stores key and value pairs.

## Examples to add

- `01-basics`: add, read, and delete, plus the `value, ok` check
- `02-loop-order`: the loop order is random, so sort the keys when order matters
- `03-sets`: use `map[T]struct{}` as a set
- `04-maps-package`: `maps.Clone`, `maps.Copy`, and `clear`

## Things to remember

- Reading from a nil map is fine. Writing to it panics.
- The order of a `for range` loop over a map is random on purpose.
- Maps aren't safe when two goroutines write at the same time. The program crashes, and `recover` can't stop it.

## Go version notes

- The `maps` package and `clear` came in Go 1.21.
