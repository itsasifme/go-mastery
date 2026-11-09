# Generics

> 📝 **Planned.** No code yet. This is what I plan to add here.

Generics let you write one function or type that works with many types.

## Examples to add

- `01-functions`: generic `Map` and `Filter`
- `02-constraints`: `any`, `comparable`, `cmp.Ordered`, and your own constraints
- `03-types`: a generic `Stack[T]`

## Things to remember

- **Use generics when the code is the same for every type. Use interfaces when the behavior is different.**
- `~int` means "any type built on `int`".
- A method can't have its own type parameters.

## Go version notes

- Generics came in Go 1.18. The `cmp` package came in Go 1.21.
- Generic type aliases need Go 1.24, so they can't be used here yet.
