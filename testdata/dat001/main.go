package dat001

import "fmt"

type (
	Cat struct {
		Name string
		Age  int
	}

	Toy interface {
		Play(c *Cat)
	}

	Cats []*Cat
)

func NewCat(name string, age int) *Cat {
	return &Cat{
		Name: name,
		Age:  age,
	}
}

func (c *Cat) Meow() {
	fmt.Println(c.Name, `meows`)
}

func (c *Cat) String() string {
	return c.Name
}

func NextYear(cats []*Cat) {
	for _, cat := range cats {
		cat.Age++
	}
}

func (cats Cats) Youngest() *Cat {
	var youngest *Cat
	for _, c := range cats {
		if youngest == nil || youngest.Age < c.Age {
			youngest = c
		}
	}
	return youngest
}

func Pet(c *Cat) {
	fmt.Println(`petting`, c.Name)
}
