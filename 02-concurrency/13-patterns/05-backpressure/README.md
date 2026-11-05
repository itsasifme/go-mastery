# Backpressure

> 📝 **Planned.** No code yet. This is what I plan to add here.

Ways to stop a fast producer from flooding a slow consumer.

## Examples to add

- `01-blocking`: a small buffer that makes the producer wait
- `02-drop`: `select` with `default` to drop work when the buffer is full

## Things to remember

- A queue with no limit just grows until the program runs out of memory.
- Choose on purpose: wait, drop, or reject.
