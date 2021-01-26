package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

const (
	// channelSize is the amount of messages that can be
	// in a channel at one time before blocking a thread.
	channelSize = 1000
)

type (
	// RunningParser is the interface for handling a parse
	// which is running asynchronously.
	RunningParser interface {

		// Cancel will stop the parse from running.
		Cancel()

		// Await will block until the parse has finished.
		Await()

		// AwaitChn will return a channel that can be selected on
		// and will be closed when the parse has finished.
		AwaitCh() chan bool
	}

	// errorMessage is a message to pass an error through the synchronizing channel.
	errorMessage struct {
		Filename string
		Error    error
	}

	// finishedMessage a message that indicate a parse of a file
	// has finished through the synchronizing channel.
	finishedMessage struct {
		Filename string
	}

	// dataMessage a message to pass data through the synchronizing channel.
	dataMessage struct {
		Filename string
		Values   []interface{}
	}

	// parserRunner is a tool for asynchronously parsing several Go source files.
	parserRunner struct {
		Parser
		fileSet   *token.FileSet
		msgCh     chan interface{}
		doneCh    chan bool
		cancelCh  chan bool
		finished  int
		cancelled bool
	}
)

// Checks that the parserRunning implements the RunningParser interface.
var _ RunningParser = (*parserRunner)(nil)

// recoveryError converts a recovered panic into an error.
func recoveryError(r interface{}) error {
	switch v := r.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return errors.New(fmt.Sprint(r))
	}
}

// newParserRunner constructs a new parser with the given parser configuration.
func newParserRunner(config Parser) *parserRunner {
	return &parserRunner{
		Parser:    config,
		fileSet:   token.NewFileSet(),
		msgCh:     make(chan interface{}, channelSize),
		doneCh:    make(chan bool),
		cancelCh:  make(chan bool),
		finished:  0,
		cancelled: false,
	}
}

// Cancel will stop the parse from running.
func (p *parserRunner) Cancel() {
	p.cancelled = true
	p.cancelCh <- true
}

// Await will block until the parse has finished.
func (p *parserRunner) Await() {
	<-p.doneCh
}

// AwaitChn will return a channel that can be selected on
// and will be closed when the parse has finished.
func (p *parserRunner) AwaitCh() chan bool {
	return p.doneCh
}

// Start will start parsing the files asynchronously.
// The returned channel can be used to await the parse to finish.
func (p *parserRunner) Start() {
	for i, filename := range p.filenames {
		source := p.sources[i]
		handlers := []NodeHandler{}
		for _, factory := range p.hndlFactories {
			handlers = append(handlers, factory(filename)...)
		}
		go p.startParse(filename, source, handlers)
	}
	go p.startCollection()
}

// handleError handles an error which occurred
func (p *parserRunner) handleError(filename string, err error) {
	if p.onErrorHndl == nil {
		DefaultOnError(filename, err)
		return
	}
	p.onErrorHndl(err)
}

// handleUpdateProgress handles the progress being updated.
func (p *parserRunner) handleUpdateProgress(finished, total int) {
	if p.updateProgress == nil {
		DefaultUpdateProgress(finished, total)
		return
	}
	p.updateProgress(finished, total)
}

// startCollection will initialize the channel reader
// to synchronize calls into the processor.
func (p *parserRunner) startCollection() {
	defer func() {
		if r := recover(); r != nil {
			p.handleError("parser", recoveryError(r))
		}
		p.doneCh <- true
	}()

	p.handleUpdateProgress(0, len(p.filenames))
	for !p.cancelled {
		select {
		case msg := <-p.msgCh:
			if p.processMessage(msg) {
				return
			}
		case <-p.cancelCh:
			return
		}
	}
}

// processMessage will process a received message from the processors.
func (p *parserRunner) processMessage(msg interface{}) bool {
	switch m := msg.(type) {
	case *errorMessage:
		p.handleError(m.Filename, m.Error)

	case *finishedMessage:
		p.finished++
		total := len(p.filenames)
		p.handleUpdateProgress(p.finished, total)
		if p.finished >= total {
			return true
		}

	case *dataMessage:
		for _, data := range m.Values {
			for _, proc := range p.processors {
				proc(m.Filename, data)
			}
		}
	}
	return false
}

// startParse starts an asynchronous parse of the given
// filename and optional source code.
func (p *parserRunner) startParse(filename string, src interface{}, handlers []NodeHandler) {
	defer func() {
		if r := recover(); r != nil {
			p.msgCh <- &errorMessage{
				Filename: filename,
				Error:    recoveryError(r),
			}
		}
		p.msgCh <- &finishedMessage{
			Filename: filename,
		}
	}()

	f, err := parser.ParseFile(p.fileSet, filename, src, 0)
	if err != nil {
		panic(err)
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if n != nil {
			for _, hndl := range handlers {
				if values := hndl(n); len(values) > 0 {
					p.msgCh <- &dataMessage{
						Filename: filename,
						Values:   values,
					}
				}
			}
		}
		return !p.cancelled
	})
}
