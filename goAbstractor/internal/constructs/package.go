package constructs

type Package struct {
	Path    string      `json:"path,omitempty"`
	Imports []string    `json:"imports,omitempty"`
	Types   []*TypeDef  `json:"types,omitempty"`
	Values  []*ValueDef `json:"values,omitempty"`
	Methods []*Method   `json:"methods,omitempty"`
}
