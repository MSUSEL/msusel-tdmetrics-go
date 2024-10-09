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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      reads: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      reads: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      reads: [`,
		`        argument1,`,
		`        argument2,`,
		`        argument3,`,
		`        basic1`,
		`      ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`    { name: $deref, signature: 1, exported: true },`,
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
		`      writes: [ argument3 ],`,
		`      reads:  [ basic1 ]`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  8,`,
		`      reads: [ basic1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      invokes: [`,
		`        selection1,`,
		`        selection2,`,
		`        tempDeclRef1`,
		`      ],`,
		`      reads:  [ tempReference1 ],`,
		`      writes: [ tempReference1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: close,   origin: tempReference1 },`,
		`    { name: doStuff, origin: tempReference1 }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: open, packagePath: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: Bar,  packagePath: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  7,`,
		`      complexity: 2,`,
		`      indents:    6,`,
		`      lineCount:  7,`,
		`      sideEffect: true,`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  10,`,
		`      complexity:  2,`,
		`      indents:    11,`,
		`      lineCount:  10,`,
		`      sideEffect: true,`,
		`      reads:  [ basic1 ],`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  9,`,
		`      complexity: 3,`,
		`      indents:    9,`,
		`      lineCount:  9,`,
		`      sideEffect: true,`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  11,`,
		`      complexity:  3,`,
		`      indents:    12,`,
		`      lineCount:  11,`,
		`      sideEffect: true,`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  12,`,
		`      complexity:  3,`,
		`      indents:    15,`,
		`      lineCount:  12,`,
		`      sideEffect: true,`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      sideEffect: true`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      sideEffect: true`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  11,`,
		`      complexity:  2,`,
		`      indents:    16,`,
		`      lineCount:  11,`,
		`      sideEffect: true,`,
		`      reads:  [ interfaceDesc1 ],`,
		`      writes: [ interfaceDesc1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      sideEffect: true`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`      sideEffect: true`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
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
		`  abstracts: [`,
		`    { name: $len,  signature: 1, exported: true },`,
		`    { name: $recv, signature: 2, exported: true },`,
		`    { name: $recv, signature: 3, exported: true },`,
		`    { name: $send, signature: 4, exported: true },`,
		`    { name: $send, signature: 5, exported: true }`,
		`  ],`,
		`  arguments: [`,
		`    {              type: basic2 },`,
		`    { name: okay,  type: basic1 },`,
		`    { name: value, type: basic1 },`,
		`    { name: value, type: typeParam1 }`,
		`  ],`,
		`  basics: [ bool, int ],`,
		`  interfaceDecls: [`,
		`    {`,
		`      name: Chan, package: 1, interface: 3, exported: true,`,
		`      instances: [ 1 ], typeParams: [ 1 ]`,
		`    },`,
		`    { name: any, package: 1, interface: 1, exported: true }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {},`,
		`    { abstracts: [ 1, 2, 4 ] },`,
		`    { abstracts: [ 1, 3, 5 ] }`,
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
		`      codeCount:  15,`,
		`      complexity:  5,`,
		`      indents:    17,`,
		`      lineCount:  15,`,
		`      sideEffect: true,`,
		`      reads:   [ basic1, interfaceInst1 ],`,
		`      writes:  [ basic1, interfaceInst1 ]`,
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
		`    { results: [ 1    ] },`,
		`    { results: [ 3, 2 ] },`,
		`    { results: [ 4, 2 ] },`,
		`    { params:  [ 3    ] },`,
		`    { params:  [ 4    ] }`,
		`  ],`,
		`  typeParams: [`,
		`    { name: T, type: interfaceDecl2 }`,
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
		`      reads: [`,
		`        selection1,`,
		`        tempReference1`,
		`      ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { origin: tempReference1, name: x }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: Bar }`,
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
		`      reads: [ tempDeclRef1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { packagePath: test, name: bar }`,
		`  ]`,
		`}`)
}

func Test_GetterWithParentheses(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo() int {`,
		`	return ((*(bar)))`,
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
		`      reads: [ tempDeclRef1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { packagePath: test, name: bar }`,
		`  ]`,
		`}`)
}

func Test_GetterWithNamedReturn(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo() (x int) {`,
		`	x = *bar`,
		`   return`,
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
		`      codeCount:  4,`,
		`      complexity: 1,`,
		`      indents:    4,`,
		`      lineCount:  4,`,
		`      #getter: false,`, // Not recognized as a getter
		`      reads:  [ tempDeclRef1 ],`,
		`      writes: [ argument1 ],`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { packagePath: test, name: bar }`,
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
		`      reads:  [ selection1, tempReference1 ],`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: x, origin: tempReference1 }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: Bar }`,
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
		`      reads:  [ argument1, tempReference1 ],`,
		`      writes: [ selection1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: x, origin: tempReference1 }`,
		`  ],`,
		`  tempReferences: [`,
		`    { packagePath: test, name: Bar }`,
		`  ]`,
		`}`)
}

func Test_SetterWithReference(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x int) {`,
		`	bar = &x`,
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
		`      setter:     true,`,
		`      sideEffect: true,`,
		`      reads:  [ argument1 ],`,
		`      writes: [ tempDeclRef1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: bar, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_SetterWithParentheses(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x int) {`,
		`	(bar) = ((&(x)))`,
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
		`      setter:     true,`,
		`      sideEffect: true,`,
		`      reads:  [ argument1 ],`,
		`      writes: [ tempDeclRef1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: bar, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_NotReverseSetter(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo(x *int) {`,
		`	*x = *bar`,
		`}`)
	tt.checkProj(
		`{`,
		`  abstracts: [`,
		`    { name: $deref, signature: 1, exported: true },`,
		`    { name: $deref, signature: 2, exported: true }`,
		`  ],`,
		`  arguments: [`,
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
		`      loc:        2,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    1,`,
		`      lineCount:  3,`,
		`      reads:  [ tempDeclRef1 ],`,
		`      writes: [ argument3 ]`,
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
		`  tempDeclRef: [`,
		`    { name: bar, packagePath: test }`,
		`  ],`,
		`  typeParams: [`,
		`    { name: T, type: interfaceDecl2 }`,
		`  ]`,
		`}`)
}

func Test_NamedResults(t *testing.T) {
	tt := parseDecl(t, `Foo`,
		`var bar *int`,
		`func Foo() (x, y int) {`,
		`	x = 10`,
		`	y = 24`,
		`	return`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 },`,
		`    { name: y, type: basic1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        2,`,
		`      codeCount:  5,`,
		`      complexity: 1,`,
		`      indents:    3,`,
		`      lineCount:  5,`,
		`      writes: [ argument1, argument2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_TwoFuncLitInit(t *testing.T) {
	tt := parseDecl(t, `val`,
		`var val = []int{`,
		`	func() int { return 12 }(),`,
		`	func() int { return 24 }(),`,
		` }`)
	tt.checkProj(
		`{`,
		`  abstracts: [`,
		`    { name: $get, signature: 2, exported: true },`,
		`    { name: $get, signature: 3, exported: true },`,
		`    { name: $len, signature: 1, exported: true },`,
		`    { name: $set, signature: 4, exported: true },`,
		`    { name: $set, signature: 5, exported: true }`,
		`  ],`,
		`  arguments: [`,
		`    {              type: basic1 },`,
		`    { name: index, type: basic1 },`,
		`    { name: value, type: basic1 },`,
		`    { name: value, type: typeParam1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  interfaceDecls: [`,
		`    {`,
		`      name: List, package: 1, interface: 3, exported: true,`,
		`      instances: [ 1 ], typeParams: [ 1 ]`,
		`    },`,
		`    { interface: 1, name: any, package: 1, exported: true }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {},`,
		`    { abstracts: [ 3, 1, 4 ] },`,
		`    { abstracts: [ 3, 2, 5 ] }`,
		`  ],`,
		`  interfaceInst: [`,
		`    {`,
		`      generic: 1, resolved: 2, `,
		`      instanceTypes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  4,`,
		`      complexity: 1,`,
		`      indents:    3,`,
		`      lineCount:  4,`,
		`      writes: [ interfaceInst1 ]`,
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
		`    { params: [ 2 ], results: [ 3 ] },`,
		`    { params: [ 2 ], results: [ 4 ] },`,
		`    { params: [ 2, 3 ] },`,
		`    { params: [ 2, 4 ] }`,
		`  ],`,
		`  typeParams: [`,
		`    { name: T, type: interfaceDecl2 }`,
		`  ]`,
		`}`)
}

func Test_BasicLitInit(t *testing.T) {
	tt := parseDecl(t, `val`,
		`var val = 10`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_ConstructInit(t *testing.T) {
	tt := parseDecl(t, `val`,
		`type Foo struct{ x, y int }`,
		`var val = Foo{`,
		`	x: 14,`,
		`	y: 42,`,
		`}`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        2,`,
		`      codeCount:  4,`,
		`      complexity: 1,`,
		`      indents:    6,`,
		`      lineCount:  4,`,
		`      reads:   [ basic1 ],`,
		`      writes:  [ tempReference1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: Foo, packagePath: test }`,
		`  ],`,
		`  tempReferences: [`,
		`    { name: Foo, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_SingletonMethodCallInit(t *testing.T) {
	tt := parseDecl(t, `val`,
		`type Foo struct{}`,
		`func (f Foo) f() float64 { return 3.14 }`,
		`var singleton = Foo{}`,
		`var val = singleton.f()`)
	tt.checkProj(
		`{`,
		`  basics: [ float64 ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        4,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      invokes: [ selection1 ],`,
		`      reads:   [ tempDeclRef1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: f, origin: tempDeclRef1 }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: singleton, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_SingletonFieldInit(t *testing.T) {
	tt := parseDecl(t, `val`,
		`type Foo struct{ x, y int}`,
		`var singleton = Foo{ x: 12, y: 24 }`,
		`var val = singleton.y`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        3,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      reads:   [ selection1, tempDeclRef1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: y, origin: tempDeclRef1 }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: singleton, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_MultipleInitMultipleValues(t *testing.T) {
	tt := parseDecl(t, `y`,
		`var x, y, z = "hello", 3.14, false`)
	tt.checkProj(
		`{`,
		`  basics: [ bool, float64, string ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      writes: [ basic1, basic2, basic3 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_MultipleInitSingleValue(t *testing.T) {
	tt := parseDecl(t, `y`,
		`var x, y = func()(int, int) {`,
		`  return 12, 24`,
		`}()`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  3,`,
		`      complexity: 1,`,
		`      indents:    6,`,
		`      lineCount:  3,`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_EnUnCapsulate(t *testing.T) {
	tt := parseDecl(t, `x`,
		`var x = struct{ y int }{ y: 24 }.y`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      reads:  [ basic1 ],`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_SelectOnUnnamedResultValue(t *testing.T) {
	tt := parseDecl(t, `val`,
		`func bar() struct{ y int }{`,
		`	return struct{ y int }{ y: 24 }`,
		`}`,
		`var val = bar().y`)
	tt.checkProj(
		`{`,
		`  basics: [ int ],`,
		`  fields: [`,
		`    { name: y, type: basic1 }`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        4,`,
		`      codeCount:  1,`,
		`      complexity: 1,`,
		`      lineCount:  1,`,
		`      invokes: [ tempDeclRef1 ],`,
		`      reads:   [ selection1, structDesc1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`  selections: [`,
		`    { name: y, origin: structDesc1 }`,
		`  ],`,
		`  structDescs: [`,
		`    { fields: [ 1 ] }`,
		`  ],`,
		`  tempDeclRef: [`,
		`    { name: bar, packagePath: test }`,
		`  ]`,
		`}`)
}

func Test_IncDec(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x int) int {`,
		`	x++`,
		`	y := x`,
		`	y--`,
		`	return y`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 },`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:        1,`,
		`      codeCount:  6,`,
		`      complexity: 1,`,
		`      indents:    4,`,
		`      lineCount:  6,`,
		`      reads:   [ argument1, basic1 ],`,
		`      writes:  [ argument1, basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_AssignInStatements(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x int) (y int, z int) {`,
		`	for y = range x {`,
		`		if z = y*x; z > 10 {`,
		`			break`,
		`		}`,
		`	}`,
		`	return`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 },`,
		`    { name: y, type: basic1 },`,
		`    { name: z, type: basic1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:         1,`,
		`      codeCount:   8,`,
		`      complexity:  3,`,
		`      indents:    10,`,
		`      lineCount:   8,`,
		`      reads:   [ argument1, argument2, argument3 ],`,
		`      writes:  [ argument2, argument3 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_DefineInForRange(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x int) (int, int) {`,
		`	for y := range x {`,
		`		if z := y*x; z > 10 {`,
		`			return y, z`,
		`		}`,
		`	}`,
		`	return -1, -1`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 },`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:         1,`,
		`      codeCount:   8,`,
		`      complexity:  3,`,
		`      indents:    10,`,
		`      lineCount:   8,`,
		`      reads:   [ argument1, basic1 ],`,
		`      writes:  [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_AssignInForLoop(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x int) (y int) {`,
		`	for y = 0; y < x; y++ {`,
		`		if y*x > 10 {`,
		`			return y`,
		`		}`,
		`	}`,
		`	return -1`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: basic1 },`,
		`    { name: y, type: basic1 }`,
		`  ],`,
		`  basics: [ int ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:         1,`,
		`      codeCount:   8,`,
		`      complexity:  3,`,
		`      indents:    10,`,
		`      lineCount:   8,`,
		`      reads:   [ argument1, argument2 ],`,
		`      writes:  [ argument2 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_InterFuncStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type u struct {`,
		`      name string`,
		`   }`,
		`	return x.(u).name`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: interfaceDesc1 }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      loc:         1,`,
		`      codeCount:   6,`,
		`      complexity:  1,`,
		`      indents:    13,`,
		`      lineCount:   6,`,
		`      reads: [ argument1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_ReferencingAnInterFuncStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type u struct {`,
		`      name string`,
		`   }`,
		`   type t struct {`,
		`      user u`,
		`   }`,
		`	return x.(t).user.name`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: interfaceDesc1 }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:   9,`,
		`      complexity:  1,`,
		`      indents:    25,`,
		`      lineCount:   9,`,
		`      loc:         1,`,
		`      reads: [ argument1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_SelfReferencingInterFuncStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type person struct {`,
		`      name string`,
		`      child *person`,
		`   }`,
		`	return x.(person).name`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: interfaceDesc1 }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:  7,`,
		`      complexity: 1,`,
		`      indents:   19,`,
		`      lineCount:  7,`,
		`      loc:        1,`,
		`      reads: [ argument1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_ReferencingAnInterFuncNestedStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type t struct {`,
		`      user struct {`,
		`         name string`,
		`      }`,
		`   }`,
		`	return x.(t).user.name`,
		`}`)
	tt.checkProj(
		`{`,
		`  arguments: [`,
		`    { name: x, type: interfaceDesc1 }`,
		`  ],`,
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:   8,`,
		`      complexity:  1,`,
		`      indents:    28,`,
		`      lineCount:   8,`,
		`      loc:         1,`,
		`      reads: [ argument1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ]`,
		`}`)
}

func Test_ReferencingAnInterFuncNestedComplexStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type t struct {`,
		`      user []*struct {`,
		`         name string`,
		`      }`,
		`   }`,
		`	return x.(t).user[0].name`,
		`}`)
	tt.checkProj(
		`{`,
		// TODO: Finish
		`}`)
}

func Test_ReferencingAnInterFuncUnnamedStruct(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`	return x.(struct {`,
		`      age  int`,
		`      name string`,
		`   }).name`,
		`}`)
	tt.checkProj(
		`{`,
		// TODO: Finish
		`}`)
}

func Test_InterFuncInterface(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   type stringer interface { String() string }`,
		`   if s, ok := x.(stringer); ok {`,
		`      return s.String()`,
		`   }`,
		`	return "unnamed"`,
		`}`)
	tt.checkProj(
		`{`,
		// TODO: Fixed
		`  arguments: [`,
		`    { type: basic2 },`,
		`    { name: x, type: interfaceDesc1 }`,
		`  ],`,
		`  basics: [ bool, string ],`,
		`  interfaceDescs: [`,
		`    {}`,
		`  ],`,
		`  language: go,`,
		`  metrics: [`,
		`    {`,
		`      codeCount:   7,`,
		`      complexity:  2,`,
		`      indents:    16,`,
		`      lineCount:   7,`,
		`      loc:         1,`,
		`      reads:  [ argument2, basic1 ],`,
		`      writes: [ basic1 ]`,
		`    }`,
		`  ],`,
		`  packages: [`,
		`    { name: test, path: test }`,
		`  ],`,
		`}`)
}

func Test_InterFuncUnnamedInterface(t *testing.T) {
	tt := parseDecl(t, `foo`,
		`func foo(x any) string {`,
		`   if s, ok := x.(interface { String() string }); ok {`,
		`      return s.String()`,
		`   }`,
		`	return "unnamed"`,
		`}`)
	tt.checkProj(
		`{`,
		// TODO: Finish
		`}`)
}

type testTool struct {
	t      *testing.T
	proj   constructs.Project
	curPkg constructs.Package
	baker  baker.Baker
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
	baker := baker.New(proj)
	conv := converter.New(baker, proj, curPkg, nil)
	return &testTool{
		t:      t,
		proj:   proj,
		curPkg: curPkg,
		baker:  baker,
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

	tt.m = Analyze(tt.info, tt.proj, tt.curPkg, tt.baker, tt.conv, expr)
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

	tt.m = Analyze(tt.info, tt.proj, tt.curPkg, tt.baker, tt.conv, target)
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
