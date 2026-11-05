# Error Handling

> 📝 **Planned.** No code yet. This is what I plan to add here.

In Go, an error is just a value that a function returns.

## Examples to add

- `01-basics`: return an error and check `if err != nil`
- `02-wrapping`: add context with `%w`, then check it with `errors.Is` and `errors.As`
- `03-custom-errors`: a fixed error value vs your own error type
- `04-join`: combine errors with `errors.Join`

## Things to remember

- Handle the error or return it. Don't do both.
- After wrapping, use `errors.Is` instead of `==`.
- A nil pointer stored in an `error` is **not** nil. This is a classic trap.
- Use `panic` for bugs, not for normal errors.

## Go version notes

- `errors.Join` came in Go 1.20.
