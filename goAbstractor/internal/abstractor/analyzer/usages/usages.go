package usages

import (
	"go/ast"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
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

func newUsage() Usages {
	return Usages{
		Reads:   sortedSet.New(usage.Comparer()),
		Writes:  sortedSet.New(usage.Comparer()),
		Invokes: sortedSet.New(usage.Comparer()),
	}
}

type usagesImp struct {
	info   *types.Info
	proj   constructs.Project
	curPkg constructs.Package
	conv   converter.Converter

	pending   constructs.Usage
	localDefs map[types.Object]constructs.Usage
	usages    Usages
}

func Calculate(info *types.Info, proj constructs.Project, curPkg constructs.Package, conv converter.Converter, node ast.Node) Usages {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotNil(`info.Defs`, info.Defs)
	assert.ArgNotNil(`info.Instances`, info.Instances)
	assert.ArgNotNil(`info.Selections`, info.Selections)
	assert.ArgNotNil(`info.Uses`, info.Uses)
	assert.ArgNotNil(`proj`, proj)
	assert.ArgNotNil(`conv`, conv)
	assert.ArgNotNil(`node`, node)
	assert.ArgNotNil(`curPkg`, curPkg)

	ui := &usagesImp{
		info:   info,
		proj:   proj,
		curPkg: curPkg,
		conv:   conv,

		pending:   nil,
		localDefs: map[types.Object]constructs.Usage{},
		usages:    newUsage(),
	}

	ui.processNode(node)

	return ui.usages
}

func (ui *usagesImp) newPending(target, origin constructs.Construct) {
	ui.flushPendingToRead()
	ui.pending = ui.proj.NewUsage(constructs.UsageArgs{
		Target: target,
		Origin: origin,
	})
}

func (ui *usagesImp) takePending() constructs.Usage {
	pending := ui.pending
	ui.pending = nil
	return pending
}

func (ui *usagesImp) flushPendingToRead() {
	if !utils.IsNil(ui.pending) {
		// If the usage hasn't been consumed it is assumed
		// to have been read from, e.g. `a + b`.
		ui.usages.Reads.Add(ui.pending)
	}
	ui.pending = nil
}

func (ui *usagesImp) createTempRefForObj(obj types.Object) constructs.TempReference {
	assert.ArgNotNil(`object`, obj)

	pkgPath := ``
	if obj.Pkg() != nil {
		pkgPath = obj.Pkg().Path()
	}

	var instType []constructs.TypeDesc
	typ := obj.Type()
	if pointer, ok := typ.(*types.Pointer); ok {
		typ = pointer.Elem()
	}
	if named, ok := typ.(*types.Named); ok {
		instType = ui.conv.ConvertInstanceTypes(named.TypeArgs())
	}

	return ui.proj.NewTempReference(constructs.TempReferenceArgs{
		RealType:      obj.Type(),
		PackagePath:   pkgPath,
		Name:          obj.Name(),
		InstanceTypes: instType,
		Package:       ui.curPkg.Source(),
	})
}

func (ui *usagesImp) processNode(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case nil:
			return true
		case *ast.AssignStmt:
			ui.processAssign(t)
		case *ast.CallExpr:
			ui.processCall(t)
		case *ast.FuncType:
			ui.processFunc(t)
		case *ast.FuncDecl:
			ui.processFunDecl(t)
		case *ast.Ident:
			ui.processIdent(t)
		case *ast.IncDecStmt:
			ui.processIncDec(t)
		case *ast.IndexExpr:
			ui.processIndex(t)
		case *ast.IndexListExpr:
			ui.processIndexList(t)
		case *ast.SelectorExpr:
			ui.processSelector(t)
		default:
			return true
		}
		return false
	})
}

func (ui *usagesImp) processAssign(assign *ast.AssignStmt) {
	// Process the right hand side (RHS) of the assignment.
	for _, exp := range assign.Rhs {
		ui.processNode(exp)
	}
	ui.flushPendingToRead()

	// Process the left hand side (LHS) of the assignment.
	// Any usage returned will be the usage that is being assigned to or nil.
	// The usage is nil if not resolvable (e.g. `*foo() = 10` where
	// `func foo() *int`) or if assignment to a local type (e.g. `x := 10`).
	for _, exp := range assign.Lhs {
		ui.processNode(exp)
		if last := ui.takePending(); !utils.IsNil(last) {
			ui.usages.Writes.Add(last)
		}
	}
}

func (ui *usagesImp) processCall(call *ast.CallExpr) {
	// Process arguments for the call.
	for _, arg := range call.Args {
		ui.processNode(arg)
	}
	ui.flushPendingToRead()

	// Process the invocation target,
	// e.g. `Bar` in `(*foo.Bar)( ** )` or `foo.Bar( ** )`.
	ui.processNode(call.Fun)
	if target := ui.takePending(); !utils.IsNil(target) {
		if tx, ok := ui.info.Types[call.Fun]; ok && tx.IsType() {
			// Explicit cast (conversion), e.g. `int(f.x)`
			ui.usages.Writes.Add(target)
		} else {
			// Function invocation, e.g. `println(f.x)`, `fmt.Println(f.x)`
			ui.usages.Invokes.Add(target)
		}
	}
}

func (ui *usagesImp) processFunc(fn *ast.FuncType) {
	// This is part of a `ast.FuncLit` or `ast.FuncDecl`.
	// Skip parameters and returns
}

func (ui *usagesImp) processFunDecl(fn *ast.FuncDecl) {
	ui.processNode(fn.Body)
}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	// Check if this identifier is part of a local definition.
	if def, ok := ui.info.Defs[id]; ok {
		ref := ui.createTempRefForObj(def)
		ui.newPending(ref, ui.takePending())
		ui.usages.Writes.Add(ui.pending)
		ui.localDefs[def] = ui.pending
		return
	}

	obj, ok := ui.info.Uses[id]
	if !ok {
		return
	}
	if usage, ok := ui.localDefs[obj]; ok {
		ui.flushPendingToRead()
		ui.pending = usage
		return
	}

	ui.newPending(ui.createTempRefForObj(obj), nil)
}

func (ui *usagesImp) processIncDec(stmt *ast.IncDecStmt) {
	ui.processNode(stmt.X)
	if target := ui.takePending(); !utils.IsNil(target) {
		ui.usages.Writes.Add(target)
	} else {
		// What about `*(func())++` where the increment is on
		// the returned type, or `mapFoo["cat"]++`.
		if typ, ok := ui.info.Types[stmt.X]; ok {
			desc := ui.conv.ConvertType(typ.Type)
			ui.newPending(desc, nil)
			ui.usages.Writes.Add(ui.takePending())
		}
	}
}

func (ui *usagesImp) processIndex(expr *ast.IndexExpr) {
	// Process index, e.g. `**[ i+3 ]`
	ui.processNode(expr.Index)
	ui.flushPendingToRead()

	// Process target of indexing, e.g. `cats[ ** ]`
	ui.processNode(expr.X)
	ui.usages.Reads.Add(ui.pending)
	target := ui.takePending()

	// Prepare the result of the index.
	if elem, ok := ui.info.Types[expr]; ok {
		ui.newPending(ui.conv.ConvertType(elem.Type), target)
	}
}

func (ui *usagesImp) processIndexList(expr *ast.IndexListExpr) {
	// Process indices, e.g. `**[ i+4 : length+2 ]`
	for _, indices := range expr.Indices {
		ui.processNode(indices)
	}
	ui.flushPendingToRead()

	// Process the object being indexed then
	// take the pending to be set after the indices.
	ui.processNode(expr.X)
	ui.usages.Reads.Add(ui.pending)
}

func (ui *usagesImp) processSelector(sel *ast.SelectorExpr) {
	ui.processNode(sel.X)

	var lhs constructs.Construct = ui.takePending()
	selection, ok := ui.info.Selections[sel]
	if !ok {
		panic(terror.New(`expected selection info but n found`).
			With(`expr`, sel))
	}

	if utils.IsNil(lhs) {
		// The left hand side is empty so this wasn't like `foo.Bar`,
		// it was instead like `foo().Bar` where some expression results
		// in a selection on a type.
		lhs = ui.conv.ConvertType(selection.Recv())
	}

	ui.newPending(ui.createTempRefForObj(selection.Obj()), lhs)
}
