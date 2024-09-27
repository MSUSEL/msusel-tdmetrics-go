package hint

// Hint is attached to an interface to indicate that the interface is a
// placeholder for a Go construct that is not directly abstracted.
// The hint is used when generating real Go types when not given.
// The hint is not needed in the output.
type Hint string

const (
	None       Hint = ``
	Pointer    Hint = `pointer`
	List       Hint = `list`
	Map        Hint = `map`
	Chan       Hint = `chan`
	Complex64  Hint = `complex64`
	Complex128 Hint = `complex128`
	Comparable Hint = `comparable`
)
