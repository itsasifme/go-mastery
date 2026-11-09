# Data Types

> 📝 **Planned.** No code yet. This is what I plan to add here.

Go's basic types: numbers, booleans, and strings.

## Examples to add

- `01-numbers`: int sizes, overflow, float precision
- `02-strings-runes`: bytes vs runes, and looping over a string
- `03-conversions`: converting between types

## Things to remember

- `int` is 32 or 64 bits, depending on the machine.
- **Go never converts number types for you.** Even `int` to `int64` needs a cast.
- **Integer overflow wraps around without any error.**
- A string is a list of bytes. `s[i]` gives a byte, but `for range s` gives runes (characters).
