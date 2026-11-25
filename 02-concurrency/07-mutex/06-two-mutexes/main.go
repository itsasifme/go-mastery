// Two mutexes: Lock only waits for the same mutex, not for the variable.
//
// A mutex doesn't lock a variable or a line of code. It only locks itself.
// mu.Lock() waits while another goroutine holds that same mu. It doesn't
// matter where that goroutine is: the same function, another method or
// another file. Code that locks a different mutex, or no mutex, doesn't
// wait at all.
//
// Here web and app handlers count views on the same counter. Every
// handler locks before views++. But web handlers lock mu, and app handlers
// lock appMu. Web handlers wait for each other. App handlers wait for each
// other. But a web handler and an app handler never wait for each other.
//
// Problem: the count is lower than expected. "go run -race" reports a
// data race, even though every views++ is between Lock and Unlock.
//
// Fix: give one piece of data one mutex. Share it by pointer, and lock it
// for every read and write. Comment out the lines marked "problem", then
// uncomment the "fix:" line.
//
// Try: in the fixed version, change countViews to take "mu sync.Mutex"
// (no *) and pass mu instead of &mu. Now each call locks its own copy, so
// the race is back. "go vet" warns that it "copies lock value".
//
// Note: compare with 01-race, where nothing locked. Here everything locks,
// and it's still a race. To avoid this, keep the mutex inside a type, next
// to the field it guards (see 02-safe-type).
//
// Real use:
//   - someone adds a method to a type and gives it a new mutex. The old
//     methods still guard the same field with the old mutex.
//   - a method with a value receiver, like func (c counter) Inc(), copies
//     the struct and its mutex. So every call locks a fresh copy.
//
// Run: go run ./02-concurrency/07-mutex/06-two-mutexes
// Race check: go run -race ./02-concurrency/07-mutex/06-two-mutexes
package main

import (
	"fmt"
	"sync"
)

// viewsPerHandler is how many views each goroutine counts.
const viewsPerHandler = 1000

func main() {
	const handlerCount = 50 // of each kind: web and app

	views := 0
	var mu sync.Mutex    // guards views
	var appMu sync.Mutex // problem: a second mutex for the same views

	var wg sync.WaitGroup
	wg.Add(2 * handlerCount) // web + app
	for range handlerCount {
		go countViews(&views, &mu, &wg)    // web handler
		go countViews(&views, &appMu, &wg) // problem: app handler locks a different mutex
		// go countViews(&views, &mu, &wg) // fix: app handler locks the same mu
	}
	wg.Wait()

	fmt.Println("expected", 2*handlerCount*viewsPerHandler)
	fmt.Println("got     ", views)
}

// countViews adds viewsPerHandler to the shared counter, locking mu each
// time. It only keeps out goroutines that lock the same mu.
func countViews(views *int, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	for range viewsPerHandler {
		mu.Lock()
		*views++
		mu.Unlock()
	}
}
