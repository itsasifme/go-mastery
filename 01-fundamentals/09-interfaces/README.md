# Interfaces

An interface is a list of methods. Any type that has those methods fits the interface. There's no `implements` keyword.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-shapes` | `go run ./01-fundamentals/09-interfaces/01-shapes` | `Rectangle` and `Triangle` both fit a `Shape` interface |
| `02-separate-contracts` | `go run ./01-fundamentals/09-interfaces/02-separate-contracts` | One type fits two small interfaces |
| `03-combined-contract` | `go run ./01-fundamentals/09-interfaces/03-combined-contract` | One interface with both methods |

## Things to remember

- Small interfaces are easier to reuse.
- The two author examples don't print a new line at the end.
