package jsonify

import "maps"

type contextKey int

const (
	keyMinimized contextKey = iota
	keyShort
	keyOnlyIndex
	keySkipDuplicates
	keyDebugKind
	keyDebugIndex
	keyDebugReceiver
	keyDebugFullLoc
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

// OnlyIndex indicates that objects should output only an index
// as a reference to the rest of the object defined elsewhere.
func (c *Context) OnlyIndex() *Context {
	return c.copyAndSet(keyOnlyIndex, true).copyAndSet(keyShort, true)
}

// Short indicates that objects should output only an identifier
// (i.e. kind and index) as a reference to the rest of the object
// defined elsewhere.
func (c *Context) Short() *Context {
	return c.copyAndSet(keyOnlyIndex, false).copyAndSet(keyShort, true)
}

// Full indicates that objects should output the whole object,
// not a shortened version.
func (c *Context) Full() *Context {
	return c.copyAndSet(keyOnlyIndex, false).copyAndSet(keyShort, false)
}

// IsShort indicates that objects should output only an identifier
// or index as a reference to the whole object defined elsewhere.
func (c *Context) IsShort() bool {
	return c.state[keyShort]
}

// IsOnlyIndex indicates that objects should output only index
// as a reference to the whole object defined elsewhere.
func (c *Context) IsOnlyIndex() bool {
	return c.state[keyOnlyIndex]
}

// SetSkipDuplicates sets the skip duplicate flag.
func (c *Context) SetSkipDuplicates(skip bool) *Context {
	return c.copyAndSet(keySkipDuplicates, skip)
}

// SkipDuplicates indicates that any object marked as a duplicate should
// return a null JSON node instead of a full output, so that the object is
// not outputted when full. The index and short should still output.
func (c *Context) SkipDuplicates() bool {
	return c.state[keySkipDuplicates]
}

// IncludeDebugKind indicates that the kind field should be included
// to the output model for debugging.
func (c *Context) IncludeDebugKind(include bool) *Context {
	return c.copyAndSet(keyDebugKind, include)
}

// IsDebugKindIncluded indicates that the kind field should be included
// to the output model for debugging.
func (c *Context) IsDebugKindIncluded() bool {
	return c.state[keyDebugKind]
}

// IncludeDebugIndex indicates that the index should be included
// to the output model for debugging.
func (c *Context) IncludeDebugIndex(include bool) *Context {
	return c.copyAndSet(keyDebugIndex, include)
}

// IsDebugIndexIncluded indicates that the index should be included
// to the output model for debugging.
func (c *Context) IsDebugIndexIncluded() bool {
	return c.state[keyDebugIndex]
}

// IncludeDebugReceiver sets if the methods should include
// receiver information in the object model for debugging.
func (c *Context) IncludeDebugReceiver(include bool) *Context {
	return c.copyAndSet(keyDebugReceiver, include)
}

// IsDebugReceiverIncluded indicates that methods should include
// receiver information in the object model for debugging.
func (c *Context) IsDebugReceiverIncluded() bool {
	return c.state[keyDebugReceiver]
}

// IncludeDebugFullLoc sets if the full location information should
// be included in the object model for debugging.
func (c *Context) IncludeDebugFullLoc(include bool) *Context {
	return c.copyAndSet(keyDebugFullLoc, include)
}

// IsDebugFullLocIncluded indicates that the full location information should
// be included in the object model for debugging.
func (c *Context) IsDebugFullLocIncluded() bool {
	return c.state[keyDebugFullLoc]
}
