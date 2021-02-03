package predeclared

// IsAny indicates if the given string is a predeclared type, function, or constant.
func IsAny(val string) bool {
	return types[val] || functions[val] || constants[val]
}

// IsType indicates if the given string is a predeclared type.
func IsType(val string) bool {
	return types[val]
}

// IsFunction indicates if the given string is a predeclared function.
func IsFunction(val string) bool {
	return functions[val]
}

// IsConstant indicates if the given string is a predeclared constant.
func IsConstant(val string) bool {
	return constants[val]
}

// createSet prepares a string set for predeclared types, functions, or constants.
func createSet(defs ...string) map[string]bool {
	result := map[string]bool{}
	for _, def := range defs {
		result[def] = true
	}
	return result
}

// types is the string set of predeclared types.
var types = createSet(
	`bool`,
	`byte`,
	`complex64`,
	`complex128`,
	`error`,
	`float32`,
	`float64`,
	`int`,
	`int8`,
	`int16`,
	`int32`,
	`int64`,
	`rune`,
	`string`,
	`uint`,
	`uint8`,
	`uint16`,
	`uint32`,
	`uint64`,
	`uintptr`)

// functions is the string set of predeclared functions.
var functions = createSet(
	`append`,
	`cap`,
	`close`,
	`complex`,
	`copy`,
	`delete`,
	`imag`,
	`len`,
	`make`,
	`new`,
	`panic`,
	`print`,
	`println`,
	`real`,
	`recover`)

// constants is the string set of predeclared constants.
var constants = createSet(
	`false`,
	`iota`,
	`nil`,
	`true`)
