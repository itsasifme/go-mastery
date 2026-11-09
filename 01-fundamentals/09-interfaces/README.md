# Interfaces

An interface is a list of methods. Any type that has those methods fits the interface. There's no `implements` keyword.

## Why

**Problem:** a function that takes a `Rectangle` can't take a `Triangle`. Code gets tied to concrete types, and it's hard to swap in a fake for tests.

**Solution:** an interface lists methods. Any type that has them fits automatically, so one function works with all of them.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `type I interface { M() }` | Define a set of methods | [01-shapes](01-shapes/) |
| (no keyword) | A type fits just by having the methods | [01-shapes](01-shapes/) |
| `var i I = v` | Store any value that fits | [02-separate-contracts](02-separate-contracts/) |
| several small interfaces | One type fits all of them | [02-separate-contracts](02-separate-contracts/) |
| one interface, many methods | A type must have all of them | [03-combined-contract](03-combined-contract/) |
| `v.(T)` | Type assertion: get the concrete type back | none |
| `any` | The empty interface: holds any value | none |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-shapes` | `go run ./01-fundamentals/09-interfaces/01-shapes` | `Rectangle` and `Triangle` both fit a `Shape` interface |
| `02-separate-contracts` | `go run ./01-fundamentals/09-interfaces/02-separate-contracts` | One type fits two small interfaces |
| `03-combined-contract` | `go run ./01-fundamentals/09-interfaces/03-combined-contract` | One interface with both methods |

## Things to remember

- **Small interfaces are easier to reuse.**
- The two author examples don't print a new line at the end.
