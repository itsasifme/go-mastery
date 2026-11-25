// Combined contract: one interface that requires several methods.
//
// Author needs both details() and articles(). author has both, so one
// variable can call either.
//
// Real use:
//   - io.ReadWriter and io.ReadCloser combine smaller interfaces. An HTTP
//     response body is an io.ReadCloser.
//   - a Repository interface with Get, Save and Delete for one kind of
//     data.
//
// Run: go run ./01-fundamentals/09-interfaces/03-combined-contract
package main

import "fmt"

// Author needs both details() and articles().
type Author interface {
	details()
	articles()
}

// author has both methods, so it fits Author.
type author struct {
	name              string
	branch            string
	college           string
	year              int
	salary            int
	publishedArticles int
	totalArticles     int
}

// details prints the author's info.
func (a author) details() {

	fmt.Printf("Author Name: %s", a.name)
	fmt.Printf("\nBranch: %s and passing year: %d", a.branch, a.year)
	fmt.Printf("\nCollege Name: %s", a.college)
	fmt.Printf("\nSalary: %d", a.salary)
	fmt.Printf("\nPublished articles: %d", a.publishedArticles)

}

// articles prints how many articles are still pending.
func (a author) articles() {

	pendingArticles := a.totalArticles - a.publishedArticles
	fmt.Printf("\nPending articles: %d", pendingArticles)
}

func main() {

	// Fill in an author by field name
	author := author{
		name:              "Mickey",
		branch:            "Computer science",
		college:           "XYZ",
		year:              2012,
		salary:            50000,
		publishedArticles: 209,
		totalArticles:     309,
	}

	// One Author variable can call both methods
	var i Author = author
	i.details()
	i.articles()
}
