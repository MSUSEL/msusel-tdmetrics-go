package constructs

type Value interface {
	Declaration
	IsValue()
}
