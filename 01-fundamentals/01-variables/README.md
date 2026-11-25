# Variables

> 📝 **Planned.** No code yet. This is what I plan to add here.

How to create variables and constants in Go, and where they can be used.

## Examples to add

- `01-declaration`: `var` vs `:=`, and zero values
- `02-constants-iota`: constants and `iota`
- `03-scope-shadowing`: scope, and how `:=` can hide an outer variable

## Things to remember

- **Every type has a zero value.** A new variable is never "empty" or random.
- `:=` only works inside functions.
- Using `:=` inside an `if` block **can create a new `err` that hides the outer one**. `go vet` won't warn you about it.
- An unused local variable is a compile error.
