// Go program to illustrate the
// concept of multiple interfaces
package main

import "fmt"

// Interface
type Author interface {
	details()
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

// Implementing details() method of the interface
func (a author) details() {

	fmt.Printf("Author Name: %s", a.name)
	fmt.Printf("\nBranch: %s and passing year: %d", a.branch, a.year)
	fmt.Printf("\nCollege Name: %s", a.college)
	fmt.Printf("\nSalary: %d", a.salary)
	fmt.Printf("\nPublished articles: %d", a.publishedArticles)

}

// Implementing articles() method of the interface
func (a author) articles() {

	pendingArticles := a.totalArticles - a.publishedArticles
	fmt.Printf("\nPending articles: %d", pendingArticles)
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

	// Accessing the methods of the interface
	var i Author = author
	i.details()
	i.articles()
}
