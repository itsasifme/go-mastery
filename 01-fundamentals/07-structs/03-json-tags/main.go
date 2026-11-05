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
		"House# 54, Road# 29, Block# D, Section# 12, Pallabi",
		"Dhaka",
	}

	user := User{
		"Asif",
		"asif@gmail.com",
		"+8801910-241412",
		address,
	}

	serialized, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	fmt.Println("Serialized User JSON:", string(serialized))
}
