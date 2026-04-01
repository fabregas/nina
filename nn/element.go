package nn

import (
	"strings"
	"syscall/js"
)

// Element represent any HTML tag
type Element struct {
	tag       string
	id        string
	classes   string
	attrs     map[string]string
	listeners map[string]func(Event) // event listeners
	children  []Node
	key       string

	// reference to real HTML element in browser
	// this value will be setup after first render
	domNode js.Value

	activeCallbacks []js.Func
}

func (e *Element) isNode() {}

func (e *Element) getKey() string {
	return e.key
}

func (e *Element) isNil() bool {
	return e == nil
}

func (e *Element) addAttr(key, val string) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}
	e.attrs[key] = val
}

func (e *Element) addListener(event string, lf func(Event)) {
	if e.listeners == nil {
		e.listeners = make(map[string]func(Event))
	}

	e.listeners[event] = lf
}

func (e *Element) ID(id string) *Element {
	e.id = id
	return e
}

func (e *Element) Key(key string) *Element {
	e.key = key
	return e
}

func (e *Element) Class(classes ...string) *Element {
	if len(classes) > 0 {
		if e.classes != "" {
			e.classes += " "
		}
		e.classes += strings.Join(classes, " ")
	}

	return e
}

func (e *Element) Attrs(key, value string) *Element {
	e.addAttr(key, value)
	return e
}

func (e *Element) Children(children ...Node) *Element {
	for _, n := range children {
		if n != nil {
			e.children = append(e.children, n)
		}
	}
	return e
}

func (e *Element) Text(text string) *Element {
	e.children = append(e.children, &TextNode{value: text})
	return e
}

func (e *Element) OnClick(handler func()) *Element {
	e.addListener("click", func(event Event) { handler() })
	return e
}

func (e *Element) OnClickEvent(handler func(Event)) *Element {
	e.addListener("click", handler)
	return e
}

func (e *Element) Value(val string) *Element {
	e.addAttr("value", val)
	return e
}

func (e *Element) OnInput(handler func(string)) *Element {
	e.addListener("input", func(event Event) {
		val := event.TargetValue()

		handler(val)
	})

	return e
}

// TextNode for general text
type TextNode struct {
	value string

	// reference to real HTML element in browser
	// this value will be setup after first render
	domNode js.Value
}

func (t *TextNode) isNode() {}

func (t *TextNode) getKey() string {
	return ""
}

func (t *TextNode) isNil() bool {
	return t == nil
}

////////////

func Text(v string) *TextNode {
	return &TextNode{value: v}
}

func Div() *Element {
	return &Element{tag: "div"}
}

func H1() *Element {
	return &Element{tag: "h1"}
}

func Input() *Element {
	return &Element{tag: "input"}
}

func A() *Element {
	return &Element{tag: "a"}
}

func Button() *Element {
	return &Element{tag: "button"}
}
