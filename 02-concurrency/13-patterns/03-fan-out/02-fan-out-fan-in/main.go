// Fan-out, fan-in: share the work, then collect the results in one place.
//
// A price comparison site checks many stores for the same product.
// Fan-out: several fetchers read store names from one channel, so the
// slow network calls overlap. Each fetcher sends its quotes on its own
// channel. Fan-in: merge combines those channels, and main picks the
// cheapest quote.
//
// Note: this is 01-basic plus fan-in/01-merge, put together.
//
// Try: set fetcherCount to 1. It takes about as long as all the fetches
// added up.
//
// Real use:
//   - price or flight search across many providers, then sorting the
//     results for one page.
//   - a crawler: fetchers download pages in parallel, and one place
//     collects the links they find.
//
// Run: go run ./02-concurrency/13-patterns/03-fan-out/02-fan-out-fan-in
package main

import (
	"fmt"
	"sync"
	"time"
)

// quote is one store's price for the product.
type quote struct {
	store string
	price int // in cents
}

const (
	fetchTime = 100 * time.Millisecond
	roundTo   = 10 * time.Millisecond
)

// prices is the fake answer from each store's API.
var prices = map[string]int{
	"alpha-shop": 2599,
	"bazaar":     2399,
	"cheap-mart": 2199,
	"deal-hub":   2450,
	"e-store":    2650,
	"fair-price": 2299,
}

func main() {
	const fetcherCount = 3

	start := time.Now()
	stores := make(chan string) // queueStores -> fetchers

	go queueStores(stores)

	outs := make([]<-chan quote, fetcherCount)
	for i := range fetcherCount {
		out := make(chan quote) // this fetcher -> merge
		go fetchQuotes(stores, out)
		outs[i] = out
	}

	best := quote{}
	for q := range merge(outs...) {
		fmt.Printf("%-10s %d cents\n", q.store, q.price)
		if best.store == "" || q.price < best.price {
			best = q
		}
	}
	fmt.Printf("cheapest: %s at %d cents (took %v)\n", best.store, best.price, time.Since(start).Round(roundTo))
}

// queueStores sends every store name, then closes stores.
func queueStores(stores chan<- string) {
	defer close(stores)
	for store := range prices {
		stores <- store
	}
}

// fetchQuotes asks each store it gets for a price, until stores is
// closed. Then it closes out.
func fetchQuotes(stores <-chan string, out chan<- quote) {
	defer close(out)
	for store := range stores {
		time.Sleep(fetchTime) // simulate the store's API call
		out <- quote{store: store, price: prices[store]}
	}
}

// merge returns one channel that gets every value from inputs. It's
// closed after every input is closed.
func merge(inputs ...<-chan quote) <-chan quote {
	out := make(chan quote)

	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, in := range inputs {
		go forward(in, out, &wg)
	}
	go closeWhenDone(&wg, out)

	return out
}

// forward copies every value from in to out, until in is closed.
func forward(in <-chan quote, out chan<- quote, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		out <- v
	}
}

// closeWhenDone waits for every forward, then closes out.
func closeWhenDone(wg *sync.WaitGroup, out chan<- quote) {
	wg.Wait()
	close(out)
}
