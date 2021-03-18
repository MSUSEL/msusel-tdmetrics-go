package dat002

import "errors"

type (
	CatBreed int

	ID string

	Cat struct {
		Name  string
		Breed CatBreed
		Age   int
	}

	CatTable map[ID]Cat
)

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

func (ct Table) Add(id ID, c Cat) {
	if _, ok := ct[id]; ok {
		panic(errors.New("cat identifier, " + id + ", already exists in a table"))
	}
	ct[id] = c
}
