package handlers

import (
	"fmt"
	"go/ast"
	"path"
	"strconv"
	"strings"

	"../parser"
)

// TypeResolver is a tool for resolving types.
type TypeResolver struct {
	filename string
	dirname  string
	imports  map[string]string
}

// NewTypeResolver creates a new type resolver for the given filename.
func NewTypeResolver(filename string) *TypeResolver {
	filename = path.Clean(filename)
	dir := path.Dir(filename)
	if dir == `.` {
		dir = ``
	}
	return &TypeResolver{
		filename: filename,
		dirname:  dir,
		imports:  map[string]string{},
	}
}

// Handler gets the node handler to collect type information.
func (t *TypeResolver) Handler() parser.NodeHandler {
	return func(n ast.Node) []interface{} {
		switch x := n.(type) {
		case *ast.ImportSpec:
			t.addImport(x)
		}
		return nil
	}
}

// addImport will add an import definition to the type resolver.
func (t *TypeResolver) addImport(n *ast.ImportSpec) {
	if n != nil {
		var key string
		if n.Name != nil {
			key = n.Name.Name
		}

		value, err := strconv.Unquote(n.Path.Value)
		if err != nil {
			panic(err)
		}

		value = path.Clean(path.Join(t.dirname, value))
		if len(key) <= 0 {
			_, key = path.Split(value)
		}

		if key == `.` {
			panic(`Currently can not handle dot imports`)
		}

		t.imports[key] = value
	}
}

// ReadTypes gets the full type from the given expression.
func (t *TypeResolver) ReadTypes(n ast.Expr) []string {
	parts := t.readTypeParts(n)
	extended := false
	if len(parts) > 0 {
		if len(parts) == 1 {
			switch parts[0] {
			case `int`, `int8`, `int16`, `int32`, `int64`,
				`uint`, `uint8`, `uint16`, `uint32`, `uint64`,
				`float32`, `float64`, `bool`, `byte`:
				extended = true
			}
		} else {
			if value, ok := t.imports[parts[0]]; ok {
				parts[0] = value
				extended = true
			}
		}
		if (!extended) && (len(t.dirname) > 0) {
			parts = append([]string{t.dirname}, parts...)
		}
	}
	return []string{strings.Join(parts, `.`)}
}

// readTypeParts will read the parts of an identifier.
func (t *TypeResolver) readTypeParts(n ast.Expr) []string {
	switch x := n.(type) {
	case *ast.StarExpr:
		return t.readTypeParts(x.X)
	case *ast.ArrayType:
		return t.readTypeParts(x.Elt)
	case *ast.Ellipsis:
		return t.readTypeParts(x.Elt)
	case *ast.Ident:
		return []string{x.Name}
	case *ast.SelectorExpr:
		return append(t.readTypeParts(x.X), x.Sel.Name)
	default:
		panic(fmt.Sprintf(`Unexpected expression type: %T`, x))
	}
}
