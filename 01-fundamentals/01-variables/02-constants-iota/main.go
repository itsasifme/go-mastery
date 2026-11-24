// Constants and iota: values that never change, numbered for you.
//
// const makes a value that's fixed when the program is compiled. Trying to
// change it is a compile error.
//
// iota counts the lines in one const block, starting at 0. The order
// statuses below get 0, 1, 2, 3 without writing a number. A line with no
// "= ..." repeats the line above it with the next iota, so the byte sizes
// each shift 10 more bits: 1024, then 1024 * 1024, and so on. "_" skips
// a value you don't want.
//
// Try: add statusRefunded between statusPaid and statusShipped. Every
// status after it moves up by one. That's why values that are saved in a
// database or sent to another service should be written out by hand, like
// http.StatusNotFound = 404, not made with iota.
//
// Real use:
//   - a fixed list of states that only lives inside the program, like
//     order statuses or log levels. time.Weekday (Sunday = 0) uses iota.
//   - size and time limits built from other constants, like
//     maxUploadSize = 10 * MB.
//
// Run: go run ./01-fundamentals/01-variables/02-constants-iota
package main

import "fmt"

// Order statuses. iota starts at 0 and goes up by one each line.
const (
	statusPending   = iota // 0
	statusPaid             // 1
	statusShipped          // 2
	statusDelivered        // 3
)

// Byte sizes. "_" skips iota's 0, so KB is 1 << 10.
const (
	_  = iota
	KB = 1 << (10 * iota) // 1 << 10
	MB                    // 1 << 20: the same formula with the next iota
	GB                    // 1 << 30
)

// maxUploadSize is made from another constant. Go works it out when it
// compiles.
const maxUploadSize = 10 * MB

func main() {
	for status := range statusDelivered + 1 {
		fmt.Println(status, statusName(status))
	}

	fmt.Println("KB:", KB, "MB:", MB, "GB:", GB)

	uploads := []int{2 * MB, 25 * MB} // sizes of two uploaded files
	for _, size := range uploads {
		fmt.Printf("upload of %d MB allowed: %t\n", size/MB, size <= maxUploadSize)
	}

	// statusPaid = 5 // compile error: cannot assign to statusPaid
}

// statusName returns the name shown to the user for a status.
func statusName(status int) string {
	switch status {
	case statusPending:
		return "pending"
	case statusPaid:
		return "paid"
	case statusShipped:
		return "shipped"
	case statusDelivered:
		return "delivered"
	}
	return "unknown"
}
