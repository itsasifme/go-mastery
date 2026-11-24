// Scope and shadowing: := inside a block can hide an outer variable.
//
// A variable exists only inside the { } block where it's declared. A block
// can declare a new variable with the same name as an outer one. Inside
// the block, the new one hides the outer one. That's called shadowing.
//
// checkout validates the cart into err. For express orders it also books a
// courier, inside an if block.
//
// Problem: "err := bookCourier()" declares a new err that only lives
// inside the if. The courier error goes into that new err, and the outer
// err is still nil. checkout returns nil, so the order looks placed even
// though no courier was booked. The program compiles and "go vet" says
// nothing.
//
// Fix: use "=" instead of ":=", so the error goes into the outer err.
// Use the "fix:" line instead of the line above it.
//
// Real use:
//   - "if x, err := f(); err != nil" is fine when err is only needed
//     inside the if. The bug is := on a variable you meant to keep, often
//     err or a result like user or total.
//
// Run: go run ./01-fundamentals/01-variables/03-scope-shadowing
package main

import (
	"errors"
	"fmt"
)

// errNoCourier is what bookCourier returns when every courier is busy.
var errNoCourier = errors.New("no courier free right now")

func main() {
	if err := checkout(true); err != nil {
		fmt.Println("checkout failed:", err)
		return
	}
	fmt.Println("order placed")
}

// checkout validates the cart, and books a courier for express orders. It
// returns the first error.
func checkout(express bool) error {
	err := validateCart() // outer err: nil, the cart is fine

	if express {
		err := bookCourier() // problem: a new err that hides the outer one
		// err = bookCourier() // fix: set the outer err
		fmt.Println("courier booked:", err == nil)
	}

	return err // the outer err
}

// validateCart pretends to check the cart. It always passes.
func validateCart() error {
	return nil
}

// bookCourier pretends to book a courier. Every courier is busy.
func bookCourier() error {
	return errNoCourier
}
