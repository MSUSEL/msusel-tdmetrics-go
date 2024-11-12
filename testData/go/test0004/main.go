//go:build test

package main

type (
	Animal interface {
		Age() int
	}

	Bird interface {
		Animal
		Fly()
	}

	Mammal interface {
		Animal
		Temp() float64
	}

	Bat interface {
		Mammal
		Fly()
	}

	Flier interface {
		Fly()
	}
)

var (
	_ Animal = Bird(nil)
	_ Flier  = Bat(nil)
)

func main() {
	println(`okay`)
}
