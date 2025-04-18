package accessor

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Accessor struct {
	Getter bool
	Setter bool
}

func Calculate(log *logger.Logger, info *types.Info, node ast.Node) Accessor {
	assert.ArgNotNil(`info`, info)
	assert.ArgNotNil(`info.Defs`, info.Defs)
	assert.ArgNotNil(`info.Types`, info.Types)
	assert.ArgNotNil(`info.Uses`, info.Uses)
	assert.ArgNotNil(`node`, node)

	log.Logf(`accessor`)

	k := Accessor{}
	if funcType, funcBody, ok := getTypeAndBody(node); ok {
		switch {
		case isGetter(info, funcType, funcBody):
			k.Getter = true
		case isSetter(info, funcType, funcBody):
			k.Setter = true
		}
	}
	return k
}

func getTypeAndBody(node ast.Node) (*ast.FuncType, *ast.BlockStmt, bool) {
	switch t := node.(type) {
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
		if !valid {
			return false
		}
		switch t := n.(type) {
		case nil, *ast.Ident, *ast.SelectorExpr, *ast.ParenExpr,
			*ast.BasicLit, *ast.StarExpr, *ast.TypeAssertExpr:
			// Check for identifiers (`foo`), selectors (`f.x`), literals (`3.24`),
			// dereference (`*f`), type assert (`f.(int)`), parenthesis, and ignore nils.
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
		if found {
			return false
		}
		id, ok := n.(*ast.Ident)
		if ok && info.Uses[id] == obj {
			found = true
			return false
		}
		return true
	})
	return found
}

// isGetter determines if this is code for a getter.
// See MetricsArgs.Getter in constructs/metrics.go for more info.
// The info must be populated with `Types`.
//
// Check that there is only one statement that is a return statement,
// one result, no parameters, and is a simple fetch for the result.
//
// The function does not need a receiver and may return a package level
// variable or constant, since this would work like a static getter
// on a static class.
func isGetter(info *types.Info, funcType *ast.FuncType, funcBody *ast.BlockStmt) bool {
	var ret *ast.ReturnStmt
	return len(funcType.Params.List) == 0 &&
		funcType.Results != nil &&
		len(funcType.Results.List) == 1 &&
		funcBody != nil &&
		len(funcBody.List) == 1 &&
		utils.Is(funcBody.List[0], &ret) &&
		len(ret.Results) == 1 &&
		isSimpleFetch(info, ret.Results[0])
}

// isSetter determines if this is code for a setter.
// See MetricsArgs.Setter in constructs/metrics.go for more info.
// The info must be populated with `Uses`, `Defs`, and `Types`.
func isSetter(info *types.Info, funcType *ast.FuncType, funcBody *ast.BlockStmt) bool {
	var assign *ast.AssignStmt
	if len(funcType.Params.List) > 1 ||
		funcType.Results != nil ||
		funcBody == nil ||
		len(funcBody.List) != 1 ||
		!utils.Is(funcBody.List[0], &assign) ||
		len(assign.Lhs) != 1 ||
		len(assign.Rhs) != 1 ||
		!isSimpleFetch(info, assign.Lhs[0]) ||
		!isSimpleFetch(info, assign.Rhs[0]) {
		// Check that there is zero or one parameters, zero results,
		// only statement in the body of the function, that there is
		// only one assignment, and both the left and right hand sides
		// must be simple fetches.
		return false
	}

	if len(funcType.Params.List) == 0 {
		// Setters may have no parameters for assigning a literal value.
		// e.g. `func(b *Bar) Hide() { b.visible = false }`
		return true
	}

	if len(funcType.Params.List[0].Names) != 1 {
		// Check that the type group in the single parameter type is only
		// for one parameter. e.g. not `func(x, y int)`.
		return false
	}

	paramId := funcType.Params.List[0].Names[0]
	if constructs.BlankName(paramId.Name) {
		// Check if single parameter is blank and therefore not used.
		return true
	}

	// Make sure the parameter isn't used on the left hand side as in a
	// reversed setter, e.g. `func (b Bar) GetX(x *int) { x* = b.x }`,
	// The parameter may be used on the right hand side or not at all.
	// The parameter may not be used at all if the setter is part of an
	// interface requirement but the value assigned is to a default value.
	return !isObjectUsed(info.Defs[paramId], info, assign.Lhs[0])
}
