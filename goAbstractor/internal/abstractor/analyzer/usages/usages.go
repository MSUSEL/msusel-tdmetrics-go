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
)

type Usages struct {
	Reads   collections.SortedSet[constructs.Construct]
	Writes  collections.SortedSet[constructs.Construct]
	Invokes collections.SortedSet[constructs.Construct]
}

func newUsage() Usages {
	cmp := constructs.Comparer[constructs.Construct]()
	return Usages{
		Reads:   sortedSet.New(cmp),
		Writes:  sortedSet.New(cmp),
		Invokes: sortedSet.New(cmp),
	}
}

type usagesImp struct {
	info   *types.Info
	proj   constructs.Project
	curPkg constructs.Package
	conv   converter.Converter

	pending   constructs.Construct
	localDefs map[types.Object]constructs.TypeDesc
	params    map[types.Object]struct{}
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
		localDefs: map[types.Object]constructs.TypeDesc{},
		params:    map[types.Object]struct{}{},
		usages:    newUsage(),
	}

	ui.processNode(node)

	return ui.usages
}

func (ui *usagesImp) setPending(pending constructs.Construct) {
	ui.flushPendingToRead()
	ui.pending = pending
}

func (ui *usagesImp) takePending() constructs.Construct {
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

func (ui *usagesImp) flushPendingToWrite() {
	if !utils.IsNil(ui.pending) {
		ui.usages.Writes.Add(ui.pending)
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
		case *ast.ReturnStmt:
			ui.processReturn(t)
		case *ast.SelectorExpr:
			ui.processSelector(t)
		case *ast.ValueSpec:
			ui.processValueSpec(t)
		default:
			//fmt.Printf("usagesImp processNode unhandled (%[1]T) %[1]v\n", t)
			return true
		}
		return false
	})
}

func (ui *usagesImp) processAssign(assign *ast.AssignStmt) {
	// Process the right hand side (RHS) of the assignment.
	for _, exp := range assign.Rhs {
		ui.processNode(exp)
		ui.flushPendingToRead()
	}

	// Process the left hand side (LHS) of the assignment.
	// Any usage returned will be the usage that is being assigned to or nil.
	// The usage is nil if not resolvable (e.g. `*foo() = 10` where
	// `func foo() *int`) or if assignment to a local type (e.g. `x := 10`).
	for _, exp := range assign.Lhs {
		ui.processNode(exp)
		ui.flushPendingToWrite()
	}
}

func (ui *usagesImp) processCall(call *ast.CallExpr) {
	// Process arguments for the call.
	for _, arg := range call.Args {
		ui.processNode(arg)
		ui.flushPendingToRead()
	}

	// Process the invocation target,
	// e.g. `Bar` in `(*foo.Bar)( ⋯ )` or `foo.Bar( ⋯ )`.
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
	if fn.Params != nil {
		ui.addLocalParams(fn.Params.List)
	}
	if fn.Results != nil {
		ui.addLocalParams(fn.Results.List)
	}
	if fn.TypeParams != nil {
		ui.addLocalParams(fn.TypeParams.List)
	}
}

func (ui *usagesImp) addLocalParams(params []*ast.Field) {
	for _, p := range params {
		ui.addLocalParam(p)
	}
}

func (ui *usagesImp) addLocalParam(p *ast.Field) {
	for _, id := range p.Names {
		if obj := ui.info.Defs[id]; obj != nil {
			ui.params[obj] = struct{}{}
		}
	}
}

func (ui *usagesImp) processFunDecl(fn *ast.FuncDecl) {
	if fn.Recv != nil {
		ui.addLocalParams(fn.Recv.List)
	}
	ui.processFunc(fn.Type)
	ui.processNode(fn.Body)
}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	// Check if this identifier is part of a local definition.
	if def, ok := ui.info.Defs[id]; ok {
		if def == nil {
			return
		}

		// TODO: Rework
		ref := ui.createTempRefForObj(def)
		ui.flushPendingToRead()
		ui.newPending(ref)
		ui.usages.Writes.Add(ui.pending)
		ui.localDefs[def] = ui.pending
		return
	}

	// Check for an identifier is being used.
	obj, ok := ui.info.Uses[id]
	if !ok {
		return
	}

	// Check if defined earlier.
	if usage, ok := ui.localDefs[obj]; ok {
		ui.flushPendingToRead()
		ui.pending = usage
		return
	}

	// Check if parameter, result, or receiver that hasn't been used yet.
	if _, ok := ui.params[obj]; ok {
		ui.flushPendingToRead()
		arg := ui.proj.NewArgument(constructs.ArgumentArgs{
			Name: obj.Name(),
			Type: ui.conv.ConvertType(obj.Type()),
		})
		ui.newPending(arg, nil)
		ui.localDefs[obj] = ui.pending
	}

	// Skip over build-in constants:
	// https://pkg.go.dev/builtin#pkg-constants
	switch obj.Id() {
	case `_.true`, `_.false`, `_.nil`, `_.iota`:
		return
	}

	// Return basic types as usage.
	if basic, ok := obj.Type().(*types.Basic); ok && basic.Kind() != types.Invalid {
		ui.newPending(ui.proj.NewBasic(constructs.BasicArgs{
			RealType: basic,
		}), nil)
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
			ui.flushPendingToWrite()
		}
	}
}

func (ui *usagesImp) processIndex(expr *ast.IndexExpr) {
	// Process index, e.g. `⋯[ i+3 ]`
	ui.processNode(expr.Index)
	ui.flushPendingToRead()

	// Process target of indexing, e.g. `cats[ ⋯ ]`
	ui.processNode(expr.X)
	ui.usages.Reads.Add(ui.pending)
	target := ui.takePending()

	// Prepare the result of the index.
	if elem, ok := ui.info.Types[expr]; ok {
		ui.newPending(ui.conv.ConvertType(elem.Type), target)
	}
}

func (ui *usagesImp) processIndexList(expr *ast.IndexListExpr) {
	// Process indices, e.g. `⋯[ i+4 : length+2 ]`
	for _, indices := range expr.Indices {
		ui.processNode(indices)
		ui.flushPendingToRead()
	}

	// Process the object being indexed then
	// take the pending to be set after the indices.
	ui.processNode(expr.X)
	if !utils.IsNil(ui.pending) {
		ui.usages.Reads.Add(ui.pending)
	}
}

func (ui *usagesImp) processReturn(ret *ast.ReturnStmt) {
	for _, r := range ret.Results {
		ui.processNode(r)
		ui.flushPendingToRead()
	}
}

func (ui *usagesImp) processSelector(sel *ast.SelectorExpr) {
	ui.processNode(sel.X)
	lhs := ui.takePending()

	selection, ok := ui.info.Selections[sel]
	if !ok {
		panic(terror.New(`expected selection info but n found`).
			With(`expr`, sel))
	}

	if utils.IsNil(lhs) {
		// The left hand side is empty so this wasn't like `foo.Bar`,
		// it was instead like `foo().Bar` where some expression results
		// in a selection on a type.
		lhs = ui.proj.NewUsage(constructs.UsageArgs{
			Target: ui.conv.ConvertType(selection.Recv()),
		})
	}

	if utils.IsNil(lhs) {
		ui.newPending(ui.createTempRefForObj(selection.Obj()), nil)
		return
	}

	ui.usages.Reads.Add(lhs)
	// TODO: Handle better selection
	ui.newPending(ui.createTempRefForObj(selection.Obj()), lhs)
}

func (ui *usagesImp) processValueSpec(spec *ast.ValueSpec) {
	for _, name := range spec.Names {
		ui.processNode(name)
		ui.flushPendingToWrite()
	}

	ui.processNode(spec.Type)
	ui.flushPendingToRead()

	for _, value := range spec.Values {
		ui.processNode(value)
		ui.flushPendingToRead()
	}
}
