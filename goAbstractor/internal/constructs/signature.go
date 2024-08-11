package constructs

type Signature interface {
	TypeDesc
	IsSignature()

	// IsVacant indicates there are no parameters and no results,
	// i.e. `func()()`.
	IsVacant() bool
}
