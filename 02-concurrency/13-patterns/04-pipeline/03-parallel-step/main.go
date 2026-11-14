// Parallel step: give the slow step more workers.
//
// A delivery app turns addresses into map points before saving them:
//
//	readAddresses -> geocode (slow, several workers) -> main saves
//
// A pipeline is only as fast as its slowest step. Here geocode calls a
// slow map API, so it gets geocodeWorkers goroutines. They all read from
// the same channel (fan-out) and write to the same channel. The last
// worker to finish closes it, using a WaitGroup.
//
// The results come out in a different order than they went in. If order
// matters, carry an index with each item and sort at the end.
//
// Note: compare with 01-basic, where every step had one goroutine.
//
// Try: set geocodeWorkers to 1. It takes about as long as all the API
// calls added up.
//
// Real use:
//   - a pipeline where one step calls a slow API or does heavy CPU work,
//     like geocoding, image resizing or hashing.
//
// Run: go run ./02-concurrency/13-patterns/04-pipeline/03-parallel-step
package main

import (
	"fmt"
	"sync"
	"time"
)

// point is an address with its map position.
type point struct {
	address string
	lat     float64
	lng     float64
}

const (
	geocodeTime = 100 * time.Millisecond
	roundTo     = 10 * time.Millisecond
)

// coords is the fake map API's data: address -> lat, lng.
var coords = map[string][2]float64{
	"Gulshan 2":    {23.7925, 90.4078},
	"Dhanmondi 27": {23.7561, 90.3742},
	"Banani 11":    {23.7937, 90.4066},
	"Uttara 7":     {23.8706, 90.3985},
	"Mirpur 10":    {23.8069, 90.3687},
	"Motijheel":    {23.7330, 90.4172},
}

func main() {
	const geocodeWorkers = 3

	start := time.Now()

	addresses := make(chan string) // readAddresses -> geocode workers
	points := make(chan point)     // geocode workers -> main

	go readAddresses(addresses)

	var wg sync.WaitGroup
	wg.Add(geocodeWorkers)
	for range geocodeWorkers {
		go geocode(addresses, points, &wg)
	}
	go closeWhenDone(&wg, points)

	for p := range points {
		fmt.Printf("saved %-13s %.4f, %.4f\n", p.address, p.lat, p.lng)
	}
	fmt.Println("all addresses saved in", time.Since(start).Round(roundTo))
}

// readAddresses sends every address, then closes out.
func readAddresses(out chan<- string) {
	defer close(out)
	for address := range coords {
		out <- address
	}
}

// geocode looks up each address on the fake map API.
func geocode(in <-chan string, out chan<- point, wg *sync.WaitGroup) {
	defer wg.Done()
	for address := range in {
		time.Sleep(geocodeTime) // simulate the API call
		c := coords[address]
		out <- point{address: address, lat: c[0], lng: c[1]}
	}
}

// closeWhenDone waits for every geocode worker, then closes out.
func closeWhenDone(wg *sync.WaitGroup, out chan<- point) {
	wg.Wait()
	close(out)
}
