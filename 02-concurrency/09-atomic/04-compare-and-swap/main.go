// Compare and swap: take one from a limited stock, without going below zero.
//
// A shop gives away a few coupons, and many users try to claim one at
// once. claimUnsafe checks that coupons are left, saves the coupon to the
// user's account, then takes one from the stock with Add(-1). Each step is
// atomic, but they're separate steps.
//
// Problem: while the first users are saving, the others pass the same
// check. More users get a coupon than there are coupons, and left goes
// below zero. "go run -race" finds nothing, because every access is
// atomic. With nothing slow between the check and the take, the gap is
// tiny and the bug hides, until real work goes there.
//
// Fix: takeOne reads left. If it's zero, it gives up. If not, it calls
// CompareAndSwap(n, n-1): "set it to n-1, but only if it's still n". If
// another user took one in between, the swap fails and takeOne tries again
// with the new value. claimSafe takes the coupon first, then saves it.
// Use the "fix:" line instead of the line above it.
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
// Race check: go run -race ./02-concurrency/09-atomic/04-compare-and-swap
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// saveTime is how long saving a coupon to a user's account takes.
const saveTime = 10 * time.Millisecond

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
		go claimUnsafe(id+1, &left, &claimed, &wg) // problem: check and take are separate steps
		// go claimSafe(id+1, &left, &claimed, &wg) // fix: take first with CompareAndSwap, then save
	}
	wg.Wait()

	fmt.Println("claimed:", claimed.Load(), "of", couponCount)
	fmt.Println("left:   ", left.Load())
}

// claimUnsafe checks for a coupon, saves it to the user's account, then
// takes it from the stock. Other users can pass the same check while the
// coupon is being saved.
func claimUnsafe(user int, left, claimed *atomic.Int64, wg *sync.WaitGroup) {
	defer wg.Done()

	if left.Load() == 0 {
		return
	}
	time.Sleep(saveTime) // save the coupon to the user's account
	left.Add(-1)
	claimed.Add(1)
	fmt.Println("user", user, "got a coupon")
}

// claimSafe takes a coupon from the stock first, then saves it to the
// user's account.
func claimSafe(user int, left, claimed *atomic.Int64, wg *sync.WaitGroup) {
	defer wg.Done()

	if !takeOne(left) {
		return
	}
	time.Sleep(saveTime) // save the coupon to the user's account
	claimed.Add(1)
	fmt.Println("user", user, "got a coupon")
}

// takeOne takes one from left if it's above zero. It reports whether it
// got one.
func takeOne(left *atomic.Int64) bool {
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
