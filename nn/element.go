package nn

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

// Element represent any HTML tag
type Element struct {
	tag       string
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

func (e *Element) ClassFunc(f func() string) *Element {
	return e.Class(f())
}

func (e *Element) Attr(key, value string) *Element {
	e.addAttr(key, value)
	return e
}

func (e *Element) BoolAttr(name string, condition bool) *Element {
	if condition {
		return e.Attr(name, "")
	}
	return e
}

func (e *Element) ID(id string) *Element         { return e.Attr("id", id) }
func (e *Element) Href(href string) *Element     { return e.Attr("href", href) }
func (e *Element) Src(src string) *Element       { return e.Attr("src", src) }
func (e *Element) Alt(alt string) *Element       { return e.Attr("alt", alt) }
func (e *Element) Type(t string) *Element        { return e.Attr("type", t) }
func (e *Element) Value(v string) *Element       { return e.Attr("value", v) }
func (e *Element) Placeholder(p string) *Element { return e.Attr("placeholder", p) }
func (e *Element) Disabled(d bool) *Element      { return e.BoolAttr("disabled", d) }
func (e *Element) Checked(c bool) *Element       { return e.BoolAttr("checked", c) }

// for passing inline-styles as string "color: red; margin: 10px;"
func (e *Element) Style(css string) *Element { return e.Attr("style", css) }

func (e *Element) On(eventName string, handler func(Event)) *Element {
	e.addListener(eventName, handler)
	return e
}

func (e *Element) OnInput(handler func(Event)) *Element      { return e.On("input", handler) }
func (e *Element) OnChange(handler func(Event)) *Element     { return e.On("change", handler) }
func (e *Element) OnSubmit(handler func(Event)) *Element     { return e.On("submit", handler) }
func (e *Element) OnKeyDown(handler func(Event)) *Element    { return e.On("keydown", handler) }
func (e *Element) OnKeyUp(handler func(Event)) *Element      { return e.On("keyup", handler) }
func (e *Element) OnMouseEnter(handler func(Event)) *Element { return e.On("mouseenter", handler) }
func (e *Element) OnMouseLeave(handler func(Event)) *Element { return e.On("mouseleave", handler) }
func (e *Element) OnClick(handler func(Event)) *Element      { return e.On("click", handler) }

func (e *Element) Children(children ...Node) *Element {
	for _, n := range children {
		if n != nil {
			e.children = append(e.children, n)
		}
	}
	return e
}

func (e *Element) Bind(target any) *Element {
	switch ptr := target.(type) {
	case *string:
		e.Value(*ptr)
		e.OnInput(func(ev Event) {
			*ptr = ev.TargetValue()
		})

	case *int:
		e.Value(fmt.Sprintf("%d", *ptr))
		e.OnInput(func(ev Event) {
			val, err := strconv.Atoi(ev.TargetValue())
			if err == nil {
				*ptr = val
			}
		})

	case *bool:
		e.Checked(*ptr)
		e.OnChange(func(ev Event) {
			*ptr = ev.TargetChecked()
		})

	default:
		// FAIL FAST
		panic(fmt.Sprintf(
			"[Nina Framework] Critical Error in Bind(): expected a pointer to a basic type (*string, *int, *bool), but got %T",
			target,
		))
	}

	return e
}

func (e *Element) Text(text string) *Element {
	e.children = append(e.children, &TextNode{value: text})
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
