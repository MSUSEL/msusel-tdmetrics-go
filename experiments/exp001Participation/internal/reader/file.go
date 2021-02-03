package reader

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"strconv"
)

type fileDef struct {
	pkg     *readPackageDef
	file    *ast.File
	imports map[string]string
}

func newFile(pkg *readPackageDef, f *ast.File) *fileDef {
	return &fileDef{
		pkg:     pkg,
		file:    f,
		imports: map[string]string{},
	}
}

func (r *fileDef) read() {
	for _, decl := range r.file.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				r.readImports(d)
			case token.VAR:
				r.readValues(d, false)
			case token.CONST:
				r.readValues(d, true)
			case token.TYPE:
				r.readTypeDef(d)
			}
		}
	}
}

// readImports reads all the imports.
func (r *fileDef) readImports(d *ast.GenDecl) {
	for _, spec := range d.Specs {
		if s, ok := spec.(*ast.ImportSpec); ok {
			r.readImport(s)
		}
	}
}

// readImport reads the given import and stores it by the import key.
func (r *fileDef) readImport(s *ast.ImportSpec) {
	value, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		panic(err)
	}

	var key string
	if s.Name != nil {
		key = s.Name.Name
	}
	if len(key) <= 0 {
		_, key = path.Split(value)
	}

	r.imports[key] = value
}

// readValue read package scope constants and values.
func (r *fileDef) readValues(d *ast.GenDecl, constant bool) {
	for _, spec := range d.Specs {
		if s, ok := spec.(*ast.ValueSpec); ok {
			r.readValue(s, constant)
		}
	}
}

func (r *fileDef) readValue(s *ast.ValueSpec, constant bool) {
	if s.Type != nil {
		types.Eval(r.pkg.fileset, r.pkg.pkg, s.Type.Pos)

	} else if constant && len(s.Values) <= 0 {
		// A constant may have no type or value is defined.
		// In this case it will inherit the previous type (see Go's iota).

		// TODO: Implement
	}
}

func (r *fileDef) readTypeDef(d *ast.GenDecl) {
	// TODO: Implement
}
