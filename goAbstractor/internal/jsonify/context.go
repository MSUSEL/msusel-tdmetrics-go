package jsonify

import "maps"

type contextKey int

const (
	keyMinimized contextKey = iota
	keyShort
	keyKindShown
	keyIdShown
	keyReceiverShown
	keyReferenceShown
	keyInheritorsShown
	keyFullLoc
)

type Context struct {
	state map[contextKey]bool
}

func NewContext() *Context {
	return &Context{state: nil}
}

func (c *Context) copyAndSet(key contextKey, state bool) *Context {
	c2 := &Context{
		state: maps.Clone(c.state),
	}
	if c2.state == nil {
		c2.state = map[contextKey]bool{}
	}
	c2.state[key] = state
	return c2
}

// SetMinimize sets if the output JSON should be minimized.
func (c *Context) SetMinimize(min bool) *Context {
	return c.copyAndSet(keyMinimized, min)
}

// IsMinimized indicates that the output JSON should be minimized.
func (c *Context) IsMinimized() bool {
	return c.state[keyMinimized]
}

// Short indicates that objects should output only an identifier
// as a reference to the rest of the object defined elsewhere.
func (c *Context) Short() *Context {
	return c.copyAndSet(keyShort, true)
}

// Long indicates that objects should output the whole object,
// not the shortened version.
func (c *Context) Long() *Context {
	return c.copyAndSet(keyShort, false)
}

// IsShort indicates that objects should output only an identifier
// as a reference to the rest of the object defined elsewhere.
func (c *Context) IsShort() bool {
	return c.state[keyShort]
}

// ShowKind indicates that the kind field should be added to the output model.
func (c *Context) ShowKind() *Context {
	return c.copyAndSet(keyKindShown, true)
}

// HideKind indicates that the kind field can be skipped.
func (c *Context) HideKind() *Context {
	return c.copyAndSet(keyKindShown, false)
}

// IsKindShown indicates that the kind field should be added to the output model.
func (c *Context) IsKindShown() bool {
	return c.state[keyKindShown]
}

// ShowId indicates that the identifier should be added to the output model.
func (c *Context) ShowId() *Context {
	return c.copyAndSet(keyIdShown, true)
}

// HideId indicates that the identifier can be skipped,
// unless output is short and only the index is outputted.
func (c *Context) HideId() *Context {
	return c.copyAndSet(keyIdShown, false)
}

// IsIdShown indicates that the identifier should be
// added to the output model when not needed.
func (c *Context) IsIdShown() bool {
	return c.state[keyIdShown]
}

// ShowReceiver sets if the methods should include receiver information
// in the object model. This is for debugging purposes.
func (c *Context) ShowReceiver(show bool) *Context {
	return c.copyAndSet(keyReceiverShown, show)
}

// IsReceiverShown indicates that methods should include receiver information
// in the object model. This is for debugging purposes.
func (c *Context) IsReceiverShown() bool {
	return c.state[keyReceiverShown]
}

// ShowReference sets if the typeDef reference should include the
// reference name in the object model. This is for debugging purposes.
func (c *Context) ShowReference(show bool) *Context {
	return c.copyAndSet(keyReferenceShown, show)
}

// IsReferenceShown indicates that typeDef reference should include the
// reference name in the object model. This is for debugging purposes.
func (c *Context) IsReferenceShown() bool {
	return c.state[keyReferenceShown]
}

// ShowInheritors sets if the interfaces inheritors should
// be included in the object model. This is for debugging purposes.
func (c *Context) ShowInheritors(show bool) *Context {
	return c.copyAndSet(keyInheritorsShown, show)
}

// IsInheritorsShown indicates that interfaces inheritors should be
// included in the object model. This is for debugging purposes.
func (c *Context) IsInheritorsShown() bool {
	return c.state[keyInheritorsShown]
}

// ShowFullLocation sets if the full location information should
// be included in the object model. This is for debugging purposes.
func (c *Context) ShowFullLocation(show bool) *Context {
	return c.copyAndSet(keyFullLoc, show)
}

// IsFullLocationShown indicates that the full location information should
// be included in the object model. This is for debugging purposes.
func (c *Context) IsFullLocationShown() bool {
	return c.state[keyFullLoc]
}
