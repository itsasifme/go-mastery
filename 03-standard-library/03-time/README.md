# time

> 📝 **Planned.** No code yet. This is what I plan to add here.

Work with dates, durations, time zones, timers, and tickers.

## Examples to add

- `01-durations`: add and compare durations, and use `time.Since`
- `02-formatting`: format dates using Go's reference time `2006-01-02 15:04:05`
- `03-timezones`: load a time zone and convert between zones
- `04-timers-tickers`: `time.NewTimer` and `time.NewTicker`

## Things to remember

- **Go formats dates with a sample date, not with `YYYY-MM-DD`.**
- **Compare times with `Equal`, not `==`.**
- Always `Stop` a ticker when you're done with it.

## Go version notes

- Timers work better since Go 1.23. Unused timers get cleaned up, and `Stop` and `Reset` don't leave old values behind.
