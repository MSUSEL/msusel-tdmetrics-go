package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"math"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/usage"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// TODO: Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

func Analyze(locs locs.Set, info *types.Info, factory constructs.MetricsFactory, node ast.Node) constructs.Metrics {
	return factory.NewMetrics(newAnalyzer(locs, info).Analyze(node).GetMetricsArgs())
}

type analyzerImp struct {
	locs locs.Set
	info *types.Info
	loc  locs.Loc

	complexity int
	minLine    int
	maxLine    int
	indents    int
	minColumn  map[int]int
	getter     bool
	setter     bool

	reads   collections.SortedSet[constructs.Usage]
	writes  collections.SortedSet[constructs.Usage]
	invokes collections.SortedSet[constructs.Usage]
	defines collections.SortedSet[constructs.Usage]
}

func newAnalyzer(locs locs.Set, info *types.Info) *analyzerImp {
	return &analyzerImp{
		locs: locs,
		info: info,
		loc:  nil,

		complexity: 1,
		maxLine:    0,
		minLine:    math.MaxInt,
		indents:    0,
		minColumn:  map[int]int{},

		reads:   sortedSet.New(usage.Comparer()),
		writes:  sortedSet.New(usage.Comparer()),
		invokes: sortedSet.New(usage.Comparer()),
		defines: sortedSet.New(usage.Comparer()),
	}
}

func (a *analyzerImp) Analyze(node ast.Node) *analyzerImp {
	if utils.IsNil(a.loc) {
		a.loc = a.locs.NewLoc(node.Pos())
	}
	// gather positional information for indents and cyclomatic complexity.
	ast.Inspect(node, a.addCodePosForNode)
	a.getUsages()
	a.getter = a.checkForGetter(node)
	a.setter = a.checkForSetter(node)
	return a
}

func (a *analyzerImp) GetMetricsArgs() constructs.MetricsArgs {
	return constructs.MetricsArgs{
		Location:   a.loc,
		Complexity: a.complexity,
		LineCount:  a.maxLine - a.minLine + 1,
		CodeCount:  len(a.minColumn),
		Indents:    a.calcIndents(),
		Getter:     a.getter,
		Setter:     a.setter,
		Reads:      a.reads,
		Writes:     a.writes,
		Invokes:    a.invokes,
		Defines:    a.defines,
	}
}

func (a *analyzerImp) calcIndents() int {
	leftMostColumn := math.MaxInt
	indentSum := 0
	for _, ind := range a.minColumn {
		leftMostColumn = min(ind, leftMostColumn)
		indentSum += ind
	}
	return indentSum - len(a.minColumn)*leftMostColumn
}

func (a *analyzerImp) incComplexity(check bool) {
	if check {
		a.complexity++
	}
}

func (a *analyzerImp) addCodePos(pos token.Pos, isEnd bool) {
	p := a.locs.FileSet().PositionFor(pos, false)
	lineNo, column := p.Line, p.Column
	a.maxLine = max(a.maxLine, lineNo)
	a.minLine = min(a.minLine, lineNo)
	if isEnd {
		column--
	}
	if otherCol, ok := a.minColumn[lineNo]; ok {
		column = min(column, otherCol)
	}
	a.minColumn[lineNo] = column
}

func (a *analyzerImp) addCodePosForNode(n ast.Node) bool {
	switch t := n.(type) {
	case nil, *ast.Comment, *ast.CommentGroup:
		return true
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.GoStmt:
		a.incComplexity(true)
	case *ast.CaseClause:
		a.incComplexity(t.List != nil)
	case *ast.CommClause:
		a.incComplexity(t.Comm != nil)
	case *ast.BinaryExpr:
		a.incComplexity(t.Op == token.LAND || t.Op == token.LOR)
	}

	a.addCodePos(n.Pos(), false)
	if ended, has := n.(interface{ End() token.Pos }); has {
		a.addCodePos(ended.End(), true)
	}
	return true
}

func (a *analyzerImp) getUsages() {

	//TODO: implement
}

func getTypeAndBody(n ast.Node) (*ast.FuncType, *ast.BlockStmt, bool) {
	switch t := n.(type) {
	case *ast.FuncDecl:
		return t.Type, t.Body, true
	case *ast.FuncLit:
		return t.Type, t.Body, true
	default:
		return nil, nil, false
	}
}

// isSimpleFetch determines if the given node is a simple access of data
// without method calls and modification. The info must be populated
// with `Types` so that explicit casts (conversions) can be distinctly
// determined from method/function calls.
func isSimpleFetch(info *types.Info, node ast.Node) bool {
	valid := true
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case nil, *ast.Ident, *ast.SelectorExpr, *ast.BasicLit, *ast.StarExpr, *ast.TypeAssertExpr:
			// Check for identifiers (`foo`), selectors (`f.x`), literals (`3.24`),
			// dereference (`*f`), type assert (`f.(int)`), and ignore nils.
			break
		case *ast.CallExpr:
			// Check for explicit cast (conversion), e.g. `int(f.x)`
			tx, ok := info.Types[t.Fun]
			valid = valid && ok && tx.IsType()
		case *ast.UnaryExpr:
			// Check for reference, e.g. `&pointer`
			valid = valid && t.Op == token.AND
		default:
			// Anything else (add, subtract, indexer, bitwise-Xor)
			// means not a simple fetch.
			valid = false
		}
		return valid
	})
	return valid
}

// isObjectUsed determines if the given obj is used somewhere in the
// branch of the AST starting at the given node. The info must be populated
// with `Uses` so that any usage of an identifier can be compared to the object.
func isObjectUsed(obj types.Object, info *types.Info, node ast.Node) bool {
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		id, ok := n.(*ast.Ident)
		found = found || (ok && info.Uses[id] == obj)
		return !found
	})
	return found
}

// checkForGetter determines if this is code for a getter.
// See MetricsArgs.Getter in constructs/metrics.go for more info.
func (a *analyzerImp) checkForGetter(n ast.Node) bool {
	funcType, funcBody, ok := getTypeAndBody(n)
	if !ok ||
		len(funcType.Params.List) != 0 ||
		funcType.Results == nil ||
		len(funcType.Results.List) != 1 ||
		len(funcBody.List) != 1 {
		return false
	}

	ret, ok := funcBody.List[0].(*ast.ReturnStmt)
	if !ok ||
		len(ret.Results) != 1 ||
		!isSimpleFetch(a.info, ret.Results[0]) {
		return false
	}

	return true
}

func (a *analyzerImp) checkForSetter(n ast.Node) bool {
	funcType, funcBody, ok := getTypeAndBody(n)
	if !ok ||
		len(funcType.Params.List) > 1 ||
		funcType.Results != nil ||
		len(funcBody.List) != 1 {
		return false
	}

	assign, ok := funcBody.List[0].(*ast.AssignStmt)
	if !ok ||
		len(assign.Lhs) != 1 ||
		len(assign.Rhs) != 1 ||
		!isSimpleFetch(a.info, assign.Lhs[0]) ||
		!isSimpleFetch(a.info, assign.Rhs[0]) {
		return false
	}

	if len(funcType.Params.List) == 0 {
		return true
	}

	if len(funcType.Params.List[0].Names) != 1 {
		return false
	}

	paramId := funcType.Params.List[0].Names[0]
	if constructs.BlankName(paramId.Name) {
		return true
	}

	// Make sure the parameter isn't used on the left hand side as in a
	// reversed setter, e.g. `func (b Bar) GetX(x *int) { x* = b.x }`,
	// The parameter may be used on the right hand side or not at all.
	// The parameter may not be used at all if the setter is part of an
	// interface requirement but the value assigned is to a default value.
	return !isObjectUsed(a.info.Defs[paramId], a.info, assign.Lhs[0])
}
