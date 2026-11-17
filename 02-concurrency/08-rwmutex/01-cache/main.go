// Cache: many readers at once, one writer at a time.
//
// A shop shows prices in the visitor's currency, so every request reads an
// exchange rate. The rates change only now and then. RLock lets all the
// readers in at the same time. Lock waits for the readers to leave, then
// lets one writer in alone. The readers keep working while the rate
// changes, and they never see a half-written map.
//
// Note: with a plain Mutex (see mutex/02-safe-type), readers would wait for
// each other, even though reading doesn't change anything.
//
// Real use:
//   - exchange rates, feature flags or a list of allowed API keys that
//     every request reads and a background job refreshes.
//   - a DNS or service discovery cache: lookups use RLock, refreshes use Lock.
//
// Run: go run ./02-concurrency/08-rwmutex/01-cache
// Race check: go run -race ./02-concurrency/08-rwmutex/01-cache
package main

import (
	"fmt"
	"sync"
	"time"
)

// rateCache holds exchange rates. It's safe to use from many goroutines.
type rateCache struct {
	mu    sync.RWMutex // guards rates
	rates map[string]float64
}

// Get returns the rate for a currency pair. Many Gets can run at once.
func (c *rateCache) Get(pair string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rates[pair]
}

// Set changes the rate for a currency pair. It waits until no one is reading.
func (c *rateCache) Set(pair string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rates[pair] = rate
}

const (
	pair       = "USD/BDT"
	readTime   = 10 * time.Millisecond // time between two lookups by one reader
	updateTime = 100 * time.Millisecond
)

func main() {
	const (
		readerCount = 3
		runTime     = 350 * time.Millisecond
	)

	cache := &rateCache{rates: map[string]float64{pair: 121.0}}
	stop := make(chan struct{}) // main -> readers and updater: "please stop"

	var wg sync.WaitGroup
	wg.Add(readerCount + 1) // +1 for the updater
	for id := range readerCount {
		go readRates(id+1, cache, stop, &wg)
	}
	go updateRates(cache, stop, &wg)

	time.Sleep(runTime)
	close(stop)
	wg.Wait()
}

// readRates looks up the rate again and again, like requests coming in,
// until stop is closed.
func readRates(id int, c *rateCache, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	lookups := 0
	var last float64
	for {
		select {
		case <-time.After(readTime):
			last = c.Get(pair)
			lookups++
		case <-stop:
			fmt.Printf("reader %d: %d lookups, last %s = %.1f\n", id, lookups, pair, last)
			return
		}
	}
}

// updateRates sets a new rate every updateTime, like a job that pulls
// rates from a bank, until stop is closed.
func updateRates(c *rateCache, stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	newRates := []float64{121.5, 122.0, 122.5, 123.0} // sample rates from the "bank"
	updates := 0
	for {
		select {
		case <-time.After(updateTime):
			c.Set(pair, newRates[updates%len(newRates)])
			updates++
		case <-stop:
			fmt.Println("updater:", updates, "updates")
			return
		}
	}
}
