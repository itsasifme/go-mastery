// Basic: run several jobs at once and get the first error back.
//
// At startup, an app checks that the database, cache and queue are
// reachable. Each check runs in g.Go. g.Wait waits for all of them, like a
// WaitGroup, and returns the first error. Here the cache check fails, so
// the app doesn't start.
//
// g.Go only takes a func() error, so each call is a one-line closure that
// calls a named function. All the work stays in the named function.
//
// Note: compare with waitgroup/02-parallel-fetch. A WaitGroup can't return
// an error, so you'd need your own error variable and a lock.
//
// Try: change the cache's error to nil, and remove the "errors" import.
// Wait returns nil and the app starts.
//
// Real use:
//   - startup checks: db.PingContext, a Redis PING and a queue connection,
//     all at once, so a slow one doesn't add to the others.
//   - deleting a user's files from several storage buckets and reporting
//     the first one that failed.
//
// Run: go run ./02-concurrency/06-errgroup/01-basic
package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// checkTime is how long one fake check takes.
const checkTime = 100 * time.Millisecond

func main() {
	services := map[string]error{ // service -> what its fake check returns
		"database": nil,
		"cache":    errors.New("connection refused"),
		"queue":    nil,
	}

	var g errgroup.Group
	for name, result := range services {
		g.Go(func() error { return ping(name, result) })
	}

	if err := g.Wait(); err != nil {
		fmt.Println("startup failed:", err)
		return
	}
	fmt.Println("all services ready, starting app")
}

// ping pretends to check one service. It returns result as the error.
func ping(name string, result error) error {
	time.Sleep(checkTime)
	if result != nil {
		return fmt.Errorf("%s: %w", name, result)
	}
	fmt.Println(name, "ok")
	return nil
}
