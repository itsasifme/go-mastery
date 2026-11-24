# Variables

How to create variables and constants in Go, and where they can be used.

## Why

**Problem:** a program needs names for its values: some that change, like a cart total, and some that never do, like an upload limit. In some languages a new variable holds random garbage until you set it, and a typo can quietly create a second variable.

**Solution:** Go gives every new variable its type's zero value, so it's never random. `var` and `:=` declare variables, `const` declares values that can't change, and `iota` numbers a list of constants for you. Each variable lives only inside the block where it's declared.

## API

| Syntax | Kind | What it does | Example |
| --- | --- | --- | --- |
| `var x T` | declaration | A variable set to the zero value of `T` | [01-declaration](01-declaration/) |
| `var x = v` | declaration | A variable with a value. Go works out the type. Also works outside functions | [01-declaration](01-declaration/) |
| `x := v` | declaration | Declare and set in one step. Only inside functions | [01-declaration](01-declaration/) |
| `a, b := v1, v2` | declaration | Declare several at once | [01-declaration](01-declaration/) |
| `x = v` | operator | Change a variable that already exists | [01-declaration](01-declaration/), [03-scope-shadowing](03-scope-shadowing/) |
| zero value | rule | `0`, `""`, `false`, or `nil` for pointers, slices and maps | [01-declaration](01-declaration/) |
| `const x = v` | declaration | A value fixed at compile time. Can't be changed | [02-constants-iota](02-constants-iota/) |
| `iota` | value | Counts the lines of a `const` block: 0, 1, 2... | [02-constants-iota](02-constants-iota/) |
| `_` | value | The blank identifier: skips a value you don't need | [02-constants-iota](02-constants-iota/) |
| `{ ... }` | rule | A variable exists only inside the block where it's declared | [03-scope-shadowing](03-scope-shadowing/) |
| `:=` in an inner block | rule | Makes a new variable that hides an outer one with the same name | [03-scope-shadowing](03-scope-shadowing/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-declaration` | `go run ./01-fundamentals/01-variables/01-declaration` | A shopping cart: the zero value of each type, then `:=` and its types |
| `02-constants-iota` | `go run ./01-fundamentals/01-variables/02-constants-iota` | Order statuses numbered with `iota`, and an upload limit built from byte sizes |
| `03-scope-shadowing` | `go run ./01-fundamentals/01-variables/03-scope-shadowing` | Problem: `:=` inside an `if` hides the outer `err`, so a failed checkout looks fine |

## Things to remember

- **Every type has a zero value.** A new variable is never "empty" or random (`01`).
- `:=` only works inside functions. Package-level variables need `var`.
- **Don't use `iota` for values that are saved or sent somewhere.** Adding a line in the middle changes every number after it (`02`).
- Using `:=` inside an `if` block **can create a new `err` that hides the outer one**. `go vet` won't warn you about it (`03`).
- An unused local variable is a compile error.
