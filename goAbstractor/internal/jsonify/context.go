package jsonify

import (
	"maps"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Context struct {
	data map[string]any
}

func NewContext() *Context {
	return &Context{
		data: map[string]any{},
	}
}

func (c *Context) Copy() *Context {
	return &Context{
		data: maps.Clone(c.data),
	}
}

func read[T any](c *Context, key string) (T, bool) {
	value, exists := c.data[key]
	if exists {
		if t, ok := value.(T); ok {
			return t, true
		}
	}
	return utils.Zero[T](), false
}

func (c *Context) Remove(key string) *Context {
	delete(c.data, key)
	return c
}

func (c *Context) Set(key string, value any) *Context {
	c.data[key] = value
	return c
}

func has[T any](c *Context, key string) bool {
	_, ok := read[T](c, key)
	return ok
}

func (c *Context) HasAny(key string) bool {
	_, ok := c.data[key]
	return ok
}

func (c *Context) HasBool(key string) bool {
	return has[bool](c, key)
}

func (c *Context) HasInt(key string) bool {
	return has[int](c, key)
}

func (c *Context) HasString(key string) bool {
	return has[string](c, key)
}

func get[T any](c *Context, key string) T {
	value, _ := read[T](c, key)
	return value
}

func (c *Context) GetAny(key string) any {
	return c.data
}

func (c *Context) GetBool(key string) bool {
	return get[bool](c, key)
}

func (c *Context) GetInt(key string) int {
	return get[int](c, key)
}

func (c *Context) GetString(key string) string {
	return get[string](c, key)
}
