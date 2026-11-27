// Compare: the same work with a map + Mutex, a map + RWMutex, and
// sync.Map.
//
// Each round runs three times, once with each map, and prints how long it
// took. All three maps have the same Load and Store methods, so the same
// code drives them.
//
// Round 1: keys are stored once at the start, then every goroutine reads
// them, like a product list that doesn't change.
//   - Mutex: readers take turns on one lock, so it's the slowest.
//   - RWMutex: readers can hold RLock together, but it's only a little
//     faster. Every RLock still updates one shared reader count, and the
//     CPU cores fight over it.
//   - sync.Map: reads take no lock at all, so it's the fastest. This is
//     the case sync.Map is built for.
//
// Round 2: every goroutine writes the same key again and again, like "who
// opened the home page last". Writers take turns on all three, so nothing
// helps here. sync.Map is the slowest, because it does more work per
// write. RWMutex can be a little faster than Mutex, even though its Lock
// does more work. With many writers waiting, most of the time goes into
// waking them up, and that depends on the machine and the Go version.
//
// sync.Map's other case, goroutines that each write their own keys, is
// in between. It's a little faster, and the gap changes a lot between
// runs, so it isn't shown here.
//
// Note: the times depend on your machine and how many CPU cores it has.
// Measure your own case with a benchmark (go test -bench) before you
// switch.
//
// Try: set workerCount to 1. Nothing waits on a lock now. Round 1 takes
// about the same time on all three, and sync.Map is clearly the slowest
// in round 2.
//
// Real use:
//   - deciding how to store a product catalog, a feature flag list, or a
//     per-host client registry that every request reads. Only switch from
//     a map with a Mutex if the numbers show it helps.
//
// Run: go run ./02-concurrency/13-sync-map/04-compare
// Race check: go run -race ./02-concurrency/13-sync-map/04-compare
package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// store is what each round needs from a map.
type store interface {
	Load(key string) (value int)
	Store(key string, value int)
}

// lockedMap is the usual way: a plain map with a Mutex.
type lockedMap struct {
	mu sync.Mutex
	m  map[string]int
}

// Load returns the value for key.
func (l *lockedMap) Load(key string) (value int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.m[key]
}

// Store sets the value for key.
func (l *lockedMap) Store(key string, value int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.m[key] = value
}

// rwLockedMap is a plain map with an RWMutex: RLock for reads, Lock for
// writes.
type rwLockedMap struct {
	mu sync.RWMutex
	m  map[string]int
}

// Load returns the value for key. Many goroutines can be in here at once.
func (r *rwLockedMap) Load(key string) (value int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.m[key]
}

// Store sets the value for key.
func (r *rwLockedMap) Store(key string, value int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[key] = value
}

// syncedMap is a sync.Map with int values.
type syncedMap struct {
	m sync.Map
}

// Load returns the value for key, or 0 if it's missing.
func (s *syncedMap) Load(key string) (value int) {
	v, _ := s.m.Load(key)
	value, _ = v.(int)
	return value
}

// Store sets the value for key.
func (s *syncedMap) Store(key string, value int) {
	s.m.Store(key, value)
}

const (
	workerCount  = 8       // goroutines that use the map at the same time
	opsPerWorker = 200_000 // reads or writes each goroutine does
	keyCount     = 100     // keys stored before round 1
	sameKey      = "/"     // the one key every goroutine writes in round 2
	roundTo      = time.Millisecond
)

func main() {
	compareReads()
	compareSameKeyWrites()
}

// compareReads is round 1: it stores keyCount keys in each map, then times
// the goroutines reading them.
func compareReads() {
	locked := &lockedMap{m: map[string]int{}}
	rwLocked := &rwLockedMap{m: map[string]int{}}
	synced := &syncedMap{}

	var keys []string
	for i := range keyCount {
		key := "p-" + strconv.Itoa(i) // a product id
		keys = append(keys, key)
		locked.Store(key, i)
		rwLocked.Store(key, i)
		synced.Store(key, i)
	}

	fmt.Println("round 1, read keys that were stored once:")
	fmt.Println("  map + Mutex:  ", timeWorkers(func() { readKeys(locked, keys) }))
	fmt.Println("  map + RWMutex:", timeWorkers(func() { readKeys(rwLocked, keys) }))
	fmt.Println("  sync.Map:     ", timeWorkers(func() { readKeys(synced, keys) }))
}

// compareSameKeyWrites is round 2: it times the goroutines writing
// sameKey in three empty maps.
func compareSameKeyWrites() {
	// Made once, outside the closures, so every goroutine shares the map.
	locked := &lockedMap{m: map[string]int{}}
	rwLocked := &rwLockedMap{m: map[string]int{}}
	synced := &syncedMap{}

	fmt.Println("round 2, write the same key:")
	fmt.Println("  map + Mutex:  ", timeWorkers(func() { writeSameKey(locked) }))
	fmt.Println("  map + RWMutex:", timeWorkers(func() { writeSameKey(rwLocked) }))
	fmt.Println("  sync.Map:     ", timeWorkers(func() { writeSameKey(synced) }))
}

// timeWorkers runs work in workerCount goroutines at once, and returns
// how long they took together.
func timeWorkers(work func()) (took time.Duration) {
	start := time.Now()
	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(work)
	}
	wg.Wait()
	return time.Since(start).Round(roundTo)
}

// readKeys reads the keys in turn from m, opsPerWorker times.
func readKeys(m store, keys []string) {
	for i := range opsPerWorker {
		m.Load(keys[i%len(keys)])
	}
}

// writeSameKey writes sameKey in m, opsPerWorker times.
func writeSameKey(m store) {
	for i := range opsPerWorker {
		m.Store(sameKey, i)
	}
}
