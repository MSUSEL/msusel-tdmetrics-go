package analyzer

import (
	"fmt"
	"go/parser"
	"go/token"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

func Test_Empty(t *testing.T) {
	m := parseExpr(t,
		`func() {}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  1,
		Complexity: 1,
		Indents:    0,
		LineCount:  1,
	})
}

func Test_Simple(t *testing.T) {
	m := parseExpr(t,
		`func() int {`,
		`	return max(1, 3, 5)`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  3,
		Complexity: 1,
		Indents:    1,
		LineCount:  3,
	})
}

func Test_SimpleWithExtraIndent(t *testing.T) {
	m := parseExpr(t,
		`		func() int {`,
		`			return max(1, 3, 5)`,
		`		}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  3,
		Complexity: 1,
		Indents:    1,
		LineCount:  3,
	})
}

func Test_SimpleParams(t *testing.T) {
	m := parseExpr(t,
		`func(a int,`,
		`	  b int,`,
		`	  c int) int {`,
		`	return max(a, b, c)`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  5,
		Complexity: 1,
		Indents:    7,
		LineCount:  5,
	})
}

func Test_SimpleWithReturn(t *testing.T) {
	m := parseExpr(t,
		`func() int {`,
		`	x := max(1, 3, 5)`,
		`	return x`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  4,
		Complexity: 1,
		Indents:    2,
		LineCount:  4,
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  3,
		Complexity: 1,
		Indents:    1,
		LineCount:  8,
	})
}

func Test_SimpleWithDefer(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := open()`,
		`	defer x.close()`,
		`	x.doStuff()`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  5,
		Complexity: 1,
		Indents:    3,
		LineCount:  5,
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  7,
		Complexity: 2,
		Indents:    6,
		LineCount:  7,
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  10,
		Complexity: 2,
		Indents:    11,
		LineCount:  10,
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  9,
		Complexity: 3,
		Indents:    9,
		LineCount:  9,
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  11,
		Complexity: 3,
		Indents:    12,
		LineCount:  11,
	})
}

func Test_SimpleSwitch(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`   switch {`,
		`	case x > 7:`,
		`		x = 4`,
		`	case x > 4:`,
		`		x = 3`,
		`	default:`,
		`		x = 2`,
		`	}`,
		`	println(x)`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  12,
		Complexity: 3,
		Indents:    15,
		LineCount:  12,
	})
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  14,
		Complexity: 1,
		Indents:    19,
		LineCount:  14,
	})
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  14,
		Complexity: 1,
		Indents:    19,
		LineCount:  14,
	})
}

func Test_DeferWithComplexity(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	print("A ")`,
		`	defer func() {`,
		`		if r := recover(); r != nil {`,
		`			print("B ")`,
		`			return`,
		`		}`,
		`		print("C ")`,
		`	}()`,
		`	print("D ")`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  11,
		Complexity: 2,
		Indents:    16,
		LineCount:  11,
	})
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
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  11,
		Complexity: 2,
		Indents:    15,
		LineCount:  11,
	})
}

func Test_GoStatement(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	go func() {`,
		`		print("A ")`,
		`	}()`,
		`	print("B ")`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  6,
		Complexity: 2,
		Indents:    5,
		LineCount:  6,
	})
}

func Test_SelectStatement(t *testing.T) {
	m := parseExpr(t,
		`func() {`,
		`	var A, B chan bool`,
		`	go func() {`,
		`		A <- true`,
		`	}()`,
		`	go func() {`,
		`		B <- true`,
		`	}()`,
		`	select {`,
		`	case <- A:`,
		`		print("A ")`,
		`	case b := <- B:`,
		`		print("B ", b)`,
		`	}`,
		`}`)
	checkMetrics(t, m, constructs.MetricsArgs{
		CodeCount:  15,
		Complexity: 5,
		Indents:    17,
		LineCount:  15,
	})
}

// TODO: Test joining metrics:
// var val = []int{
//   func() int { ** }(),
//   func() int { ** }(),
// }

// TODO: Test nothing variable:
// var val = 10

// TODO: Test reading metrics with only read reference:
// var val = singleton.f()

// TODO: Test reading metrics with only read reference but without a function:
// var val = singleton.value

// TODO: Test reading metrics with read reference as parameter:
// var val = func(f Foo) int { ** }(singleton.f)

// TODO: Test reading metrics with read reference in typed call:
// var val = Foo[int](singleton)

func parseExpr(t *testing.T, lines ...string) constructs.MetricsArgs {
	code := strings.Join(lines, "\n")
	fSet := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)

	an := New(locs.NewSet(fSet))
	an.Analyze(expr)
	metrics := an.GetMetrics()
	metrics.Location = nil // ignore locations
	return metrics
}

func checkMetrics(t *testing.T, m, exp constructs.MetricsArgs) {
	ctx := jsonify.NewContext()

	mMet := metrics.New().NewMetrics(m)
	gotData, err := jsonify.Marshal(ctx, mMet)
	check.NoError(t).Assert(err)

	expMet := metrics.New().NewMetrics(exp)
	expData, err := jsonify.Marshal(ctx, expMet)
	check.NoError(t).Assert(err)

	if !slices.Equal(gotData, expData) {
		gotLines := strings.Split(string(gotData), "\n")
		expLines := strings.Split(string(expData), "\n")
		d := diff.Default().PlusMinus(gotLines, expLines)
		fmt.Println(strings.Join(d, "\n"))
		t.Fail()
	}
}
