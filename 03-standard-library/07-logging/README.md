# Logging

> 📝 **Planned.** No code yet. This is what I plan to add here.

The simple `log` package, and structured logs with `log/slog`.

## Examples to add

- `01-log`: `log.Println`, flags, and prefixes
- `02-slog`: text and JSON output with `slog`
- `03-slog-fields`: levels and extra fields

## Things to remember

- `log.Fatal` exits the program right away, so deferred code doesn't run.
- Logs with key and value pairs are easier to search.

## Go version notes

- `log/slog` came in Go 1.21.
