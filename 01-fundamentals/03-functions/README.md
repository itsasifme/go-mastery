# Functions

> 📝 **Planned.** No code yet. This is what I plan to add here.

Functions in Go can return many values, and you can pass them around like any other value.

## Examples to add

- `01-multiple-returns`: return a value and an error
- `02-variadic`: functions that take any number of arguments
- `03-closures`: a function that remembers a variable
- `04-function-values`: pass a function to another function

## Things to remember

- Go always passes a copy. But a copy of a slice, map, or pointer **still points to the same data**.
- **A closure keeps the variable itself, not a copy of its value.**
- Since Go 1.22, every loop round gets its own loop variable. Closures made in a loop no longer share one variable.
