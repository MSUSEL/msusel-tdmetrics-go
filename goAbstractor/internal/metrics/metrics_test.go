package metrics

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Empty(t *testing.T) {
	m := parseExpr(t,
		`func() {}`)
	checkMetrics(t, m, Metrics{
		Complexity: 1,
		LineCount:  1,
		CodeCount:  1,
		Indents:    0,
	})
}

func Test_Simple(t *testing.T) {
	m := parseExpr(t,
		`func() int {`,
		`	return max(1, 3, 5)`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 1,
		LineCount:  3,
		CodeCount:  1,
		Indents:    1,
	})
}

func Test_SimpleWithReturn(t *testing.T) {
	m := parseExpr(t,
		`func() int {`,
		`	x := max(1, 3, 5)`,
		`	return x`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 1,
		LineCount:  4,
		CodeCount:  2,
		Indents:    2,
	})
}

func Test_SimpleWithSpace(t *testing.T) {
	m := parseExpr(t,
		`func() int {`,
		`   // Bacon is tasty`,
		`   `,
		`	return max(1, 3, 5)`,
		`   /* This is not a comment`,
		`	   it is a sandwich */`,
		`   `,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 1,
		LineCount:  8,
		CodeCount:  1,
		Indents:    1,
	})
}

func Test_SimpleWithDefer(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := open()`,
		`	defer x.close()`,
		`	x.doStuff()`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 1,
		LineCount:  5,
		CodeCount:  3,
		Indents:    3,
	})
}

func Test_SimpleIf(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	}`,
		`	println(x)`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 2,
		LineCount:  7,
		CodeCount:  5,
		Indents:    5,
	})
}

func Test_SimpleIfElse(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	} else {`,
		`		x = 2`,
		`		print("cat")`,
		`	}`,
		`	println(x)`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 2,
		LineCount:  10,
		CodeCount:  4,
		Indents:    11,
	})
}

func Test_SimpleIfElseIf(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	} else if x > 4 {`,
		`		x = 3`,
		`	}`,
		`	println(x)`,
		`}`)
	checkMetrics(t, m, Metrics{
		Complexity: 3,
		LineCount:  9,
		CodeCount:  6,
		Indents:    9,
	})
}

func Test_SimpleIfElseIfElse(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	} else if x > 4 {`,
		`		x = 3`,
		`	} else {`,
		`		x = 2`,
		`	}`,
		`	println(x)`,
		`}`)
	checkMetrics(t, m, Metrics{})
}

func Test_DeferInBlock(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	print("A ")`,
		`	defer func() {`,
		`		print("B ")`,
		`	}()`,
		`	{`,
		`		print("C ")`,
		`		defer func() {`,
		`			print("D ")`,
		`		}()`,
		`		print("E ")`,
		`	}`,
		`	print("F ")`,
		`}`)
	checkMetrics(t, m, Metrics{})
}

func Test_DeferInFuncLiteral(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	print("A ")`,
		`	defer func() {`,
		`		print("B ")`,
		`	}()`,
		`	func() {`,
		`		print("C ")`,
		`		defer func() {`,
		`			print("D ")`,
		`		}()`,
		`		print("E ")`,
		`	}()`,
		`	print("F ")`,
		`}`)
	checkMetrics(t, m, Metrics{})
}

func Test_DeferWithComplexity(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	print("A ")`,
		`	defer func() {`,
		`		 if r := recover(); r != nil {`,
		`			print("B ")`,
		`			return`,
		`		 }`,
		`		print("C ")`,
		`	}()`,
		`	print("D ")`,
		`}`)
	checkMetrics(t, m, Metrics{})
}

func Test_ForRangeWithDefer(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	print("A ")`,
		`	for _ = range 4 {`,
		`		print("B ")`,
		`		defer func() {`,
		`			print("C ")`,
		`		}()`,
		`		print("D ")`,
		`	}`,
		`	print("E ")`,
		`}`)
	checkMetrics(t, m, Metrics{})
}

func parseExpr(t *testing.T, lines ...string) Metrics {
	code := strings.Join(lines, "\n")
	fSet := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)
	block := expr.(*ast.FuncLit).Body
	return New(fSet, block)
}

func checkMetrics(t *testing.T, m, exp Metrics) {
	gotLines := m.String()
	expLines := exp.String()
	if gotLines != expLines {
		diff := diff.Default().PlusMinus(strings.Split(gotLines, "\n"), strings.Split(expLines, "\n"))
		fmt.Println(strings.Join(diff, "\n"))
		t.Fail()
	}
}
