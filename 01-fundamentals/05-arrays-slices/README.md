# Arrays and Slices

> 📝 **Planned.** No code yet. This is what I plan to add here.

Arrays have a fixed size. Slices are flexible and are what you use most of the time.

## Examples to add

- `01-arrays`: arrays are copied when you assign them
- `02-slice-basics`: `len`, `cap`, `append`, `make`
- `03-shared-memory`: two slices sharing the same data
- `04-slices-package`: sort and search with the `slices` package

## Things to remember

- Always use the result of `append`. It may create a new array behind the scenes.
- Slicing a slice doesn't copy it. Both slices share the same data.
- A nil slice and an empty slice both have length 0, but JSON shows them as `null` and `[]`.

## Go version notes

- The `slices` package came in Go 1.21.
- `for i := range 10` needs Go 1.22.
