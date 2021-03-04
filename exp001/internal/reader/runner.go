package reader

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
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
	fileSet := token.NewFileSet()
	files := []*ast.File{}
	for _, filename := range filenames {
		source := sources[filename]
		f, err := parser.ParseFile(fileSet, filename, source, parser.ParseComments)
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
	imp := importer.ForCompiler(fileSet, "source", nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check(basePath, fileSet, files, info)
	if err != nil {
		log.Fatal("Type Check Failed: ", err)
	}

	// Gather up read results to be returned.
	return &Project{
		BasePath:   basePath,
		FileSet:    fileSet,
		Package:    pkg,
		Types:      info.Types,
		Defs:       info.Defs,
		Uses:       info.Uses,
		Implicits:  info.Implicits,
		Selections: info.Selections,
	}
}
