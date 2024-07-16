package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

type Definition interface {
	TypeDesc
	Name() string
	Package() Package
	Location() locs.Loc
}
