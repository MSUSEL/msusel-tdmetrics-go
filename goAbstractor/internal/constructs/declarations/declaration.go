package declarations

import (
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Package interface {
	constructs.Construct
	Source() *packages.Package
}

type Declaration interface {
	constructs.Construct

	Package() Package
	Name() string
	Location() locs.Loc
}
