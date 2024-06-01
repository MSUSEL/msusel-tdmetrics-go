package abstractor

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func (ab *abstractor) determineReceiver(m constructs.Method, src *packages.Package, decl *ast.FuncDecl) {
	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(fmt.Errorf(`function declaration has unexpected receiver fields: %s`, pos(src, decl.Pos())))
		}
		recv := src.TypesInfo.Types[decl.Recv.List[0].Type].Type
		noCopyRecv := false
		if p, ok := recv.(*types.Pointer); ok {
			noCopyRecv = true
			recv = p.Elem()
		}
		n, ok := recv.(*types.Named)
		if !ok {
			panic(fmt.Errorf(`function declaration has unexpected receiver type: %T: %s`, recv, pos(src, decl.Pos())))
		}
		name := n.String()
		if index := strings.Index(name, `[`); index >= 0 {
			name = name[:index]
		}
		if index := strings.LastIndexAny(name, `/.`); index >= 0 {
			name = name[index+1:]
		}
		m.SetReceiver(noCopyRecv, name)
	}
}

func (ab *abstractor) resolveReceivers() {
	ab.log(`resolve receivers`)
	for _, pkg := range ab.proj.Packages() {
		resolveReceiversInPackage(pkg)
	}
}

func resolveReceiversInPackage(pkg constructs.Package) {
	pkgChanged := false
	methods := pkg.Methods()
	for i, m := range methods {
		if len(m.Receiver()) > 0 {

			t := pkg.FindTypeForReceiver(m.Receiver())
			if t == nil {
				panic(fmt.Errorf(`failed to find receiver for %s`, m.Receiver()))
			}

			pkgChanged = true
			t.AppendMethod(m)
			methods[i] = nil
		}
	}
	if pkgChanged {
		pkg.SetMethods(utils.RemoveZeros(methods))
	}
}
