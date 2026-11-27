# Data Types

Go's basic types: numbers, booleans, and strings.

## Why

**Problem:** every value has a type, and each type has limits. An integer can only get so big. A float can't store most decimals exactly. A string isn't always one byte per character. If you don't know the limits, you get wrong totals and broken text, with no error.

**Solution:** know what each basic type can hold, and convert between types on purpose. Go helps here: it never converts number types for you, so every conversion is visible in the code.

## API

| Syntax | Kind | What it does | Example |
| --- | --- | --- | --- |
| `int`, `int8` ... `int64` | type | Whole numbers. `int` is 64 bits on most computers | [01-numbers](01-numbers/) |
| `uint`, `uint8` ... `uint64` | type | Whole numbers that can't go below zero | [01-numbers](01-numbers/) |
| `float64`, `float32` | type | Numbers with a fraction, stored approximately | [01-numbers](01-numbers/) |
| `bool` | type | `true` or `false` | [01-numbers](01-numbers/) |
| past the biggest value | rule | An integer wraps around to the smallest value, with no error | [01-numbers](01-numbers/) |
| `string` | type | Text, stored as UTF-8 bytes. Can't be changed | [02-strings-runes](02-strings-runes/) |
| `byte`, `rune` | type | `byte` is `uint8`. `rune` is one Unicode character (`int32`) | [02-strings-runes](02-strings-runes/) |
| `len(s)` | built-in | The number of bytes in a string, not characters | [02-strings-runes](02-strings-runes/) |
| `s[i]` | operator | The byte at position `i` | [02-strings-runes](02-strings-runes/) |
| `for i, r := range s` | statement | Each rune, and the byte where it starts | [02-strings-runes](02-strings-runes/) |
| `utf8.RuneCountInString(s)` | function | The number of runes in a string | [02-strings-runes](02-strings-runes/) |
| `T(v)` | operator | Convert `v` to type `T`, like `float64(n)` | [03-conversions](03-conversions/) |
| `[]rune(s)`, `[]byte(s)`, `string(b)` | operator | Convert a string to runes or bytes, and back | [02-strings-runes](02-strings-runes/), [03-conversions](03-conversions/) |

## Examples

| Example | Run | What it shows |
| --- | --- | --- |
| `01-numbers` | `go run ./01-fundamentals/02-data-types/01-numbers` | Integer sizes, overflow with no error, `0.1 + 0.2`, and money as whole cents |
| `02-strings-runes` | `go run ./01-fundamentals/02-data-types/02-strings-runes` | Problem: cutting a Bangla name by bytes breaks a letter. Fix: cut by runes |
| `03-conversions` | `go run ./01-fundamentals/02-data-types/03-conversions` | Problem: an average rating loses its fraction. Fix: convert to `float64` before you divide |

## Things to remember

- `int` is 32 or 64 bits, depending on the machine.
- **When an integer goes past its biggest value, it jumps to the smallest one, with no error.** A `uint` below zero becomes huge (`01`).
- **Don't compare floats with `==`, and don't store money in a float.** Use whole cents in an `int64` (`01`).
- A string is a list of bytes. `s[i]` gives a byte, but `for range s` gives runes. **`len(s)` counts bytes, not characters** (`02`).
- **Go never converts number types for you.** Even `int` to `int64` needs a conversion (`03`).
- `int(f)` drops the fraction. It doesn't round (`03`).
