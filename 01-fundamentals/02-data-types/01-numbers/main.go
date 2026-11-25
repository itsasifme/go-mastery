// Numbers: int sizes, overflow, and float precision.
//
// Go has integers of fixed sizes, like int8 (-128 to 127) and int64. Plain
// int is 64 bits on most computers today, and 32 bits on some small ones.
//
// When an integer goes past its biggest value, it wraps around to the
// smallest one. There's no error. An unsigned integer (uint) can't go
// below zero, so 0 - 1 wraps around to its biggest value.
//
// A float64 stores most decimals only approximately. So 0.1 + 0.2 isn't
// exactly 0.3, and == says false. That's why money is usually stored as
// whole cents in an int64.
//
// Try: change stock to an int. Now 0 - 1 is -1, and you can check for it.
//
// Real use:
//   - prices as whole cents in an int64. Payment APIs like Stripe take
//     amounts this way: 1099 means 10.99.
//   - a uint size or count that goes below zero turns into a huge number.
//     That's why Go uses int for len and most counts.
//   - comparing floats with a small tolerance, never with ==.
//
// Run: go run ./01-fundamentals/02-data-types/01-numbers
package main

import (
	"fmt"
	"math"
	"math/bits"
)

const (
	centsPerTaka = 100
	tolerance    = 1e-9 // how close two floats must be to count as equal
)

func main() {
	fmt.Println("int has", bits.UintSize, "bits on this computer")
	fmt.Println("int8 range: ", math.MinInt8, "to", math.MaxInt8)
	fmt.Println("int64 range:", math.MinInt64, "to", math.MaxInt64)

	// Overflow: past the biggest value, an integer wraps to the smallest.
	var retries int8 = math.MaxInt8
	retries++
	fmt.Println("int8 max + 1:", retries)

	// A uint can't go below zero. It wraps to its biggest value.
	var stock uint // 0
	stock--
	fmt.Println("uint 0 - 1:", stock)

	// Floats are close, not exact.
	a, b := 0.1, 0.2
	fmt.Println("0.1 + 0.2 =", a+b)
	fmt.Println("equal to 0.3:", a+b == 0.3)
	fmt.Println("close to 0.3:", math.Abs(a+b-0.3) < tolerance)

	// Money as whole cents in an int64 stays exact.
	var bookCents, penCents, stickerCents int64 = 1099, 250, 1 // 10.99, 2.50 and 0.01
	totalCents := bookCents + penCents + stickerCents
	fmt.Printf("total: %d.%02d Taka\n", totalCents/centsPerTaka, totalCents%centsPerTaka)
}
