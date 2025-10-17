package constructs

type Project struct {
	Abstracts  *AbstractFactory
	Arguments  *ArgumentFactory
	Signatures *SignatureFactory
}

func NewProject() *Project {
	p := &Project{}
	p.Abstracts = newAbstractFactory(p)
	p.Arguments = newArgumentFactory(p)
	p.Signatures = newSignatureFactory(p)
	return p
}

func (p *Project) AllFactories() []any {
	return []any{
		p.Abstracts,
		p.Arguments,
		p.Signatures,
	}
}
