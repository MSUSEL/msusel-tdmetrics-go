package handlers

import (
	"fmt"
	"strings"
	"testing"

	"../parser"
)

func lns(lines ...string) string {
	return strings.Join(lines, "\n")
}

func testFuncs(t *testing.T, filename, source, exp string) {
	buf := &strings.Builder{}
	parser.NewParser().
		UpdateProgress(parser.MuteProgress).
		OnError(parser.RecordError(buf)).
		AddProcessor(parser.RecordProcess(buf)).
		AddHandler(JustFunctionParameterTypes).
		AddSource(filename, source).
		Start().
		Await()

	result := buf.String()
	if result != exp {
		t.Error(fmt.Sprint("Error in function parameter handler:\n",
			"Expected: ", strings.ReplaceAll(exp, "\n", "\n          "), "\n",
			"Result:   ", strings.ReplaceAll(result, "\n", "\n          ")))
	}
}

func Test_Parser_ParamBasics(t *testing.T) {
	testFuncs(t, `main.go`,
		lns(`package main`,
			`func noParamsFunc() {}`),
		lns(`noParamsFunc: []`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`func twoIntParams(a, b int) {}`),
		lns(`twoIntParams: [int int]`))

	testFuncs(t, `animal.go`,
		lns(`package animal`,
			`func localParam(a Cat) {}`),
		lns(`localParam: [Cat]`))

	testFuncs(t, `animal.go`,
		lns(`package animal`,
			`func localPointerParam(a *Cat) {}`),
		lns(`localPointerParam: [Cat]`))

	testFuncs(t, `object/animal/cat.go`,
		lns(`package animal`,
			`func localPointerParam(a *Cat) {}`),
		lns(`localPointerParam: [object/animal.Cat]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import "animal"`,
			`func importedPointerParam(a *animal.Cat) {}`),
		lns(`importedPointerParam: [animal.Cat]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import "object/animal"`,
			`func importedPointerParam(a *animal.Cat) {}`),
		lns(`importedPointerParam: [object/animal.Cat]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import ani "object/animal"`,
			`func importedPointerParam(a *ani.Cat) {}`),
		lns(`importedPointerParam: [object/animal.Cat]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import ani "object/animal"`,
			`func importedPointerParam(a ***ani.Cat) {}`),
		lns(`importedPointerParam: [object/animal.Cat]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import ani "object/animal"`,
			`func importedPointerParam(a ani.Cat) {}`),
		lns(`importedPointerParam: [object/animal.Cat]`))
}

func Test_Parser_ParamUnhandled(t *testing.T) {
	testFuncs(t, `main.go`,
		lns(`package main`,
			`import . "animal"`,
			`func importedPointerParam(a *animal.Cat) {}`),
		lns(`Error: Currently can not handle dot imports`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`func paramTypeDef(a struct{`,
			`	first string`,
			`   last string`,
			`}) {}`),
		lns(`Error: Unexpected expression type: *ast.StructType`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`func paramTypeDef(a interface{`,
			`	Name() string`,
			`}) {}`),
		lns(`Error: Unexpected expression type: *ast.InterfaceType`))
}

func Test_Parser_ParamMultiples(t *testing.T) {
	testFuncs(t, `main.go`,
		lns(`package main`,
			`import ani "object/animal"`,
			`func meow(a *ani.Cat) {}`,
			`func bark(a *ani.Dog) {}`,
			`func chase(a *ani.Cat, b *ani.Dog, direction bool) {}`),
		lns(`meow: [object/animal.Cat]`,
			`bark: [object/animal.Dog]`,
			`chase: [object/animal.Cat object/animal.Dog bool]`))

	testFuncs(t, `main.go`,
		lns(`package main`,
			`import cats "object/animal/cats"`,
			`import dogs "object/animal/dogs"`,
			`func meow(a *cats.Cat) {}`,
			`func bark(a *dogs.Dog) {}`,
			`func chase(a *cats.Cat, b *dogs.Dog, direction bool) {}`),
		lns(`meow: [object/animal/cats.Cat]`,
			`bark: [object/animal/dogs.Dog]`,
			`chase: [object/animal/cats.Cat object/animal/dogs.Dog bool]`))
}
