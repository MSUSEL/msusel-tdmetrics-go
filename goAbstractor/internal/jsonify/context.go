package jsonify

import "maps"

type Context struct {
	state map[string]bool
}

func NewContext() *Context {
	return &Context{
		state: map[string]bool{},
	}
}

func (c *Context) copyAndSet(name string, state bool) *Context {
	c2 := &Context{
		state: maps.Clone(c.state),
	}
	if c2.state == nil {
		c2.state = map[string]bool{}
	}
	c2.state[name] = state
	return c2
}

// SetMinimize sets if the output JSON should be minimized.
func (c *Context) SetMinimize(min bool) *Context {
	return c.copyAndSet(`minimized`, min)
}

// IsMinimized indicates that the output JSON should be minimized.
func (c *Context) IsMinimized() bool {
	return c.state[`minimized`]
}

// Short indicates that objects should output only an index or name
// as a reference to the rest of the object defined elsewhere.
func (c *Context) Short() *Context {
	return c.copyAndSet(`short`, true)
}

// Long indicates that objects should output the whole object,
// not the shortened version.
func (c *Context) Long() *Context {
	return c.copyAndSet(`short`, false)
}

// IsShort indicates that objects should output only an index or name
// as a reference to the rest of the object defined elsewhere.
func (c *Context) IsShort() bool {
	return c.state[`short`]
}

// ShowKind indicates that the kind field should be added to the output model.
func (c *Context) ShowKind() *Context {
	return c.copyAndSet(`kindShown`, true)
}

// HideKind indicates that the "kind" field can be skipped.
func (c *Context) HideKind() *Context {
	return c.copyAndSet(`kindShown`, false)
}

// IsKindShown indicates that the kind field should be added to the output model.
func (c *Context) IsKindShown() bool {
	return c.state[`kindShown`]
}

// ShowReceiver sets if the methods should include receiver information
// in the object model. This is for debugging purposes.
func (c *Context) ShowReceiver(show bool) *Context {
	return c.copyAndSet(`receiverShown`, show)
}

// IsReceiverShown indicates that methods should include receiver information
// in the object model. This is for debugging purposes.
func (c *Context) IsReceiverShown() bool {
	return c.state[`receiverShown`]
}

// ShowReference sets if the typeDef reference should include the
// reference name in the object model. This is for debugging purposes.
func (c *Context) ShowReference(show bool) *Context {
	return c.copyAndSet(`referenceShown`, show)
}

// IsReferenceShown indicates that typeDef reference should include the
// reference name in the object model. This is for debugging purposes.
func (c *Context) IsReferenceShown() bool {
	return c.state[`referenceShown`]
}
