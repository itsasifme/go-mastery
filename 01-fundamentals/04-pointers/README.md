# Pointers

A pointer holds the memory address of a value. You can use it to read or change that value.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./01-fundamentals/04-pointers/01-basic` | Change an int and a struct through a pointer, and print the addresses |

## Things to remember

- The printed addresses are different on every run.
- The last line turns the pointer into a `uintptr` using `unsafe`. That's only to look at the raw address. Normal Go code shouldn't do this.
