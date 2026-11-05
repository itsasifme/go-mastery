package main

import (
	"fmt"
)

func main() {
	go helloworld()
	// time.Sleep(1 * time.Second)
	goodbye()
}

func helloworld() {
	fmt.Println("Hello World!")
}

func goodbye() {
	fmt.Println("Good Bye!")
}
