# singleflight

> 📝 **Planned.** No code yet. This is what I plan to add here.

`singleflight` makes many calls for the same key run only once. All callers get the same result.

## Examples to add

- `01-basic`: many goroutines ask for the same key, and the function runs once
- `02-cache`: protect a slow database when a cache entry expires

## Things to remember

- Callers share the error too, not just the result.
- `Forget` removes a key, so the next call runs again.

## Package note

It comes from `golang.org/x/sync`, not the standard library.
