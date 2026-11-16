// Nested RLock: a method that holds RLock calls another method that takes
// RLock again.
//
// A shop keeps feature flags in memory. CheckoutPage holds RLock while it
// builds the page, and calls IsEnabled, which takes RLock again. Readers
// share the lock, so this usually works. But a refresh job calls Lock in
// between. A waiting writer blocks new readers, so the second RLock waits
// for the writer, and the writer waits for our first RLock.
//
// Problem: the page and the refresh wait for each other forever. Nothing
// else is running, so Go stops with "fatal error: all goroutines are
// asleep - deadlock!". In a real server, requests just hang.
//
// Fix: same as mutex/04-nested-lock. Use a helper that doesn't lock
// (isEnabledLocked) when you already hold the lock. Comment out the line
// marked "problem", then uncomment the "fix:" line.
//
// Try: set refreshDelay to 150ms, so the refresh comes after the page is
// built. The problem line works fine. That's why tests rarely catch this
// bug: it only happens when a write lands in the small gap between the
// two RLocks.
//
// Real use:
//   - a config, flags or permissions type where one method reads several
//     values and calls other public methods that also lock.
//
// Run: go run ./02-concurrency/08-rwmutex/04-nested-rlock
package main

import (
	"fmt"
	"sync"
	"time"
)

// featureFlags holds which features are on. A background job refreshes it.
type featureFlags struct {
	mu sync.RWMutex // guards on
	on map[string]bool
}

// IsEnabled reports whether a feature is on.
func (f *featureFlags) IsEnabled(name string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.isEnabledLocked(name)
}

// isEnabledLocked does the work of IsEnabled. The caller must hold f.mu.
func (f *featureFlags) isEnabledLocked(name string) bool {
	return f.on[name]
}

// CheckoutPage builds the checkout page from the flags, under one RLock so
// every flag comes from the same version.
func (f *featureFlags) CheckoutPage() string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	time.Sleep(buildTime) // build the page header and cart with the flags

	newCheckout := f.IsEnabled("new-checkout") // problem: takes RLock again while we hold it
	// newCheckout := f.isEnabledLocked("new-checkout") // fix: we already hold the lock
	if newCheckout {
		return "new checkout"
	}
	return "old checkout"
}

// Set turns a feature on or off.
func (f *featureFlags) Set(name string, on bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.on[name] = on
}

const (
	buildTime    = 100 * time.Millisecond // how long building the page takes
	refreshDelay = 50 * time.Millisecond  // when the refresh job runs
)

func main() {
	flags := &featureFlags{on: map[string]bool{"new-checkout": true}}

	var wg sync.WaitGroup
	wg.Add(2) // the page request and the refresh job
	go showCheckout(flags, &wg)
	go refreshFlags(flags, &wg)
	wg.Wait() // a stuck page and refresh show up here as a deadlock
}

// showCheckout is one request for the checkout page.
func showCheckout(f *featureFlags, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("page:", f.CheckoutPage())
}

// refreshFlags is the background job that loads new flags, after
// refreshDelay.
func refreshFlags(f *featureFlags, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(refreshDelay)
	f.Set("new-checkout", false)
	fmt.Println("refresh: flags updated")
}
