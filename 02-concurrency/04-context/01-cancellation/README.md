# Context Cancellation

> 📝 **Planned.** No code yet. This is what I plan to add here.

Stop work by hand with `context.WithCancel`.

## Examples to add

- `01-basic`: a goroutine that stops when `ctx.Done()` closes
- `02-stop-workers`: one `cancel` call stops many workers
- `03-cancel-cause`: say why it was cancelled with `WithCancelCause`

## Things to remember

- Cancelling a parent also cancels all its children, but not the other way around.
- A goroutine has to check `ctx.Done()` itself. Go won't stop it for you.
