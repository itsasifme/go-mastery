# Structs

A struct groups related fields into one type.

## Why

**Problem:** related values, like a user's name, email and address, kept as separate variables are easy to mix up and painful to pass around.

**Solution:** a struct groups them into one named type, with a name for each field.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `type T struct { ... }` | Define a struct | [01-basic-fields](01-basic-fields/) |
| `T{a, b}` | Create a value, fields in order | [01-basic-fields](01-basic-fields/) |
| `T{Name: a}` | Create a value, fields by name | [02-function-field](02-function-field/) |
| `s.field` | Read or change a field | [01-basic-fields](01-basic-fields/) |
| a `func` field | Store behaviour in a field | [02-function-field](02-function-field/) |
| `` `json:"name"` `` | A field tag, read by packages like `encoding/json` | [03-json-tags](03-json-tags/) |
| embedded struct | Its fields are promoted to the outer struct | [03-json-tags](03-json-tags/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-basic-fields` | `go run ./01-fundamentals/07-structs/01-basic-fields` | Read and change struct fields |
| `02-function-field` | `go run ./01-fundamentals/07-structs/02-function-field` | A struct field that holds a function |
| `03-json-tags` | `go run ./01-fundamentals/07-structs/03-json-tags` | Turn a struct into JSON using field tags |

## Things to remember

- **Only fields that start with a capital letter go into JSON.**
- In `03-json-tags`, `Address` is embedded in `User`, so its fields show up at the top level of the JSON.
