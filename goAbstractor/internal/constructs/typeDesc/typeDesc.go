package typeDesc

// TypeDesc is an interface for all type descriptors.
type TypeDesc interface {

	// _isTypeDesc is to prevent arbitrary things duck-typing to a type.
	_isTypeDesc()
}
