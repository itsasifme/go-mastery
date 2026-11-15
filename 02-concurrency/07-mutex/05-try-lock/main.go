// Try lock: skip a job that's already running, instead of waiting for it.
//
// Several admins click "sync products" at once. TryLock lets the first one
// run the sync. The others get false right away and skip. Once the sync is
// done, the lock is free again, so a later click starts a new one.
//
// If TryLock returns false, you don't hold the lock, so don't call Unlock.
//
// Try: use s.mu.Lock() instead. Every click waits and runs its own sync.
//
// Note: this only works inside one server. Across servers, use a database
// lock like PostgreSQL's pg_try_advisory_xact_lock. If the other callers
// need the result, use singleflight instead (see singleflight/01-basic).
//
// Real use:
//   - a "sync now" or "export CSV" handler that answers 409 Conflict on a
//     double click.
//   - a time.Ticker job that skips a tick while the last flush is running.
//   - a cache refresh where one request reloads the data and the others
//     return the old data right away.
//
// Run: go run ./02-concurrency/07-mutex/05-try-lock
package main

import (
	"fmt"
	"sync"
	"time"
)

// syncer runs the product sync, one at a time.
type syncer struct {
	mu   sync.Mutex // held while a sync is running; also guards runs
	runs int
}

// SyncNow runs a sync and returns true, or returns false right away if a
// sync is already running.
func (s *syncer) SyncNow() bool {
	if !s.mu.TryLock() {
		return false // someone else holds the lock, so don't Unlock
	}
	defer s.mu.Unlock()

	time.Sleep(syncTime) // simulate the sync
	s.runs++
	return true
}

// Runs returns how many syncs have finished.
func (s *syncer) Runs() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.runs
}

// syncTime is how long one sync takes. It's long enough that every click
// in the first group arrives while the sync is still running.
const syncTime = 200 * time.Millisecond

func main() {
	const adminCount = 5

	s := &syncer{}

	// Everyone clicks at once: one sync runs, the rest skip.
	var wg sync.WaitGroup
	wg.Add(adminCount)
	for id := range adminCount {
		go click(id+1, s, &wg)
	}
	wg.Wait()

	// One more click after the sync is finished: the lock is free again.
	lateAdmin := adminCount + 1
	wg.Add(1)
	go click(lateAdmin, s, &wg)
	wg.Wait()

	fmt.Println("syncs run:", s.Runs())
}

// click is one admin pressing "sync products".
func click(admin int, s *syncer, wg *sync.WaitGroup) {
	defer wg.Done()

	if !s.SyncNow() {
		fmt.Println("admin", admin, "sync already running, skipped")
		return
	}
	fmt.Println("admin", admin, "sync finished")
}
