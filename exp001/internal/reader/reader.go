package reader

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/data"
	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/filter"
	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/utils"
)

// Reader is the information used to setup how the reader
// will run and the data to run inside the reader.
type Reader struct {

	// basePath is the base path to use while reading the data.
	basePath string

	// sources is the filename and associated code source.
	// The code source can be nil to read the file, otherwise the source is used.
	sources map[string]interface{}
}

// New constructs a new reader.
func New() *Reader {
	return &Reader{
		basePath: ".",
		sources:  make(map[string]interface{}, 0),
	}
}

// SetBasePath is the base path to use while reading the data.
func (r *Reader) SetBasePath(basePath string) *Reader {
	r.basePath = basePath
	return r
}

// BasePath gets the base path to use while reading the data.
func (r *Reader) BasePath() string {
	return r.basePath
}

// AddFiles adds one or more new files to the reader to this configuration.
func (r *Reader) AddFiles(filenames ...string) *Reader {
	for _, filename := range filenames {
		r.sources[filename] = nil
	}
	return r
}

// AddSource adds a string of source code to read directly.
// The filename is used to identify this source code string.
func (r *Reader) AddSource(filename string, source string) *Reader {
	r.sources[filename] = source
	return r
}

// AddDir adds all the files in the given directory.
func (r *Reader) AddDir(foldername string) *Reader {
	fileInfo, err := ioutil.ReadDir(foldername)
	if err != nil {
		panic(err)
	}
	files := make([]string, len(fileInfo))
	for i, info := range fileInfo {
		files[i] = path.Clean(path.Join(foldername, info.Name()))
	}
	return r.AddFiles(files...)
}

// AddDirRecursively adds all the files in the given directory and
// all the files in all children folders.
func (r *Reader) AddDirRecursively(dirname string) *Reader {
	var files []string
	err := filepath.Walk(dirname,
		func(filepath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, filepath)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	return r.AddFiles(files...)
}

// FilterFilenames will remove all the filenames and paired sources which
// match the given handler. If the matcher returns true the file is removed.
func (r *Reader) FilterFilenames(matcher filter.Matcher) *Reader {
	for filename := range r.sources {
		if matcher(filename) {
			delete(r.sources, filename)
		}
	}
	return r
}

// Filenames get the list of filenames that have been set to this reader
// and should be processed when started.
func (r *Reader) Filenames() []string {
	filenames := make([]string, 0, len(r.sources))
	for filename := range r.sources {
		filenames = append(filenames, filename)
	}
	sort.Strings(filenames)
	return filenames
}

// sortedFilenames gets all the filenames sorted.
func (r *Reader) sortedFilenames() []string {
	names := utils.NewStringSet()
	for filename := range r.sources {
		names.Add(filename)
	}
	return names.Values()
}

// Read and parse all the sources and groups by package.
func (r *Reader) parseSources() (*token.FileSet, map[string][]*ast.File) {
	fileSet := token.NewFileSet()
	filesByPackage := map[string][]*ast.File{}
	filenames := r.sortedFilenames()
	for _, filename := range filenames {
		source := r.sources[filename]
		f, err := parser.ParseFile(fileSet, filename, source, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		// Using the package name to differentiate packages may not
		// work if the same package name is used in two placed in a project.
		pkg := f.Name.Name
		list := filesByPackage[pkg]
		if len(list) <= 0 {
			list = []*ast.File{}
		}
		list = append(list, f)
		filesByPackage[pkg] = list
	}
	return fileSet, filesByPackage
}

// readPackage constructs a new package with data pulled from the Go type checker.
func (r *Reader) readPackage(proj *data.Project, pkgName string, files []*ast.File) *data.Package {
	// Prepare the info for collecting data.
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	// Resolve types in the packages.
	imp := importer.ForCompiler(proj.FileSet, "source", nil)
	conf := types.Config{Importer: imp}
	pkg, err := conf.Check(r.basePath, proj.FileSet, files, info)
	if err != nil {
		log.Fatal("Type Check Failed: ", err)
	}

	// Gather up read results to be returned.
	return &data.Package{
		Name:       pkgName,
		Project:    proj,
		Package:    pkg,
		Types:      info.Types,
		Defs:       info.Defs,
		Uses:       info.Uses,
		Implicits:  info.Implicits,
		Selections: info.Selections,
	}
}

// Read will perform the read of the data.
func (r *Reader) Read() *data.Project {
	fileSet, filesByPackage := r.parseSources()

	proj := &data.Project{
		BasePath: r.basePath,
		FileSet:  fileSet,
		Packages: nil,
	}

	for pkgName, files := range filesByPackage {
		pkg := r.readPackage(proj, pkgName, files)
		proj.Packages = append(proj.Packages, pkg)
	}
	return proj
}
