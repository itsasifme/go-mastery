# Context Timeout

> 📝 **Planned.** No code yet. This is what I plan to add here.

Stop work after some time with `context.WithTimeout`.

## Examples to add

- `01-basic`: a slow task that hits the timeout
- `02-result-or-timeout`: use `select` to wait for a result or the timeout

## Things to remember

- Still call `cancel`, even when the timeout already fired.
- Make sure a slow worker can still finish and exit after the timeout. Otherwise it leaks.
