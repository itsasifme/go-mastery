// Pointer: swap a whole config at once while requests keep reading it.
//
// Handlers read the config on every request. When the config file
// changes, main builds a new config and swaps it in with Store. Handlers
// call Load and get either the old config or the new one, never a mix of
// both. No lock, so readers never wait.
//
// The rule: never change a config after it's stored. Build a new one and
// Store it. Handlers may still be using the old one.
//
// Real use:
//   - reloading a config file on SIGHUP without restarting the server.
//   - swapping a routing table or a TLS certificate while requests are
//     being served (tls.Config's GetCertificate often does this).
//
// Run: go run ./02-concurrency/09-atomic/03-pointer
// Race check: go run -race ./02-concurrency/09-atomic/03-pointer
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// config is never changed after it's stored. A reload makes a new one.
type config struct {
	version   int
	rateLimit int // requests per minute per user
	banner    string
}

// requestTime is how long one fake request takes.
const requestTime = 20 * time.Millisecond

func main() {
	const (
		handlerCount = 2
		reloadAfter  = 100 * time.Millisecond
		runTime      = 200 * time.Millisecond
	)

	var current atomic.Pointer[config]
	current.Store(&config{version: 1, rateLimit: 60, banner: "Welcome!"})

	stop := make(chan struct{}) // main -> handlers: "please stop"

	var wg sync.WaitGroup
	wg.Add(handlerCount)
	for id := range handlerCount {
		go handle(id+1, &current, stop, &wg)
	}

	time.Sleep(reloadAfter)
	fmt.Println("main: config file changed, reloading")
	current.Store(&config{version: 2, rateLimit: 120, banner: "Big sale today!"})

	time.Sleep(runTime - reloadAfter)
	close(stop)
	wg.Wait()
}

// handle serves requests until stop is closed. It counts how many requests
// used each config version.
func handle(id int, current *atomic.Pointer[config], stop <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	perVersion := map[int]int{} // config version -> requests served with it
	var last *config
	for {
		select {
		case <-time.After(requestTime):
			last = current.Load() // one Load per request: a request sees one config from start to end
			perVersion[last.version]++
		case <-stop:
			fmt.Printf("handler %d: requests per config version %v, banner now %q\n", id, perVersion, last.banner)
			return
		}
	}
}
