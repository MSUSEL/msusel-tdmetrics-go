package processors

import (
	"fmt"
	"sort"
	"strings"

	"../handlers"
	"../parser"
)

// FuncUsage counts the number of times an interface or structure
// is a parameter in a function.
type FuncUsage struct {
	usage map[string]int
}

// NewFuncUsage constructs a new function usage processor.
func NewFuncUsage() *FuncUsage {
	return &FuncUsage{
		usage: map[string]int{},
	}
}

// ProcessFunction processes function information pulled from the parsed files.
func (p *FuncUsage) ProcessFunction() parser.ProcessHandler {
	return func(filename string, data interface{}) {
		switch x := data.(type) {
		case *handlers.FuncData:
			p.funcDataType(x)
		}
	}
}

// funcDataType process the function data by counting the parameter and receiver types.
func (p *FuncUsage) funcDataType(d *handlers.FuncData) {
	for _, t := range d.ParamTypes {
		p.usage[t]++
	}
}

// String will output the result from this processor as a string.
func (p *FuncUsage) String() string {
	countMaxWidth := 0
	for _, c := range p.usage {
		width := len(fmt.Sprintf(`%d`, c))
		if width > countMaxWidth {
			countMaxWidth = width
		}
	}

	parts := make([]string, len(p.usage))
	i := 0
	format := fmt.Sprintf(`%%%dd %%s`, countMaxWidth)
	for t, c := range p.usage {
		parts[i] = fmt.Sprintf(format, c, t)
		i++
	}

	sort.Strings(parts)
	return strings.Join(parts, "\n")
}
