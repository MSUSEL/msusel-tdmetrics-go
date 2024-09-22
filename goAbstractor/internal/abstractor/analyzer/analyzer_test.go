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
	"golang.org/x/tools/go/packages"
	"gopkg.in/yaml.v3"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/project"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

func Test_Empty(t *testing.T) {
	tt := parseExpr(t,
		`func() {}`)
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      #indents:   0,`,
		`      lineCount:  1`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`}`)
}

func Test_Simple(t *testing.T) {
	tt := parseExpr(t,
		`func() int {`,
		`	return max(1, 3, 5)`,
		`}`)
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      invokes: [ tempReference1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ]`,
		`}`)
}

func Test_SimpleWithExtraIndent(t *testing.T) {
	tt := parseExpr(t,
		`		func() int {`,
		`			return max(1, 3, 5)`,
		`		}`)
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      invokes: [ tempReference1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ]`,
		`}`)
}

func Test_SimpleParams(t *testing.T) {
	tt := parseExpr(t,
		`func(a int,`,
		`	  b int,`,
		`	  c int) int {`,
		`	return max(a, b, c)`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: a, type: basic1 },`,
		`    { name: b, type: basic1 },`,
		`    { name: c, type: basic1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  5,`,
		`      complexity: 1,`,
		`      indents:    7,`,
		`      lineCount:  5,`,
		`      invokes: [ tempReference1 ],`,
		`      reads: [`,
		`        argument1,`,
		`        argument2,`,
		`        argument3`,
		`      ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ]`,
		`}`)
}

func Test_SimpleWithReturn(t *testing.T) {
	tt := parseExpr(t,
		`func() int {`,
		`	x := max(1, 3, 5)`,
		`	return x`,
		`}`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  4,`,
		`      complexity: 1,`,
		`      indents:    2,`,
		`      lineCount:  4,`,
		`      invokes: [ tempReference1 ],`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ]`,
		`}`)
}

func Test_SimpleWithOutParam(t *testing.T) {
	tt := parseExpr(t,
		`func(x *int) {`,
		`	*x = max(1, 3, 5)`,
		`}`)
	tt.checkProj(
		`{`,
		`  abstracts: [`,
		`    { name: $deref, signature: 1 },`, // TODO: Shouldn't this be exported?
		`    { name: $deref, signature: 2, exported: true }`,
		`  ],`,
		`  "arguments": [`,
		`    {          type: basic1 },`,
		`    {          type: typeParam1 },`,
		`    { name: x, type: interfaceInst1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  interfaceDecls: [`,
		`    {`,
		`      name: Pointer, package: 1, interface: 3, exported: true,`,
		`      instances: [ 1 ], typeParams: [ 1 ]`,
		`    },`,
		`    { name: any, package: 1, interface: 1, exported: true }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {},`,
		`    { abstracts: [ 1 ] },`,
		`    { abstracts: [ 2 ] }`,
		`  ],`,
		`  interfaceInst: [`,
		`    {`,
		`      generic: 1, resolved: 2,`,
		`      instanceTypes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      invokes: [ tempReference1 ],`,
		`      writes:  [ argument3 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    {`,
		`      name: $builtin, path: $builtin,`,
		`      interfaces: [ 1, 2 ]`,
		`    },`,
		`    { name: test, path: test }`,
		`  ],`,
		`  signatures: [`,
		`    { results: [ 1 ] },`,
		`    { results: [ 2 ] }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ],`,
		`  typeParams: [`,
		`    { name: T, type: interfaceDecl2 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  8,`,
		`      invokes: [ tempReference1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: max }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        5,`,
		`      codeCount:  5,`,
		`      complexity: 1,`,
		`      indents:    3,`,
		`      lineCount:  5,`,
		`      invokes:  [ 1, 2, 3 ],`,
		`      writes:   [ 4 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: close   },`,
		`    { packagePath: test, name: doStuff },`,
		`    { packagePath: test, name: open    },`,
		`    { packagePath: test, name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1, origin: 4 },`,
		`    { target: tempReference2, origin: 4 },`,
		`    { target: tempReference3 },`,
		`    { target: tempReference4 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  7,`,
		`      complexity: 2,`,
		`      indents:    6,`,
		`      lineCount:  7,`,
		`      invokes:  [ 1 ],`,
		`      reads:    [ 2 ],`,
		`      writes:   [ 2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: println },`,
		`    { name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  10,`,
		`      complexity:  2,`,
		`      indents:    11,`,
		`      lineCount:  10,`,
		`      invokes:  [ 1, 2 ],`,
		`      reads:    [ 3 ],`,
		`      writes:   [ 3 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print   },`,
		`    { name: println },`,
		`    { name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 },`,
		`    { target: tempReference3 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  9,`,
		`      complexity: 3,`,
		`      indents:    9,`,
		`      lineCount:  9,`,
		`      invokes:  [ 1 ],`,
		`      reads:    [ 2 ],`,
		`      writes:   [ 2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: println },`,
		`    { name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  11,`,
		`      complexity:  3,`,
		`      indents:    12,`,
		`      lineCount:  11,`,
		`      invokes:  [ 1 ],`,
		`      reads:    [ 2 ],`,
		`      writes:   [ 2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: println },`,
		`    { name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  12,`,
		`      complexity:  3,`,
		`      indents:    15,`,
		`      lineCount:  12,`,
		`      invokes:  [ 1 ],`,
		`      reads:    [ 2 ],`,
		`      writes:   [ 2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: println },`,
		`    { name: x       }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  14,`,
		`      complexity:  1,`,
		`      indents:    19,`,
		`      lineCount:  14,`,
		`      invokes:  [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  14,`,
		`      complexity:  1,`,
		`      indents:    19,`,
		`      lineCount:  14,`,
		`      invokes:  [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  11,`,
		`      complexity:  2,`,
		`      indents:    16,`,
		`      lineCount:  11,`,
		`      invokes:  [ 1, 3 ],`,
		`      reads:    [ 2 ],`,
		`      writes:   [ 2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print   },`,
		`    { name: r       },`,
		`    { name: recover }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 },`,
		`    { target: tempReference3 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  11,`,
		`      complexity:  2,`,
		`      indents:    15,`,
		`      lineCount:  11,`,
		`      invokes:  [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  6,`,
		`      complexity: 2,`,
		`      indents:    5,`,
		`      lineCount:  6,`,
		`      invokes:  [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: print }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 }`,
		`  ]`,
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
	tt.checkProj(
		`{`,
		`  basics: [ bool ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  15,`,
		`      complexity:  5,`,
		`      indents:    17,`,
		`      lineCount:  15,`,
		`      invokes:  [ 2, 3, 5 ],`,
		`      reads:    [ 1, 2, 3, 4 ],`,
		`      writes:   [ 2, 3, 4 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: A },`,
		`    { name: B },`,
		`    { name: b },`,
		`    { name: print }`,
		`  ],`,
		`  usages: [`,
		`    { target: basic1 },`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2 },`,
		`    { target: tempReference3 },`,
		`    { target: tempReference4 }`,
		`  ]`,
		`}`)
}

func Test_GetterWithSelect(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x int }`,
		`func (b Bar) Foo() int {`,
		`	return b.x`,
		`}`)
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`	   loc:        2,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      getter:     true,`,
		`      reads:    [ 1, 2 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: b },`,
		`    { packagePath: test, name: x }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2, origin: 1 }`,
		`  ]`,
		`}`)
}

func Test_GetterWithDereference(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo() int {`,
		`	return *bar`,
		`}`)
	tt.checkProj(
		`{`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        2,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      getter:  true,`,
		`      reads:    [ 1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: bar }`,
		`  ],`,
		`  usages: [`,
		`    { target: tempReference1 }`,
		`  ]`,
		`}`)
}

func Test_GetterWithConvert(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x float64 }`,
		`func (b Bar) Foo() int {`,
		`	return int(b.x)`,
		`}`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        2,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      getter:  true,`,
		`      reads:    [ 2, 3 ],`,
		`      writes:   [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: b },`,
		`    { packagePath: test, name: x }`,
		`  ],`,
		`  usages: [`,
		`    { target: basic1 },`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2, origin: 2 }`,
		`  ]`,
		`}`)
}

func Test_SetterWithSelect(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`type Bar struct { x int }`,
		`func (b Bar) Foo(x int) {`,
		`	b.x = x`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        2,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      setter:  true,`,
		`      reads:    [ 1 ],`,
		`      writes:   [ 1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: Bar },`,
		`    { packagePath: test, name: b },`,
		`    { packagePath: test, name: "x WHAT? No, this is b.x, the ref won't resolve" }`,
		`  ],`,
		`  usages: [`,
		`    { target: basic1 },`,
		`    { target: tempReference1 },`,
		`    { target: tempReference2, origin: 2 }`,
		`  ]`,
		`}`)
}

func Test_SetterWithReference(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x int) {`,
		`	bar = &x`,
		`}`)
	//tt.checkMetrics(
	//	`{`,
	//	`	loc:        2,`,
	//	`	codeCount:  3,`,
	//	`	complexity: 1,`,
	//	`	indents:    1,`,
	//	`	lineCount:  3,`,
	//	`   setter:  true,`,
	//	`}`)
	tt.checkProj()
}

func Test_NotReverseSetter(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x *int) {`,
		`	*x = *bar`,
		`}`)
	//tt.checkMetrics(
	//	`{`,
	//	`	loc:        2,`,
	//	`	codeCount:  3,`,
	//	`	complexity: 1,`,
	//	`	indents:    1,`,
	//	`	lineCount:  3,`,
	//	`}`)
	tt.checkProj()
}

// TODO: Test a named return
// func() (x, y int) { x = 10; y = 24; return }

// TODO: Test joining metrics:
// var val = []int{
//   func() int { ⋯ }(),
//   func() int { ⋯ }(),
// }

// TODO: Test nothing variable:
// var val = 10

// TODO: Test reading metrics with only read reference:
// var val = singleton.f()

// TODO: Test reading metrics with only read reference but without a function:
// var val = singleton.value

// TODO: Test reading metrics with read reference as parameter:
// var val = func(f Foo) int { ⋯ }(singleton.f)

// TODO: Test reading metrics with read reference in typed call:
// var val = Foo[int](singleton)

// TODO: Test parentheses in getters and setters.

// TODO: Test assignment of returned pointer:
// func() { *(getIntPointer()) = 10 }

// TODO: Test assigning named result:
// func() (x int) { x = 10; return }

// TODO: Test multiple assignments:
// x, y := 1, 2  and  x, y := func()(int, int) { ⋯ }

// TODO: Test local encapsulation of type:
// x := struct{y externalType}{y: ext}.y

// TODO: Test selection from return value.
// x := foo().y

// TODO: Test inc and dec also work as assignment.

// TODO: Test literal cast and call
// type foo int; func(f foo) bar { ⋯ }; foo(6).bar()

// TODO: Test the assignment in a for-loop or if-statement
// are picked up as writes, `for i := 0; ⋯`

type testTool struct {
	t      *testing.T
	proj   constructs.Project
	curPkg constructs.Package
	conv   converter.Converter
	fSet   *token.FileSet
	info   *types.Info
	m      constructs.Metrics
}

func newTestTool(t *testing.T) *testTool {
	pkgPath := `test`
	pkgName := `test`
	info := &types.Info{
		Defs:       make(map[*ast.Ident]types.Object),
		Instances:  make(map[*ast.Ident]types.Instance),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Uses:       make(map[*ast.Ident]types.Object),
	}
	fSet := token.NewFileSet()
	proj := project.New(locs.NewSet(fSet))
	curPkg := proj.NewPackage(constructs.PackageArgs{
		RealPkg: &packages.Package{
			PkgPath:   pkgPath,
			Name:      pkgName,
			Types:     types.NewPackage(pkgPath, pkgName),
			TypesInfo: info,
		},
		Path: pkgPath,
		Name: pkgName,
	})
	conv := converter.New(baker.New(proj), proj, curPkg, nil)
	return &testTool{
		t:      t,
		proj:   proj,
		curPkg: curPkg,
		conv:   conv,
		fSet:   fSet,
		info:   info,
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
	tt := newTestTool(t)

	code := strings.Join(lines, "\n")
	expr, err := parser.ParseExprFrom(tt.fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)

	err = types.CheckExpr(tt.fSet, nil, token.NoPos, expr, tt.info)
	check.NoError(t).Require(err)

	tt.m = Analyze(tt.info, tt.proj, tt.curPkg, tt.conv, expr)
	tt.proj.UpdateIndices()
	return tt
}

func parseDecl(t *testing.T, name string, lines ...string) *testTool {
	tt := newTestTool(t)

	code := "package test\n" + strings.Join(lines, "\n")
	file, err := parser.ParseFile(tt.fSet, ``, []byte(code), parser.ParseComments)
	check.NoError(t).Require(err)

	var conf types.Config
	_, err = conf.Check("test", tt.fSet, []*ast.File{file}, tt.info)
	check.NoError(t).Require(err)

	target := findNode(file, name)
	check.NotNil(t).Name(`found name`).With(`name`, name).Assert(target)

	tt.m = Analyze(tt.info, tt.proj, tt.curPkg, tt.conv, target)
	tt.proj.UpdateIndices()
	return tt
}

func (tt *testTool) checkProj(expLines ...string) {
	tt.check(tt.proj, expLines...)
}

func (tt *testTool) check(data any, expLines ...string) {
	ctx := jsonify.NewContext().Full()
	gotData, err := jsonify.Marshal(ctx, data)
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
		d := diff.Default().PlusMinus(expLines, gotLines)
		fmt.Println(strings.Join(d, "\n"))
		tt.t.Fail()
	}
}
