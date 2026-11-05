# Context Deadline

> 📝 **Planned.** No code yet. This is what I plan to add here.

Stop work at a fixed time with `context.WithDeadline`.

## Examples to add

- `01-basic`: a deadline at a set time
- `02-parent-deadline`: a child can make the deadline earlier, but never later

## Things to remember

- `WithTimeout` is the same thing as `WithDeadline` with `time.Now().Add(d)`.
- `ctx.Deadline()` tells you if there's a deadline, so you can skip work that won't finish in time.
