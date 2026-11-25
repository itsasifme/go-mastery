// Function field: a struct field can hold a function.
//
// rect holds a function that computes the area, and s.rect(...) calls it.
//
// Real use:
//   - hooks and callbacks, like http.Server's ConnState or
//     http.Transport's Proxy field.
//   - swapping behaviour in tests: a Clock struct with a Now func field
//     that a test replaces with a fixed time.
//
// Run: go run ./01-fundamentals/07-structs/02-function-field
package main

import "fmt"

// Rectangle is a function type. It takes a length and a breadth and
// returns the area.
type Rectangle func(int, int) int

// Shapes holds a size, and a function that works out the area.
type Shapes struct {
	length  int
	breadth int

	// A field that holds a function
	rect Rectangle
}

func main() {

	// Set the fields by name. rect gets a function that multiplies the sides.
	s := Shapes{
		length:  10,
		breadth: 20,
		rect: func(length int, breadth int) int {
			return length * breadth
		},
	}

	fmt.Println("Area of Rectangle: ", s.rect(s.length, s.breadth))
}
