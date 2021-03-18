package reader

import (
	"go/types"
)

type BaseTypes struct {
	cache map[types.Type]map[types.Type]bool
}

func NewBaseTypes() *BaseTypes {
	return &BaseTypes{
		cache: map[types.Type]map[types.Type]bool{},
	}
}

// getBaseTypes gets the base types from the given type and add them to the given map.
func (bt *BaseTypes) getBaseTypes(t types.Type) map[types.Type]bool {
	touchedTypes := map[types.Type]bool{}
	baseTypes := map[types.Type]bool{}

	bt.recurseBaseTypes(t, touchedTypes, baseTypes)

	// Now that this type has been fully computed, cache it.
	if bt != nil {
		bt.cache[t] = baseTypes
	}
	return baseTypes
}

// recurseBaseTypes is the internal recursive method for `getBaseTypes`.
func (bt *BaseTypes) recurseBaseTypes(t types.Type, touchedTypes, baseTypes map[types.Type]bool) {
	if t == nil {
		return
	}

	// Check if this type has been fully computed before.
	if bt != nil {
		if b, ok := bt.cache[t]; ok {
			for key := range b {
				baseTypes[key] = true
			}
			return
		}
	}

	// Check if this type has been reached yet.
	if _, touched := touchedTypes[t]; touched {
		return
	}
	touchedTypes[t] = true

	// Determine which type
	switch t2 := t.(type) {
	case *types.Array:
		bt.recurseBaseTypes(t2.Elem(), touchedTypes, baseTypes)

	case *types.Chan:
		bt.recurseBaseTypes(t2.Elem(), touchedTypes, baseTypes)

	case *types.Interface:
		// Since interfaces don't have internal types,
		// only the types referenced in the functions of the interface,
		// just add the interface and no inner function types.
		baseTypes[t2] = true

	case *types.Map:
		bt.recurseBaseTypes(t2.Key(), touchedTypes, baseTypes)
		bt.recurseBaseTypes(t2.Elem(), touchedTypes, baseTypes)

	case *types.Named:
		// For named types the top type and its underlying type should both be added.
		baseTypes[t2] = true
		bt.recurseBaseTypes(t2.Underlying(), touchedTypes, baseTypes)

	case *types.Pointer:
		bt.recurseBaseTypes(t2.Elem(), touchedTypes, baseTypes)

	case *types.Signature:
		if recv := t2.Recv(); recv != nil {
			bt.recurseBaseTypes(recv.Type(), touchedTypes, baseTypes)
		}
		params := t2.Params()
		for i := 0; i < params.Len(); i++ {
			param := params.At(i)
			bt.recurseBaseTypes(param.Type(), touchedTypes, baseTypes)
		}

	case *types.Slice:
		bt.recurseBaseTypes(t2.Elem(), touchedTypes, baseTypes)

	case *types.Struct:
		for i := 0; i < t2.NumFields(); i++ {
			bt.recurseBaseTypes(t2.Field(i).Type(), touchedTypes, baseTypes)
		}

	default:
		baseTypes[t] = true
	}
}
