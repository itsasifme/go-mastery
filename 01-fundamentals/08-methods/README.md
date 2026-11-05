# Methods

A method is a function that belongs to a type.

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-value-receiver` | `go run ./01-fundamentals/08-methods/01-value-receiver` | Calculate an area with a value receiver |
| `02-pointer-receiver` | `go run ./01-fundamentals/08-methods/02-pointer-receiver` | Change a field with a pointer receiver |

## Things to remember

- A value receiver works on a copy, so changes inside the method are lost.
- A pointer receiver changes the original value.
- You can call a pointer method on a normal variable (`rect.ChangeLength(20)`). Go takes the address for you.
