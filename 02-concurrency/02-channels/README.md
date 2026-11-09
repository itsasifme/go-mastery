# Channels

A channel lets goroutines send values to each other. It also makes them wait for each other.

1. [Directional channels](01-directional/): how to pass a channel to a function as send-only or receive-only
2. [Closing channels](02-closing/): how to say "no more values", and how to signal many goroutines at once
3. [Unbuffered channels](03-unbuffered/): a send waits for a receiver
4. [Buffered channels](04-buffered/): a send only waits when the buffer is full

Directional channels and closing come first, because every example after them passes channels as `chan<- int` / `<-chan int` and uses `close`.

## Why

**Problem:** goroutines that share variables need locks, and it's easy to get a race. You also need a way to wait for another goroutine's result.

**Solution:** "Don't communicate by sharing memory; share memory by communicating." A channel passes a value from one goroutine to another, and makes them wait for each other.

## API

| Syntax | What it does | Example |
| --- | --- | --- |
| `make(chan T)` | Unbuffered channel | [03-unbuffered](03-unbuffered/) |
| `make(chan T, n)` | Buffered channel | [04-buffered](04-buffered/) |
| `ch <- v` | Send | [03-unbuffered/01-basic](03-unbuffered/01-basic/) |
| `v := <-ch` | Receive | [03-unbuffered/01-basic](03-unbuffered/01-basic/) |
| `v, ok := <-ch` | Receive. `ok == false` once closed and empty | [04-buffered/06-drain-after-close](04-buffered/06-drain-after-close/) |
| `close(ch)` | No more values | [02-closing/01-range](02-closing/01-range/) |
| `for v := range ch` | Read until closed | [02-closing/01-range](02-closing/01-range/) |
| `chan<- T`, `<-chan T` | Send-only, receive-only | [01-directional](01-directional/) |
| `len(ch)`, `cap(ch)` | Values waiting, buffer size | [04-buffered/01-basic](04-buffered/01-basic/) |

## Things to remember

- **Signal channels have the same names in every example:** `stop` tells a goroutine to stop (`main` → goroutine), and `done` says a goroutine finished (goroutine → `main`).
- Both are `chan struct{}`. **A signal carries no data**, and `struct{}` takes no memory.
- Close `stop` to reach every goroutine at once ([02-closing/02-broadcast](02-closing/02-broadcast/)). Close `done` when one goroutine reports, or send `done <- struct{}{}` once per goroutine when several report on the same channel.
