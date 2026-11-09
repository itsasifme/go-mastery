# io

> 📝 **Planned.** No code yet. This is what I plan to add here.

`io.Reader` and `io.Writer` are used all over Go. This topic shows how they work and fit together.

## Examples to add

- `01-reader-writer`: write a small `io.Reader` yourself
- `02-copy`: `io.Copy`, `io.TeeReader`, `io.MultiWriter`
- `03-bufio`: read line by line with `bufio.Scanner`
- `04-pipe`: connect a writer and a reader with `io.Pipe`

## Things to remember

- **`Read` can return some data and an error at the same time.** Use the data first.
- `bufio.Scanner` fails on lines longer than 64 KB by default.
- **Always call `Flush` on a `bufio.Writer`.**
