# Concurrency Patterns

Common ways to build things with goroutines and channels.

1. [Worker pool](01-worker-pool/)
2. [Fan-in](02-fan-in/)
3. [Fan-out](03-fan-out/)
4. [Pipeline](04-pipeline/)
5. [Backpressure](05-backpressure/)

## Things to remember

- **Every goroutine needs a clear way to stop.** Otherwise it leaks.
- Only the sender closes a channel, and only after all sending is done.
