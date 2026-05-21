package nn

import (
	"fmt"
	"strconv"
	"strings"
)

// Element represent any HTML tag
type Element struct {
	tag       string
	classes   string
	attrs     map[string]string
	children  []Node
	rawHTML   string
	key       string
	listeners *listenersInfo

	// reference to real HTML element in browser
	// this value will be setup after first render
	domNode NativeNode
	refs    []*Ref
}

type eventInfo struct {
	name     string
	isGlobal bool
}

type listenersInfo struct {
	events          map[eventInfo]func(Event)
	activeCallbacks map[eventInfo]func()
	parentComponent Component // for re-render from element's callbacks
}

var globalClassPreprocessor func(string) string

func SetGlobalClassPreprocessor(f func(string) string) {
	globalClassPreprocessor = f
}

func (e *Element) isNode() {}

func (e *Element) getKey() string {
	return e.key
}

func (e *Element) isNil() bool {
	return e == nil
}

func (e *Element) Ref(r *Ref) *Element {
	for _, ro := range e.refs {
		if ro == r {
			return e
		}
	}
	e.refs = append(e.refs, r)

	return e
}

func (e *Element) addAttr(key, val string) {
	if e.attrs == nil {
		e.attrs = make(map[string]string)
	}
	e.attrs[key] = val
}

func (e *Element) addListener(event string, lf func(Event), isGlobal bool) {
	if e.listeners == nil {
		e.listeners = &listenersInfo{
			events:          make(map[eventInfo]func(Event)),
			activeCallbacks: make(map[eventInfo]func()),
		}
	}

	ei := eventInfo{name: event, isGlobal: isGlobal}
	if handler, ok := e.listeners.events[ei]; ok {
		newHandler := func(e Event) {
			handler(e)
			lf(e)
		}
		e.listeners.events[ei] = newHandler
	} else {
		e.listeners.events[ei] = lf
	}
}

func (e *Element) compClasses() {
	if globalClassPreprocessor == nil {
		return
	}

	e.classes = globalClassPreprocessor(e.classes)
}

func (e *Element) InnerHTML(html string) *Element {
	e.rawHTML = html
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

func (e *Element) Clone() *Element {
	clone := &Element{
		tag:     e.tag,
		classes: e.classes,
		key:     e.key,
		rawHTML: e.rawHTML,
		refs:    e.refs,
	}

	if e.attrs != nil {
		clone.attrs = make(map[string]string, len(e.attrs))
		for k, v := range e.attrs {
			clone.attrs[k] = v
		}
	}

	if e.listeners != nil {
		for einfo, handler := range e.listeners.events {
			clone.addListener(einfo.name, handler, einfo.isGlobal)
		}
	}

	if len(e.children) > 0 {
		clone.children = append([]Node{}, e.children...)
	}

	return clone
}

func (e *Element) Merge(other *Element) *Element {
	if other.classes != "" {
		e.Class(other.classes)
	}

	for k, v := range other.attrs {
		e.Attr(k, v)
	}

	if other.listeners != nil {
		for ei, handler := range other.listeners.events {
			e.addListener(ei.name, handler, ei.isGlobal)
		}
	}

	if other.key != "" {
		e.key = other.key

	}
	for _, r := range other.refs {
		e.Ref(r)
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

func (e *Element) AttrIf(condition bool, key, value string) *Element {
	if condition {
		e.addAttr(key, value)
	}
	return e
}

func (e *Element) GetAttr(key string) string {
	return e.attrs[key]
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
	e.addListener(eventName, handler, false)
	return e
}

func (e *Element) OnGlobal(eventName string, handler func(Event)) *Element {
	e.addListener(eventName, handler, true)
	return e
}

func (e *Element) OnInput(handler func(Event)) *Element      { return e.On("input", handler) }
func (e *Element) OnChange(handler func(Event)) *Element     { return e.On("change", handler) }
func (e *Element) OnSubmit(handler func(Event)) *Element     { return e.On("submit", handler) }
func (e *Element) OnKeyDown(handler func(Event)) *Element    { return e.On("keydown", handler) }
func (e *Element) OnKeyUp(handler func(Event)) *Element      { return e.On("keyup", handler) }
func (e *Element) OnMouseEnter(handler func(Event)) *Element { return e.On("mouseenter", handler) }
func (e *Element) OnMouseLeave(handler func(Event)) *Element { return e.On("mouseleave", handler) }
func (e *Element) OnMouseOver(handler func(Event)) *Element  { return e.On("mouseover", handler) }
func (e *Element) OnMouseOut(handler func(Event)) *Element   { return e.On("mouseout", handler) }
func (e *Element) OnClick(handler func(Event)) *Element      { return e.On("click", handler) }
func (e *Element) OnResize(handler func(Event)) *Element     { return e.OnGlobal("resize-el", handler) }

func (e *Element) Empty() bool {
	return len(e.children) == 0
}

func (e *Element) Children(children ...AsNode) *Element {
	newC := intoNodesList(children)
	if len(newC) > 0 {
		e.children = append(e.children, newC...)
	}

	return e
}

func intoNodesList(children []AsNode) []Node {
	var ret []Node

	for _, n := range children {
		if n == nil {
			continue
		}

		cn := n.AsNode()
		if isNilNode(cn) {
			continue
		}

		ret = append(ret, cn)
	}

	return ret
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
			"[Nina] Critical Error in Bind(): expected a pointer to a basic type (*string, *int, *bool), but got %T",
			target,
		))
	}

	return e
}

func (e *Element) Text(text string) *Element {
	e.children = append(e.children, &TextNode{value: text})
	return e
}
func (e *Element) AsNode() Node {
	return e
}

// TextNode for general text
type TextNode struct {
	value string

	// reference to real HTML element in browser
	// this value will be setup after first render
	domNode NativeNode
}

func (t *TextNode) isNode() {}

func (t *TextNode) getKey() string {
	return ""
}

func (t *TextNode) isNil() bool {
	return t == nil
}

func (t *TextNode) AsNode() Node {
	return t
}
