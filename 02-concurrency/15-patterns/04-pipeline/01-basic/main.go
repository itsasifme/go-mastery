// Basic: a chain of steps connected by channels.
//
// An order report runs in three steps, each in its own goroutine:
//
//	generate -> addTax -> format -> main prints
//
// Each step reads from the channel before it and writes to the one after.
// While format works on order 1, addTax can already work on order 2. When
// a step's input is closed, it closes its own output. The close moves down
// the chain, and main's range loop ends.
//
// main makes the channels and starts each step, so you can see how they
// connect.
//
// Real use:
//   - ETL jobs: read rows from a CSV or database, clean them, then write
//     them somewhere else.
//   - image processing: decode -> resize -> watermark -> upload.
//
// Run: go run ./02-concurrency/15-patterns/04-pipeline/01-basic
package main

import "fmt"

// order is one line in the report.
type order struct {
	id       string
	subtotal int // in cents
	tax      int // in cents
}

const (
	taxPercent = 15  // sales tax added to every order
	percent    = 100 // taxPercent is out of this
)

func main() {
	orders := []order{
		{id: "A-1", subtotal: 1000},
		{id: "A-2", subtotal: 2500},
		{id: "A-3", subtotal: 400},
	}

	raw := make(chan order)    // generate -> addTax
	taxed := make(chan order)  // addTax -> format
	lines := make(chan string) // format -> main

	go generate(orders, raw)
	go addTax(raw, taxed)
	go format(taxed, lines)

	for line := range lines {
		fmt.Println(line)
	}
}

// generate sends every order, then closes out.
func generate(orders []order, out chan<- order) {
	defer close(out)
	for _, o := range orders {
		out <- o
	}
}

// addTax fills in the tax for each order.
func addTax(in <-chan order, out chan<- order) {
	defer close(out)
	for o := range in {
		o.tax = o.subtotal * taxPercent / percent
		out <- o
	}
}

// format turns each order into a line of text.
func format(in <-chan order, out chan<- string) {
	defer close(out)
	for o := range in {
		out <- fmt.Sprintf("%s  subtotal %5d  tax %4d  total %5d", o.id, o.subtotal, o.tax, o.subtotal+o.tax)
	}
}
