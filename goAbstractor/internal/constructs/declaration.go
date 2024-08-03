package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

type Declaration interface {
	TypeDesc
	Name() string
	Package() Package
	Location() locs.Loc
}
