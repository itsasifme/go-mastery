// Compare: the same read-only work with a Mutex and with an RWMutex.
//
// Several readers each do some lookups. Each lookup holds the lock for a
// short time, like copying a big value out of a cache. With a Mutex, only
// one reader can hold the lock, so the lookups run one after another. With
// an RWMutex, every reader can hold RLock at once, so they overlap.
//
// Note: here every goroutine only reads, which is the best case for
// RWMutex. With many writes, or with very short lookups, RWMutex can be
// slower than a Mutex, because it does more bookkeeping. Measure your own
// case with a benchmark (go test -bench) before you switch.
//
// Try: set holdTime to 0. Both take almost no time, and the gap is gone.
//
// Real use:
//   - deciding between Mutex and RWMutex for a cache that every request
//     reads. Only switch if the numbers show it helps.
//
// Run: go run ./02-concurrency/08-rwmutex/02-compare
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	readerCount      = 10
	lookupsPerReader = 10
	holdTime         = time.Millisecond // how long one lookup holds the lock
	roundTo          = 10 * time.Millisecond
)

func main() {
	fmt.Println("Mutex:  ", timeMutexReaders())
	fmt.Println("RWMutex:", timeRWMutexReaders())
}

// timeMutexReaders runs every reader with a shared Mutex and returns how
// long they took.
func timeMutexReaders() time.Duration {
	start := time.Now()

	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(readerCount)
	for range readerCount {
		go readWithMutex(&mu, &wg)
	}
	wg.Wait()

	return time.Since(start).Round(roundTo)
}

// readWithMutex does its lookups while holding the Mutex.
func readWithMutex(mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for range lookupsPerReader {
		mu.Lock()
		time.Sleep(holdTime) // simulate the lookup
		mu.Unlock()
	}
}

// timeRWMutexReaders runs every reader with a shared RWMutex and returns
// how long they took.
func timeRWMutexReaders() time.Duration {
	start := time.Now()

	var mu sync.RWMutex
	var wg sync.WaitGroup
	wg.Add(readerCount)
	for range readerCount {
		go readWithRWMutex(&mu, &wg)
	}
	wg.Wait()

	return time.Since(start).Round(roundTo)
}

// readWithRWMutex does its lookups while holding the read lock.
func readWithRWMutex(mu *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	for range lookupsPerReader {
		mu.RLock()
		time.Sleep(holdTime) // simulate the lookup
		mu.RUnlock()
	}
}
