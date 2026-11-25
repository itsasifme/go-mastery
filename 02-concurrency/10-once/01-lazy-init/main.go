// Lazy init: set something up on first use, only once.
//
// Opening a database pool and checking it with Ping is slow, so the app
// does it the first time a request needs it, not at startup. Several
// requests arrive at once and all call getDB. once.Do runs openPool only
// for the first caller. The others wait until openPool is finished, then
// all of them get the same pool.
//
// Without Once, each request would check "pool == nil" at the same time,
// and several of them would open their own pool.
//
// Note: the server shares one pool, not one connection. A real pool like
// *sql.DB holds many connections. It lends one to each query, then takes
// it back for the next query. "pool 1" in the output is the pool's ID, so
// every request got the same pool. The fake pool here has no connections
// inside, because this lesson is about Once.
//
// Real use:
//   - opening a sql.DB, a Redis client or an S3 client on first use.
//   - compiling a big regexp.MustCompile or loading a template only when
//     it's first needed.
//
// Run: go run ./02-concurrency/10-once/01-lazy-init
package main

import (
	"fmt"
	"sync"
	"time"
)

// dbPool stands in for a real connection pool like *sql.DB.
type dbPool struct {
	id int
}

// openTime is how long opening the fake pool and checking it takes.
const openTime = 100 * time.Millisecond

var (
	dbOnce sync.Once
	pool   *dbPool // set by openPool, only once
)

func main() {
	const requestCount = 5

	var wg sync.WaitGroup
	wg.Add(requestCount)
	for id := range requestCount {
		go handle(id+1, &wg)
	}
	wg.Wait()
}

// handle is a request that needs the database.
func handle(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	db := getDB()
	fmt.Printf("request %d: using pool %d\n", id, db.id)
}

// getDB returns the shared pool, opening it on the first call.
func getDB() *dbPool {
	dbOnce.Do(openPool)
	return pool
}

// openPool pretends to open a connection pool and check it with Ping.
// With a real database, db, err := sql.Open(...) is fast and doesn't
// connect. db.PingContext(ctx) makes the first real connection.
func openPool() {
	fmt.Println("opening the database pool and checking it with Ping...")
	time.Sleep(openTime)
	pool = &dbPool{id: 1}
}
