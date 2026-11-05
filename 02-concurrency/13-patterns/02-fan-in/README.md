# Fan-In

> 📝 **Planned.** No code yet. This is what I plan to add here.

Merge many channels into one.

## Examples to add

- `01-merge`: a `merge` function that combines several channels

## Things to remember

- Use one goroutine per input channel. Close the output when all of them are done.
- Values from different inputs can come out in any order.
