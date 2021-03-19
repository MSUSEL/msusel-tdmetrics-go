package dat002

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	alreadyExistsInTable = "cat identifier, %s, already exists in the table"
	doesNotExistInTable  = "no cat by the given identifier, %s, in the table"
)

type (
	// CatBreed is a enumerator type for predefined cat breeds.
	CatBreed int

	// ID is a unique identifier used to specify cats.
	ID string

	// IDSlice is a slice of IDs.
	IDSlice []ID

	// Cat stores a single cat's information.
	Cat struct {
		Name  string
		Breed CatBreed
		Age   int
	}

	// CatTable is a collection of identified set of cats.
	CatTable map[ID]*Cat
)

// The enumeration of all predefined cat breeds.
const (
	BritishShorthair = CatBreed(iota)
	Persian
	MaineCoon
	Siamese
	Bengal
	AmericanShorthair
	Ragdoll
	Abyssinian
	ExoticShorthair
	NorwegianForest
	ScottishFold
	Burmese
	Siberian
	Savannah
	RussianBlue
	JapaneseBobtail
	Manx
	DevonRex
	Himalayan
	Bombay
	CornishRex
	EgyptianMau
	Munchkin
	Balinese
)

// String gets the name for a cat breed.
func (cb CatBreed) String() string {
	switch cb {
	case BritishShorthair:
		return "British Shorthair"
	case Persian:
		return "Persian"
	case MaineCoon:
		return "Maine Coon"
	case Siamese:
		return "Siamese"
	case Bengal:
		return "Bengal"
	case AmericanShorthair:
		return "American Shorthair"
	case Ragdoll:
		return "Ragdoll"
	case Abyssinian:
		return "Abyssinian"
	case ExoticShorthair:
		return "Exotic Shorthair"
	case NorwegianForest:
		return "Norwegian Forest"
	case ScottishFold:
		return "Scottish Fold"
	case Burmese:
		return "Burmese"
	case Siberian:
		return "Siberian"
	case Savannah:
		return "Savannah"
	case RussianBlue:
		return "Russian Blue"
	case JapaneseBobtail:
		return "Japanese Bobtail"
	case Manx:
		return "Manx"
	case DevonRex:
		return "Devon Rex"
	case Himalayan:
		return "Himalayan"
	case Bombay:
		return "Bombay"
	case CornishRex:
		return "Cornish Rex"
	case EgyptianMau:
		return "Egyptian Mau"
	case Munchkin:
		return "Munchkin"
	case Balinese:
		return "Balinese"
	}
	return "Unknown"
}

// nextIDValue is the next value to use in constructing an ID
var nextIDValue = 0

// NewID creates a new cat ID for the table.
func NewID() ID {
	val := nextIDValue
	nextIDValue++
	return ID(fmt.Sprintf("cat-%d", val))
}

// Len gets the number of IDs in the slice used for sorting.
func (ids IDSlice) Len() int {
	return len(ids)
}

// Less determines if the `i`th ID is less than the `j`th ID used for sorting.
func (ids IDSlice) Less(i, j int) bool {
	return ids[i] < ids[j]
}

// Swap will swap the `i`th ID with the `j`th ID.
func (ids IDSlice) Swap(i, j int) {
	ids[j], ids[i] = ids[i], ids[j]
}

// String gets a string of all ths IDs in the slice used for debugging.
func (ids IDSlice) String() string {
	var s strings.Builder
	for i, id := range ids {
		if i > 0 {
			if _, err := s.WriteString(`, `); err != nil {
				panic(err)
			}
		}
		if _, err := s.WriteString(string(id)); err != nil {
			panic(err)
		}
	}
	return s.String()
}

// Has determines if there is a cat with the given ID.
func (ct CatTable) Has(id ID) bool {
	_, ok := ct[id]
	return ok
}

// Adds a new cat with the given information and ID.
func (ct CatTable) Add(id ID, name string, breed CatBreed, age int) {
	if ct.Has(id) {
		panic(fmt.Errorf(alreadyExistsInTable, id))
	}
	ct[id] = &Cat{
		Name:  name,
		Breed: breed,
		Age:   age,
	}
}

// Adds a new cat with a generated ID. The ID will be returned.
func (ct CatTable) AddNew(name string, breed CatBreed, age int) ID {
	id := NewID()
	ct.Add(id, name, breed, age)
	return id
}

// Name gets the name of the cat with the given ID.
func (ct CatTable) Name(id ID) string {
	if !ct.Has(id) {
		panic(fmt.Errorf(doesNotExistInTable, id))
	}
	return ct[id].Name
}

// Breed gets the breed of the cat with the given ID.
func (ct CatTable) Breed(id ID) CatBreed {
	if !ct.Has(id) {
		panic(fmt.Errorf(doesNotExistInTable, id))
	}
	return ct[id].Breed
}

// Age gets the age of the cat with the given ID.
func (ct CatTable) Age(id ID) int {
	if !ct.Has(id) {
		panic(fmt.Errorf(doesNotExistInTable, id))
	}
	return ct[id].Age
}

// AllIDs get a slice of all the IDs in the table.
func (ct CatTable) AllIDs() IDSlice {
	ids := IDSlice{}
	for id := range ct {
		ids = append(ids, id)
	}
	sort.Sort(ids)
	return ids
}

// AgeRange gets the minimum and maximum ages in this table.
func (ct CatTable) AgeRange() (int, int) {
	min, max, first := -1, -1, true
	for _, c := range ct {
		if first {
			min, max, first = c.Age, c.Age, false
		} else {
			if min > c.Age {
				min = c.Age
			}
			if max < c.Age {
				max = c.Age
			}
		}
	}
	return min, max
}

// AllBreeds gets the set all the breeds in this table.
func (ct CatTable) AllBreeds() []CatBreed {
	intMap := map[int]bool{}
	for _, c := range ct {
		intMap[int(c.Breed)] = true
	}

	index, intSlice := 0, make([]int, len(intMap))
	for b := range intMap {
		intSlice[index] = b
		index++
	}
	sort.Ints(intSlice)

	breeds := make([]CatBreed, len(intSlice))
	for i, b := range intSlice {
		breeds[i] = CatBreed(b)
	}
	return breeds
}

// CatsWithBreed gets the IDs for all the cats of the given breed type.
func (ct CatTable) CatsWithBreed(breed CatBreed) IDSlice {
	ids := IDSlice{}
	for id, c := range ct {
		if c.Breed == breed {
			ids = append(ids, id)
		}
	}
	sort.Sort(ids)
	return ids
}

// CatsWithName gets the IDs for all the cats of the given name.
func (ct CatTable) CatsWithName(name string) IDSlice {
	ids := IDSlice{}
	for id, c := range ct {
		if c.Name == name {
			ids = append(ids, id)
		}
	}
	sort.Sort(ids)
	return ids
}

// String returns a formatted table view of the cat table for debugging.
func (ct CatTable) String() string {
	idTitle, nameTitle, breedTitle, ageTitle := `ID`, `Name`, `Breed`, `Age`
	idLenMax, nameLenMax, breedLenMax, ageLenMax :=
		len(idTitle), len(nameTitle), len(breedTitle), len(ageTitle)
	for id, c := range ct {
		if idLen := len(id); idLen > idLenMax {
			idLenMax = idLen
		}
		if nameLen := len(c.Name); nameLen > nameLenMax {
			nameLenMax = nameLen
		}
		if breedLen := len(c.Breed.String()); breedLen > breedLenMax {
			breedLenMax = breedLen
		}
		if ageLen := len(strconv.Itoa(c.Age)); ageLen > ageLenMax {
			ageLenMax = ageLen
		}
	}

	s := &strings.Builder{}
	fmt.Fprintf(s, "|%-*s|%-*s|%-*s|%*s|",
		idLenMax, idTitle, nameLenMax, nameTitle,
		breedLenMax, breedTitle, ageLenMax, ageTitle)
	fmt.Fprintf(s, "\n|:%s|:%s|:%s|%s:|",
		strings.Repeat(`-`, idLenMax-1), strings.Repeat(`-`, nameLenMax-1),
		strings.Repeat(`-`, breedLenMax-1), strings.Repeat(`-`, ageLenMax-1))
	ids := ct.AllIDs()
	for _, id := range ids {
		c := ct[id]
		fmt.Fprintf(s, "\n|%-*s|%-*s|%-*s|%*s|",
			idLenMax, id, nameLenMax, c.Name,
			breedLenMax, c.Breed.String(), ageLenMax, strconv.Itoa(c.Age))
	}
	return s.String()
}
