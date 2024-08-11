package constructs

type Reference interface {
	TypeDesc
	IsReference()

	PackagePath() string
	Name() string
	Resolved() bool
	SetType(typ TypeDesc)
}
