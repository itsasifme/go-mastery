// Cache: stop a crowd of cache misses from hitting the database.
//
// A product page is cached. When the cache entry expires, every request
// that comes in at that moment misses at once. Without singleflight, each
// one queries the database for the same row. This is called a cache
// stampede, and it can take a database down.
//
// Get wraps the "miss, then load" part in group.Do. The first miss loads
// the row and fills the cache, and the other misses join that load. Inside
// the shared call, the cache is checked once more, in case a load finished
// just before.
//
// The same requests run twice, once without and once with singleflight,
// and it prints how many queries reached the database.
//
// Real use:
//   - a Redis or in-memory cache in front of Postgres: wrap the "miss ->
//     load -> set" part in g.Do(key, ...).
//   - refreshing an OAuth token that many requests need at the same time.
//
// Run: go run ./02-concurrency/12-singleflight/02-cache
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// productStore is a cache in front of a fake database.
type productStore struct {
	mu        sync.Mutex        // guards cache
	cache     map[string]string // product id -> product name
	group     singleflight.Group
	dbQueries atomic.Int64
}

// GetNoDedup returns the product's name. It checks the cache, and on a miss
// queries the database itself. "Dedup" means removing duplicates: here,
// nothing stops them.
func (s *productStore) GetNoDedup(id string) (name string) {
	if name, found := s.cached(id); found {
		return name
	}
	return s.load(id)
}

// Get returns the product's name too, but misses for the same id share one
// call. Inside it, GetNoDedup checks the cache again before it loads.
func (s *productStore) Get(id string) (name string) {
	if name, found := s.cached(id); found {
		return name // most requests stop here, without singleflight
	}
	v, _, _ := s.group.Do(id, func() (any, error) { return s.GetNoDedup(id), nil })
	return v.(string) // Do returns any, so turn it back into the name
}

// cached returns the product's name from the cache. found is false on a miss.
func (s *productStore) cached(id string) (name string, found bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	name, found = s.cache[id]
	return name, found
}

// load returns the product's name from the fake database, and saves it in
// the cache.
func (s *productStore) load(id string) (name string) {
	s.dbQueries.Add(1)
	time.Sleep(queryTime) // simulate the query
	name = "Desk lamp"

	s.mu.Lock()
	s.cache[id] = name
	s.mu.Unlock()
	return name
}

// queryTime is how long the fake database query takes.
const queryTime = 100 * time.Millisecond

func main() {
	const userCount = 20 // users who open the product page at the same moment

	plain := &productStore{cache: map[string]string{}} // empty cache, like right after it expired
	// plain.GetNoDedup is a method value: a func that remembers plain. So each
	// round uses its own store and counts its own queries.
	simulateStampede(userCount, plain.GetNoDedup)
	fmt.Println("without singleflight:", plain.dbQueries.Load(), "database queries")

	deduped := &productStore{cache: map[string]string{}}
	simulateStampede(userCount, deduped.Get)
	fmt.Println("with singleflight:   ", deduped.dbQueries.Load(), "database queries")
}

// simulateStampede makes userCount users get the same product at once.
func simulateStampede(userCount int, getProduct func(id string) (name string)) {
	var wg sync.WaitGroup
	for range userCount {
		wg.Go(func() { getProduct("p-7") })
	}
	wg.Wait()
}
