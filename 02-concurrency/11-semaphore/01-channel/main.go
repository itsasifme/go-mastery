// Channel: a buffered channel used as a semaphore.
//
// A downloader has many files to get, but the server allows only
// maxParallel connections per client. sem has that many slots. main puts a
// value in before it starts each download, and the download takes it out
// when it ends. When sem is full, main blocks until a download finishes.
//
// main takes the slot, not the goroutine. That way only maxParallel
// goroutines exist at a time, instead of all of them waiting at once.
//
// Note: compare with errgroup/02-limit. SetLimit does the same, and also
// returns the first error.
//
// Try: set maxParallel to 1. The files download one by one.
//
// Real use:
//   - limiting http.Get calls to one host, so it doesn't block you.
//   - limiting how many files are open at once, to stay under the OS limit.
//
// Run: go run ./02-concurrency/11-semaphore/01-channel
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	downloadTime = 100 * time.Millisecond
	roundTo      = 10 * time.Millisecond
)

func main() {
	const maxParallel = 3

	files := []string{"1.zip", "2.zip", "3.zip", "4.zip", "5.zip", "6.zip", "7.zip"}
	sem := make(chan struct{}, maxParallel) // one slot per download allowed at once
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		sem <- struct{}{} // take a slot. This blocks while all slots are taken
		go download(file, start, sem, &wg)
	}
	wg.Wait()

	fmt.Println("all files downloaded in", time.Since(start).Round(roundTo))
}

// download pretends to download one file, then gives its slot back.
func download(file string, start time.Time, sem <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer release(sem) // give the slot back, even on a panic

	fmt.Printf("%-6v start %s\n", time.Since(start).Round(roundTo), file)
	time.Sleep(downloadTime)
}

// release gives one slot back to the semaphore.
func release(sem <-chan struct{}) {
	<-sem
}
