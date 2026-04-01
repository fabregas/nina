package nn

import "syscall/js"

// any struct that has View method is Component
type Component interface {
	View() *Element
}

// ComponentNode — is a special це tree node that contain whole component
type ComponentNode struct {
	comp Component

	key        string
	hash       string
	lastRender *Element
	parentDOM  js.Value
}

func (c *ComponentNode) isNode() {}

func (c *ComponentNode) getKey() string {
	return c.key
}
func (c *ComponentNode) isNil() bool {
	return c == nil
}

func (c *ComponentNode) Key(key string) *ComponentNode {
	c.key = key
	return c
}

func C(comp Component) *ComponentNode {
	return &ComponentNode{
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
