// JSON tags: control the key names when a struct becomes JSON.
//
// json:"name" turns the Name field into a "name" key. Address is
// embedded, so its fields appear at the top level of the JSON, not nested
// under "Address".
//
// Real use:
//   - API request and response bodies, where the JSON uses lowercase or
//     snake_case keys.
//   - config files, and messages sent through queues like Kafka.
//   - json:"-" to hide a field like a password, and json:",omitempty" to
//     skip empty values.
//
// Run: go run ./01-fundamentals/07-structs/03-json-tags
package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Address
}

func main() {
	address := Address{
		"House 1, Road 2, Example Area",
		"Dhaka",
	}

	user := User{
		"Asif",
		"asif@example.com",
		"+880 1700-000000",
		address,
	}

	serialized, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Serialized User JSON:", string(serialized))
}
