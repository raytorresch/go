package domain

import "math"

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 { // Implement the Area method for Circle
	return math.Pi * math.Pow(c.Radius, 2)
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 { // Implement the Area method for Rectangle
	return r.Width * r.Height
}

func Less(s1, s2 Sizer) Sizer {
	if s1.Area() < s2.Area() {
		return s1
	}
	return s2
}
