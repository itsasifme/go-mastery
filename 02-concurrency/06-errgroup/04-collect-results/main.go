// Collect results: fetch many values at once, or fail as a whole.
//
// A checkout page needs the price of every item in the cart. Each price
// comes from a pricing service, so all of them are fetched at once. Each
// goroutine writes to its own index of prices, so no lock is needed. main
// reads prices only after Wait returns nil.
//
// If one item fails, ctx is cancelled, the other fetches stop, and the
// page shows an error instead of a wrong total.
//
// Try: add "ghost" to the cart. It isn't in the catalog, so it fails right
// away. ctx is cancelled, the other fetches print "stopped", and Wait
// returns ghost's error. No total is printed.
//
// Note: cart and prices are matched only by index. prices[i] is the price
// of cart[i]. Two slices kept in step like this are usually a bad sign. If
// you sort or filter one, the match breaks. Here it's fine. It's the usual
// Go way to collect results, and errgroup's own docs do it too. Make the
// slice first, and let each goroutine write only to its own index. Then
// no lock is needed, and the results stay in the cart's order.
//
// Real use:
//   - an API that loads several things for one response (prices, stock,
//     shipping cost) and returns 500 if any of them fails.
//   - running several db.QueryContext(ctx, ...) calls for one page, and
//     cancelling the rest when one fails.
//
// Run: go run ./02-concurrency/06-errgroup/04-collect-results
package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// fetchTime is how long the fake pricing service takes to answer.
const fetchTime = 100 * time.Millisecond

// catalog is the fake pricing service's data: product -> price in cents.
// Goroutines only read it, so it's safe without a lock.
var catalog = map[string]int{
	"book": 1250,
	"pen":  199,
	"lamp": 3400,
}

func main() {
	cart := []string{"book", "pen", "lamp"}
	prices := make([]int, len(cart)) // each goroutine writes only its own index

	g, ctx := errgroup.WithContext(context.Background())
	for i, product := range cart {
		g.Go(func() error { return fetchPrice(ctx, product, &prices[i]) })
	}

	if err := g.Wait(); err != nil {
		fmt.Println("checkout failed:", err)
		return
	}

	total := 0
	for i, product := range cart {
		fmt.Printf("%-5s %d cents\n", product, prices[i])
		total += prices[i]
	}
	fmt.Println("total", total, "cents")
}

// fetchPrice asks the fake pricing service for one product's price and
// writes it to out. An unknown product fails right away (like a quick 404).
// A known one takes fetchTime, and stops early if ctx is cancelled.
func fetchPrice(ctx context.Context, product string, out *int) error {
	price, ok := catalog[product]
	if !ok {
		return fmt.Errorf("no price for %q", product)
	}

	select {
	case <-time.After(fetchTime): // simulate the network call
	case <-ctx.Done():
		fmt.Println(product, "stopped: checkout already failed")
		return ctx.Err()
	}

	*out = price
	return nil
}
