package types

type Counter struct {
	N int
}

func (c Counter) IncrementValue() {
	c.N++
}

func (c *Counter) IncrementPointer() {
	c.N++
}
