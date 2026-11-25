// With context: stop the other jobs when one fails.
//
// errgroup.WithContext gives a ctx that is cancelled as soon as any job
// returns an error. Three regions download their part of a report in
// chunks. The "eu" download fails early, so ctx is cancelled and the other
// downloads stop. They don't waste time on a report that can't be built.
//
// Wait returns the first error (from eu), not the context.Canceled that
// the other jobs return after it.
//
// Note: compare with 01-basic. There, the other jobs kept running until
// they finished, and Wait waited for all of them.
//
// Real use:
//   - building one report from several regions with
//     http.NewRequestWithContext(ctx, ...). If one region fails, the other
//     requests are cancelled.
//   - uploading a big file in parts. One part fails, so stop the others.
//
// Run: go run ./02-concurrency/06-errgroup/03-with-context
package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	chunkTime  = 50 * time.Millisecond
	chunkCount = 10 // chunks in a full download
	neverFail  = -1 // failAt value for a download that works
)

func main() {
	const euFailAt = 3 // eu fails when it reaches this chunk

	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return download(ctx, "us", neverFail) })
	g.Go(func() error { return download(ctx, "eu", euFailAt) })
	g.Go(func() error { return download(ctx, "asia", neverFail) })

	if err := g.Wait(); err != nil {
		fmt.Println("report failed:", err)
		return
	}
	fmt.Println("report ready")
}

// download fetches one region's data chunk by chunk. It fails when it
// reaches chunk failAt, and stops early if ctx is cancelled.
func download(ctx context.Context, region string, failAt int) error {
	for chunk := range chunkCount {
		if chunk == failAt {
			return fmt.Errorf("%s: server returned 500 at chunk %d", region, chunk)
		}

		select {
		case <-time.After(chunkTime): // simulate downloading one chunk
		case <-ctx.Done():
			fmt.Printf("%s: stopped after %d of %d chunks\n", region, chunk, chunkCount)
			return ctx.Err()
		}
	}
	fmt.Println(region, "done")
	return nil
}
