// Errors: stop the whole pipeline when one step fails.
//
// A product import reads a CSV file, parses each row, and saves it:
//
//	readRows -> parse -> save
//
// Each step runs with errgroup's g.Go. One row has a bad price. So parse
// returns an error, and errgroup cancels ctx. readRows sees ctx.Done() and
// stops sending. save stops when parse closes its output. g.Wait returns
// the parse error.
//
// The import is all or nothing. save writes into a transaction. main
// commits only when g.Wait returns nil. Otherwise it rolls back, so half
// an import never reaches the database.
//
// Note: compare with 02-cancel. There, main stopped the pipeline when it
// had enough. Here, a failing step stops it (see errgroup/03-with-context).
//
// Try: change "N/A" in csvFile to a number. Every row is saved and the
// import is committed.
//
// Real use:
//   - importing a CSV or JSON file in one transaction: db.BeginTx, then
//     tx.Commit if every row worked, or tx.Rollback if one failed.
//   - a report built from several steps, where one failed step makes the
//     whole report wrong.
//
// Run: go run ./02-concurrency/15-patterns/04-pipeline/04-errors
package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

// product is one parsed row.
type product struct {
	name  string
	price int // in cents
}

// tx stands in for a database transaction like *sql.Tx.
type tx struct {
	pending []product // saved, but not committed yet
}

// Insert adds a product to the transaction.
func (t *tx) Insert(p product) {
	time.Sleep(insertTime) // simulate the INSERT
	t.pending = append(t.pending, p)
}

// insertTime is how long one fake INSERT takes.
const insertTime = 20 * time.Millisecond

// csvFile is the sample file to import, one "name,price" per line. One
// price is wrong.
var csvFile = []string{
	"Desk lamp,2599",
	"Chair,4999",
	"Mug,N/A",
	"Notebook,350",
	"Backpack,3999",
	"Pen,120",
}

func main() {
	t := &tx{}
	g, ctx := errgroup.WithContext(context.Background())

	rows := make(chan string)      // readRows -> parse
	products := make(chan product) // parse -> save

	g.Go(func() error { return readRows(ctx, csvFile, rows) })
	g.Go(func() error { return parse(ctx, rows, products) })
	g.Go(func() error { return save(products, t) })

	if err := g.Wait(); err != nil {
		fmt.Println("import failed:", err)
		fmt.Println("rolled back", len(t.pending), "saved products")
		return
	}
	fmt.Println("committed", len(t.pending), "products")
}

// readRows sends each line of the file, until it runs out or ctx is done.
func readRows(ctx context.Context, lines []string, out chan<- string) error {
	defer close(out)
	for i, line := range lines {
		select {
		case out <- line:
		case <-ctx.Done():
			fmt.Printf("readRows: stopped after %d of %d rows\n", i, len(lines))
			return ctx.Err()
		}
	}
	return nil
}

// parse turns each line into a product. It returns an error on the first
// bad row.
func parse(ctx context.Context, in <-chan string, out chan<- product) error {
	defer close(out)

	row := 0
	for line := range in {
		row++
		name, priceText, ok := strings.Cut(line, ",")
		if !ok {
			return fmt.Errorf("row %d: want name,price, got %q", row, line)
		}
		price, err := strconv.Atoi(priceText)
		if err != nil {
			return fmt.Errorf("row %d: bad price %q: %w", row, priceText, err)
		}

		select {
		case out <- product{name: name, price: price}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

// save inserts each product into the transaction, until in is closed.
func save(in <-chan product, t *tx) error {
	for p := range in {
		t.Insert(p)
	}
	fmt.Println("save: input closed after", len(t.pending), "products")
	return nil
}
