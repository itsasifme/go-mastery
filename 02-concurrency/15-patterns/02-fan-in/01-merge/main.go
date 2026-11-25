// Merge: combine several channels into one.
//
// Three services each send log lines on their own channel, at their own
// speed. merge starts one forward goroutine per input. Each one copies
// lines from its input to the shared output. When every input is closed
// and every forward has returned, closeWhenDone closes the output. Then
// main's range loop ends.
//
// Lines from different services mix together in the order they arrive.
// Lines from one service stay in their own order.
//
// Real use:
//   - collecting logs or events from several sources into one stream,
//     like "kubectl logs" for many pods.
//   - reading from several queues or Kafka partitions and handling all
//     messages in one loop.
//
// Run: go run ./02-concurrency/15-patterns/02-fan-in/01-merge
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	lineCount = 3 // lines each service sends
	apiEvery  = 30 * time.Millisecond
	dbEvery   = 50 * time.Millisecond
	authEvery = 70 * time.Millisecond
)

func main() {
	api := make(chan string)
	db := make(chan string)
	auth := make(chan string)

	go emitLogs("api", apiEvery, api)
	go emitLogs("db", dbEvery, db)
	go emitLogs("auth", authEvery, auth)

	for line := range merge(api, db, auth) {
		fmt.Println(line)
	}
	fmt.Println("all services finished")
}

// emitLogs sends lineCount log lines, one every interval, then closes out.
func emitLogs(service string, every time.Duration, out chan<- string) {
	defer close(out)
	for i := range lineCount {
		time.Sleep(every)
		out <- fmt.Sprintf("[%s] line %d", service, i+1)
	}
}

// merge returns one channel that gets every value from inputs. It's
// closed after every input is closed.
func merge(inputs ...<-chan string) <-chan string {
	out := make(chan string)

	var wg sync.WaitGroup
	wg.Add(len(inputs))
	for _, in := range inputs {
		go forward(in, out, &wg)
	}
	go closeWhenDone(&wg, out)

	return out
}

// forward copies every value from in to out, until in is closed.
func forward(in <-chan string, out chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range in {
		out <- v
	}
}

// closeWhenDone waits for every forward, then closes out.
func closeWhenDone(wg *sync.WaitGroup, out chan<- string) {
	wg.Wait()
	close(out)
}
