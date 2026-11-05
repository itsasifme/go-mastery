# Select

> 📝 **Planned.** No code yet. This is what I plan to add here.

`select` waits on several channels at once and runs the first one that's ready.

## Examples to add

- `01-basic`: read from whichever channel is ready first
- `02-default`: send or receive without waiting
- `03-timeout`: give up after some time with `time.After`
- `04-nil-channels`: turn off a case by setting its channel to `nil`

## Things to remember

- If more than one case is ready, Go picks one at random.
- `select {}` waits forever.
- A closed channel is always ready. Set it to `nil` to stop picking it.
