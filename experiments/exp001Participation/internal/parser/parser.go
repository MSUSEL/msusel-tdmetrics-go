package parser

import (
	"go/ast"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type (
	// OnErrorHander is called when the parser has an error.
	OnErrorHander func(err error)

	// UpdateProgressHandler is called periodically to indicate the current parser progress.
	UpdateProgressHandler func(finished, total int)

	// ProcessHandler processes data pulled from the parsed files.
	// These calls will have been synchronizes.
	ProcessHandler func(filename string, data interface{})

	// NodeHandler will process the given asynchronously node
	// and return any number of data objects.
	NodeHandler func(n ast.Node) []interface{}

	// NodeHandlerFactory creates a new node handler for each file.
	NodeHandlerFactory func(filename string) []NodeHandler

	// FilenameMatcher is a handler for comparing if a filename matches some criteria.
	FilenameMatcher func(filename string) bool

	// Parser is the information used to setup how the parser
	// will run and the data to run inside the parser.
	Parser struct {
		onErrorHndl    OnErrorHander
		updateProgress UpdateProgressHandler
		processors     []ProcessHandler
		hndlFactories  []NodeHandlerFactory
		filenames      []string
		sources        []interface{}
	}
)

// NewParser constructs a new configuration to parse with.
func NewParser() *Parser {
	return &Parser{
		onErrorHndl:    nil,
		updateProgress: nil,
		processors:     make([]ProcessHandler, 0),
		hndlFactories:  make([]NodeHandlerFactory, 0),
		filenames:      make([]string, 0),
		sources:        make([]interface{}, 0),
	}
}

// OnError will set the function handler to call to when an error occurs.
// If nil then this will print out to the terminal.
func (c *Parser) OnError(hndl OnErrorHander) *Parser {
	c.onErrorHndl = hndl
	return c
}

// UpdateProgress will set the function handler to call to when progress is updated.
// If nil then this will print out to the terminal.
func (c *Parser) UpdateProgress(hndl UpdateProgressHandler) *Parser {
	c.updateProgress = hndl
	return c
}

// AddProcessor adds the processors that the parser will send
// all of its gathered information to.
func (c *Parser) AddProcessor(proc ProcessHandler) *Parser {
	c.processors = append(c.processors, proc)
	return c
}

// AddHandler will insert a new node handler which will be called during the parse.
func (c *Parser) AddHandler(hndl NodeHandlerFactory) *Parser {
	c.hndlFactories = append(c.hndlFactories, hndl)
	return c
}

// AddFile adds a new file to parser to this configuration.
func (c *Parser) AddFile(filenames ...string) *Parser {
	c.filenames = append(c.filenames, filenames...)
	count := len(filenames)
	sources := make([]interface{}, count)
	c.sources = append(c.sources, sources...)
	return c
}

// AddSource adds a string of source code to parse directly.
// The filename is used to identify this source code string.
func (c *Parser) AddSource(filename, source string) *Parser {
	c.filenames = append(c.filenames, filename)
	c.sources = append(c.sources, source)
	return c
}

// AddDir adds all the files in the given directory.
func (c *Parser) AddDir(foldername string, recursive bool) *Parser {
	fileInfo, err := ioutil.ReadDir(foldername)
	if err != nil {
		panic(err)
	}
	files := make([]string, len(fileInfo))
	for i, info := range fileInfo {
		files[i] = path.Clean(path.Join(foldername, info.Name()))
	}
	return c.AddFile(files...)
}

// AddDirRecursively adds all the files in the given directory and
// all the files in all children folders.
func (c *Parser) AddDirRecursively(dirname string) *Parser {
	var files []string
	err := filepath.Walk(dirname,
		func(filepath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				filename := filepath
				if !path.IsAbs(filename) {
					filename = path.Clean(path.Join(dirname, filepath))
				}
				files = append(files, filename)
			}
			return nil
		})
	if err != nil {
		panic(err)
	}
	return c.AddFile(files...)
}

// FilterFilenames will remove all the filenames and paired sources which
// match the given handler. If the mater returns true the file is removed.
func (c *Parser) FilterFilenames(matcher FilenameMatcher) *Parser {
	remove := map[int]bool{}
	for i, filename := range c.filenames {
		if matcher(filename) {
			remove[i] = true
		}
	}

	count := len(c.filenames) - len(remove)
	filenames := make([]string, count)
	sources := make([]interface{}, count)
	j := 0
	for i, filename := range c.filenames {
		if !remove[i] {
			filenames[j] = filename
			sources[j] = c.sources[i]
			j++
		}
	}

	c.filenames = filenames
	c.sources = sources
	return c
}

// Filenames get the list of filenames that have been set to this parser
// and should be processed when started.
func (c *Parser) Filenames() []string {
	filenames := make([]string, len(c.filenames))
	copy(filenames, c.filenames)
	return filenames
}

// Start will start the asynchronous parse and
// return a control to manage the running parse.
func (c *Parser) Start() RunningParser {
	p := newParserRunner(*c)
	p.Start()
	return p
}
