package constructs

import (
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Package struct {
	pkg *packages.Package

	Path    string
	Imports []*Package
	Types   []*TypeDef
	Values  []*ValueDef
	Methods []*Method

	Index       int
	ImportPaths []string
}

func NewPackage(pkg *packages.Package, path string, importPaths []string) *Package {
	return &Package{
		pkg:         pkg,
		Path:        path,
		ImportPaths: importPaths,
	}
}

func (p *Package) Source() *packages.Package {
	return p.pkg
}

func (p *Package) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, p.Index)
	}

	ctx2 := ctx.ShowKind().Short()
	return jsonify.NewMap().
		AddNonZero(ctx2, `path`, p.Path).
		AddNonZero(ctx2, `imports`, p.Imports).
		AddNonZero(ctx2, `types`, p.Types).
		AddNonZero(ctx2, `values`, p.Values).
		AddNonZero(ctx2, `methods`, p.Methods)
}

func (p *Package) String() string {
	return jsonify.ToString(p)
}

func (p *Package) FindTypeForReceiver(receiver string) *TypeDef {
	for _, t := range p.Types {
		if receiver == t.Name {
			return t
		}
	}
	return nil
}
