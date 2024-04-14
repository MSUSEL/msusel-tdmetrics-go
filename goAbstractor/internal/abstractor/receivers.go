package abstractor

import (
	"fmt"
	"go/ast"
	"go/types"
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func squeeze[E any, S ~[]E](s S) S {
	dest := 0
	for src, count := 0, len(s); src < count; src++ {
		if !utils.IsNil(s[src]) {
			s[dest], s[src] = s[src], s[dest]
			dest++
		}
	}
	return slices.Clip(s[:dest])
}

func (ab *abstractor) determineReceiver(m *constructs.Method, src *packages.Package, decl *ast.FuncDecl) {
	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(fmt.Errorf(`function declaration has unexpected receiver fields: %s`, pos(src, decl.Pos())))
		}
		recv := src.TypesInfo.Types[decl.Recv.List[0].Type].Type
		if p, ok := recv.(*types.Pointer); ok {
			m.NoCopyRecv = true
			recv = p.Elem()
		}
		n, ok := recv.(*types.Named)
		if !ok {
			panic(fmt.Errorf(`function declaration has unexpected receiver type: %T: %s`, recv, pos(src, decl.Pos())))
		}
		name := n.String()
		if index := strings.LastIndexAny(name, `/.`); index >= 0 {
			name = name[index+1:]
		}
		m.Receiver = name
	}
}

func (ab *abstractor) resolveReceivers() {
	for _, pkg := range ab.proj.Packages {
		pkgChanged := false
		for i, m := range pkg.Methods {
			if len(m.Receiver) > 0 {
				for _, t := range pkg.Types {
					if t.Name == m.Receiver {
						pkgChanged = true
						t.Methods = append(t.Methods, m)
						pkg.Methods[i] = nil
						break
					}
				}
			}
		}
		if pkgChanged {
			pkg.Methods = squeeze(pkg.Methods)
		}
	}
}
