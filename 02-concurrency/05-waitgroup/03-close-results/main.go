// Close results: many workers send to one channel, close it when all are done.
//
// A health checker checks several sites at once. Each check runs in its own
// goroutine and sends its result on one channel. main reads with range, so
// someone must close the channel, but only after every check is finished.
// closeWhenDone does that: it waits on the WaitGroup, then closes results.
//
// Results arrive in the order the checks finish, not the order of sites.
//
// Note: compare with channels/02-closing/03-many-senders. There, a done
// channel and a counting loop decided when to close. WaitGroup does the
// counting for you.
//
// Real use:
//   - an uptime checker that runs http.Get on many URLs and prints each
//     status as soon as it comes back.
//   - resizing uploaded images in parallel and streaming each finished
//     file name back to the client.
//
// Run: go run ./02-concurrency/05-waitgroup/03-close-results
package main

import (
	"fmt"
	"sync"
	"time"
)

// result is what one check sends back.
type result struct {
	site   string
	status int
}

// checkTime is how long one fake check takes.
const checkTime = 100 * time.Millisecond

func main() {
	sites := map[string]int{ // site -> the status our fake check returns
		"example.com": 200,
		"shop.dev":    200,
		"old-api.io":  404,
		"slow.net":    503,
	}

	results := make(chan result) // checks -> main: one result each

	var wg sync.WaitGroup
	wg.Add(len(sites)) // one per check, before any goroutine starts
	for site, status := range sites {
		go check(site, status, results, &wg)
	}

	go closeWhenDone(&wg, results)

	for r := range results { // stops when closeWhenDone closes results
		fmt.Printf("%-12s %d\n", r.site, r.status)
	}
	fmt.Println("all sites checked")
}

// check pretends to send a request to site and sends back the status.
func check(site string, status int, out chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(checkTime)
	out <- result{site: site, status: status}
}

// closeWhenDone waits for every check, then closes the channel once.
func closeWhenDone(wg *sync.WaitGroup, out chan<- result) {
	wg.Wait()
	close(out)
}
