package metrics

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strings"
	"testing"

	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Simple(t *testing.T) {
	parseFunc(t, -1,
		`x := max(1, 3, 5)`)
	// ─┬─[enter_8]
	//  └───[exit_29]
}

func Test_SimpleWithReturn(t *testing.T) {
	parseFunc(t, -1,
		`x := max(1, 3, 5)`,
		`return x`)
	// ─┬─[enter_8]
	//  └───[exit_38]
}

func Test_SimpleWithDefer(t *testing.T) {
	parseFunc(t, -1,
		`x := open()`,
		`defer x.close()`,
		`x.doStuff()`)
	// ─┬─[enter_8]
	//  └─┬─[defer_22]
	//    └───[exit_51]
}

func Test_SimpleIf(t *testing.T) {
	parseFunc(t, -1,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`}`,
		`println(x)`)
	// ─┬─[enter_8]
	//  └─┬─[if_17]
	//    ├─┬─[ifBody_26]
	//    │ └─┬─[endIf_37]
	//    │   └───[exit_50]
	//    └───<endIf_37>
}

func Test_SimpleIfElse(t *testing.T) {
	parseFunc(t, -1,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else {`,
		`  x = 2`,
		`  print("cat")`,
		`}`,
		`println(x)`)
	// ─┬─[enter_8]
	//  └─┬─[if_17]
	//    ├─┬─[ifBody_26]
	//    │ └─┬─[endIf_69]
	//    │   └───[exit_82]
	//    └─┬─[elseBody_43]
	//      └───<endIf_69>
}

func Test_SimpleIfElseIf(t *testing.T) {
	parseFunc(t, -1,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else if x > 4 {`,
		`  x = 3`,
		`}`,
		`println(x)`)
	// ─┬─[enter_8]
	//  └─┬─[if_17]
	//    ├─┬─[ifBody_26]
	//    │ └─┬─[endIf_63]
	//    │   └───[exit_76]
	//    └─┬─[elseBody_43]
	//      └─┬─[if_43]
	//        ├─┬─[ifBody_52]
	//        │ └───<endIf_63>
	//        └───<endIf_63>
}

func Test_SimpleIfElseIfElse(t *testing.T) {
	parseFunc(t, -1,
		`x := 9`,
		`if x > 7 {`,
		`  x = 4`,
		`} else if x > 4 {`,
		`  x = 3`,
		`} else {`,
		`  x = 2`,
		`}`,
		`println(x)`)
	// ─┬─[enter_8]
	//  └─┬─[if_17]
	//    ├─┬─[ifBody_26]
	//    │ └─┬─[endIf_80]
	//    │   └───[exit_93]
	//    └─┬─[elseBody_43]
	//      └─┬─[if_43]
	//        ├─┬─[ifBody_52]
	//        │ └───<endIf_80>
	//        └─┬─[elseBody_69]
	//          └───<endIf_80>
}

func Test_DeferInBlock(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`defer func() {`,
		`	print("B ")`,
		`}()`,
		`{`,
		`	print("C ")`,
		`	defer func() {`,
		`		print("D ")`,
		`	}()`,
		`	print("E ")`,
		`}`,
		`print("F ")`)
	// Output: A C E F D B
	// ─┬─[enter_8]
	//  └─┬─[defer_70]
	//    └─┬─[defer_22]
	//      └───[exit_132]
}

func Test_DeferInFuncLiteral(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`defer func() {`,
		`	print("B ")`,
		`}()`,
		`func() {`,
		`	print("C ")`,
		`	defer func() {`,
		`		print("D ")`,
		`	}()`,
		`	print("E ")`,
		`}()`,
		`print("F ")`)
	// Output: A C E D F B
}

func Test_DeferWithComplexity(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`defer func() {`,
		`   if r := recover(); r != nil {`,
		`		print("B ")`,
		`       return`,
		`   }`,
		`	print("C ")`,
		`}()`,
		`print("D ")`)
}

func Test_ForRangeWithDefer(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`for _ = range 4 {`,
		`	print("B ")`,
		`	defer func() {`,
		`		print("C ")`,
		`	}()`,
		`	print("D ")`,
		`}`,
		`print("E ")`)
	// Output: A B D B D B D B D E C C C C
}

func Test_InfiniteGotoLoop(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`BEANS:`,
		`print("B ")`,
		`goto BEANS`,
		`print("C ")`) // Unreachable
	// Output: A B B B...
	// ─┬─[enter_8]
	//  └─┬─[label_22]
	//    └─┬─[goto_41]
	//      └───<label_22>
}

func Test_SkippingGoto(t *testing.T) {
	parseFunc(t, -1,
		`print("A ")`,
		`goto BEANS`,
		`print("B ")`, // Unreachable
		`BEANS:`,
		`print("C ")`)
	// Output: A C
	// ─┬─[enter_8]
	//  └─┬─[goto_22]
	//    └─┬─[label_45]
	//      └───[exit_65]
}

func Test_GotoWithIf(t *testing.T) {
	parseFunc(t, -1,
		`x := 10`,
		`TOP:`,
		`if x <= 0 {`,
		`  goto BOTTOM`,
		`}`,
		`print(x)`,
		`x--`,
		`goto TOP`,
		`BOTTOM:`,
		`print("Done")`)
}

func parseFunc(t *testing.T, exp int, lines ...string) {
	code := fmt.Sprintf("func() {\n%s\n}\n", strings.Join(lines, "\n"))
	expr, err := parser.ParseExpr(code)
	check.NoError(t).Require(err)
	block := expr.(*ast.FuncLit).Body
	m := New(block)
	check.Equal(t, exp).Assert(m.Complexity)
}
