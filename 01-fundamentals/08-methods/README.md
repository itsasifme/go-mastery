# Methods

A method is a function that belongs to a type.

## Why

**Problem:** plain functions like `area(rect)` live apart from the data. With many types, you end up with `areaRect`, `areaTriangle`, and so on.

**Solution:** a method belongs to a type: `rect.Area()`. Methods are also how a type fits an interface ([09-interfaces](../09-interfaces/)).

## API

| Syntax | Kind | What it does | Example |
| --- | --- | --- | --- |
| `func (r T) M()` | declaration | Value receiver: works on a copy | [01-value-receiver](01-value-receiver/) |
| `func (r *T) M()` | declaration | Pointer receiver: changes the original | [02-pointer-receiver](02-pointer-receiver/) |
| `v.M()` with a pointer receiver | rule | Go takes `&v` for you | [02-pointer-receiver](02-pointer-receiver/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-value-receiver` | `go run ./01-fundamentals/08-methods/01-value-receiver` | Calculate an area with a value receiver |
| `02-pointer-receiver` | `go run ./01-fundamentals/08-methods/02-pointer-receiver` | Change a field with a pointer receiver |

## Method or function?

Both are fine in Go. Pick by what the code needs.

| Use | When | Example |
| --- | --- | --- |
| Method | It needs the type's fields, or it's an action the type does | `rect.Area()` reads `length` and `breadth`. `conn.Send(msg)` sends on that connection |
| Function | It only uses its arguments | `formatCM(n int) string` turns a number into text |

A method calling a function is normal. It keeps each piece small:

```go
// Describe uses rect's fields, so it's a method.
func (rect Rectangle) Describe() string {
	return formatCM(rect.length) + " x " + formatCM(rect.breadth)
}

// formatCM only needs n, so it's a function. Any type can use it.
func formatCM(n int) string {
	return strconv.Itoa(n) + " cm"
}
```

Don't make something a method just to group it with a type. A method that never uses its receiver makes the reader think it does.

## Things to remember

- A value receiver works on a copy, so **changes inside the method are lost**.
- A pointer receiver changes the original value.
- You can call a pointer method on a normal variable (`rect.ChangeLength(20)`). Go takes the address for you.
- **Make it a method only if it uses the type's data.** If it only uses its arguments, a plain function is better, and a method can call it.
