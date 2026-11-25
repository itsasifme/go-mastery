// Struct fields: group related values and read or change them by name.
//
// rect.length reads a field, and rect.breadth = 15 changes one.
//
// Real use:
//   - modelling data: a User, an Order, a Config loaded at startup.
//   - request and response bodies in an HTTP API.
//
// Run: go run ./01-fundamentals/07-structs/01-basic-fields
package main

import "fmt"

func main() {

	// Declare a struct
	type Rectangle struct {
		length  int
		breadth int
	}

	// Make a Rectangle. The values go in field order: length, then breadth.
	rect := Rectangle{22, 12}

	// Read the length field
	fmt.Println("Length:", rect.length)

	rect.breadth = 15

	// Read the breadth field
	fmt.Println("Breadth:", rect.breadth)

	area := rect.length * rect.breadth
	fmt.Println("Area:", area)
	fmt.Println("Rect:", rect)
}
