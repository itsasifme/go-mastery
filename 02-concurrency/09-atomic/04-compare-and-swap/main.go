// Compare and swap: take one from a limited stock, without going below zero.
//
// A shop gives away a few coupons, and many users try to claim one at
// once. claim reads how many are left. If none, it gives up. If some, it
// calls CompareAndSwap(n, n-1): "set it to n-1, but only if it's still n".
// If another user claimed in between, the swap fails and claim tries again
// with the new value.
//
// A plain left.Add(-1) followed by a check isn't enough. Several users can
// subtract at once, and left goes below zero even though they get turned
// away.
//
// Note: compare with mutex/03-check-then-act. There, a mutex kept the
// check and the change together. Here, CompareAndSwap does it for one value.
// This works inside one server. With many servers, let the database do it:
// UPDATE coupons SET left = left - 1 WHERE id = $1 AND left > 0.
//
// Real use:
//   - limited coupons, free trial slots, or "first 100 users" offers.
//   - a simple rate limiter or connection counter with a maximum.
//
// Run: go run ./02-concurrency/09-atomic/04-compare-and-swap
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	const (
		couponCount = 5
		userCount   = 50
	)

	var left atomic.Int64
	left.Store(couponCount)
	var claimed atomic.Int64

	var wg sync.WaitGroup
	wg.Add(userCount)
	for id := range userCount {
		go tryClaim(id+1, &left, &claimed, &wg)
	}
	wg.Wait()

	fmt.Println("claimed:", claimed.Load(), "of", couponCount)
	fmt.Println("left:   ", left.Load())
}

// tryClaim tries to claim one coupon for a user and prints who got one.
func tryClaim(user int, left, claimed *atomic.Int64, wg *sync.WaitGroup) {
	defer wg.Done()

	if claim(left) {
		claimed.Add(1)
		fmt.Println("user", user, "got a coupon")
	}
}

// claim takes one from left if it's above zero. It reports whether it got one.
func claim(left *atomic.Int64) bool {
	for {
		n := left.Load()
		if n == 0 {
			return false
		}
		if left.CompareAndSwap(n, n-1) {
			return true
		}
		// someone else changed left after our Load, so try again
	}
}
