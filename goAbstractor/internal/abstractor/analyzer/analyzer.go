package analyzer

import (
	"go/ast"
	"go/token"
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

func Analyze(locs locs.Set, factory constructs.MetricsFactory, node ast.Node) constructs.Metrics {
	return factory.NewMetrics(newAnalyzer(locs).Analyze(node).GetMetricsArgs())
}

type analyzerImp struct {
	locs locs.Set
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

func newAnalyzer(locs locs.Set) *analyzerImp {
	return &analyzerImp{
		locs: locs,
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
	a.getter = checkForGetter(node)
	a.setter = checkForSetter(node)
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

func isSimpleFetch(n ast.Node) bool {
	valid := true
	ast.Inspect(n, func(n2 ast.Node) bool {
		switch t := n2.(type) {
		case nil, *ast.Ident, *ast.SelectorExpr, *ast.BasicLit, *ast.StarExpr, *ast.TypeAssertExpr:
			return valid
		// TODO: Implement explicit casts (type conversions)
		//case *ast.CallExpr:
		case *ast.UnaryExpr:
			valid = valid && t.Op == token.AND
		default:
			valid = false
		}
		return valid
	})
	return valid
}

// checkForGetter determines if this is code for a getter.
// See MetricsArgs.Getter in constructs/metrics.go for more info.
func checkForGetter(n ast.Node) bool {
	funcType, funcBody, ok := getTypeAndBody(n)
	if !ok || len(funcType.Params.List) != 0 ||
		len(funcType.Results.List) != 1 || len(funcBody.List) != 1 {
		return false
	}

	ret, ok := funcBody.List[0].(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 || !isSimpleFetch(ret.Results[0]) {
		return false
	}

	return true
}

func checkForSetter(n ast.Node) bool {
	funcType, funcBody, ok := getTypeAndBody(n)
	if !ok || len(funcType.Params.List) > 1 ||
		len(funcType.Results.List) != 0 || len(funcBody.List) != 1 {
		return false
	}

	assign, ok := funcBody.List[0].(*ast.AssignStmt)
	if !ok || len(assign.Lhs) != 1 || len(assign.Rhs) != 1 ||
		!isSimpleFetch(assign.Lhs[0]) || !isSimpleFetch(assign.Rhs[0]) {
		return false
	}

	if len(funcType.Params.List) == 0 {
		return true
	}

	if len(funcType.Params.List[0].Names) != 1 {
		return false
	}

	paramName := funcType.Params.List[0].Names[0].Name
	if constructs.BlankName(paramName) {
		return true
	}

	// TODO: Finish checking if `paramName is in right side only and failing
	//       with the left side (`func Foo(b* Bar) { b.x = b.y}`).

	return true
}
