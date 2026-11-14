// Limit: run only a few jobs at a time with SetLimit.
//
// Six photos need resizing, but resizing uses a lot of memory, so only
// maxParallel can run at once. After g.SetLimit, g.Go waits until a slot
// is free before it starts the next job. The start times show the photos
// going in groups.
//
// Note: this is a semaphore with error handling built in (see
// semaphore/01-channel).
//
// Try: set maxParallel to len(photos). Every photo starts at 0s.
//
// Real use:
//   - resizing uploaded images or videos without running out of memory.
//   - calling a partner API that allows only a few requests at once.
//
// Run: go run ./02-concurrency/06-errgroup/02-limit
package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	resizeTime = 100 * time.Millisecond
	roundTo    = 10 * time.Millisecond
)

func main() {
	const maxParallel = 2

	photos := []string{"a.jpg", "b.jpg", "c.jpg", "d.jpg", "e.jpg", "f.jpg"}
	start := time.Now()

	var g errgroup.Group
	g.SetLimit(maxParallel)
	for _, photo := range photos {
		g.Go(func() error { return resize(photo, start) }) // waits here while every slot is busy
	}

	if err := g.Wait(); err != nil {
		fmt.Println("resize failed:", err)
		return
	}
	fmt.Println("all photos resized in", time.Since(start).Round(roundTo))
}

// resize pretends to resize one photo, and prints when it started.
func resize(photo string, start time.Time) error {
	fmt.Printf("%-6v start %s\n", time.Since(start).Round(roundTo), photo)
	time.Sleep(resizeTime)
	return nil
}
