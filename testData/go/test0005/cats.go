//go:build test

package cats

type (
	// Cat represents an evil murder machine.
	Cat struct {
		Name string
		Age  int
	}

	// Toy represents anything to distract cats from knocking things off of counters.
	Toy interface {
		// Play uses this toy on the given cat.
		Play(c *Cat)
	}

	// Cats is a pride of murder machines.
	Cats []*Cat
)

// log will print messages, overwrite to log out to a different place.
var log = func(value string) { println(value) }

// NewCat creates a new cat instance with the given name and age.
func NewCat(name string, age int) *Cat {
	return &Cat{
		Name: name,
		Age:  age,
	}
}

// Meow is the cat's way to tell its human servant that it
// must have food now or the servant will suffer consequences.
func (c *Cat) Meow() {
	log(c.Name + ` meows`)
}

// String gets the name of the cat.
func (c *Cat) String() string {
	return c.Name
}

// NextYear increments the age of all the given cats.
func NextYear(cats ...*Cat) {
	for _, cat := range cats {
		cat.Age++
	}
}

// Youngest finds the first of the youngest cats in the given list of cats.
func (cats Cats) Youngest() *Cat {
	var youngest *Cat
	for _, c := range cats {
		if youngest == nil || youngest.Age > c.Age {
			youngest = c
		}
	}
	return youngest
}

// Pet puts the human servant in danger of being
// murdered while trying to please its kitty master.
func Pet(c *Cat) {
	log(`petting ` + c.Name)
}
