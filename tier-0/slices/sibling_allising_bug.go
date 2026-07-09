package main

import "fmt"

func procesarLote(lote []int) []int {
	// add an element to the slice, which may cause a reallocation if the capacity is exceeded
	return append(lote, -1)
}

func siblingAliasingBug() {
	//shared buffer, typical of a batch processing pipeline.
	buffer := make([]int, 6)
	for i := range buffer {
		buffer[i] = i + 1
	}
	fmt.Println("buffer original:", buffer) // [1 2 3 4 5 6]

	//two "batches" derived from the SAME buffer, they share the underlying array.
	lote1 := buffer[0:3] // [1 2 3], cap = 6 (legacy buffer cap)
	lote2 := buffer[3:6] // [4 5 6]

	fmt.Println("lote1 antes:", lote1, "cap:", cap(lote1))
	fmt.Println("lote2 antes:", lote2)

	resultado1 := procesarLote(lote1)

	fmt.Println("resultado1:", resultado1)
	fmt.Println("lote2 DESPUES de procesar lote1:", lote2)
}
