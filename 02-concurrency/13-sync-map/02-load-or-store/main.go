// Load or store: get the value for a key, or add it if it's missing, in
// one step.
//
// An API gateway keeps one client per partner host, made on first use.
// Five requests for the same new host arrive together. LoadOrStore stores
// the first client and hands that same client to every other caller
// (loaded is true for them). Load followed by Store would be a
// check-then-act bug (see mutex/03-check-then-act): several requests see
// "missing", and each stores its own client.
//
// Problem: LoadOrStore takes a ready value, so every request builds a
// client first, and only one is kept. The others are thrown away. A real
// client holds open connections, so each thrown-away one is a leak.
//
// Fix: store a sync.OnceValue instead of the client (see
// once/02-once-value). Making the func is cheap. Only the stored one ever
// runs, so exactly one client is built, and the other requests wait for
// it. Use the "fix:" line instead of the line above it.
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
		fmt.Printf("  client #%d thrown away\n", c.id)
	}
	return v.(*client)
}

// ClientOnce returns the client for host. Only the first caller's build
// runs, and the others wait for it.
func (g *gateway) ClientOnce(host string) *client {
	build := sync.OnceValue(func() *client { return g.newClient(host) })
	v, _ := g.clients.LoadOrStore(host, build)
	return v.(func() *client)()
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

	fmt.Println("clients built:", g.built.Load(), "- kept: 1")
}

// handle is one request that needs the client for host.
func handle(g *gateway, request int, host string) {
	c := g.Client(host) // problem: builds a client even when one is already stored
	// c := g.ClientOnce(host) // fix: only the stored build runs
	fmt.Printf("request %d uses client #%d\n", request, c.id)
}
