package main

import "fmt"

func agregar(s []int) {
	fmt.Printf("antes del append -> ptr=%p len=%d cap=%d\n", s, len(s), cap(s))
	s = append(s, 999)
	fmt.Printf("despues del append -> ptr=%p len=%d cap=%d\n", s, len(s), cap(s))
	s[0] = 111
}

func capacityCheck() {
	original := make([]int, 3, 10) // len=3, cap=10
	original[0], original[1], original[2] = 1, 2, 3
	fmt.Printf("original antes -> ptr=%p len=%d cap=%d\n", original, len(original), cap(original))
	agregar(original)
	fmt.Println("original despues (len=3, lo unico 'visible'):", original)

	// re-slice to "see" original all lenngth,
	// but also up to the real capacity of the underlying array.
	extended := original[:cap(original)]
	fmt.Println("mismo array, pero re-sliced hasta cap:", extended)
}
