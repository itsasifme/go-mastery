# strings

> 📝 **Planned.** No code yet. This is what I plan to add here.

Search, split, and build text with the `strings` package.

## Examples to add

- `01-search-split`: `Contains`, `Index`, `Split`, `Fields`, `Cut`
- `02-builder`: build a long string fast with `strings.Builder`
- `03-replace`: `ReplaceAll` and `strings.NewReplacer`

## Things to remember

- `s += x` in a loop is slow. Use `strings.Builder`.
- For non-English text, a byte index isn't the same as a character position.

## Go version notes

- `strings.Cut` came in Go 1.18. `CutPrefix` and `CutSuffix` came in Go 1.20.
