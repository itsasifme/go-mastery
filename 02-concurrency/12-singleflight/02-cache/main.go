// Cache: stop a crowd of cache misses from hitting the database.
//
// A product page is cached. When the cache entry expires, every request
// that comes in at that moment misses at once. Without singleflight, each
// of them queries the database for the same row. This is called a cache
// stampede, and it can take a database down.
//
// With singleflight, the first miss loads the row and fills the cache.
// The other misses join that load. The database sees one query. Inside
// the shared call, Get looks in the cache once more. Without that, a
// request that missed just as a load finished would start a second load.
//
// The program runs the same burst of requests twice, once with each kind
// of store, and prints how many queries reached the database.
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
	mu        sync.Mutex // guards cache
	cache     map[string]string
	group     singleflight.Group
	dbQueries atomic.Int64
}

// GetNoDedup checks the cache, and on a miss queries the database itself.
func (s *productStore) GetNoDedup(id string) string {
	if name, ok := s.cached(id); ok {
		return name
	}
	return s.loadAndCache(id)
}

// Get checks the cache, and on a miss joins any load for the same id that
// is already running.
func (s *productStore) Get(id string) string {
	if name, ok := s.cached(id); ok {
		return name
	}
	v, _, _ := s.group.Do(id, func() (any, error) { return s.cachedOrLoad(id), nil })
	return v.(string)
}

// cachedOrLoad runs inside the shared call. It looks in the cache again,
// because a load that just finished may have filled it. On a miss, it loads.
func (s *productStore) cachedOrLoad(id string) string {
	if name, ok := s.cached(id); ok {
		return name
	}
	return s.loadAndCache(id)
}

// cached looks id up in the cache.
func (s *productStore) cached(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	name, ok := s.cache[id]
	return name, ok
}

// loadAndCache queries the fake database and saves the result in the cache.
func (s *productStore) loadAndCache(id string) string {
	s.dbQueries.Add(1)
	time.Sleep(queryTime) // simulate the query
	name := "Desk lamp"

	s.mu.Lock()
	s.cache[id] = name
	s.mu.Unlock()
	return name
}

// queryTime is how long the fake database query takes.
const queryTime = 100 * time.Millisecond

func main() {
	const requestCount = 20

	plain := &productStore{cache: map[string]string{}} // empty cache, like right after it expired
	runBurst(requestCount, plain.GetNoDedup)
	fmt.Println("without singleflight:", plain.dbQueries.Load(), "database queries")

	deduped := &productStore{cache: map[string]string{}}
	runBurst(requestCount, deduped.Get)
	fmt.Println("with singleflight:   ", deduped.dbQueries.Load(), "database queries")
}

// runBurst sends requestCount requests at the same moment and waits for
// all of them.
func runBurst(requestCount int, get func(id string) string) {
	var wg sync.WaitGroup
	wg.Add(requestCount)
	for range requestCount {
		go request("p-7", get, &wg)
	}
	wg.Wait()
}

// request asks the store for one product.
func request(id string, get func(id string) string, wg *sync.WaitGroup) {
	defer wg.Done()
	get(id)
}
