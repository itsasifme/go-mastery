// Parallel fetch: load a page's data from several services at once.
//
// A dashboard needs the user's profile, orders and settings. Each one comes
// from a different service and takes some time. Calling them one by one
// adds up all the waits. Here each call runs in its own goroutine, and
// wg.Wait blocks until all of them are back. The total time is close to the
// slowest call, not the sum.
//
// Each goroutine writes to its own field of dashboard, so they never touch
// the same memory and no lock is needed. main only reads after wg.Wait.
//
// Try: make settingsTime the longest and see the total time follow it.
//
// Real use:
//   - an API handler that calls userSvc.Get, orderSvc.List and
//     settingsSvc.Get at once, then writes one JSON response.
//   - a product page that runs several db.QueryContext calls in parallel
//     (details, reviews, stock) before rendering the template.
//
// Run: go run ./02-concurrency/05-waitgroup/02-parallel-fetch
package main

import (
	"fmt"
	"sync"
	"time"
)

// dashboard holds everything the page needs. Each field is filled by a
// different goroutine.
type dashboard struct {
	profile  string
	orders   []string
	settings string
}

// How long each fake service takes to answer.
const (
	profileTime  = 100 * time.Millisecond
	ordersTime   = 300 * time.Millisecond // the slowest, so it sets the total time
	settingsTime = 200 * time.Millisecond
)

func main() {
	const (
		callCount = 3 // must match the number of "go" lines below
		roundTo   = 10 * time.Millisecond
	)

	start := time.Now()

	var d dashboard
	var wg sync.WaitGroup
	wg.Add(callCount)
	go fetchProfile(&wg, &d.profile)
	go fetchOrders(&wg, &d.orders)
	go fetchSettings(&wg, &d.settings)
	wg.Wait() // after this line, every field is filled

	fmt.Println("profile: ", d.profile)
	fmt.Println("orders:  ", d.orders)
	fmt.Println("settings:", d.settings)
	fmt.Println("took", time.Since(start).Round(roundTo), "instead of", profileTime+ordersTime+settingsTime)
}

// fetchProfile pretends to call the user service.
func fetchProfile(wg *sync.WaitGroup, out *string) {
	defer wg.Done()
	time.Sleep(profileTime)
	*out = "Asif, Dhaka"
}

// fetchOrders pretends to call the order service.
func fetchOrders(wg *sync.WaitGroup, out *[]string) {
	defer wg.Done()
	time.Sleep(ordersTime)
	*out = []string{"order-101", "order-102"}
}

// fetchSettings pretends to call the settings service.
func fetchSettings(wg *sync.WaitGroup, out *string) {
	defer wg.Done()
	time.Sleep(settingsTime)
	*out = "dark mode, English"
}
