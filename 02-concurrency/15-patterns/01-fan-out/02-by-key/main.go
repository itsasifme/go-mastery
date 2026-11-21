// By key: keep the order of events that belong together.
//
// A shop sends status events for each order: created, then paid, then
// shipped. Three workers share the work, and some are slower than others.
//
// Problem: with plain fan-out (see 01-basic), any free worker takes the
// next event. The events of one order go to different workers, and a slow
// worker can finish "created" after a fast one already applied "paid".
// The store counts every event that arrived out of order.
//
// Fix: give each worker its own channel, and send every event of an order
// to the same worker: order % workerCount picks it. One worker handles an
// order's events one by one, so they stay in order. Different orders still
// run in parallel on different workers. Use the "fix:" line instead of the
// line above it.
//
// Note: for string keys like user IDs, hash the key first (hash/fnv) and
// take that modulo workerCount. One very busy key can keep its worker busy
// while the others sit idle.
//
// Real use:
//   - Kafka sends every message with the same key to the same partition,
//     and one consumer reads each partition, so one user's or one
//     account's events are handled in order.
//   - a bank account's transactions, a chat room's messages, or an
//     order's status changes, handled in parallel across keys but in order
//     within each key.
//
// Run: go run ./02-concurrency/15-patterns/01-fan-out/02-by-key
package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

// event is one status change of an order.
type event struct {
	order  int
	status string
}

// store keeps the latest status of each order.
type store struct {
	mu         sync.Mutex     // guards status and outOfOrder
	status     map[int]string // order -> latest status
	outOfOrder int
}

// Apply records e. An event older than the order's current status arrived
// out of order: it's counted and ignored.
func (s *store) Apply(e event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.status[e.order]
	if slices.Index(steps, e.status) < slices.Index(steps, current) {
		s.outOfOrder++
		fmt.Printf("order %d: %q arrived after %q\n", e.order, e.status, current)
		return
	}
	s.status[e.order] = e.status
}

// OutOfOrder returns how many events arrived out of order.
func (s *store) OutOfOrder() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.outOfOrder
}

const (
	workerCount = 3
	orderCount  = 6
	workUnit    = 10 * time.Millisecond // worker 1's time per event. Worker 2 takes twice as long, and so on
)

// steps is the order the statuses must happen in.
var steps = []string{"created", "paid", "shipped"}

func main() {
	s := &store{status: map[int]string{}}

	fanOutShared(s) // problem: any free worker takes the next event
	// fanOutByKey(s) // fix: every event of one order goes to the same worker

	fmt.Println("events out of order:", s.OutOfOrder())
}

// fanOutShared starts workers that all read one channel, sends every
// event, and waits for the workers.
func fanOutShared(s *store) {
	events := make(chan event) // sendEvents -> any free worker

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go work(id+1, events, s, &wg)
	}

	sendEvents([]chan event{events}) // one channel, so every event goes to it
	close(events)
	wg.Wait()
}

// fanOutByKey starts workers that each read their own channel, sends every
// event to the worker picked by its order, and waits for the workers.
func fanOutByKey(s *store) {
	queues := make([]chan event, workerCount) // sendEvents -> one worker each

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		queues[id] = make(chan event)
		go work(id+1, queues[id], s, &wg)
	}

	sendEvents(queues)
	for _, q := range queues {
		close(q)
	}
	wg.Wait()
}

// sendEvents sends created, paid and shipped for every order, in that
// order and right after each other, like webhooks for one order arriving
// close together. The order number picks the queue, so one order's events
// always go to the same queue.
func sendEvents(queues []chan event) {
	for order := range orderCount {
		for _, status := range steps {
			e := event{order: order + 1, status: status}
			queues[e.order%len(queues)] <- e
		}
	}
}

// work applies events until its channel is closed. Higher ids are slower.
func work(id int, in <-chan event, s *store, wg *sync.WaitGroup) {
	defer wg.Done()

	workTime := time.Duration(id) * workUnit
	for e := range in {
		time.Sleep(workTime) // simulate handling the event
		s.Apply(e)
	}
}
