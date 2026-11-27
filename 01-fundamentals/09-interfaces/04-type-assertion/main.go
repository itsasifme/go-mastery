// Type assertion: get the real value back out of an interface.
//
// A variable of type any can hold a value of any type. To use it as that
// type again, you write x.(T). It means "I know x really holds a T, give
// it to me as a T". If you're wrong, the program panics.
//
// The two-value form, v, ok := x.(T), never panics. When x doesn't hold a
// T, ok is false and v is T's zero value. A missing map key gives nil, and
// ok is false for that too.
//
// An app keeps its settings in a map[string]any, like a config file loaded
// at start. Someone wrote the retry count as the text "3" instead of the
// number 3.
//
// Problem: retriesOrPanic uses settings["retries"].(int). The value is a
// string, so the program panics with "interface conversion: interface {}
// is string, not int".
//
// Fix: retriesOrDefault uses the two-value form. When the value isn't an
// int, it says so and uses a default. Use the "fix:" line instead of the
// line above it.
//
// Note: T can also be an interface. w.(io.Closer) asks "does w also have
// a Close method?", whatever type it really is.
//
// Real use:
//   - getting your type back from a container that stores any, like
//     sync.Pool's Get or a value from ctx.Value(key).
//   - checking for an extra method: if f, ok := w.(http.Flusher); ok
//     { f.Flush() } sends streamed data right away.
//   - reading a field from JSON decoded into map[string]any.
//
// Run: go run ./01-fundamentals/09-interfaces/04-type-assertion
package main

import "fmt"

// defaultRetries is used when the retries setting isn't an int.
const defaultRetries = 3

func main() {
	// Values of different types. "retries" was written as text by mistake.
	settings := map[string]any{
		"appName": "shop",
		"debug":   true,
		"retries": "3",
	}

	name := settings["appName"].(string) // it is a string, so this works
	fmt.Println("app:", name)

	debug, ok := settings["debug"].(bool)
	fmt.Println("debug:", debug, "ok:", ok)

	port, ok := settings["port"].(int) // missing key: nil, so ok is false
	fmt.Println("port:", port, "ok:", ok)

	retries := retriesOrPanic(settings) // problem: panics, because "retries" is a string
	// retries := retriesOrDefault(settings) // fix: checks the type and falls back to a default
	fmt.Println("retries:", retries)
}

// retriesOrPanic assumes "retries" holds an int. It panics if it doesn't.
func retriesOrPanic(settings map[string]any) int {
	return settings["retries"].(int)
}

// retriesOrDefault checks the type first, and returns defaultRetries when
// "retries" isn't an int.
func retriesOrDefault(settings map[string]any) int {
	retries, ok := settings["retries"].(int)
	if !ok {
		fmt.Printf("retries is a %T, not an int. Using %d\n", settings["retries"], defaultRetries)
		return defaultRetries
	}
	return retries
}
