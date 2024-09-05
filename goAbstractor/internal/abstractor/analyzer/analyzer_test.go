package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/differs/diff"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
	"gopkg.in/yaml.v3"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

func Test_Empty(t *testing.T) {
	tt := parseExpr(t,
		`func() {}`)
	tt.check(
		`{`,
		`	codeCount:  1,`,
		`	complexity: 1,`,
		`	#indents:   0,`,
		`	lineCount:  1`,
		`}`)
}

func Test_Simple(t *testing.T) {
	tt := parseExpr(t,
		`func() int {`,
		`	return max(1, 3, 5)`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`}`)
}

func Test_SimpleWithExtraIndent(t *testing.T) {
	tt := parseExpr(t,
		`		func() int {`,
		`			return max(1, 3, 5)`,
		`		}`)
	tt.check(
		`{`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`}`)
}

func Test_SimpleParams(t *testing.T) {
	tt := parseExpr(t,
		`func(a int,`,
		`	  b int,`,
		`	  c int) int {`,
		`	return max(a, b, c)`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  5,`,
		`	complexity: 1,`,
		`	indents:    7,`,
		`	lineCount:  5,`,
		`}`)
}

func Test_SimpleWithReturn(t *testing.T) {
	tt := parseExpr(t,
		`func() int {`,
		`	x := max(1, 3, 5)`,
		`	return x`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  4,`,
		`	complexity: 1,`,
		`	indents:    2,`,
		`	lineCount:  4,`,
		`}`)
}

func Test_SimpleWithSpace(t *testing.T) {
	tt := parseExpr(t,
		`func() int {`,
		`   // Bacon is tasty`,
		`   `,
		`	return max(1, 3, 5)`,
		`   /* This is not a comment`,
		`	   it is a sandwich */`,
		`   `,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  8,`,
		`}`)
}

func Test_SimpleWithDefer(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct {}`,
		`func (b Bar) close() { }`,
		`func (b Bar) doStuff() { }`,
		`func open() Bar { return Bar{} }`,
		`func Foo() {`,
		`	x := open()`,
		`	defer x.close()`,
		`	x.doStuff()`,
		`}`)
	tt.check(
		`{`,
		`   loc:        5,`,
		`	codeCount:  5,`,
		`	complexity: 1,`,
		`	indents:    3,`,
		`	lineCount:  5,`,
		`}`)
}

func Test_SimpleIf(t *testing.T) {
	tt := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	}`,
		`	println(x)`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  7,`,
		`	complexity: 2,`,
		`	indents:    6,`,
		`	lineCount:  7,`,
		`}`)
}

func Test_SimpleIfElse(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  10,`,
		`	complexity:  2,`,
		`	indents:    11,`,
		`	lineCount:  10,`,
		`}`)
}

func Test_SimpleIfElseIf(t *testing.T) {
	tt := parseExpr(t,
		`func() {`,
		`	x := 9`,
		`	if x > 7 {`,
		`		x = 4`,
		`	} else if x > 4 {`,
		`		x = 3`,
		`	}`,
		`	println(x)`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  9,`,
		`	complexity: 3,`,
		`	indents:    9,`,
		`	lineCount:  9,`,
		`}`)
}

func Test_SimpleIfElseIfElse(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  11,`,
		`	complexity:  3,`,
		`	indents:    12,`,
		`	lineCount:  11,`,
		`}`)
}

func Test_SimpleSwitch(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  12,`,
		`	complexity:  3,`,
		`	indents:    15,`,
		`	lineCount:  12,`,
		`}`)
}

func Test_DeferInBlock(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  14,`,
		`	complexity:  1,`,
		`	indents:    19,`,
		`	lineCount:  14,`,
		`}`)
}

func Test_DeferInFuncLiteral(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  14,`,
		`	complexity:  1,`,
		`	indents:    19,`,
		`	lineCount:  14,`,
		`}`)
}

func Test_DeferWithComplexity(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  11,`,
		`	complexity:  2,`,
		`	indents:    16,`,
		`	lineCount:  11,`,
		`}`)
}

func Test_ForRangeWithDefer(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  11,`,
		`	complexity:  2,`,
		`	indents:    15,`,
		`	lineCount:  11,`,
		`}`)
}

func Test_GoStatement(t *testing.T) {
	tt := parseExpr(t,
		`func() {`,
		`	go func() {`,
		`		print("A ")`,
		`	}()`,
		`	print("B ")`,
		`}`)
	tt.check(
		`{`,
		`	codeCount:  6,`,
		`	complexity: 2,`,
		`	indents:    5,`,
		`	lineCount:  6,`,
		`}`)
}

func Test_SelectStatement(t *testing.T) {
	tt := parseExpr(t,
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
	tt.check(
		`{`,
		`	codeCount:  15,`,
		`	complexity:  5,`,
		`	indents:    17,`,
		`	lineCount:  15,`,
		`}`)
}

func Test_GetterWithSelect(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x int }`,
		`func (b Bar) Foo() int {`,
		`	return b.x`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`   getter:  true,`,
		`}`)
}

func Test_GetterWithDereference(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo() int {`,
		`	return *bar`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`   getter:  true,`,
		`}`)
}

func Test_GetterWithConvert(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x float64 }`,
		`func (b Bar) Foo() int {`,
		`	return int(b.x)`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`   getter:  true,`,
		`}`)
}

func Test_SetterWithSelect(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x int }`,
		`func (b Bar) Foo(x int) {`,
		`	b.x = x`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`   setter:  true,`,
		`}`)
}

func Test_SetterWithReference(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x int) {`,
		`	bar = &x`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`   setter:  true,`,
		`}`)
}

func Test_NotReverseSetter(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x *int) {`,
		`	*x = *bar`,
		`}`)
	tt.check(
		`{`,
		`	loc:        2,`,
		`	codeCount:  3,`,
		`	complexity: 1,`,
		`	indents:    1,`,
		`	lineCount:  3,`,
		`}`)
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

// TODO: Test parentheses in getters and setters.

type testTool struct {
	t *testing.T
	m constructs.Metrics
}

func createInfo() *types.Info {
	return &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
}

func findNode(src ast.Node, name string) ast.Node {
	found := false
	var target ast.Node
	ast.Inspect(src, func(n ast.Node) bool {
		if found {
			return false
		}
		switch t := n.(type) {
		case *ast.Ident:
			if t.Name == name {
				found = true
				return false
			}
		case *ast.FuncDecl, *ast.TypeSpec, *ast.ValueSpec:
			target = t
		}
		return true
	})
	return target
}

func parseExpr(t *testing.T, lines ...string) *testTool {
	code := strings.Join(lines, "\n")
	fSet := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)

	info := createInfo()
	err = types.CheckExpr(fSet, nil, token.NoPos, expr, info)
	check.NoError(t).Require(err)

	m := Analyze(locs.NewSet(fSet), info, metrics.New(), expr)
	return &testTool{t: t, m: m}
}

func parseDecl(t *testing.T, name string, lines ...string) *testTool {
	code := "package test\n" + strings.Join(lines, "\n")
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)

	info := createInfo()
	var conf types.Config
	_, err = conf.Check("test", fSet, []*ast.File{file}, info)
	check.NoError(t).Require(err)

	target := findNode(file, name)
	check.NotNil(t).Name(`found name`).With(`name`, name).Assert(target)

	m := Analyze(locs.NewSet(fSet), info, metrics.New(), target)
	return &testTool{t: t, m: m}
}

func (tt *testTool) check(expLines ...string) {
	ctx := jsonify.NewContext()

	gotData, err := jsonify.Marshal(ctx, tt.m)
	check.NoError(tt.t).Require(err)

	exp := strings.Join(expLines, "\n")
	var expObj any
	err = yaml.Unmarshal([]byte(exp), &expObj)
	check.NoError(tt.t).Require(err)
	expData, err := jsonify.Marshal(ctx, expObj)
	check.NoError(tt.t).Require(err)

	if !slices.Equal(gotData, expData) {
		gotLines := strings.Split(string(gotData), "\n")
		expLines := strings.Split(string(expData), "\n")
		d := diff.Default().PlusMinus(gotLines, expLines)
		fmt.Println(strings.Join(d, "\n"))
		tt.t.Fail()
	}
}
