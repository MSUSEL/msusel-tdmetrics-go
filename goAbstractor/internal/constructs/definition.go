package constructs

type Definition interface {
	TypeDesc
	Name() string
	Package() Package
}
