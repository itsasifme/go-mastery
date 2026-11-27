# Interfaces

An interface is a list of methods. Any type that has those methods fits the interface. There's no `implements` keyword.

## Why

**Problem:** a function that takes a `Rectangle` can't take a `Triangle`. Code gets tied to concrete types, and it's hard to swap in a fake for tests.

**Solution:** an interface lists methods. Any type that has them fits automatically, so one function works with all of them.

## API

| Syntax | Kind | What it does | Example |
| --- | --- | --- | --- |
| `type I interface { M() }` | declaration | Define a set of methods | [01-shapes](01-shapes/) |
| (no keyword) | rule | A type fits just by having the methods | [01-shapes](01-shapes/) |
| `var i I = v` | declaration | Store any value that fits | [02-separate-contracts](02-separate-contracts/) |
| several small interfaces | rule | One type fits all of them | [02-separate-contracts](02-separate-contracts/) |
| one interface, many methods | rule | A type must have all of them | [03-combined-contract](03-combined-contract/) |
| `any` | type | The empty interface: holds a value of any type | [04-type-assertion](04-type-assertion/), [05-type-switch](05-type-switch/) |
| `x.(T)` | operator | Type assertion: get the value back as a `T`. Panics if it isn't one | [04-type-assertion](04-type-assertion/) |
| `v, ok := x.(T)` | operator | The same, but never panics: `ok` is false if `x` isn't a `T` | [04-type-assertion](04-type-assertion/) |
| `switch v := x.(type)` | statement | Type switch: one case per type. In each case, `v` has that type | [05-type-switch](05-type-switch/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-shapes` | `go run ./01-fundamentals/09-interfaces/01-shapes` | `Rectangle` and `Triangle` both fit a `Shape` interface |
| `02-separate-contracts` | `go run ./01-fundamentals/09-interfaces/02-separate-contracts` | One type fits two small interfaces |
| `03-combined-contract` | `go run ./01-fundamentals/09-interfaces/03-combined-contract` | One interface with both methods |
| `04-type-assertion` | `go run ./01-fundamentals/09-interfaces/04-type-assertion` | Problem: a setting saved as text makes `x.(int)` panic. Fix: check with `v, ok := x.(int)` and use a default |
| `05-type-switch` | `go run ./01-fundamentals/09-interfaces/05-type-switch` | A settings page shows text, numbers and switches, each in its own case |

## Things to remember

- **Small interfaces are easier to reuse.**
- **`x.(T)` panics when you're wrong.** Use `v, ok := x.(T)` unless you're sure, like when only your own code puts values in (`04`).
- A missing map key gives `nil`, and asserting `nil` to any type gives `ok == false` (`04`).
- **Use a type switch when there are two or more types to check** (`05`). In a case with two types, like `case int, int64`, `v` stays `any`.
- The two author examples don't print a new line at the end.
