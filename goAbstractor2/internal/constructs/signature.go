package constructs

import "go/types"

type Signature struct {
	// TODO: Implement
}

type SignatureFactory struct {
	*Factory[*types.Signature, Signature]
}

func newSignatureFactory(proj *Project) *SignatureFactory {
	return &SignatureFactory{
		Factory: NewFactory(proj, func(proj *Project, src *types.Signature, sig *Signature) {

			// TODO: Implement
		}),
	}
}

func (s Signature) Kind() string { return `signature` }
