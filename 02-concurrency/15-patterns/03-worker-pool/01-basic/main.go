// Basic: a fixed number of workers share a queue of jobs.
//
// An upload service makes thumbnails. Images go into the jobs channel.
// workerCount workers take them one at a time, make the thumbnail, and
// send a result. However many images there are, only workerCount run at
// once.
//
// Who closes what:
//   - queueImages closes jobs after the last image, so the workers' range
//     loops end.
//   - closeWhenDone closes results after every worker has returned, so
//     main's range loop ends (see waitgroup/03-close-results).
//
// Workers sharing one jobs channel is fan-out (see fan-out/01-basic).
// Every result going to one results channel is fan-in (see
// fan-in/01-merge). A worker pool puts them together with a fixed number
// of workers.
//
// Try: set workerCount to 1. One worker does every job, and it takes
// longer.
//
// Real use:
//   - processing uploads, sending emails, or resizing images with a set
//     number of workers, so a burst of work doesn't start thousands of
//     goroutines.
//   - a queue consumer (SQS, RabbitMQ, Kafka) with N workers reading from
//     one channel of messages.
//
// Run: go run ./02-concurrency/15-patterns/03-worker-pool/01-basic
package main

import (
	"fmt"
	"sync"
	"time"
)

// result says which worker made which thumbnail.
type result struct {
	image  string
	worker int
}

const (
	thumbTime = 100 * time.Millisecond
	roundTo   = 10 * time.Millisecond
)

func main() {
	const workerCount = 3

	images := []string{"cat.jpg", "dog.jpg", "sun.png", "sea.png", "car.jpg", "map.png"}
	start := time.Now()

	jobs := make(chan string)    // queueImages -> workers
	results := make(chan result) // workers -> main

	go queueImages(images, jobs)

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(id+1, jobs, results, &wg)
	}

	go closeWhenDone(&wg, results)

	for r := range results {
		fmt.Printf("worker %d made a thumbnail of %s\n", r.worker, r.image)
	}
	fmt.Println("all thumbnails done in", time.Since(start).Round(roundTo))
}

// queueImages sends every image to jobs, then closes it.
func queueImages(images []string, jobs chan<- string) {
	defer close(jobs)
	for _, img := range images {
		jobs <- img
	}
}

// worker makes thumbnails until jobs is closed.
func worker(id int, jobs <-chan string, results chan<- result, wg *sync.WaitGroup) {
	defer wg.Done()
	for img := range jobs {
		time.Sleep(thumbTime) // simulate making the thumbnail
		results <- result{image: img, worker: id}
	}
}

// closeWhenDone waits for every worker, then closes results.
func closeWhenDone(wg *sync.WaitGroup, results chan<- result) {
	wg.Wait()
	close(results)
}
