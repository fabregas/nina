package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

// ==========================================
// PROP STORE
// ==========================================

type PropStore struct {
	text    string
	classes string
	attrs   map[string]string
	events  map[string]func(nn.Event)
}

func (p *PropStore) Merge(other *PropStore) {
	if other.classes != "" {
		p.classes = p.classes + " " + other.classes
	}

	for k, v := range other.attrs {
		if p.attrs == nil {
			p.attrs = make(map[string]string)
		}
		p.attrs[k] = v
	}

	for evt, handler := range other.events {
		if p.events == nil {
			p.events = make(map[string]func(nn.Event))
		}
		p.events[evt] = handler
	}

	if other.text != "" {
		p.text = other.text
	}
}

func (p *PropStore) ApplyTo(el *nn.Element) *nn.Element {
	if p.classes != "" {
		el.Class(p.classes)
	}
	for k, v := range p.attrs {
		el.Attr(k, v)
	}
	for evt, handler := range p.events {
		el.On(evt, handler)
	}

	if p.text != "" {
		el.Text(p.text)
	}

	return el
}

type UIProps[T any] struct {
	instance T

	props         PropStore
	children      []nn.AsNode
	childOverride *nn.Element // optional for AsChild()
	isInit        bool
}

func (p *UIProps[T]) init(instance T) {
	p.instance = instance
	p.isInit = true
}

func (p *UIProps[T]) Class(c string) T {
	p.assertInit()
	if p.props.classes != "" {
		p.props.classes += " "
	}
	p.props.classes += c

	return p.instance
}

func (p *UIProps[T]) Attr(k, v string) T {
	p.assertInit()
	if p.props.attrs == nil {
		p.props.attrs = make(map[string]string)
	}
	p.props.attrs[k] = v

	return p.instance
}

func (p *UIProps[T]) On(event string, handler func(nn.Event)) T {
	p.assertInit()
	if p.props.events == nil {
		p.props.events = make(map[string]func(nn.Event))
	}
	p.props.events[event] = handler

	return p.instance
}

func (p *UIProps[T]) OnClick(handler func(nn.Event)) T {
	return p.On("click", handler)
}

func (p *UIProps[T]) Text(t string) T {
	p.assertInit()
	p.props.text = t

	return p.instance
}

func (p *UIProps[T]) assertInit() {
	if !p.isInit {
		panic(fmt.Sprintf("%T is not initialized\n", p.instance))
	}
}

func (p *UIProps[T]) Children(items ...nn.AsNode) T {
	p.assertInit()
	p.children = append(p.children, items...)

	return p.instance
}

func (p *UIProps[T]) ID(id string) T {
	p.Attr("id", id)

	return p.instance
}

func (p *UIProps[T]) For(id string) T {
	p.Attr("for", id)

	return p.instance
}

func (p *UIProps[T]) Disabled(disabled bool) T {
	if disabled {
		p.Attr("disabled", "true")
	}

	return p.instance
}

func (p *UIProps[T]) AsChild(child *nn.Element) T {
	p.childOverride = child

	return p.instance
}

// ==========================================
// baseBuilder
// ==========================================

type builder interface {
	build() *nn.Element
}

type baseBuilder[T builder] struct {
	UIProps[T]

	el    *nn.Element
	built bool
}

func (b *baseBuilder[T]) AsNode() nn.Node {
	return b.El()
}

func (b *baseBuilder[T]) El() *nn.Element {
	if b.built {
		//fmt.Printf("ALREADY BUILT: %T\n", b.instance)
		return b.el
	}

	b.el.Children(b.children...)
	b.el = b.props.ApplyTo(b.el)

	b.el = b.instance.build()

	if b.childOverride != nil {
		b.el = b.childOverride.MergeEl(b.el)
	}

	b.built = true

	return b.el
}

func base[T builder](self T, el *nn.Element) baseBuilder[T] {
	b := baseBuilder[T]{
		el: el,
	}

	b.init(self)

	return b
}

// ==========================================
// UIComponent
// ==========================================

type uiComponent[T any] struct {
	nn.BaseComponent
	UIProps[T]
}

func (c *uiComponent[T]) ApplyProps(p *nn.Element) *nn.Element {
	p.Children(c.children...)
	c.props.ApplyTo(p)
	if c.childOverride != nil {
		p = c.childOverride.MergeEl(p)
	}

	return p
}

type simpleBuilder struct {
	baseBuilder[*simpleBuilder]
}

func (s *simpleBuilder) build() *nn.Element {
	return s.el
}

func simple(el *nn.Element) *simpleBuilder {
	b := &simpleBuilder{}
	b.baseBuilder = baseBuilder[*simpleBuilder]{
		el: el,
	}
	b.init(b)

	return b
}

func noUpdate(cb func(nn.Event)) func(nn.Event) {
	return func(e nn.Event) {
		e.PreventUpdate()
		cb(e)
	}
}
