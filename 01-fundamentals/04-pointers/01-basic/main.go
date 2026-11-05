package main

import (
	"fmt"
	"unsafe"
)

func main() {
	// A plain integer value
	value := 10
	fmt.Println(value)

	// Create a integer type pointer that stores the address of value
	// var ptr *int
	// ptr = &value
	// OR
	ptr := &value
	fmt.Println(ptr)
	fmt.Printf("Memory address of ptr: %v: %v\n", ptr, *ptr)

	// Change the pointed-to value; the pointer's address stays unchanged
	*ptr = 20
	fmt.Printf("Memory address of ptr: %v: %v\n", ptr, *ptr)
	// ------------------------------------------------

	// Declare a Person struct
	type Person struct {
		name string
		age  int
	}

	// Instance of the struct Person
	person := Person{"John", 25}
	fmt.Println(person)

	// Create a struct type pointer that stores the address of person
	var ptr2 *Person
	ptr2 = &person
	fmt.Println(ptr2)

	// Print pointer address and dereference the pointer value
	fmt.Printf("Memory address of ptr2: %p: %v\n", ptr2, *ptr2)

	// Change the pointed-to value; the pointer's address stays unchanged
	*ptr2 = Person{"Asif", 30}
	fmt.Printf("Memory address of ptr2: %p: %v\n", ptr2, *ptr2)

	fmt.Printf("Exact memory address (uintptr): %x: %v\n", uintptr(unsafe.Pointer(ptr2)), *ptr2)
}
