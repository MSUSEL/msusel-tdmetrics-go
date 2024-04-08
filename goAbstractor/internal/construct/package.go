package construct

import "encoding/json"

type Package struct {
	Path    string
	Imports []string
	Types   []*TypeDef
	Values  []*ValueDef
	Methods []*Method
}

func (p *Package) MarshalJSON() ([]byte, error) {
	data := map[string]any{}
	if len(p.Path) > 0 {
		data[`path`] = p.Path
	}
	if len(p.Imports) > 0 {
		data[`imports`] = p.Imports
	}
	if len(p.Types) > 0 {
		data[`types`] = p.Types
	}
	if len(p.Values) > 0 {
		data[`values`] = p.Values
	}
	if len(p.Methods) > 0 {
		data[`methods`] = p.Methods
	}
	return json.Marshal(data)
}
