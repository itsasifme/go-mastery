// Load or store: get a key's value, or add it if the key is missing.
//
// A gateway calls partner APIs. It keeps one client per partner host, and
// makes it the first time that host is needed. Here, several requests for
// the same new host arrive at once.
//
// LoadOrStore does "get it, or add it" in one step. The first caller adds
// its client, and everyone after gets that same client back. Load, then
// Store, would be two steps. Several requests could all see "missing" and
// each add their own (see mutex/03-check-then-act).
//
// Problem: you must build the client before you call LoadOrStore. So every
// request builds one, but only one is kept. The rest are thrown away. A
// real client opens connections, so those are wasted.
//
// Fix: store a sync.OnceValue instead (see once/02-once-value). It's a
// small func that builds the client the first time it's called. Making the
// func costs almost nothing. Only the stored one is ever called, so only
// one client is built, and the other requests wait for it. It also calls
// Load first, so once the host is stored, no new func is made at all. Use
// the "fix:" line instead of the line above it.
//
// Real use:
//   - one HTTP or gRPC client per upstream host, one DB pool per tenant,
//     or one compiled template per name, made on first use.
//
// Run: go run ./02-concurrency/13-sync-map/02-load-or-store
// Race check: go run -race ./02-concurrency/13-sync-map/02-load-or-store
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// client stands in for an HTTP client with its own connection pool.
type client struct {
	id   int64
	host string
}

// gateway keeps one client per host.
type gateway struct {
	clients sync.Map     // host -> *client (fix: host -> func() *client)
	built   atomic.Int64 // how many clients were built, kept or not
}

// newClient builds a client. It's slow, like opening connections.
func (g *gateway) newClient(host string) *client {
	time.Sleep(buildTime)
	return &client{id: g.built.Add(1), host: host}
}

// Client returns the client for host. Every caller builds one first, even
// when a client is already stored.
func (g *gateway) Client(host string) *client {
	c := g.newClient(host)
	v, loaded := g.clients.LoadOrStore(host, c)
	if loaded {
		fmt.Printf("  client #%d built, then thrown away\n", c.id) // another was stored first
	} else {
		fmt.Printf("  client #%d built and kept\n", c.id)
	}
	return v.(*client)
}

// ClientOnce returns the client for host. Only the stored builder is
// called, so one client is built, and the other callers wait for it.
func (g *gateway) ClientOnce(host string) *client {
	// Fast path: the host is already stored, so nothing new is made. Most
	// calls stop here.
	if v, ok := g.clients.Load(host); ok {
		fmt.Println("  builder already stored, nothing new made")
		return v.(func() *client)()
	}

	// Slow path: the host is missing. Several requests can get here at
	// once, so use LoadOrStore, not Store. The first build is kept, and the
	// others are thrown away without ever being called.
	build := sync.OnceValue(func() *client { return g.newClient(host) })
	v, loaded := g.clients.LoadOrStore(host, build)
	if loaded {
		// This request's builder is dropped. Nobody calls it, so the
		// newClient inside it never runs. c below is the winner's client.
		fmt.Println("  my builder wasn't stored, so newClient never ran for it")
	}
	c := v.(func() *client)()
	if !loaded {
		fmt.Printf("  client #%d built and kept\n", c.id)
	}
	return c
}

// keptCount returns how many clients the map holds, one per host.
func (g *gateway) keptCount() (kept int64) {
	g.clients.Range(func(_, _ any) bool {
		kept++
		return true // keep going to the next host
	})
	return kept
}

// buildTime is how long building one client takes.
const buildTime = 50 * time.Millisecond

func main() {
	const (
		requestCount = 5
		host         = "api.partner.com"
	)

	g := &gateway{}

	var wg sync.WaitGroup
	for id := range requestCount {
		wg.Go(func() { handle(g, id+1, host) })
	}
	wg.Wait()

	built, kept := g.built.Load(), g.keptCount()
	fmt.Println("clients built:", built, "- kept:", kept, "- thrown away:", built-kept)
}

// handle is one request that needs the client for host.
func handle(g *gateway, request int, host string) {
	c := g.Client(host) // problem: builds a client even when one is already stored
	// c := g.ClientOnce(host) // fix: only the stored builder is called
	fmt.Printf("request %d uses client #%d\n", request, c.id)
}
