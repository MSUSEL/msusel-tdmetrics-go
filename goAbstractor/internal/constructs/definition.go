package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

// TODO: Add scope (exported)

type Definition interface {
	TypeDesc
	Name() string
	Package() Package
	Location() locs.Loc
}
