package main

import "fmt"

func alising() {
	s := make([]int, 0)
	prevCap := cap(s)
	for i := 0; i < 15; i++ {
		s = append(s, i)
		if cap(s) != prevCap {
			fmt.Printf("len=%d  cap=%d  (crecio de %d)\n", len(s), cap(s), prevCap)
			prevCap = cap(s)
		}
	}
}
