package abstractor

import (
	"go/ast"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

// TODO: Handle multiple init's per package.

func (ab *abstractor) setTypeParamOverrides(args *types.TypeList, params *types.TypeParamList, src *packages.Package, decl *ast.FuncDecl) {
	count := args.Len()
	if count != params.Len() {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, pos(src, decl.Pos())))
	}

	ab.typeParamReplacer = map[*types.TypeParam]*types.TypeParam{}
	for i := range count {
		ab.typeParamReplacer[args.At(i).(*types.TypeParam)] = params.At(i)
	}
}

func (ab *abstractor) clearTypeParamOverrides() {
	ab.typeParamReplacer = nil
}

func (ab *abstractor) abstractFuncDecl(pkg constructs.Package, src *packages.Package, decl *ast.FuncDecl) {
	obj := src.TypesInfo.Defs[decl.Name]

	noCopyRecv := false
	recvName := ``
	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(terror.New(`function declaration has unexpected receiver fields`).
				With(`pos`, pos(src, decl.Pos())))
		}
		recv := src.TypesInfo.Types[decl.Recv.List[0].Type].Type
		if p, ok := recv.(*types.Pointer); ok {
			noCopyRecv = true
			recv = p.Elem()
		}
		n, ok := recv.(*types.Named)
		if !ok {
			panic(terror.New(`function declaration has unexpected receiver type`).
				WithType(`receiver`, recv).
				With(`pos`, pos(src, decl.Pos())))
		}
		recvName = n.Origin().Obj().Name()
		ab.setTypeParamOverrides(n.TypeArgs(), n.TypeParams(), src, decl)
	}

	sig := ab.convertSignature(obj.Type().(*types.Signature))
	ab.clearTypeParamOverrides()

	mets := metrics.New(src.Fset, decl)
	ab.proj.NewMethod(constructs.MethodArgs{
		Package:    pkg,
		Name:       decl.Name.Name,
		Signature:  sig,
		Metrics:    mets,
		NoCopyRecv: noCopyRecv,
		Receiver:   recvName,
		Location:   ab.proj.NewLoc(decl.Pos()),
	})
}
