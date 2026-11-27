// Basic: many callers ask for the same thing, and the work runs once.
//
// A popular user's profile gets several requests at the same moment. Each
// one calls g.Do with the same key. The first caller runs loadProfile.
// The others don't start their own load. They wait for the first one and
// get its result. shared is true when the result went to more than one
// caller.
//
// Once the call is finished, the key is free again. A later request for
// the same key runs loadProfile again. singleflight joins calls that
// overlap in time. It isn't a cache.
//
// g.Do only takes a func() (any, error), so each call is a one-line
// closure that calls a named function.
//
// Note: the group lives in this server's memory. With many servers, each
// one makes its own call, so the database gets one call per server. That's
// only extra work, not a wrong answer, so it's usually fine. If it must be
// one call for the whole app, use a Redis lock, or send each key to one
// server.
//
// Real use:
//   - an API where a viral post or a popular product gets many requests
//     at once: one db.QueryRowContext instead of hundreds.
//   - fetching the same remote file or token for many goroutines at once.
//
// Run: go run ./02-concurrency/12-singleflight/01-basic
package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// loadTime is how long the fake database takes.
const loadTime = 100 * time.Millisecond

func main() {
	const requestCount = 5

	var g singleflight.Group

	var wg sync.WaitGroup
	wg.Add(requestCount)
	for id := range requestCount {
		go handle(id+1, "user:42", &g, &wg)
	}
	wg.Wait()
}

// handle is one request for a profile. Requests for the same key at the
// same time share one load.
func handle(id int, key string, g *singleflight.Group, wg *sync.WaitGroup) {
	defer wg.Done()

	v, err, shared := g.Do(key, func() (any, error) { return loadProfile(key) })
	if err != nil {
		fmt.Printf("request %d: %v\n", id, err)
		return
	}
	fmt.Printf("request %d: %s (shared=%v)\n", id, v.(string), shared)
}

// loadProfile pretends to read a profile from the database.
func loadProfile(key string) (string, error) {
	fmt.Println("loading", key, "from the database...")
	time.Sleep(loadTime)
	return "Asif, joined 2019", nil
}
