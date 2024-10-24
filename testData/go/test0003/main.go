//go:build test

package main

type Animal interface {
	Pet()
}

type Cat struct {
	Name string
}

func (c *Cat) Pet() {
	println(`Petting`, c.Name)
}

func main() {
	c := &Cat{
		Name: `Mittens`,
	}
	c.Pet()
}
