// Interfaces: one function works with any type that has the right methods.
//
// Rectangle and Triangle both have area(), so both satisfy Shape.
// calculate takes a Shape and doesn't care which one it gets. There's no
// "implements" keyword: having the method is enough.
//
// Real use:
//   - io.Reader: the same code reads from a file, an HTTP body, or a
//     string.
//   - swapping a real database for a fake one in tests, through a
//     UserStore interface.
//   - payment providers (Stripe, PayPal) behind one PaymentGateway
//     interface.
//
// Run: go run ./01-fundamentals/09-interfaces/01-shapes
package main

import "fmt"

// Interface
type Shape interface {
	area() float32
}

// Rectangle struct implements the interface
type Rectangle struct {
	length, breadth float32
}

// Rectangle provides implementation for area()
func (r Rectangle) area() float32 {
	return r.length * r.breadth
}

// Triangle struct implements the interface
type Triangle struct {
	base, height float32
}

// Triangle provides implementation for area()
func (t Triangle) area() float32 {
	return 0.5 * t.base * t.height
}

// Access method of the interface
func calculate(s Shape) float32 {
	return s.area()
}

// Main function
func main() {

	// Assigns value to struct members
	r := Rectangle{7, 4}
	t := Triangle{8, 12}

	fmt.Println("Area of Rectangle (Access via Rectangle):", r.area())

	// Call calculate() with struct variable r
	rectangleArea := calculate(r)
	fmt.Println("Area of Rectangle (Access via Shape):", rectangleArea)

	fmt.Println("Area of Triangle (Access via Triangle):", t.area())

	// Call calculate() with struct variable t
	triangleArea := calculate(t)
	fmt.Println("Area of Triangle (Access via Shape):", triangleArea)
}
