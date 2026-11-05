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
