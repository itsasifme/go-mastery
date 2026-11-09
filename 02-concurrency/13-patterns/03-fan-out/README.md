# Fan-Out

> 📝 **Planned.** No code yet. This is what I plan to add here.

Many goroutines read from the same channel to share the work.

## Examples to add

- `01-basic`: a few workers reading from one channel

## Things to remember

- **Each value goes to only one goroutine.** It shares work. It doesn't copy it to everyone.
- It's often used together with fan-in to collect the results.
