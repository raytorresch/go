package main

import (
	"fmt"
	"tier-0/pointers/types"
)

func main() {
	c := types.Counter{N: 0}
	c.IncrementValue()
	fmt.Println(c.N) // Output: 0

	// method usually used to read
	c.IncrementPointer() // pointer indirection -  solved by go compiler
	fmt.Println(c.N)     // Output: 1

	// method usually used to mutate
	(&c).IncrementPointer() // pointer
	fmt.Println(c.N)        // Output: 2

	// counters := map[string]types.Counter{
	//	"a": {N: 0},
	// }
	// this will not compile because IncrementPointer() has a pointer receiver,
	// but counters["a"] is a value of type Counter, not *Counter.
	// You need to take the address of the value in the map to call the method with a pointer receiver.
	// counters["a"].IncrementPointer() - function expects a pointer receiver, but counters["a"] is a value of type Counter
	// fmt.Println(counters["a"].N)

	counters := map[string]*types.Counter{
		"a": {N: 0},
	}

	counters["a"].IncrementPointer() // pointer indirection -  solved by go compiler
	fmt.Println(counters["a"].N)     // Output: 1
}
