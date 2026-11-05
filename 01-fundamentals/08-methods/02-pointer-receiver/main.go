package main

import "fmt"

// Taking a struct
type Rectangle struct {
	length, breadth int
}

func (rect *Rectangle) ChangeLength(length int) {
	rect.length = length
}

func (rect Rectangle) Area() int {
	return rect.length * rect.breadth
}

func main() {
	rect := Rectangle{10, 20}
	rect.ChangeLength(20)

	fmt.Println("Length and Width are:", rect)
	fmt.Println("Area of Rectangle: ", rect.Area())
}
