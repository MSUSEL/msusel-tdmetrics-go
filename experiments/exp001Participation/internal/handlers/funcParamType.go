package handlers

import (
	"fmt"
	"go/ast"

	"../parser"
)

type (
	// FuncData stores the collected information about a found function.
	FuncData struct {

		// FuncName is the name of the function.
		FuncName string

		// ParamTypes is the set of parameter types.
		ParamTypes []string
	}
)

// String will get a string for the function data.
func (d *FuncData) String() string {
	return fmt.Sprintf(`%s: %v`, d.FuncName, d.ParamTypes)
}

// JustFunctionParameterTypes returns a factory for a type resolver
// attached to a function parameter node handler.
func JustFunctionParameterTypes(filename string) []parser.NodeHandler {
	typeRev := NewTypeResolver(filename)
	return []parser.NodeHandler{
		typeRev.Handler(),
		FunctionParameterTypes(typeRev),
	}
}

// FunctionParameterTypes creates a new parser node handler which
// gets of parameter types of each function.
func FunctionParameterTypes(tres *TypeResolver) parser.NodeHandler {
	return func(n ast.Node) []interface{} {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if data := readFunc(x, tres); data != nil {
				return []interface{}{data}
			}
		}
		return nil
	}
}

// readFunc will read the information for a function from the AST node.
func readFunc(f *ast.FuncDecl, tres *TypeResolver) *FuncData {
	paramTypes := []string{}
	paramTypes = append(paramTypes, readFieldListTypes(f.Recv, tres)...)
	paramTypes = append(paramTypes, readFieldListTypes(f.Type.Params, tres)...)
	return &FuncData{
		FuncName:   f.Name.Name,
		ParamTypes: paramTypes,
	}
}

// readFieldListTypes will read the receiver or parameters for the function.
// This will only return one type for a multiple variable param.
func readFieldListTypes(fl *ast.FieldList, tres *TypeResolver) []string {
	result := []string{}
	if fl != nil && fl.List != nil {
		for _, f := range fl.List {
			if f != nil {
				result = append(result, readField(f, tres)...)
			}
		}
	}
	return result
}

// readField will read the given field and add the type the number of times
// it was used in the field based on the identifier (name) number.
func readField(f *ast.Field, tres *TypeResolver) []string {
	result := []string{}
	if ids := tres.ReadTypes(f.Type); len(ids) > 0 {
		for _, id := range ids {
			for i := 0; i < len(f.Names); i++ {
				result = append(result, id)
			}
		}
	}
	return result
}
