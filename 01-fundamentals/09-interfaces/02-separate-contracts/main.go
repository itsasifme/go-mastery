// Separate contracts: one type can satisfy several small interfaces.
//
// author has details() and articles(), so it satisfies both AuthorDetails
// and AuthorArticles. Each variable only sees the methods of its own
// interface.
//
// Real use:
//   - io.Reader and io.Writer: an os.File is both, but a function that
//     only reads asks for a Reader.
//   - a function asks for the smallest interface it needs, so it's easy
//     to test with a tiny fake.
//
// Run: go run ./01-fundamentals/09-interfaces/02-separate-contracts
package main

import "fmt"

// AuthorDetails is any type that can print author details.
type AuthorDetails interface {
	details()
}

// AuthorArticles is any type that can print pending articles.
type AuthorArticles interface {
	articles()
}

// author has both methods, so it fits both interfaces.
type author struct {
	name              string
	branch            string
	college           string
	year              int
	salary            int
	publishedArticles int
	totalArticles     int
}

// details prints the author's info. AuthorDetails asks for this method.
func (a author) details() {

	fmt.Printf("Author Name: %s", a.name)
	fmt.Printf("\nBranch: %s and passing year: %d", a.branch, a.year)
	fmt.Printf("\nCollege Name: %s", a.college)
	fmt.Printf("\nSalary: %d", a.salary)
	fmt.Printf("\nPublished articles: %d", a.publishedArticles)

}

// articles prints how many articles are still pending. AuthorArticles
// asks for this method.
func (a author) articles() {

	pendingArticles := a.totalArticles - a.publishedArticles
	fmt.Printf("\n\nPending articles: %d", pendingArticles)
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

	// i1 only sees details()
	var i1 AuthorDetails = author
	i1.details()

	// i2 only sees articles()
	var i2 AuthorArticles = author
	i2.articles()

}
