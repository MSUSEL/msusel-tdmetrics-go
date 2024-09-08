package usages

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
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

type usagesImp struct {
	info *types.Info
	proj constructs.Project
	conv converter.Converter

	localDefs    collections.Set[types.Object]
	identHandled collections.Set[*ast.Ident]

	reads   collections.SortedSet[constructs.Usage]
	writes  collections.SortedSet[constructs.Usage]
	invokes collections.SortedSet[constructs.Usage]
}

func Calculate(info *types.Info, proj constructs.Project, conv converter.Converter, node ast.Node) Usages {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotEmpty(`info.Defs`, info.Defs)
	assert.ArgNotEmpty(`info.Instances`, info.Instances)
	assert.ArgNotEmpty(`info.Selections`, info.Selections)
	assert.ArgNotEmpty(`info.Uses`, info.Uses)
	assert.ArgNotNil(`proj`, proj)
	assert.ArgNotNil(`conv`, conv)
	assert.ArgNotNil(`node`, node)

	ui := &usagesImp{
		info: info,
		proj: proj,
		conv: conv,

		localDefs:    set.New[types.Object](),
		identHandled: set.New[*ast.Ident](),

		reads:   sortedSet.New(usage.Comparer()),
		writes:  sortedSet.New(usage.Comparer()),
		invokes: sortedSet.New(usage.Comparer()),
	}

	ui.getUsages(node)

	return Usages{
		Reads:   ui.reads,
		Writes:  ui.writes,
		Invokes: ui.invokes,
	}
}

func (ui *usagesImp) getUsages(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.AssignStmt:
			ui.processAssign(t)
		case *ast.Ident:
			ui.processIdent(t)
		case *ast.SelectorExpr:
			ui.processSelector(t)
		}
		return true
	})
}

func (ui *usagesImp) processAssign(assign *ast.AssignStmt) {

	fmt.Printf(">>> assign: %v\n", assign) // TODO: REMOVE

}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	if ui.identHandled.Contains(id) {
		return
	}
	ui.identHandled.Add(id)

	if def, ok := ui.info.Defs[id]; ok {
		ui.localDefs.Add(def)
		return
	}

	// TODO: Finish implementing
	// TODO: Use local to change selection into normal target so that
	//       if someone uses a struct locally to external types then the usage
	//       of a selection on that struct are the same as just using that type.

	if useObj, ok := ui.info.Uses[id]; ok && !ui.localDefs.Contains(useObj) {
		usage := ui.createUsage(id, useObj)

		fmt.Printf(">>> (%T) %v\t\t%s\n", useObj, useObj, usage.String()) // TODO: REMOVE

		//usages[id] = usage
	}
}

func (ui *usagesImp) processSelector(sel *ast.SelectorExpr) {
	selection, ok := ui.info.Selections[sel]
	if !ok {
		panic(terror.New(`expected selection info but n found`).
			With(`expr`, sel))
	}

	fmt.Printf(">>> sel: %v, %v\n", sel, selection) // TODO: REMOVE
	fmt.Printf("  >>> %v\n", selection.Obj())
	fmt.Printf("  >>> %v\n", selection.Recv())
}

func (ui *usagesImp) createUsage(id *ast.Ident, useObj types.Object) constructs.Usage {
	pkgPath := ``
	if !utils.IsNil(useObj.Pkg()) {
		pkgPath = useObj.Pkg().Path()
	}

	var instType []constructs.TypeDesc
	if inst, ok := ui.info.Instances[id]; ok {
		instType = ui.conv.ConvertInstanceTypes(inst.TypeArgs)
	}

	target := useObj.Name()
	selection := ``
	if v, ok := useObj.(*types.Var); ok && v.IsField() {
		fmt.Printf(">>>X %v\n", v)

		if v == v.Origin() {
			fmt.Println("SAME!!!")
		}

		target = v.Origin().Name()
		selection = useObj.Name()
		// TODO: Need to get instance information for origin, maybe, create tests.
	}

	return ui.proj.NewUsage(constructs.UsageArgs{
		PackagePath:   pkgPath,
		Target:        target,
		InstanceTypes: instType,
		Selection:     selection,
	})
}
