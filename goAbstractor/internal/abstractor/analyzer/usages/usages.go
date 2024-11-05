package usages

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/collections/stack"
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

	// SideEffect indicates that the usage definitely has a side effect
	// when true, however, false only means it isn't known yet.
	// Any invokes of another method that has a side effect means this
	// has a side effect too but will not be determined yet.
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

const printDebugLogging = false

func logDebug(format string, args ...any) {
	if printDebugLogging {
		fmt.Printf(format, args...)
		fmt.Println()
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

	compLits  collections.Stack[*ast.CompositeLit]
	pending   constructs.Construct
	pendingSE bool
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

	fSet := proj.Locs().FileSet()
	assert.ArgNotNil(`fSet`, fSet)

	start := root
	switch t := start.(type) {
	case *ast.FuncDecl:
		// Skip over receiver, parameter, returns, and type parameters.
		// Those will be visible in the abstraction data outside of usages.
		start = t.Body
	case *ast.FuncLit:
		start = t.Body
	}

	ui := &usagesImp{
		fSet:     fSet,
		info:     info,
		proj:     proj,
		curPkg:   curPkg,
		baker:    baker,
		conv:     conv,
		root:     root,
		usages:   newUsage(),
		compLits: stack.New[*ast.CompositeLit](),
	}

	ui.processNode(start)

	return ui.usages
}

func (ui *usagesImp) pos(pr posReader) token.Position {
	pos := token.NoPos
	if !utils.IsNil(pr) {
		pos = pr.Pos()
	}
	return ui.fSet.Position(pos)
}

func (ui *usagesImp) hasPending() bool {
	return !utils.IsNil(ui.pending)
}

func (ui *usagesImp) setPendingConstruct(c constructs.Construct) {
	ui.flushPendingToRead()
	logDebug(`  - PendingCon: %v`, c)
	ui.pending = c
}

func (ui *usagesImp) setPendingType(t types.Type) {
	ui.flushPendingToRead()
	logDebug(`  - PendingType: %v`, t)

	named := getNamed(t)
	if named == nil {
		logDebug(`    - no named`)
		return
	}

	logDebug(`    - named.Obj: %v @ %v`, named.Obj(), ui.pos(named.Obj()))
	if isLocal(ui.root, named.Obj()) {
		logDebug(`    - local`)
		return
	}

	ui.pending = ui.conv.ConvertType(named)
	logDebug(`    - converted: %v`, ui.pending)
}

func (ui *usagesImp) setPendingObject(o types.Object, instType []constructs.TypeDesc) {
	ui.flushPendingToRead()
	logDebug(`  - PendingObject: %v`, o)

	if utils.IsNil(o) {
		logDebug(`    + nil`)
		return
	}

	if _, ok := o.(*types.Label); ok {
		// Skip over labels
		logDebug(`    + label`)
		return
	}

	if _, ok := o.(*types.PkgName); ok {
		// Skip over package names
		logDebug(`    + package name`)
		return
	}

	if isLocal(ui.root, o) {
		logDebug(`    - local: set type`)
		ui.setPendingType(o.Type())
		return
	}

	if tn, ok := o.(*types.TypeName); ok {
		logDebug(`    - type name: %v: %v`, o, tn)

		pkgPath := getPkgPath(o)
		if len(pkgPath) <= 0 {
			if typ := ui.baker.TypeByName(o.Name()); !utils.IsNil(typ) {
				logDebug(`      - built-in type: %v`, typ)
				ui.setPendingConstruct(typ)
				return
			}

			if basic, ok := o.Type().(*types.Basic); ok && basic.Kind() != types.Invalid {
				logDebug(`      + basic: %v`, basic)
				return
			}

			logDebug(`      + dump built-in`)
			return
		}

		_, typ, found := ui.proj.FindType(pkgPath, o.Name(), instType, false)
		if found {
			logDebug(`      + type found: %v`, typ)
			ui.pending = typ
			return
		}

		logDebug(`      - temp ref: %v`, o)
		ui.pending = ui.proj.NewTempReference(constructs.TempReferenceArgs{
			RealType:      o.Type(),
			PackagePath:   pkgPath,
			Name:          o.Name(),
			InstanceTypes: instType,
			Package:       ui.curPkg.Source(),
		})
		return
	}

	if v, ok := o.(*types.Var); ok {
		logDebug(`    - type var: %v`, v)
		if v.IsField() {
			if compLit := ui.compLits.Peek(); !utils.IsNil(compLit) {
				compType := ui.info.Types[compLit.Type].Type
				logDebug(`      - field sel: %v => %v`, compType, v.Name())
				ui.setPendingType(compType)
				if ui.hasPending() {
					logDebug(`      - pending field selObj: %v`, ui.pending)
					ui.setPendingConstruct(ui.proj.NewSelection(constructs.SelectionArgs{
						Name:   v.Name(),
						Origin: ui.pending,
					}))
					return
				}
			}

			logDebug(`      - field without recv: %v`, v)
			ui.setPendingType(v.Type())
			return
		}

		// If not a field, then it is a global variable being either
		// read from or written two so it is a potential side effect.
		ui.pendingSE = true
	}

	pkgPath := getPkgPath(o)
	_, typ, found := ui.proj.FindDecl(pkgPath, o.Name(), instType, false)
	if found {
		logDebug(`      + decl found: %v`, typ)
		ui.pending = typ
		return
	}

	logDebug(`    - temp decl ref: %v`, o)
	ui.pending = ui.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   pkgPath,
		Name:          o.Name(),
		InstanceTypes: instType,
	})
}

func (ui *usagesImp) clearPending() {
	ui.pending = nil
	ui.pendingSE = false
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
	ui.addWrite(ui.pending, ui.pendingSE)
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
		logDebug(`  + Reads: %v`, c)
		ui.usages.Reads.Add(c)
	}
}

// addWrite adds the given construct as a write usage.
func (ui *usagesImp) addWrite(c constructs.Construct, sideEffect bool) {
	if !utils.IsNil(c) {
		logDebug(`  + Write: %v`, c)
		ui.usages.Writes.Add(c)
		if sideEffect {
			ui.usages.SideEffect = true
		}
	}
}

// addInvoke adds the given construct as an invoke usage.
func (ui *usagesImp) addInvoke(c constructs.Construct) {
	if !utils.IsNil(c) {
		logDebug(`  + Invoke: %v`, c)
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
		case *ast.KeyValueExpr:
			ui.processKeyValue(t)
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
			//logDebug(`usagesImp.processNode unhandled (%[1]T) %[1]v`, t)
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
	logDebug(`>>> processCompositeLit: %v @ %s`, comp, ui.pos(comp))
	ui.compLits.Push(comp)
	defer ui.compLits.Pop()

	for _, elem := range comp.Elts {
		ui.processNode(elem)
		ui.flushPendingToWrite()
	}

	ui.processNode(comp.Type)
	if ui.hasPending() {
		ui.addWrite(ui.pending, ui.pendingSE)
	}
}

func (ui *usagesImp) processIdent(id *ast.Ident) {
	logDebug(`>>> processIdent: %v @ %s`, id, ui.pos(id))

	// Check if this identifier is part of a definition.
	if def, ok := ui.info.Defs[id]; ok {
		logDebug(`  > processIdent: def object: %v`, def)
		ui.setPendingObject(def, nil)
		ui.addWrite(ui.pending, ui.pendingSE)
		return
	}

	// Check if the identifier is being used.
	obj, ok := ui.info.Uses[id]
	if !ok {
		logDebug(`  > processIdent: no uses`)
		return
	}

	// Skip over build-in constants:
	// https://pkg.go.dev/builtin#pkg-constants
	switch obj.Id() {
	case `_.true`, `_.false`, `_.nil`, `_.iota`:
		logDebug(`  > processIdent: build-in constants: %v`, obj.Id())
		return
	}

	// Return built-in type.
	if obj.Pkg() == nil {
		if typ := ui.baker.TypeByName(obj.Name()); !utils.IsNil(typ) {
			logDebug(`  > processIdent: build-in type: %v`, typ)
			ui.setPendingConstruct(typ)
			return
		}
	}

	// Find any type arguments for this object.
	var instType []constructs.TypeDesc
	if itList := ui.info.Instances[id].TypeArgs; !utils.IsNil(itList) {
		instType = ui.conv.ConvertInstanceTypes(itList)
	} else if itList := getInstTypes(obj); !utils.IsNil(itList) {
		instType = ui.conv.ConvertInstanceTypes(itList)
	}

	logDebug(`  > processIdent: object: %v`, obj)
	ui.setPendingObject(obj, instType)
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

func (ui *usagesImp) processKeyValue(kv *ast.KeyValueExpr) {
	logDebug(`>>> processKeyValue: %v @ %s`, kv, ui.pos(kv))
	ui.processNode(kv.Key)
	ui.flushPendingToWrite()

	ui.processNode(kv.Value)
	ui.flushPendingToRead()
}

func (ui *usagesImp) processRange(r *ast.RangeStmt) {
	logDebug(`>>> processRange: %v @ %s`, r, ui.pos(r))
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
	logDebug(`>>> processSelector: %v @ %s`, sel, ui.pos(sel))
	ui.processNode(sel.X)
	logDebug(`>>> processSelector.X: %v`, ui.pending)
	if ui.hasPending() {
		logDebug(`  > pending sel.X: %v`, ui.pending)
		ui.setPendingConstruct(ui.proj.NewSelection(constructs.SelectionArgs{
			Name:   sel.Sel.Name,
			Origin: ui.pending,
		}))
		return
	}

	selObj, ok := ui.info.Selections[sel]
	if !ok {
		logDebug(`  > no selection info: %v`, sel)
		return
	}
	logDebug(`  > selObj: %v`, selObj)

	if !isLocal(ui.root, selObj.Obj()) {
		logDebug(`  > non-local selObj: %v at %v`, selObj.Obj(), ui.pos(selObj.Obj()))
		ui.setPendingConstruct(ui.conv.ConvertType(selObj.Recv()))
		if ui.hasPending() {
			logDebug(`  > pending selObj: %v`, ui.pending)
			ui.setPendingConstruct(ui.proj.NewSelection(constructs.SelectionArgs{
				Name:   sel.Sel.Name,
				Origin: ui.pending,
			}))
			return
		}
	}

	logDebug(`  > selection fallback: %v`, selObj)
	ui.setPendingType(selObj.Recv())
	ui.flushPendingToRead()
	ui.setPendingType(selObj.Obj().Type())
}

func (ui *usagesImp) processTypeAssert(exp *ast.TypeAssertExpr) {
	ui.processNode(exp.X)
	ui.flushPendingToRead()

	if utils.IsNil(exp.Type) {
		// Type assertions for type switches, e.g `switch t := x.(type)`.
		return
	}

	tv, ok := ui.info.Types[exp.Type]
	if !ok {
		panic(terror.New(`Expected a type in a TypeAssert for usages.`).
			With(`type`, types.ExprString(exp.Type)).
			With(`x`, types.ExprString(exp.X)).
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
