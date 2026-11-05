# Structs

A struct groups related fields into one type.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic-fields` | `go run ./01-fundamentals/07-structs/01-basic-fields` | Read and change struct fields |
| `02-function-field` | `go run ./01-fundamentals/07-structs/02-function-field` | A struct field that holds a function |
| `03-json-tags` | `go run ./01-fundamentals/07-structs/03-json-tags` | Turn a struct into JSON using field tags |

## Things to remember

- Only fields that start with a capital letter go into JSON.
- In `03-json-tags`, `Address` is embedded in `User`, so its fields show up at the top level of the JSON.
