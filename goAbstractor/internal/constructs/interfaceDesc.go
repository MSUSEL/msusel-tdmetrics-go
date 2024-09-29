package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
)

// TODO: Add test for if there are unexported abstracts the interface needs to be locked to a package.

// InterfaceDesc is a named interface typically explicitly defined at the given
// location in the source code. The underlying type description
// can be a class or interface with optional parameter types.
//
// If type parameters are given then the interface is generic.
// Instances with realized versions of the interface,
// are added for each used instance in the source code. If there
// are no instances then the generic interface isn't used.
type InterfaceDesc interface {
	TypeDesc
	TempReferenceContainer
	IsInterfaceDesc()

	// Hint indicates if the interface is a placeholder for something that
	// isn't abstracted directly, such as maps and channels.
	Hint() hint.Hint

	// PinnedPackage is non-nil if the interface is tied to a specific
	// package. This only happens if any of the abstracts are unexported.
	// Unexported methods can only be implemented by something in the same
	// package as the interface is defined.
	PinnedPackage() Package

	// IsPinned indicates if this interface is pinned to a package.
	IsPinned() bool

	// Abstracts is the set of named signatures for this interface.
	Abstracts() []Abstract

	// Exact types are like `string|int|bool` where the
	// data type must match exactly.
	Exact() []TypeDesc

	// Approx types are like `~string|~int` where the data type
	// may be exact or an extension of the base type.
	Approx() []TypeDesc

	// IsGeneral indicates if there is two or more exact or approximate types.
	IsGeneral() bool

	// Implements determines if this interface implements the other interface.
	Implements(other InterfaceDesc) bool

	AddInherits(it InterfaceDesc) InterfaceDesc

	Inherits() collections.SortedSet[InterfaceDesc]
}

type InterfaceDescArgs struct {
	Hint     hint.Hint
	RealType types.Type

	// PinnedPkg is non-nil when an abstract is unexported.
	// This is the package the interface is pinned to.
	PinnedPkg Package

	// Abstracts is the set of named signatures for this interface.
	Abstracts []Abstract

	// Exact types are like `string|int|bool` where the
	// data type must match exactly.
	Exact []TypeDesc

	// Approx types are like `~string|~int` where the data type
	// may be exact or an extension of the base type.
	Approx []TypeDesc

	// Package is needed when the real type isn't given.
	// The package is used to help create the real type.
	Package *packages.Package
}

type InterfaceDescFactory interface {
	NewInterfaceDesc(args InterfaceDescArgs) InterfaceDesc
	InterfaceDescs() collections.ReadonlySortedSet[InterfaceDesc]
}
