// Struct: keep the mutex inside the type, next to the data it protects.
//
// A wallet has a balance that many requests change at once. The mutex
// lives in the struct, right above the field it guards. Every method locks
// it, so code that uses the wallet can't forget to. Withdraw checks and
// changes the balance under one lock, so it can't go below zero.
//
// Methods use a pointer receiver (*wallet). A value receiver would copy
// the mutex, and the copy wouldn't protect anything. go vet warns about it.
//
// Note: compare with 01-race. There, every caller had to lock by hand.
//
// Real use:
//   - an in-memory session store or rate limiter: a struct with a
//     sync.Mutex and a map, and methods like Get, Set and Delete.
//   - a connection pool that tracks free and busy connections.
//
// Run: go run ./02-concurrency/07-mutex/02-struct
package main

import (
	"errors"
	"fmt"
	"sync"
)

// wallet is safe to use from many goroutines.
type wallet struct {
	mu      sync.Mutex // guards balance
	balance int
}

// Deposit adds amount to the balance.
func (w *wallet) Deposit(amount int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.balance += amount
}

// Withdraw takes amount out, or returns errNoMoney. The check and the
// change happen under the same lock.
func (w *wallet) Withdraw(amount int) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.balance < amount {
		return errNoMoney
	}
	w.balance -= amount
	return nil
}

// Balance returns the current balance.
func (w *wallet) Balance() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.balance
}

// Sample amounts in cents.
const (
	depositAmount  = 500
	withdrawAmount = 300
	timesEach      = 100 // how many times each goroutine deposits or withdraws
)

// errNoMoney is returned when a withdrawal is bigger than the balance.
var errNoMoney = errors.New("not enough money")

func main() {
	const (
		depositors  = 10
		withdrawers = 10
	)

	// Start with enough money that no withdrawal fails, so the final
	// balance is the same on every run.
	start := withdrawers * timesEach * withdrawAmount
	w := &wallet{balance: start}

	var wg sync.WaitGroup
	wg.Add(depositors + withdrawers)
	for range depositors {
		go depositMany(w, &wg)
	}
	for range withdrawers {
		go withdrawMany(w, &wg)
	}
	wg.Wait()

	expected := start + depositors*timesEach*depositAmount - withdrawers*timesEach*withdrawAmount
	fmt.Println("expected", expected)
	fmt.Println("got     ", w.Balance())
}

// depositMany deposits depositAmount timesEach times.
func depositMany(w *wallet, wg *sync.WaitGroup) {
	defer wg.Done()
	for range timesEach {
		w.Deposit(depositAmount)
	}
}

// withdrawMany withdraws withdrawAmount timesEach times.
func withdrawMany(w *wallet, wg *sync.WaitGroup) {
	defer wg.Done()
	for range timesEach {
		if err := w.Withdraw(withdrawAmount); err != nil {
			fmt.Println("withdraw:", err)
		}
	}
}
