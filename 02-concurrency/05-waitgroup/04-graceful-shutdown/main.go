// Graceful shutdown: let running jobs finish before the program exits.
//
// Workers keep processing jobs one after another. When ctx times out,
// each worker finishes the job it's on, then returns. It never cuts a job
// in half. main calls wg.Wait, so it doesn't exit while a job is still
// running.
//
// Try: comment out wg.Wait(). main exits right away, and the workers never
// print "stopping".
//
// Note: here a timeout stands in for Ctrl+C. In a real program, ctx comes
// from signal.NotifyContext(context.Background(), os.Interrupt) (see
// context/01-cancellation/02-stop-workers). WithTimeoutCause works like in
// context/02-timeout/03-timeout-cause: context.Cause says why it stopped.
//
// Real use:
//   - a queue consumer that finishes the current message, acks it, then
//     exits, so no message is lost or handled twice.
//   - http.Server.Shutdown(ctx), which waits for requests in progress the
//     same way before the server stops.
//
// Run: go run ./02-concurrency/05-waitgroup/04-graceful-shutdown
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// jobTime is how long one fake job takes.
const jobTime = 150 * time.Millisecond

// errShutdown is why ctx ends. context.Cause returns it.
var errShutdown = errors.New("shutdown requested")

func main() {
	const (
		workerCount = 3
		runTime     = 400 * time.Millisecond // not a multiple of jobTime, so the timeout lands mid-job
		roundTo     = 10 * time.Millisecond
	)

	// After runTime, ctx ends and every worker sees it.
	ctx, cancel := context.WithTimeoutCause(context.Background(), runTime, errShutdown)
	defer cancel() // frees the timer if main returns early

	var wg sync.WaitGroup
	wg.Add(workerCount)
	for id := range workerCount {
		go worker(ctx, id+1, &wg)
	}

	<-ctx.Done() // the program runs until the timeout
	fmt.Println("main:", context.Cause(ctx), "- waiting for jobs in progress")

	shutdownStart := time.Now()
	wg.Wait()
	fmt.Println("main: all workers stopped after", time.Since(shutdownStart).Round(roundTo))
}

// worker processes jobs until ctx is canceled. It checks ctx only between
// jobs, so a job that started always finishes.
func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()

	jobs := 0
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: stopping after %d jobs\n", id, jobs)
			return
		default:
		}

		time.Sleep(jobTime) // simulate one job. It isn't stopped halfway
		jobs++
	}
}
