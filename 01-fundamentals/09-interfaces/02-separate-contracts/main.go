// Go program to illustrate the concept of multiple interfaces

package main

import "fmt"

// Interface 1
type AuthorDetails interface {
	details()
}

// Interface 2
type AuthorArticles interface {
	articles()
}

// Structure
type author struct {
	name              string
	branch            string
	college           string
	year              int
	salary            int
	publishedArticles int
	totalArticles     int
}

// Implementing method of the interface 1
func (a author) details() {

	fmt.Printf("Author Name: %s", a.name)
	fmt.Printf("\nBranch: %s and passing year: %d", a.branch, a.year)
	fmt.Printf("\nCollege Name: %s", a.college)
	fmt.Printf("\nSalary: %d", a.salary)
	fmt.Printf("\nPublished articles: %d", a.publishedArticles)

}

// Implementing method of the interface 2
func (a author) articles() {

	pendingArticles := a.totalArticles - a.publishedArticles
	fmt.Printf("\n\nPending articles: %d", pendingArticles)
}

func main() {

	// Assigning values to the structure
	author := author{
		name:              "Mickey",
		branch:            "Computer science",
		college:           "XYZ",
		year:              2012,
		salary:            50000,
		publishedArticles: 209,
		totalArticles:     309,
	}

	// Accessing the method of the interface 1
	var i1 AuthorDetails = author
	i1.details()

	// Accessing the method of the interface 2
	var i2 AuthorArticles = author
	i2.articles()

}
