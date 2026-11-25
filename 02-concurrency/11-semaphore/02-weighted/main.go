// Weighted: jobs that take different amounts of a shared limit.
//
// A video server has a memory budget. A short clip needs 1 unit, a movie
// needs 3, and a 4K movie needs the whole budget. semaphore.Weighted lets
// each job take as many units as it needs. Acquire blocks until there's
// enough room. Release gives the units back.
//
// Weighted is fair: jobs get their units in the order they asked. A big
// job that's waiting isn't passed by small jobs behind it, so it never
// waits forever.
//
// Note: compare with 01-channel. There, every job took exactly one slot.
//
// Real use:
//   - limiting total memory used by image or video processing, when you
//     know each job's size before it starts.
//   - limiting total bytes being uploaded at once, not just the number of
//     uploads.
//
// Run: go run ./02-concurrency/11-semaphore/02-weighted
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// job is one video to transcode and the memory units it needs.
type job struct {
	name  string
	units int64
}

const (
	memoryBudget = 4                      // total units the server can use at once
	timePerUnit  = 100 * time.Millisecond // bigger jobs take longer
	roundTo      = 10 * time.Millisecond
)

func main() {
	jobs := []job{
		{"clip-1", 1},
		{"movie-1", 3},
		{"clip-2", 1},
		{"4k-movie", memoryBudget},
		{"clip-3", 1},
	}

	sem := semaphore.NewWeighted(memoryBudget)
	ctx := context.Background()
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(jobs))
	for _, j := range jobs {
		if err := sem.Acquire(ctx, j.units); err != nil { // blocks until j.units are free
			fmt.Println("acquire:", err)
			wg.Done()
			continue
		}
		go transcode(j, start, sem, &wg)
	}
	wg.Wait()

	fmt.Println("all videos done in", time.Since(start).Round(roundTo))
}

// transcode pretends to convert one video, then gives its units back.
func transcode(j job, start time.Time, sem *semaphore.Weighted, wg *sync.WaitGroup) {
	defer wg.Done()
	defer sem.Release(j.units)

	fmt.Printf("%-6v start %-8s (%d units)\n", time.Since(start).Round(roundTo), j.name, j.units)
	time.Sleep(time.Duration(j.units) * timePerUnit)
}
