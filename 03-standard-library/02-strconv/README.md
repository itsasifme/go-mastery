# strconv

> 📝 **Planned.** No code yet. This is what I plan to add here.

Convert between strings and numbers.

## Examples to add

- `01-parse`: `Atoi`, `ParseInt`, `ParseFloat`, `ParseBool`
- `02-format`: `Itoa`, `FormatInt`, `FormatFloat`, `Quote`
- `03-errors`: what the error looks like when parsing fails

## Things to remember

- `string(65)` gives `"A"`, not `"65"`. Use `strconv.Itoa`.
- `strconv` is faster and clearer than `fmt.Sprint` for simple conversions.
