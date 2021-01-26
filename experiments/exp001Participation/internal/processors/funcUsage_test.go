package processors

import (
	"fmt"
	"strings"
	"testing"

	"../handlers"
	"../parser"
)

func lns(lines ...string) string {
	return strings.Join(lines, "\n")
}

func testFuncs(t *testing.T, filename, source, exp string) {
	proc := NewFuncUsage()
	parser.NewParser().
		UpdateProgress(parser.MuteProgress).
		AddProcessor(proc.ProcessFunction()).
		AddHandler(handlers.JustFunctionParameterTypes).
		AddSource(filename, source).
		Start().
		Await()

	result := proc.String()
	if result != exp {
		t.Error(fmt.Sprint("Error in function parameter handler:\n",
			"Expected: ", strings.ReplaceAll(exp, "\n", "\n          "), "\n",
			"Result:   ", strings.ReplaceAll(result, "\n", "\n          ")))
	}
}

func Test_FuncUsage_Basics(t *testing.T) {
	testFuncs(t, `main.go`,
		lns(`package main`,
			`import ani "object/animal"`,
			`func meow(a *ani.Cat) {}`,
			`func bark(a *ani.Dog) {}`,
			`func chase(a *ani.Cat, b *ani.Dog, direction bool) {}`),
		lns(`1 bool`,
			`2 object/animal.Cat`,
			`2 object/animal.Dog`))

}
