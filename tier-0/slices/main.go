package main

import "fmt"

func mutate(s []int) {
	s[0] = 999
}

func add(s []int) {
	s = append(s, 999) // this make a new slice with a new underlying array, so the original slice is not modified
	s[0] = 111         // modification to the new slice does not affect the original slice
}

func mutateSuccessfully() {
	original := []int{1, 2, 3}
	mutate(original)      // pass "value" of slice, but slice is a reference type, so the underlying array is modified
	fmt.Println(original) // Output: [999 2 3]
}

func failedMutateUncapSlice() {
	orgNoCap := make([]int, 3) // create a slice with length 3 and capacity 3
	orgNoCap[0], orgNoCap[1], orgNoCap[2] = 1, 2, 3
	add(orgNoCap)

	fmt.Println(orgNoCap) // Output: [1 2 3] - the underlying array is not modified because the slice was reallocated in the add function
}

func main() {
	// mutateSuccessfully()
	// failedMutateUncapSlice()

	orgNoCap := make([]int, 3, 10) // create a slice with length 3 and capacity 3
	orgNoCap[0], orgNoCap[1], orgNoCap[2] = 1, 2, 3
	add(orgNoCap)

	fmt.Println(orgNoCap)
}
