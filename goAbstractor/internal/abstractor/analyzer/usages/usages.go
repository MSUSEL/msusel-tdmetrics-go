package usages

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type Usages struct {
	Reads   collections.SortedSet[constructs.Construct]
	Writes  collections.SortedSet[constructs.Construct]
	Invokes collections.SortedSet[constructs.Construct]

	SideEffect bool
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
	fSet   *token.FileSet
	info   *types.Info
	proj   constructs.Project
	curPkg constructs.Package
	baker  baker.Baker
	conv   converter.Converter
	root   ast.Node
	usages Usages

	pending constructs.Construct
}

func Calculate(info *types.Info, proj constructs.Project, curPkg constructs.Package, baker baker.Baker, conv converter.Converter, root ast.Node) Usages {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotNil(`info.Defs`, info.Defs)
	assert.ArgNotNil(`info.Instances`, info.Instances)
	assert.ArgNotNil(`info.Selections`, info.Selections)
	assert.ArgNotNil(`info.Uses`, info.Uses)
	assert.ArgNotNil(`info.Types`, info.Types)
	assert.ArgNotNil(`proj`, proj)
	assert.ArgNotNil(`baker`, baker)
	assert.ArgNotNil(`conv`, conv)
	assert.ArgNotNil(`root`, root)
	assert.ArgNotNil(`curPkg`, curPkg)

	ui := &usagesImp{
		fSet:    curPkg.Source().Fset,
		info:    info,
		proj:    proj,
		curPkg:  curPkg,
		baker:   baker,
		conv:    conv,
		root:    root,
		usages:  newUsage(),
		pending: nil,
	}

	//ast.Print(proj.Locs().FileSet(), root)
	ui.processNode(root)

	return ui.usages
}

func (ui *usagesImp) pos(pr posReader) token.Position {
	return ui.fSet.Position(pr.Pos())
}

func (ui *usagesImp) hasPending() bool {
	return !utils.IsNil(ui.pending)
}

func (ui *usagesImp) setPendingConstruct(c constructs.Construct) {
	ui.flushPendingToRead()
	ui.pending = c
}

func (ui *usagesImp) setPendingType(t types.Type) {
	ui.flushPendingToRead()
	if utils.IsNil(t) {
		return
	}
	if named, ok := stripNamed(t); ok {
		if isLocal(ui.root, named.Obj()) {
			return
		}
	}
	ui.pending = ui.conv.ConvertType(t)
}

func (ui *usagesImp) clearPending() {
	ui.pending = nil
}

// flushPendingToRead writes any pending as a read usage
// and clears the pending.
//
// If the usage hasn't been consumed it is assumed
// to have been read from, e.g. `a + b`.
func (ui *usagesImp) flushPendingToRead() {
	ui.addRead(ui.pending)
	ui.clearPending()
}

// flushPendingToWrite writes any pending as a write usage
// and clears the pending.
func (ui *usagesImp) flushPendingToWrite() {
	ui.addWrite(ui.pending)
	ui.clearPending()
}

// flushPendingToInvoke writes any pending as an invoke usage
// and clears the pending.
func (ui *usagesImp) flushPendingToInvoke() {
	ui.addInvoke(ui.pending)
	ui.clearPending()
}

// addRead adds the given construct as a read usage.
func (ui *usagesImp) addRead(c constructs.Construct) {
	if !utils.IsNil(c) {
		ui.usages.Reads.Add(c)
	}
}

// addWrite adds the given construct as a write usage.
func (ui *usagesImp) addWrite(c constructs.Construct) {
	if !utils.IsNil(c) {
		ui.usages.Writes.Add(c)
	}
}

// addInvoke adds the given construct as an invoke usage.
func (ui *usagesImp) addInvoke(c constructs.Construct) {
	if !utils.IsNil(c) {
		ui.usages.Invokes.Add(c)
	}
}

func (ui *usagesImp) handleBuiltinCall(call *ast.CallExpr) {
	switch name := getName(ui.fSet, call.Fun); name {
	case `append`, `cap`, `complex`, `copy`, `imag`, `len`,
		`make`, `max`, `min`, `new`, `real`, `recover`,
		`unsafe.Alignof`, `unsafe.Offsetof`, `unsafe.Sizeof`,
		`unsafe.String`, `unsafe.StringData`, `unsafe.Slice`,
		`unsafe.SliceData`, `unsafe.Add`:
		if typ, ok := ui.info.Types[call]; ok {
			ui.setPendingType(typ.Type)
		}
		return

	case `clear`, `close`, `delete`, `panic`:
		return

	case `print`, `println`:
		ui.usages.SideEffect = true
		return

	default:
		panic(terror.New(`failed to get name of builtin function`).
			With(`name`, name).
			With(`position`, ui.pos(call)))
	}
}

func (ui *usagesImp) processNode(node ast.Node) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		panic(terror.New(`error processing node`, terror.RecoveredPanic(r)).
	//			With(`pos`, ui.pos(node)).
	//			WithType(`node`, node))
	//	}
	//}()

	if utils.IsNil(node) {
		return
	}
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case nil:
			return true
		case *ast.AssignStmt:
			ui.processAssign(t)
		case *ast.BlockStmt:
			ui.processBlock(t)
		case *ast.CallExpr:
			ui.processCall(t)
		case *ast.CompositeLit:
			ui.processCompositeLit(t)
		case *ast.Ident:
			ui.processIdent(t)
		case *ast.IncDecStmt:
			ui.processIncDec(t)
		case *ast.IndexExpr:
			ui.processIndex(t)
		case *ast.IndexListExpr:
			ui.processIndexList(t)
		case *ast.RangeStmt:
			ui.processRange(t)
		case *ast.ReturnStmt:
			ui.processReturn(t)
		case *ast.SendStmt:
			ui.processSend(t)
		case *ast.SelectorExpr:
			ui.processSelector(t)
		case *ast.TypeAssertExpr:
			ui.processTypeAssert(t)
		case *ast.TypeSpec:
			ui.processTypeSpec(t)
		case *ast.ValueSpec:
			ui.processValueSpec(t)
		default:
			// The following print is useful for debugging by showing
			// which nodes do not have custom handling on them yet.
			// Not all nodes need custom handling but a bug might indicate
			// one that doesn't have custom handling probably should.
			//fmt.Printf("usagesImp.processNode unhandled (%[1]T) %[1]v\n", t)
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

func (ui *usagesImp) processBlock(block *ast.BlockStmt) {
	// Each separate statement should be flushed after.
	for _, statement := range block.List {
		ui.processNode(statement)
		ui.flushPendingToRead()
	}
}

func (ui *usagesImp) processCall(call *ast.CallExpr) {
	// Process arguments for the call and flush each as read.
	for _, arg := range call.Args {
		ui.processNode(arg)
		ui.flushPendingToRead()
	}

	// Process the invocation target,
	// e.g. `Bar` in `(*foo.Bar)( ⋯ )` or `foo.Bar( ⋯ )`.
	typ, ok := ui.info.Types[call.Fun]
	if !ok {
		panic(terror.New(`failed to find type info`).
			With(`function`, call.Fun).
			With(`position`, ui.pos(call)))
	}

	// Check for explicit cast (conversion), e.g. `int(f.x)`
	if typ.IsType() {
		ui.processNode(call.Fun)
		ui.flushPendingToWrite()
		return
	}

	// Check for builtin method, e.g. `println(f.x)`
	if typ.IsBuiltin() {
		ui.handleBuiltinCall(call)
		return
	}

	// Function invocation, e.g. `fmt.Println(f.x)`
	ui.processNode(call.Fun)
	ui.flushPendingToInvoke()
}

func (ui *usagesImp) processCompositeLit(comp *ast.CompositeLit) {
	for _, elem := range comp.Elts {
		ui.processNode(elem)
		ui.flushPendingToRead()
	}

	// Skip over comp.Type, the internally defined type.
	ui.processNode(comp.Type)
}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	// Check if this identifier is part of a definition.
	if def, ok := ui.info.Defs[id]; ok {
		if def == nil {
			// Skip over `t` in `select t := x.(type)` type definitions.
			return
		}
		if _, ok := def.(*types.Label); ok {
			// Skip over labels
			return
		}

		// If the definition isn't of a local type the it has been written.
		// Add to write and leave the pending so it can also be used in
		// read or invoke.
		ui.setPendingType(def.Type())
		ui.addWrite(ui.pending)
		return
	}

	// Check if the identifier is being used.
	obj, ok := ui.info.Uses[id]
	if !ok {
		return
	}

	// Skip over build-in constants:
	// https://pkg.go.dev/builtin#pkg-constants
	switch obj.Id() {
	case `_.true`, `_.false`, `_.nil`, `_.iota`:
		return
	}

	// Check if object was defined within root and reuse type.
	if ok := isLocal(ui.root, obj); ok {
		return
	}

	// Return builtin type.
	if obj.Pkg() == nil {
		if typ := ui.baker.TypeByName(obj.Name()); !utils.IsNil(typ) {
			ui.setPendingConstruct(typ)
			return
		}
	}

	// Return basic types as usage.
	if basic, ok := obj.Type().(*types.Basic); ok && basic.Kind() != types.Invalid {
		switch basic.Kind() {
		case types.Complex64:
			ui.setPendingConstruct(ui.baker.BakeComplex64())
		case types.Complex128:
			ui.setPendingConstruct(ui.baker.BakeComplex128())
		default:
			ui.setPendingConstruct(ui.proj.NewBasic(constructs.BasicArgs{
				RealType: basic,
			}))
		}
		return
	}

	// Create a temp reference for this object.
	pkgPath := ``
	if obj.Pkg() != nil {
		pkgPath = obj.Pkg().Path()
	}

	var instType []constructs.TypeDesc
	if named, ok := stripNamed(obj.Type()); ok {
		instType = ui.conv.ConvertInstanceTypes(named.TypeArgs())
	}

	ui.setPendingConstruct(ui.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   pkgPath,
		Name:          obj.Name(),
		InstanceTypes: instType,
	}))
}

func (ui *usagesImp) processIncDec(stmt *ast.IncDecStmt) {
	ui.processNode(stmt.X)
	if ui.hasPending() {
		ui.flushPendingToWrite()
		return
	}

	// Handle `*(func())++` where the increment is on
	// the returned type, or `mapFoo["cat"]++`.
	if tv, ok := ui.info.Types[stmt.X]; ok {
		ui.setPendingType(tv.Type)
		ui.flushPendingToWrite()
	}
}

func (ui *usagesImp) processIndex(expr *ast.IndexExpr) {
	// Process index, e.g. `⋯[ i+3 ]`
	ui.processNode(expr.Index)
	ui.flushPendingToRead()

	// Process target of indexing, e.g. `cats[ ⋯ ]`
	// Add to read but leave in pending.
	ui.processNode(expr.X)
	ui.flushPendingToRead()

	// Prepare the pending after the indexing.
	if tv, ok := ui.info.Types[expr]; ok {
		if _, ok := tv.Type.(*types.Tuple); ok {
			// The indexing returned a (value, ok) tuple.
			// The results are probably used in a function parameter or an
			// assignment so the pending construct doesn't need to be set.
			return
		}

		// The indexing returned a single value.
		ui.setPendingType(tv.Type)
	}
}

func (ui *usagesImp) processIndexList(expr *ast.IndexListExpr) {
	// Process indices, e.g. `⋯[ i+4 : length+2 ]`
	for _, indices := range expr.Indices {
		ui.processNode(indices)
		ui.flushPendingToRead()
	}

	// Process the object being indexed then
	// set the pending after the index to be read from but leave
	// it in pending since the same type is going to be used.
	ui.processNode(expr.X)
	ui.addRead(ui.pending)
}

func (ui *usagesImp) processRange(r *ast.RangeStmt) {
	ui.processNode(r.Key)
	ui.flushPendingToWrite()

	ui.processNode(r.Value)
	ui.flushPendingToWrite()

	ui.processNode(r.X)
	ui.processNode(r.Body)
}

func (ui *usagesImp) processReturn(ret *ast.ReturnStmt) {
	for _, r := range ret.Results {
		ui.processNode(r)
		ui.flushPendingToRead()
	}
}

func (ui *usagesImp) processSend(send *ast.SendStmt) {
	ui.processNode(send.Value)
	ui.flushPendingToRead()

	ui.processNode(send.Chan)
	ui.flushPendingToWrite()
}

func (ui *usagesImp) processSelector(sel *ast.SelectorExpr) {
	ui.processNode(sel.X)
	if ui.hasPending() {
		ui.setPendingConstruct(ui.proj.NewSelection(constructs.SelectionArgs{
			Name:   sel.Sel.Name,
			Origin: ui.pending,
		}))
		return
	}

	selObj, ok := ui.info.Selections[sel]
	if !ok {
		panic(terror.New(`expected selection info but not found`).
			With(`expr`, sel).
			With(`position`, ui.pos(sel)))
	}

	// The left hand side is empty so this wasn't like `foo.Bar`,
	// it was instead like `foo().Bar` where some expression results
	// in a selection on a type. If a locally defined type then skip
	// the select since it won't reference anything useful.
	ui.setPendingType(selObj.Recv())
	ui.flushPendingToRead()
	ui.setPendingType(selObj.Obj().Type())
}

func (ui *usagesImp) processTypeAssert(exp *ast.TypeAssertExpr) {
	ui.processNode(exp.X)
	ui.flushPendingToRead()

	tv, ok := ui.info.Types[exp.Type]
	if !ok {
		panic(terror.New(`Expected a type in a TypeAssert for usages.`).
			With(`node`, exp.Type).
			With(`pos`, ui.pos(exp)))
	}

	ui.setPendingType(tv.Type)
}

func (ui *usagesImp) processTypeSpec(_ *ast.TypeSpec) {
	// A local type is being defined, skip the definition.
	ui.flushPendingToRead()
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
