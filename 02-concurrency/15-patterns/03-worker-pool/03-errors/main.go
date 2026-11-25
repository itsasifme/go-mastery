// Errors: keep going when some jobs fail, and collect what failed.
//
// A billing job charges a list of invoices. Some cards are declined. One
// failure shouldn't stop the others, so each result holds its own err.
// main collects the failures into a list to retry later or show to
// support.
//
// Note: compare with errgroup/03-with-context. There, the first error
// stopped everything. Pick what fits the job: "all or nothing" or "do as
// many as you can".
//
// Real use:
//   - charging invoices, sending webhooks, or importing CSV rows, where
//     each item can fail on its own.
//   - a report at the end: "98 sent, 2 failed: [id-4 id-9]".
//
// Run: go run ./02-concurrency/15-patterns/03-worker-pool/03-errors
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// invoice is one charge to make.
type invoice struct {
	id     string
	amount int // in cents
	card   string
}

// result is the outcome of one charge.
type result struct {
	invoiceID string
	err       error
}

// chargeTime is how long one fake charge takes.
const chargeTime = 50 * time.Millisecond

// errDeclined is what the fake payment provider returns for a bad card.
var errDeclined = errors.New("card declined")

func main() {
	const workerCount = 3

	invoices := []invoice{
		{"inv-1", 1200, "visa-ok"},
		{"inv-2", 4500, "visa-ok"},
		{"inv-3", 990, "expired"},
		{"inv-4", 3000, "visa-ok"},
		{"inv-5", 750, "expired"},
		{"inv-6", 2100, "visa-ok"},
	}

	jobs := make(chan invoice)
	results := make(chan result)

	go queueInvoices(invoices, jobs)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for range workerCount {
		go worker(jobs, results, &wg)
	}
	go closeWhenDone(&wg, results)

	paid := 0
	var failed []string
	for r := range results {
		if r.err != nil {
			fmt.Printf("%s: %v\n", r.invoiceID, r.err)
			failed = append(failed, r.invoiceID)
			continue
		}
		paid++
	}
	fmt.Printf("paid %d, failed %d: %v\n", paid, len(failed), failed)
}

// queueInvoices sends every invoice to jobs, then closes it.
func queueInvoices(invoices []invoice, jobs chan<- invoice) {
	defer close(jobs)
	for _, inv := range invoices {
		jobs <- inv
	}
}

// worker charges invoices until jobs is closed. A failed charge becomes a
// result with an error. It doesn't stop the worker.
func worker(jobs <-chan invoice, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for inv := range jobs {
		results <- result{invoiceID: inv.id, err: charge(inv)}
	}
}

// closeWhenDone waits for every worker, then closes results.
func closeWhenDone(wg *sync.WaitGroup, results chan<- result) {
	wg.Wait()
	close(results)
}

// charge pretends to charge a card.
func charge(inv invoice) error {
	time.Sleep(chargeTime)
	if inv.card == "expired" {
		return fmt.Errorf("charge %d cents: %w", inv.amount, errDeclined)
	}
	return nil
}
