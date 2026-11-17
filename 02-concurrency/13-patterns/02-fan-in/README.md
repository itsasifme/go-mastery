# Fan-In

Merge many channels into one.

## Why

**Problem:** values come from several sources, each on its own channel. One loop can't `range` over several channels, and a `select` needs one case per channel, written out by hand.

**Solution:** a `merge` function. It starts one goroutine per input that copies values into a shared output channel, and closes the output after every input is done. The caller reads one channel.

## API

| Idiom | What it does | Example |
| --- | --- | --- |
| `func merge(inputs ...<-chan T) <-chan T` | Takes any number of inputs, returns one output | [01-merge](01-merge/) |
| `go forward(in, out, &wg)` | One goroutine per input copies values to `out` | [01-merge](01-merge/) |
| `wg.Wait(); close(out)` | Close the output after every input is done | [01-merge](01-merge/) |
| `select { case out <- v: case <-ctx.Done(): return }` | A forward that stops when the reader leaves | [02-merge-context](02-merge-context/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-merge` | `go run ./02-concurrency/13-patterns/02-fan-in/01-merge` | Logs from three services merged into one stream |
| `02-merge-context` | `go run ./02-concurrency/13-patterns/02-fan-in/02-merge-context` | Stop reading at the first CPU alert, and prove no goroutine leaked |

## Things to remember

- Use one goroutine per input channel. **Close the output when all of them are done.**
- Values from different inputs can come out in any order. Values from one input keep their order.
- **If the reader can stop early, pass a `ctx` to merge.** Otherwise the forward goroutines block forever (`02`).
