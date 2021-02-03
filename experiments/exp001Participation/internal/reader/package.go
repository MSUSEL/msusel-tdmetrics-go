package reader

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

// packageDef will read all the files in a project and
// parse them into a structure for measuring the code.
type packageDef struct {
	fileset *token.FileSet
	pkg     *ast.Package

	files map[string]*fileDef
}

func newPackage(sources map[string]interface{}) *packageDef {
	fileset := token.NewFileSet()
	files := make(map[string]*ast.File)
	filenames := []string{}

	for filename, source := range sources {
		f, err := parser.ParseFile(fileset, filename, source, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files[filename] = f
		filenames = append(filenames, filename)
	}

	pkg, err := ast.NewPackage(fileset, files, nil, nil)
	if err != nil {
		panic(err)
	}

	sort.Strings(filenames)
	return &packageDef{
		fileset: fileset,
		pkg:     pkg,
	}
}

func (p *packageDef) read() {
	filenames := make([]string, len(p.pkg.Files))
	i := 0
	for filename := range p.pkg.Files {
		filenames[i] = filename
		i++
	}
	sort.Strings(filenames)
	for _, filename := range filenames {
		f := newFile(p, p.pkg.Files[filename])
		p.files[filename] = f
		f.read()
	}
}
