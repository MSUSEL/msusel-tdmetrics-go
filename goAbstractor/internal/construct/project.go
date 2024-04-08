package construct

import "encoding/json"

type Project struct {
	Packages []*Package
}

func (p *Project) MarshalJSON() ([]byte, error) {
	data := map[string]any{
		`language`: `go`,
	}
	if len(p.Packages) > 0 {
		data[`packages`] = p.Packages
	}
	return json.Marshal(data)
}
