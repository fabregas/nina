package nn

import "syscall/js"

// any struct that has View method is Component
type Component interface {
	View() *Element
}

// componentNode — is a special це tree node that contain whole component
type componentNode struct {
	comp Component

	key        string
	hash       string
	lastRender *Element
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

func (c *componentNode) ToNode() Node {
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

type BaseComponent struct{}

func (c BaseComponent) ToNode() Node {
	panic("system error: called ToNode at component object")
}
