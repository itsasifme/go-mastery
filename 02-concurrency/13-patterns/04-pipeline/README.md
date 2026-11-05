# Pipeline

> 📝 **Planned.** No code yet. This is what I plan to add here.

A chain of steps connected by channels. Each step takes values, changes them, and passes them on.

## Examples to add

- `01-basic`: generate → square → print
- `02-cancel`: stop every step early with a context

## Things to remember

- Each step closes its own output channel.
- If a later step stops early, the earlier steps get stuck unless they also watch `ctx.Done()`.
