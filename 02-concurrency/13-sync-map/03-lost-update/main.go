// Lost update: Load, then Store on a sync.Map, loses counts.
//
// A site counts views per page. Each request calls Load, adds 1, then
// calls Store. Each call is safe on its own, but they're separate steps.
// Two requests can Load the same count, and one of the updates is lost.
// It's mutex/01-race again, even though the map itself is safe now.
//
// Problem: the totals are lower than expected, and they change between
// runs. "go run -race" finds nothing, because sync.Map locks every call
// for you. The bug is in the logic, not the memory.
//
// Fix: keep a *atomic.Int64 per page. LoadOrStore adds it the first time
// a page is seen. After that, every view is one atomic Add (see
// atomic/01-counter). The map only changes when a new page shows up. Use
// the "fix:" line instead of the line above it.
//
// Note: new(atomic.Int64) is made on every call, even when the page
// already has a counter. In hot code, call Load first and use
// LoadOrStore only when the page is missing.
//
// Real use:
//   - per-page, per-endpoint or per-customer counters for metrics or
//     usage limits.
//
// Run: go run ./02-concurrency/13-sync-map/03-lost-update
// Race check: go run -race ./02-concurrency/13-sync-map/03-lost-update
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// viewsPerWorker is how many views each worker counts on each page.
const viewsPerWorker = 100

func main() {
	const workerCount = 100

	pages := []string{"/", "/pricing", "/blog"}
	var views sync.Map // page -> its view count

	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(func() { countViews(&views, pages) })
	}
	wg.Wait()

	fmt.Println("expected", workerCount*viewsPerWorker, "per page")
	for _, page := range pages {
		v, _ := views.Load(page)
		fmt.Printf("%-9s %d\n", page, viewCount(v))
	}
}

// countViews counts viewsPerWorker views on every page.
func countViews(views *sync.Map, pages []string) {
	for range viewsPerWorker {
		for _, page := range pages {
			addViewUnsafe(views, page) // problem: Load and Store are two separate steps
			// addView(views, page) // fix: one atomic counter per page
		}
	}
}

// addViewUnsafe reads the count, adds 1, and stores it back. Another
// request can store in between, and its update is lost.
func addViewUnsafe(views *sync.Map, page string) {
	v, _ := views.Load(page)
	n, _ := v.(int) // 0 when the page has no count yet
	views.Store(page, n+1)
}

// addView adds 1 to the page's atomic counter. If the page is new, it
// adds the counter first.
func addView(views *sync.Map, page string) {
	v, _ := views.LoadOrStore(page, new(atomic.Int64))
	v.(*atomic.Int64).Add(1)
}

// viewCount reads a stored count: an int from addViewUnsafe, or a
// *atomic.Int64 from addView.
func viewCount(v any) int64 {
	switch n := v.(type) {
	case int:
		return int64(n)
	case *atomic.Int64:
		return n.Load()
	}
	return 0
}
