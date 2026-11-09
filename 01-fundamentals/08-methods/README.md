# Methods

A method is a function that belongs to a type.

## Why

**Problem:** plain functions like `area(rect)` live apart from the data. With many types, you end up with `areaRect`, `areaTriangle`, and so on.

**Solution:** a method belongs to a type: `rect.Area()`. Methods are also how a type fits an interface ([09-interfaces](../09-interfaces/)).

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `func (r T) M()` | Value receiver: works on a copy | [01-value-receiver](01-value-receiver/) |
| `func (r *T) M()` | Pointer receiver: changes the original | [02-pointer-receiver](02-pointer-receiver/) |
| `v.M()` with a pointer receiver | Go takes `&v` for you | [02-pointer-receiver](02-pointer-receiver/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-value-receiver` | `go run ./01-fundamentals/08-methods/01-value-receiver` | Calculate an area with a value receiver |
| `02-pointer-receiver` | `go run ./01-fundamentals/08-methods/02-pointer-receiver` | Change a field with a pointer receiver |

## Things to remember

- A value receiver works on a copy, so **changes inside the method are lost**.
- A pointer receiver changes the original value.
- You can call a pointer method on a normal variable (`rect.ChangeLength(20)`). Go takes the address for you.
