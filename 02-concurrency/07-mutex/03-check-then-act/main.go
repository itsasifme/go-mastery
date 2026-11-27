// Check then act: a lock on each step isn't enough.
//
// A concert has one ticket left and several buyers click "buy" at once.
// Each method of inventory locks the mutex, so there's no data race. But
// buyUnsafe calls Available, then charges the card, then calls Take. Those
// are separate locks. Between them, other buyers see the same last ticket.
//
// Problem: more buyers get a ticket than there are tickets, and the stock
// goes below zero. "go run -race" finds nothing, because every access is
// locked. The bug is in the logic, not the memory.
//
// Fix: check and take in one method, under one lock (TryReserve). Charge
// the card after the ticket is reserved. Comment out the line marked
// "problem", then uncomment the "fix:" line.
//
// Note: when the stock is just one number, you don't need a lock. Atomic's
// CompareAndSwap can check and take in one step (see
// atomic/04-compare-and-swap). You still need a mutex when more than one
// value must change together.
//
// Note: a mutex only works inside one server. The tickets belong to the
// whole app, so with many servers each one has its own lock, and the last
// ticket gets sold again. Let the database check and take in one step:
// UPDATE tickets SET left = left - 1 WHERE id = $1 AND left > 0.
//
// Real use:
//   - selling the last item in stock, booking the last seat, or using a
//     coupon that has a limit.
//   - "if not in cache, then add": check and add under one lock, or use
//     a single call like sync.Map's LoadOrStore (see
//     sync-map/02-load-or-store).
//
// Run: go run ./02-concurrency/07-mutex/03-check-then-act
// Race check: go run -race ./02-concurrency/07-mutex/03-check-then-act
package main

import (
	"fmt"
	"sync"
	"time"
)

// inventory counts tickets left. Every method locks, so it has no data race.
type inventory struct {
	mu   sync.Mutex // guards left
	left int
}

// Available reports how many tickets are left.
func (inv *inventory) Available() int {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	return inv.left
}

// Take removes one ticket, without checking.
func (inv *inventory) Take() {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	inv.left--
}

// TryReserve takes one ticket if there is one. The check and the change
// happen under the same lock.
func (inv *inventory) TryReserve() bool {
	inv.mu.Lock()
	defer inv.mu.Unlock()
	if inv.left == 0 {
		return false
	}
	inv.left--
	return true
}

// chargeTime is how long charging a card takes.
const chargeTime = 50 * time.Millisecond

func main() {
	const (
		tickets    = 1
		buyerCount = 5
	)

	inv := &inventory{left: tickets}

	var wg sync.WaitGroup
	wg.Add(buyerCount)
	for id := range buyerCount {
		go buyUnsafe(id+1, inv, &wg) // problem: check and take are two separate locks
		// go buySafe(id+1, inv, &wg) // fix: reserve first, then charge
	}
	wg.Wait()

	fmt.Println("tickets left:", inv.Available())
}

// buyUnsafe checks for a ticket, charges the card, then takes the ticket.
// Other buyers can pass the same check while the card is being charged.
func buyUnsafe(buyer int, inv *inventory, wg *sync.WaitGroup) {
	defer wg.Done()

	if inv.Available() == 0 {
		fmt.Println("buyer", buyer, "sold out")
		return
	}
	time.Sleep(chargeTime) // charge the card
	inv.Take()
	fmt.Println("buyer", buyer, "got a ticket")
}

// buySafe reserves the ticket first, then charges the card.
func buySafe(buyer int, inv *inventory, wg *sync.WaitGroup) {
	defer wg.Done()

	if !inv.TryReserve() {
		fmt.Println("buyer", buyer, "sold out")
		return
	}
	time.Sleep(chargeTime) // charge the card
	fmt.Println("buyer", buyer, "got a ticket")
}
