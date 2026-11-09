# Pointers

A pointer holds the memory address of a value. You can use it to read or change that value.

## Why

**Problem:** Go passes everything by value. A function gets a copy, so it can't change the caller's variable, and a big struct is copied on every call.

**Solution:** pass a pointer, which is the address of the value. The function can change the original, and only the address (8 bytes) is copied.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `&x` | The address of `x` | [01-basic](01-basic/) |
| `*T` | The type "pointer to a `T`" | [01-basic](01-basic/) |
| `*p` | Read or change the value `p` points to | [01-basic](01-basic/) |
| `new(T)` | Make a zero `T` and return a `*T` | none |
| `nil` | A pointer to nothing. `*p` on it panics | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic` | `go run ./01-fundamentals/04-pointers/01-basic` | Change an int and a struct through a pointer, and print the addresses |

## Things to remember

- The printed addresses are different on every run.
- The last line turns the pointer into a `uintptr` using `unsafe`. That's only to look at the raw address. Normal Go code shouldn't do this.
