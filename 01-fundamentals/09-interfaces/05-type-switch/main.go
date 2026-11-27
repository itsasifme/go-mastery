// Type switch: handle each type an interface can hold in its own case.
//
// switch v := x.(type) checks x against several types in one go. Inside
// each case, v already has that case's type, so you can use it right away
// with no extra assertion. case nil catches an empty interface, and
// default catches every type you didn't list.
//
// A settings page shows each setting. Text, numbers and on/off switches
// each need a different look.
//
// Note: when one case lists two types, like case int, int64, v stays any,
// because Go can't know which of the two it is.
//
// Note: compare with 04-type-assertion. One assertion checks one type. A
// type switch is shorter once there are two or more.
//
// Real use:
//   - json.Unmarshal into an any gives map[string]any, []any, string,
//     float64, bool or nil. A type switch walks through it.
//   - a logger like log/slog formats each field by its type.
//   - a message handler that does something different for each kind of
//     message it receives.
//
// Run: go run ./01-fundamentals/09-interfaces/05-type-switch
package main

import "fmt"

// setting is one row on the settings page. value can be any type.
type setting struct {
	name  string
	value any
}

// percent turns a fraction like 0.15 into 15.
const percent = 100

func main() {
	settings := []setting{
		{name: "appName", value: "shop"},
		{name: "debug", value: true},
		{name: "retries", value: 3},
		{name: "discount", value: 0.15},
		{name: "maintenance", value: nil},
		{name: "adminIDs", value: []int{1, 7}},
	}

	for _, s := range settings {
		fmt.Printf("%-12s %s\n", s.name, describe(s.value))
	}
}

// describe turns one setting value into text for the page.
func describe(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v) // v is a string here
	case bool:
		if v { // v is a bool, so it works as a condition
			return "on"
		}
		return "off"
	case int:
		return fmt.Sprintf("%d times", v)
	case float64:
		return fmt.Sprintf("%.0f%%", v*percent) // v is a float64, so math works
	case nil:
		return "not set"
	default:
		return fmt.Sprintf("can't show a %T", v) // v is still any here
	}
}
