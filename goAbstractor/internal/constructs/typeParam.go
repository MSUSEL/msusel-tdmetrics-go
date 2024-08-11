package constructs

type TypeParam interface {
	TypeDesc
	IsTypeParam()

	Name() string
	Type() TypeDesc
}
