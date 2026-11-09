// Random pick: when several cases are ready, select picks one at random.
//
// Both channels always have a value waiting. main runs select 1000 times
// and counts which case won. The split is close to 50/50 and changes on
// every run. The order of the cases in the code doesn't matter.
//
// Note: never count on select picking a case first. If one case must
// win, see 09-priority.
//
// Real use:
//   - fairness: under load, one busy channel can't starve the others,
//     because every ready case gets picked sometimes.
//
// Run: go run ./02-concurrency/03-select/02-random-pick
package main

import "fmt"

func main() {
	a := make(chan int, 1)
	b := make(chan int, 1)
	wins := map[string]int{}

	for range 1000 {
		a <- 1
		b <- 1 // now both cases are ready

		select {
		case <-a:
			wins["a"]++
			<-b // empty the other channel for the next round
		case <-b:
			wins["b"]++
			<-a
		}
	}

	fmt.Println("a won:", wins["a"], "b won:", wins["b"])
}
