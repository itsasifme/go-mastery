// Value receiver: a method that works on a copy of the struct.
//
// Area only reads length and breadth, so a value receiver is enough.
//
// Real use:
//   - read-only methods on small types, like Money.String() or
//     Point.Distance().
//   - methods on values that never change, like most of time.Time's
//     methods.
//
// Run: go run ./01-fundamentals/08-methods/01-value-receiver
package main

import "fmt"

// Taking a struct
type Rectangle struct {
	length, breadth int
}

func (rect Rectangle) Area() int {
	return rect.length * rect.breadth
}

func main() {
	rect := Rectangle{10, 20}

	fmt.Println("Length and Width are:", rect)
	fmt.Println("Area of Rectangle: ", rect.Area())
}
