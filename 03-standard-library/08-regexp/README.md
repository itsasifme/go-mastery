# regexp

> 📝 **Planned.** No code yet. This is what I plan to add here.

Match and replace text with regular expressions.

## Examples to add

- `01-match`: `MatchString`, `FindString`, `FindAllString`
- `02-groups`: capture groups and named groups
- `03-replace`: `ReplaceAllString`

## Things to remember

- **Compile the pattern once with `regexp.MustCompile`, not inside a loop.**
- Go's regexp is always fast, but it has no back-references and no lookahead.
- If plain text search is enough, use the `strings` package instead.
