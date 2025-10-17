package constructs

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
)

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
	// This does not contain any additional abstracts, only the original
	// that this interface was defined with.
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

	// AdditionalAbstracts the set of additional named signatures for this
	// interface that were resolved from an underlying type.
	// The additional abstracts are not used in a comparison since they are
	// resolved from an underlying type they should match anyway once resolved.
	AdditionalAbstracts() []Abstract

	// SetAdditionalAbstracts overrides any prior additional abstracts with
	// the given set of abstracts.
	SetAdditionalAbstracts(abstracts []Abstract)

	// AddInherits tries to add an interface this interface inherits from.
	// Returns the given interface or the equivalent interface that already existed.
	AddInherits(it InterfaceDesc) InterfaceDesc

	// Inherits is the set of interfaces this interface inherits from.
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
	Factory
	NewInterfaceDesc(args InterfaceDescArgs) InterfaceDesc
	InterfaceDescs() collections.ReadonlySortedSet[InterfaceDesc]
}
