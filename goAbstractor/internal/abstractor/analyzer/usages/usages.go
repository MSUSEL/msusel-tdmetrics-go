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
	info   *types.Info
	proj   constructs.Project
	curPkg constructs.Package
	baker  baker.Baker
	conv   converter.Converter

	pendingEff bool
	pendingCon constructs.Construct

	localDefs map[types.Object]constructs.Construct
	usages    Usages
}

func Calculate(info *types.Info, proj constructs.Project, curPkg constructs.Package, baker baker.Baker, conv converter.Converter, node ast.Node) Usages {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotNil(`info.Defs`, info.Defs)
	assert.ArgNotNil(`info.Instances`, info.Instances)
	assert.ArgNotNil(`info.Selections`, info.Selections)
	assert.ArgNotNil(`info.Uses`, info.Uses)
	assert.ArgNotNil(`info.Types`, info.Types)
	assert.ArgNotNil(`proj`, proj)
	assert.ArgNotNil(`baker`, baker)
	assert.ArgNotNil(`conv`, conv)
	assert.ArgNotNil(`node`, node)
	assert.ArgNotNil(`curPkg`, curPkg)

	ui := &usagesImp{
		info:   info,
		proj:   proj,
		curPkg: curPkg,
		baker:  baker,
		conv:   conv,

		pendingEff: false,
		pendingCon: nil,

		localDefs: map[types.Object]constructs.Construct{},
		usages:    newUsage(),
	}

	ui.processNode(node)

	return ui.usages
}

func (ui *usagesImp) clearPending() {
	ui.pendingCon = nil
	ui.pendingEff = false
}

// flushPendingToRead writes any pending as a read usage
// and clears the pending.
//
// If the usage hasn't been consumed it is assumed
// to have been read from, e.g. `a + b`.
func (ui *usagesImp) flushPendingToRead() {
	ui.addRead(ui.pendingCon)
	ui.clearPending()
}

// flushPendingToWrite writes any pending as a write usage
// and clears the pending.
func (ui *usagesImp) flushPendingToWrite() {
	ui.addWrite(ui.pendingCon, ui.pendingEff)
	ui.clearPending()
}

// addRead adds the given construct as a read usage.
func (ui *usagesImp) addRead(c constructs.Construct) {
	if !utils.IsNil(c) {
		ui.usages.Reads.Add(c)
	}
}

// addWrite adds the given construct as a write usage.
func (ui *usagesImp) addWrite(c constructs.Construct, eff bool) {
	if !utils.IsNil(c) {
		ui.usages.Writes.Add(c)
		if eff {
			ui.usages.SideEffect = true
		}
	}
}

func (ui *usagesImp) processNode(node ast.Node) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		panic(terror.New(`error processing node`, terror.RecoveredPanic(r)).
	//			With(`pos`, ui.proj.Locs().FileSet().Position(node.Pos())).
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
		case *ast.RangeStmt:
			ui.processRange(t)
		case *ast.ReturnStmt:
			ui.processReturn(t)
		case *ast.SendStmt:
			ui.processSend(t)
		case *ast.SelectorExpr:
			ui.processSelector(t)
		case *ast.ValueSpec:
			ui.processValueSpec(t)
		default:
			// The following print is useful for debugging by showing
			// which nodes do not have custom handling on them yet.
			// Not all nodes need custom handling but a bug might indicate
			// one that doesn't have custom handling probably should.
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

func (ui *usagesImp) processBlock(block *ast.BlockStmt) {
	for _, statement := range block.List {
		ui.processNode(statement)
		ui.flushPendingToRead()
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
	typ, ok := ui.info.Types[call.Fun]
	if !ok {
		panic(terror.New(`failed to find type info`).
			With(`function`, call.Fun))
	}

	// Check for explicit cast (conversion), e.g. `int(f.x)`
	if typ.IsType() {
		ui.processNode(call.Fun)
		ui.flushPendingToWrite()
		return
	}

	// Check for builtin method, e.g. `println(f.x)`
	if typ.IsBuiltin() {
		ui.processBuiltinCall(call)
		return
	}

	// Function invocation, e.g. `fmt.Println(f.x)`
	ui.processNode(call.Fun)
	if !utils.IsNil(ui.pendingCon) {
		ui.usages.Invokes.Add(ui.pendingCon)
		ui.clearPending()
	}
}

func (ui *usagesImp) processBuiltinCall(call *ast.CallExpr) {
	exp := call.Fun
	if p, ok := exp.(*ast.ParenExpr); ok {
		exp = p.X
	}
	if id, ok := exp.(*ast.Ident); ok {
		ui.clearPending()

		switch id.Name {
		case `append`, `cap`, `complex`, `copy`, `imag`, `len`,
			`make`, `max`, `min`, `new`, `real`, `recover`:
			if typ, ok := ui.info.Types[call]; ok {
				ui.pendingCon = ui.conv.ConvertType(typ.Type)
				ui.pendingEff = false
			}
			return

		case `clear`, `close`, `delete`, `panic`:
			return

		case `print`, `println`:
			ui.usages.SideEffect = true
			return
		}
	}
	panic(terror.New(`failed to get name of builtin function`).
		With(`expression`, exp))
}

func (ui *usagesImp) processCompositeLit(comp *ast.CompositeLit) {
	for _, elem := range comp.Elts {
		ui.processNode(elem)
		ui.flushPendingToRead()
	}
	ui.processNode(comp.Type)
	ui.clearPending() // flush the internally defined type
}

func (ui *usagesImp) processFunc(fn *ast.FuncType) {
	// This is part of a `ast.FuncLit` or `ast.FuncDecl`.
	if fn.Params != nil {
		ui.addArguments(fn.Params.List)
	}
	if fn.Results != nil {
		ui.addArguments(fn.Results.List)
	}
	if fn.TypeParams != nil {
		ui.addTypeParam(fn.TypeParams.List)
	}
}

func (ui *usagesImp) addArguments(args []*ast.Field) {
	for _, arg := range args {
		for _, id := range arg.Names {
			if obj := ui.info.Defs[id]; obj != nil {
				ui.localDefs[obj] = ui.proj.NewArgument(constructs.ArgumentArgs{
					Name: obj.Name(),
					Type: ui.conv.ConvertType(obj.Type()),
				})
			}
		}
	}
}

func (ui *usagesImp) addTypeParam(tps []*ast.Field) {
	for _, tp := range tps {
		for _, id := range tp.Names {
			if obj := ui.info.Defs[id]; obj != nil {
				ui.localDefs[obj] = ui.proj.NewTypeParam(constructs.TypeParamArgs{
					Name: obj.Name(),
					Type: ui.conv.ConvertType(obj.Type()),
				})
			}
		}
	}
}

func (ui *usagesImp) processFunDecl(fn *ast.FuncDecl) {
	if fn.Recv != nil {
		for _, recv := range fn.Recv.List {
			for _, id := range recv.Names {
				if obj := ui.info.Defs[id]; obj != nil {
					ui.localDefs[obj] = ui.conv.ConvertType(obj.Type())
				}
			}
		}
	}

	ui.processFunc(fn.Type)
	ui.processNode(fn.Body)
}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	// Check if this identifier is part of a local definition.
	if def, ok := ui.info.Defs[id]; ok {
		if def == nil {
			// Skip over `t` in `select `t := x.(type)` type definitions.
			return
		}
		if _, ok := def.(*types.Label); ok {
			// Skip over labels
			return
		}

		ui.flushPendingToRead()
		ui.pendingCon = ui.conv.ConvertType(def.Type())
		ui.addWrite(ui.pendingCon, false)
		ui.localDefs[def] = ui.pendingCon
		return
	}

	// Check for an identifier is being used.
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

	// Check if defined earlier and reuse defined type.
	if usage, ok := ui.localDefs[obj]; ok {
		ui.flushPendingToRead()
		ui.pendingCon = usage
		return
	}

	// Return builtin type.
	if obj.Pkg() == nil {
		if typ := ui.baker.TypeByName(obj.Name()); !utils.IsNil(typ) {
			ui.pendingCon = typ
			ui.pendingEff = false
			return
		}
	}

	// Return basic types as usage.
	if basic, ok := obj.Type().(*types.Basic); ok && basic.Kind() != types.Invalid {
		switch basic.Kind() {
		case types.Complex64:
			ui.pendingCon = ui.baker.BakeComplex64()
		case types.Complex128:
			ui.pendingCon = ui.baker.BakeComplex128()
		default:
			ui.pendingCon = ui.proj.NewBasic(constructs.BasicArgs{
				RealType: basic,
			})
		}
		ui.pendingEff = false
		return
	}

	// Create a temp reference for this object.
	ui.pendingEff = true
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

	ui.pendingCon = ui.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   pkgPath,
		Name:          obj.Name(),
		InstanceTypes: instType,
	})
}

func (ui *usagesImp) processIncDec(stmt *ast.IncDecStmt) {
	ui.processNode(stmt.X)
	if !utils.IsNil(ui.pendingCon) {
		ui.addWrite(ui.pendingCon, false)
		ui.clearPending()
	} else if typ, ok := ui.info.Types[stmt.X]; ok {
		// Handle `*(func())++` where the increment is on
		// the returned type, or `mapFoo["cat"]++`.
		ui.addWrite(ui.conv.ConvertType(typ.Type), false)
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
	ui.pendingEff = false
	if elem, ok := ui.info.Types[expr]; ok {
		if _, ok := elem.Type.(*types.Tuple); ok {
			// The indexing returned a (value, ok) tuple.
			// The results are probably used in a function parameter or an
			// assignment so the pending construct doesn't need to be set.
			ui.pendingCon = nil
			return
		}

		// The indexing returned a single value.
		ui.pendingCon = ui.conv.ConvertType(elem.Type)
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
	if !utils.IsNil(ui.pendingCon) {
		ui.usages.Reads.Add(ui.pendingCon)
	}
}

func (ui *usagesImp) processRange(r *ast.RangeStmt) {
	if r.Key != nil {
		if r.Tok == token.DEFINE {
			ui.processNode(r.Key)
			ui.processNode(r.Value)
		} else { // r.Tok == token.ASSIGN
			ui.processNode(r.Key)
			ui.flushPendingToWrite()
			ui.processNode(r.Value)
			ui.flushPendingToWrite()
		}
	}
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
	if utils.IsNil(ui.pendingCon) {
		selObj, ok := ui.info.Selections[sel]
		if !ok {
			panic(terror.New(`expected selection info but not found`).
				With(`expr`, sel))
		}

		// The left hand side is empty so this wasn't like `foo.Bar`,
		// it was instead like `foo().Bar` where some expression results
		// in a selection on a type. If a locally defined type then skip
		// the select since it won't reference anything useful.
		ui.pendingEff = false
		if _, ok := ui.localDefs[selObj.Obj()]; ok {
			ui.pendingCon = ui.conv.ConvertType(selObj.Type())
			return
		}
		ui.pendingCon = ui.conv.ConvertType(selObj.Recv())
	}

	ui.addRead(ui.pendingCon)
	ui.pendingCon = ui.proj.NewSelection(constructs.SelectionArgs{
		Name:   sel.Sel.Name,
		Origin: ui.pendingCon,
	})
	ui.pendingEff = false
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
