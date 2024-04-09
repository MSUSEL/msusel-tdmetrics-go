package typeKind

type TypeKind string

const (
	Array   TypeKind = `array`
	List    TypeKind = `list`
	Chan    TypeKind = `chan`
	Pointer TypeKind = `pointer`
)
