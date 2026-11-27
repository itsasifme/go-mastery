// Declaration: var, :=, and zero values.
//
// A shopping cart starts empty. "var name type" makes a variable set to
// its type's zero value: 0 for numbers, "" for strings, false for bools,
// and nil for pointers, slices and maps. A nil slice still works: len is 0
// and append adds to it.
//
// ":=" declares a variable and sets it in one step. Go works out the type
// from the value, so "price := 1250" is an int.
//
// Use var when the zero value is the right start, or outside functions.
// Use := inside functions when you already have the value.
//
// Try: move "store := ..." above func main. It won't compile: := only
// works inside functions. Package-level variables need var.
//
// Real use:
//   - "var total int", then add to it in a loop. No need for "= 0".
//   - "var mu sync.Mutex" or "var buf bytes.Buffer": their zero values are
//     ready to use.
//   - "user, err := loadUser(id)" in a handler: the types come from the
//     function.
//
// Run: go run ./01-fundamentals/01-variables/01-declaration
package main

import "fmt"

// currency is set once for the whole program. Package-level variables
// need var.
var currency = "BDT"

func main() {
	// var: every variable starts at its zero value.
	var itemCount int        // 0
	var coupon string        // ""
	var giftWrap bool        // false
	var discount *int        // nil: no discount
	var items []string       // nil, but len and append still work
	var stock map[string]int // nil: reading works, writing would panic

	fmt.Printf("count=%d coupon=%q giftWrap=%t discount=%v items=%v len=%d stock=%v\n",
		itemCount, coupon, giftWrap, discount, items, len(items), stock)

	// :=: declare and set in one step. Go picks the type from the value.
	store := "Dhaka Books"
	price := 1250             // int
	weight := 0.4             // float64
	name, qty := "Go book", 2 // two at once
	fmt.Printf("store=%T price=%T weight=%T name=%T qty=%T\n", store, price, weight, name, qty)

	// = changes a variable that already exists.
	items = append(items, name)
	itemCount = len(items)
	total := price * qty
	fmt.Println(store+":", itemCount, "item,", items, "total", total, currency)
}
