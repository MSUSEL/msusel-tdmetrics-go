package typeDesc

type Register struct {
	BasicSet     []Basic
	InterfaceSet []Interface
	NamedSet     []Named
	SignatureSet []Signature
	SolidSet     []Solid
	StructSet    []Struct
	UnionSet     []Union
}

func (r *Register) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	index := 0
	index = setIndices(index, r.BasicSet)
	index = setIndices(index, r.InterfaceSet)
	index = setIndices(index, r.NamedSet)
	index = setIndices(index, r.SignatureSet)
	index = setIndices(index, r.SolidSet)
	index = setIndices(index, r.StructSet)
	setIndices(index, r.UnionSet)
}

func setIndices[T TypeDesc](index int, s []T) int {
	for _, t := range s {
		t.SetIndex(index)
		index++
	}
	return index
}

func (r *Register) RegisterBasic(t Basic) Basic {
	return registerType(t, &r.BasicSet)
}

func (r *Register) RegisterInterface(t Interface) Interface {
	return registerType(t, &r.InterfaceSet)
}

func (r *Register) Named(t Named) Named {
	return registerType(t, &r.NamedSet)
}

func (r *Register) RegisterSignature(t Signature) Signature {
	return registerType(t, &r.SignatureSet)
}

func (r *Register) RegisterSolid(t Solid) Solid {
	return registerType(t, &r.SolidSet)
}

func (r *Register) RegisterStruct(t Struct) Struct {
	return registerType(t, &r.StructSet)
}

func (r *Register) RegisterUnion(t Union) Union {
	return registerType(t, &r.UnionSet)
}

func registerType[T TypeDesc](t T, s *[]T) T {
	for _, t2 := range *s {
		if t.Equal(t2) {
			return t2
		}
	}
	*s = append(*s, t)
	return t
}
