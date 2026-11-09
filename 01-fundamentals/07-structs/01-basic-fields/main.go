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

	// Declare instance rect1 and defining the struct
	rect := Rectangle{22, 12}

	// Access the length of the struct
	fmt.Println("Length:", rect.length)

	rect.breadth = 15

	// Access the breadth of the struct
	fmt.Println("Breadth:", rect.breadth)

	area := rect.length * rect.breadth
	fmt.Println("Area:", area)
	fmt.Println("Rect:", rect)
}
