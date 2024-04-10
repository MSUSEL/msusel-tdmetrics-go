package wrapKind

type WrapKind string

const (
	Array   WrapKind = `array`
	List    WrapKind = `list`
	Chan    WrapKind = `chan`
	Pointer WrapKind = `pointer`
)
