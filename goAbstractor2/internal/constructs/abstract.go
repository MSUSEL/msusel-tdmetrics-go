package constructs

import "go/types"

type Abstract struct {
	Name      string     `json:"name"`
	Signature *Signature `json:"signature"`
}

type AbstractFactory struct {
	*Factory[*types.Func, Abstract]
}

func newAbstractFactory(proj *Project) *AbstractFactory {
	return &AbstractFactory{
		Factory: NewFactory(proj, func(proj *Project, src *types.Func, ab *Abstract) {
			ab.Name = src.Name()
			ab.Signature = proj.Signatures.New(src.Signature())
		}),
	}
}

func (a Abstract) Kind() string { return `abstract` }
