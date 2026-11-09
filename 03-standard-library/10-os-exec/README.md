# os/exec

> 📝 **Planned.** No code yet. This is what I plan to add here.

Run other programs from Go.

## Examples to add

- `01-output`: run a command and read its output
- `02-exit-code`: get the exit code when a command fails
- `03-timeout`: kill a slow command with `exec.CommandContext`

## Things to remember

- **There's no shell**, so things like `*` or quotes aren't processed.
- Commands are different on Windows and Linux. `go version` works on both.

## Go version notes

- Since Go 1.19, Go won't run a program from the current folder by name alone.
