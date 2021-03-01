package reader

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"runtime"
	"sort"
)

type runner struct {
}

func newRunner() *runner {
	return &runner{}
}

func (r *runner) run(basePath string, sources map[string]interface{}) *Project {
	// Sort the file names
	filenames := []string{}
	for filename := range sources {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)

	// Read and parse all the sources
	fileset := token.NewFileSet()
	files := []*ast.File{}
	for _, filename := range filenames {
		source := sources[filename]
		f, err := parser.ParseFile(fileset, filename, source, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}

	// Prepare the info for collecting data.
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	// Resolve types in the packages.
	imp := importer.ForCompiler(fileset, runtime.Compiler, nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check(basePath, fileset, files, info)
	if err != nil {
		log.Fatal(err)
	}

	// Gather up read results to be returned.
	return &Project{
		Fileset:    fileset,
		Package:    pkg,
		Types:      info.Types,
		Defs:       info.Defs,
		Uses:       info.Uses,
		Implicits:  info.Implicits,
		Selections: info.Selections,
	}
}
