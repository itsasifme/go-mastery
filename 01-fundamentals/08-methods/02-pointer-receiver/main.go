// Pointer receiver: a method that changes the original struct.
//
// ChangeLength has a *Rectangle receiver, so it changes rect itself, not
// a copy. Go takes &rect for you when you call rect.ChangeLength(20).
//
// Real use:
//   - methods that update state, like account.Deposit(amount) or
//     cache.Set(key, value).
//   - any struct that holds a sync.Mutex: copying it breaks the lock, so
//     its methods must use a pointer receiver.
//
// Run: go run ./01-fundamentals/08-methods/02-pointer-receiver
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
