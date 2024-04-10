package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"

type Method struct {
	Name      string              `json:"name"`
	Signature *typeDesc.Signature `json:"signature"`
	Receiver  typeDesc.TypeDesc   `json:"receiver,omitempty"`
}
