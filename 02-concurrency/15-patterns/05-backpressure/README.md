# Backpressure

Ways to stop a fast producer from flooding a slow consumer.

## Why

**Problem:** when values come in faster than they can be handled, they pile up somewhere. With no limit, that pile grows until the program runs out of memory.

**Solution:** put a buffer with a size limit between the two sides. Then decide what happens when it's full: the producer waits, the new value is dropped, or the caller is told "busy".

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `make(chan T, n)` + `ch <- v` | Wait when the buffer is full | [01-blocking](01-blocking/) |
| `select { case ch <- v: default: }` | Drop the value when the buffer is full | [02-drop](02-drop/) |
| `select { case ch <- v: case <-time.After(d): }` | Wait up to `d`, then give up and reject | [03-reject](03-reject/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-blocking` | `go run ./02-concurrency/15-patterns/05-backpressure/01-blocking` | A full log buffer slows the app down to the shipper's speed |
| `02-drop` | `go run ./02-concurrency/15-patterns/05-backpressure/02-drop` | Metric samples are dropped and counted when the exporter can't keep up |
| `03-reject` | `go run ./02-concurrency/15-patterns/05-backpressure/03-reject` | Uploads wait a little for room in the queue, then get a 503 |

## Things to remember

- **A queue with no limit just grows until the program runs out of memory.**
- Choose on purpose: wait (`01`), drop (`02`), or reject (`03`).
- **When you drop, count the drops,** so you know it's happening.
- **When you reject, tell the caller** (503 or 429), so it can retry later.
