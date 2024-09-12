package usages

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/usage"
)

type Usages struct {
	Reads   collections.SortedSet[constructs.Usage]
	Writes  collections.SortedSet[constructs.Usage]
	Invokes collections.SortedSet[constructs.Usage]
}

func newUsage() Usages {
	return Usages{
		Reads:   sortedSet.New(usage.Comparer()),
		Writes:  sortedSet.New(usage.Comparer()),
		Invokes: sortedSet.New(usage.Comparer()),
	}
}

type usagesImp struct {
	info *types.Info
	proj constructs.Project
	conv converter.Converter

	localDefs   collections.Set[types.Object]
	alreadyDefs map[types.Object]constructs.Usage
	usages      Usages
}

func Calculate(info *types.Info, proj constructs.Project, conv converter.Converter, node ast.Node) Usages {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotNil(`info.Defs`, info.Defs)
	assert.ArgNotNil(`info.Instances`, info.Instances)
	assert.ArgNotNil(`info.Selections`, info.Selections)
	assert.ArgNotNil(`info.Uses`, info.Uses)
	assert.ArgNotNil(`proj`, proj)
	assert.ArgNotNil(`conv`, conv)
	assert.ArgNotNil(`node`, node)

	ui := &usagesImp{
		info: info,
		proj: proj,
		conv: conv,

		localDefs:   set.New[types.Object](),
		alreadyDefs: make(map[types.Object]constructs.Usage),
		usages:      newUsage(),
	}

	ui.processNode(node)

	return ui.usages
}

func (ui *usagesImp) processNode(node ast.Node) constructs.Usage {
	var last constructs.Usage
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.AssignStmt:
			ui.processAssign(t)
			last = nil
		case *ast.CallExpr:
			ui.processCall(t)
		case *ast.Ident:
			last = ui.processIdent(t)
		case *ast.SelectorExpr:
			last = ui.processSelector(t)
		default:
			return true
		}
		return false
	})
	return last
}

func (ui *usagesImp) processAssign(assign *ast.AssignStmt) {
	// Process the left hand side (Lhs) of the assignment.
	// Any usage returned will be the usage that is being assigned to or nil.
	// The usage is nil if not resolvable (e.g. `*foo() = 10` where
	// `func foo() *int`) or if assignment to a local type (e.g. `x := 10`).
	for _, exp := range assign.Lhs {
		if last := ui.processNode(exp); !utils.IsNil(last) {
			ui.usages.Writes.Add(last)
		}
	}

	// Process the right hand side (Rhs) of the assignment.
	for _, exp := range assign.Rhs {
		ui.processNode(exp)
	}
}

func (ui *usagesImp) processCall(call *ast.CallExpr) {
	if last := ui.processNode(call.Fun); !utils.IsNil(last) {
		if tx, ok := ui.info.Types[call.Fun]; ok && tx.IsType() {
			// Explicit cast (conversion), e.g. `int(f.x)`
			ui.usages.Writes.Add(last)
		} else {
			ui.usages.Invokes.Add(last)
		}
	}
}

func (ui *usagesImp) processIdent(id *ast.Ident) constructs.Usage {
	// Check if this identifier is part of a local definition.
	// Don't create a usage for the local definition, they should be skipped.
	if def, ok := ui.info.Defs[id]; ok {
		ui.localDefs.Add(def)
		return nil
	}

	obj, ok := ui.info.Uses[id]
	if !ok || ui.localDefs.Contains(obj) {
		return nil
	}

	return ui.createUsage(id, obj, nil)
}

func (ui *usagesImp) processSelector(sel *ast.SelectorExpr) constructs.Usage {
	last := ui.processNode(sel.X)
	fmt.Printf(">>> last: %v\n", last)
	fmt.Printf(">>> sel:  %v\n", sel.Sel)

	// TODO: Finish implementing

	/*
		selection, ok := ui.info.Selections[sel]
		if !ok {
			panic(terror.New(`expected selection info but n found`).
				With(`expr`, sel))
		}
		fmt.Printf(">>> sel: %v, %v\n", sel, selection) // TODO: REMOVE
		fmt.Printf("  >>> %v\n", selection.Obj())
		fmt.Printf("  >>> %v\n", selection.Recv())
	*/
	return nil
}

func (ui *usagesImp) createUsage(id *ast.Ident, obj types.Object, origin constructs.Construct) constructs.Usage {
	if usage, ok := ui.alreadyDefs[obj]; ok {
		return usage
	}

	pkgPath := ``
	if !utils.IsNil(obj.Pkg()) {
		pkgPath = obj.Pkg().Path()
	}

	var instType []constructs.TypeDesc
	if inst, ok := ui.info.Instances[id]; ok {
		instType = ui.conv.ConvertInstanceTypes(inst.TypeArgs)
	}

	//if utils.IsNil(origin) {
	//	if v, ok := obj.(*types.Var); ok && v.IsField() {
	// TODO: Finish implementing by looking up the origin type
	//       if not given. Note `v.Origin` may equal `v`.
	//origin = ui.conv.ConvertVar(v.Origin())
	//	}
	//}

	usage := ui.proj.NewUsage(constructs.UsageArgs{
		PackagePath:   pkgPath,
		Name:          obj.Name(),
		InstanceTypes: instType,
		Origin:        origin,
	})
	ui.alreadyDefs[obj] = usage
	ui.usages.Reads.Add(usage)

	fmt.Printf("Created: %s\n", usage.String()) // TODO: remove

	return usage
}
