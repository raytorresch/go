// Package main demonstrates idiomatic Go interfaces (implicit satisfaction,
// small interfaces) with a minimal domain layer (VOs + port).
package main

import (
	"fmt"
	"tier-0/interfaces/domain"
)

func main() {
	c := domain.Circle{Radius: 5}
	r := domain.Rectangle{Width: 4, Height: 6}

	l := domain.Less(c, r)
	fmt.Printf("The shape with the smaller area is: %T with area %.2f\n", l, l.Area())
}
