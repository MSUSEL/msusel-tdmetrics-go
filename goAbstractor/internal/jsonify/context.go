package jsonify

type Context struct {
	Minimize      bool
	OnlyIndex     bool
	NoKind        bool
	ShowReceivers bool
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Copy() *Context {
	b := *c
	return &b
}
