package constructs

import "go/types"

type Argument struct {
	Name string `json:"name,omitempty"`
	Type Type   `json:"type"`
}

type ArgumentFactory struct {
	*Factory[*types.Signature, Argument]
}

func newArgumentFactory(proj *Project) *ArgumentFactory {
	return &ArgumentFactory{
		Factory: NewFactory(proj, func(proj *Project, src *types.Signature, arg *Argument) {

			// TODO: Implement
		}),
	}
}

func (s Argument) Kind() string { return `argument` }
