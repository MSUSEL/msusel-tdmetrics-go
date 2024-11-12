package innate

// These are constant construct names innate to the abstractor and not
// part of Go. These are prepended with a `$` to avoid duck-typing
// with user-defined types.
const (
	Builtin = `$builtin` // Name for builtin package.
	Compare = `$compare` // Name for compare method in comparable.
	Data    = `$data`    // Name for synthetic field in an object.
	Deref   = `$deref`   // Name for a pointer method to get underlying type.
	Get     = `$get`     // Name for getting a values from slice, map, etc.
	Imag    = `$image`   // Name for getting imaginary part of a complex number.
	Len     = `$len`     // Name for getting a length from slice, map, etc.
	Real    = `$real`    // Name for getting real part of a complex number
	Recv    = `$recv`    // Name for receiving a value from a channel.
	Send    = `$send`    // Name for sending a value to a channel.
	Set     = `$set`     // Name for setting a value in a slice, map, etc.
)

// Is determines if the given string is an innate name.
func Is(name string) bool {
	switch name {
	case Builtin, Compare, Data, Deref, Get, Imag, Len, Real, Recv, Send, Set:
		return true
	}
	return false
}
