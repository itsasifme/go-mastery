// Basic: several workers read from one channel to share the work.
//
// A warehouse has a stream of orders and three packers. All packers read
// from the same orders channel. Each order goes to exactly one packer,
// whichever is free first. Nobody assigns the work. A faster packer just
// comes back for more sooner, so it ends up packing more orders.
//
// Try: in pack, set packTime to packTimeUnit, so every packer has the
// same speed. The orders split about evenly.
//
// Real use:
//   - several consumers on one queue, like workers reading the same SQS
//     queue or Kafka consumer group: each message is handled once.
//   - spreading CPU-heavy work (hashing, compressing) over runtime.NumCPU()
//     goroutines.
//
// Run: go run ./02-concurrency/13-patterns/03-fan-out/01-basic
package main

import (
	"fmt"
	"sync"
	"time"
)

// packTimeUnit is packer 1's time per order. Packer 2 takes twice as
// long, and so on.
const packTimeUnit = 20 * time.Millisecond

func main() {
	const (
		packerCount = 3
		orderCount  = 12
	)

	orders := make(chan int) // queueOrders -> packers: order numbers

	go queueOrders(orderCount, orders)

	packed := make([]int, packerCount) // each packer writes only its own index
	var wg sync.WaitGroup
	wg.Add(packerCount)
	for id := range packerCount {
		go pack(id+1, orders, &packed[id], &wg)
	}
	wg.Wait()

	for id, n := range packed {
		fmt.Printf("packer %d packed %d orders\n", id+1, n)
	}
}

// queueOrders sends order numbers 1 to count, then closes orders.
func queueOrders(count int, orders chan<- int) {
	defer close(orders)
	for i := range count {
		orders <- i + 1
	}
}

// pack packs orders until the channel is closed, and counts them in packed.
func pack(id int, orders <-chan int, packed *int, wg *sync.WaitGroup) {
	defer wg.Done()

	packTime := time.Duration(id) * packTimeUnit // higher id = slower packer
	for order := range orders {
		time.Sleep(packTime)
		fmt.Printf("packer %d packed order %d\n", id, order)
		*packed++
	}
}
