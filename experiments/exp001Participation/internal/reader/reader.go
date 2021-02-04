package reader

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"../filter"
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

// Read will perform the read of the data.
func (r *Reader) Read() {
	runner := newRunner()
	runner.run(r.basePath, r.sources)
}
