// Program to access the individual elements of struct

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
