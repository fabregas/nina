package nn

import (
	"reflect"
	"syscall/js"
)

type Component interface {
	View() Node

	AddCleanup(func())
	destroy()

	parent() Component
	setParent(Component)

	setContext(typ reflect.Type, val any)
	getContext(typ reflect.Type) (any, bool)
	exportContext() any
	importContext(any)
}

// componentNode — is a special це tree node that contain whole component
type componentNode struct {
	comp Component

	key        string
	hash       string
	lastRender Node
	parentDOM  js.Value
}

func (c *componentNode) isNode() {}

func (c *componentNode) getKey() string {
	return c.key
}
func (c *componentNode) isNil() bool {
	return c == nil
}

func (c *componentNode) Key(key string) *componentNode {
	c.key = key
	return c
}

func (c *componentNode) AsNode() Node {
	return c
}

func Comp(comp Component) *componentNode {
	return &componentNode{
		comp: comp,
	}
}

// Pure — implement this interface if the component is heavy
// and you want to control when it is redrawn
type Pure interface {
	// Hash should return a string. If the string has not changed since the last time,
	// the framework WILL NOT call View() and will skip this component.
	Hash() string
}

type BaseComponent struct {
	cleanups   []func()
	contextMap map[reflect.Type]any
	parentC    Component
}

func (c BaseComponent) AsNode() Node {
	panic("system error: called ToNode at component object")
}

func (c *BaseComponent) AddCleanup(fn func()) {
	c.cleanups = append(c.cleanups, fn)
}

func (c *BaseComponent) destroy() {
	for _, cleanup := range c.cleanups {
		cleanup()
	}
	c.cleanups = nil
}

func (c *BaseComponent) setParent(p Component) {
	c.parentC = p
}

func (c *BaseComponent) parent() Component {
	return c.parentC
}

func (b *BaseComponent) setContext(typ reflect.Type, val any) {
	if b.contextMap == nil {
		b.contextMap = make(map[reflect.Type]any)
	}
	b.contextMap[typ] = val
}

func (b *BaseComponent) getContext(typ reflect.Type) (any, bool) {
	if b.contextMap == nil {
		return nil, false
	}
	val, ok := b.contextMap[typ]
	return val, ok
}

func (b *BaseComponent) importContext(c any) {
	newCtx := c.(map[reflect.Type]any)
	b.contextMap = newCtx
}

func (b *BaseComponent) exportContext() any {
	return b.contextMap
}

func ProvideContext[T any](c Component, value T) {
	var dummy *T

	typ := reflect.TypeOf(dummy).Elem()

	c.setContext(typ, value)
}

func GetContext[T any](startNode Component) T {
	var dummy *T
	typ := reflect.TypeOf(dummy).Elem()

	curr := startNode

	for curr != nil {
		if val, ok := curr.getContext(typ); ok {
			return val.(T)
		}
		curr = curr.parent()
	}

	var zero T
	return zero
}
