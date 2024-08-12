package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

type TypeParam interface {
	TypeDesc
	IsTypeParam()

	Name() string
	Type() TypeDesc
}

type TypeParamArgs struct {
	Name string
	Type TypeDesc
}

type TypeParamFactory interface {
	NewTypeParam(args TypeParamArgs) TypeParam
	TypeParams() collections.ReadonlySortedSet[TypeParam]
}
