// Conversions: Go never converts types for you. You write T(v).
//
// float64(n) turns an int into a float. int(f) turns a float into an int
// and drops the fraction. It doesn't round. Even int to int64 needs a
// conversion. A conversion to a smaller type, like int64 to int32, wraps
// big values with no error.
//
// A product page shows the average of its star ratings.
//
// Problem: averageIntDivide divides two ints, so the fraction is gone
// before float64 sees it. The page shows 4.00 instead of 4.33.
//
// Fix: averageFloatDivide converts both to float64 first, then divides.
// main runs both, so you can compare.
//
// Note: string(65) gives "A", not "65", because it's a rune. "go vet"
// warns about it. To turn a number into text, use strconv.Itoa or
// fmt.Sprint (see the strconv topic).
//
// Real use:
//   - averages, percentages and ratios: convert to float64 before you
//     divide.
//   - math.Round(f) before int(f), when you want the closest whole number.
//   - []byte(s) to write a string to a file or an HTTP body.
//
// Run: go run ./01-fundamentals/02-data-types/03-conversions
package main

import (
	"fmt"
	"math"
)

// ratingCount is how many customers rated the product.
const ratingCount = 3

func main() {
	first, second, third := 5, 4, 4 // each customer's stars
	sum := first + second + third

	fmt.Printf("problem, average: %.2f\n", averageIntDivide(sum, ratingCount))
	fmt.Printf("fix, average:     %.2f\n", averageFloatDivide(sum, ratingCount))

	// float to int drops the fraction. Round first if you want the closest.
	price := 4.99
	fmt.Println("int(4.99):", int(price), " rounded:", int(math.Round(price)))

	// A smaller type wraps big values, with no error.
	var orderID int64 = math.MaxInt32 + 1
	fmt.Println("int64:", orderID, " as int32:", int32(orderID))

	// A string and a []byte convert both ways. Each makes a copy.
	body := []byte("hello")
	fmt.Println("bytes:", body, " back to string:", string(body))

	// fmt.Println(string(65)) // go vet: conversion from int to string yields a string of one rune
}

// averageIntDivide divides first, as ints, so the fraction is lost.
func averageIntDivide(sum, count int) float64 {
	return float64(sum / count)
}

// averageFloatDivide converts to float64 first, so the fraction is kept.
func averageFloatDivide(sum, count int) float64 {
	return float64(sum) / float64(count)
}
