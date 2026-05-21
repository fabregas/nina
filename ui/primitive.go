package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

type UIProps[T any] struct {
	instance T

	props         nn.Element
	children      []nn.AsNode
	childOverride nn.AsNode // optional for RenderAs()

	isInit      bool
	isRendering bool
}

func (p *UIProps[T]) init(instance T) {
	p.instance = instance
	p.isInit = true
}

func (p *UIProps[T]) MergeProps(other *nn.Element) {
	p.props.Merge(other)
}

func (p *UIProps[T]) Class(c string) T {
	p.assert()
	p.props.Class(c)

	return p.instance
}

func (p *UIProps[T]) Style(style string) T {
	p.assert()
	p.props.Style(style)

	return p.instance
}

func (p *UIProps[T]) Attr(k, v string) T {
	p.assert()
	p.props.Attr(k, v)

	return p.instance
}

func (p *UIProps[T]) AttrIf(condition bool, k, v string) T {
	p.assert()
	if condition {
		p.props.Attr(k, v)
	}

	return p.instance
}

func (p *UIProps[T]) Key(key string) T {
	p.assert()
	p.props.Key(key)

	return p.instance
}

func (p *UIProps[T]) Ref(ref *nn.Ref) T {
	p.assert()
	p.props.Ref(ref)

	return p.instance
}

func (p *UIProps[T]) On(event string, handler func(nn.Event)) T {
	p.assert()
	p.props.On(event, handler)

	return p.instance
}

func (p *UIProps[T]) OnGlobal(event string, handler func(nn.Event)) T {
	p.assert()
	p.props.OnGlobal(event, handler)

	return p.instance
}

func (p *UIProps[T]) OnClick(handler func(nn.Event)) T {
	return p.On("click", handler)
}

func (p *UIProps[T]) assert() {
	if !p.isInit {
		panic(fmt.Sprintf("%T is not initialized\n", p.instance))
	}

	if p.isRendering {
		panic("DX Error: you try to mutate builder on rendering")
	}
}

func (p *UIProps[T]) Children(items ...nn.AsNode) T {
	p.assert()
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

func (p *UIProps[T]) Text(t string) T {
	return p.Children(nn.Text(t))
}

func (p *UIProps[T]) RenderAs(child nn.AsNode) T {
	p.assert()
	p.childOverride = child

	return p.instance
}

// ==========================================
// baseBuilder
// ==========================================

type buildContext struct {
	Props    *nn.Element
	Children []nn.AsNode
}

type builder interface {
	build(*buildContext)
}

type baseBuilder[T builder] struct {
	UIProps[T]

	tag string
}

func (b *baseBuilder[T]) AsNode() nn.Node {
	b.isRendering = true
	defer func() { b.isRendering = false }()

	ctx := &buildContext{
		Props:    b.props.Clone(),
		Children: append([]nn.AsNode(nil), b.children...),
	}

	b.instance.build(ctx)

	var coreNode nn.Node

	if b.childOverride != nil {
		if receiver, ok := b.childOverride.(interface{ MergeProps(*nn.Element) }); ok {
			receiver.MergeProps(ctx.Props)
			coreNode = b.childOverride.AsNode()
		} else {
			node := b.childOverride.AsNode()
			if el, ok := node.(*nn.Element); ok {
				el.Merge(ctx.Props)
				coreNode = el
			}
		}
	} else {
		freshEl := nn.Tag(b.tag)
		freshEl.Merge(ctx.Props)
		freshEl.Children(ctx.Children...)

		coreNode = freshEl
	}

	if nw, ok := any(b.instance).(interface{ wrap(nn.Node) nn.Node }); ok {
		return nw.wrap(coreNode)
	}

	return coreNode
}

func base[T builder](self T, tag string) baseBuilder[T] {
	b := baseBuilder[T]{
		tag: tag,
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

func (c *uiComponent[T]) AsNode() nn.Node {
	if comp, ok := any(c.instance).(nn.Component); ok {
		return nn.Comp(comp)
	}

	panic(fmt.Sprintf("Method View() must be implemented for %T!", c.instance))
}

func (c *uiComponent[T]) ApplyProps(defaultRoot nn.Node) nn.Node {
	customerProps := c.props.Clone()

	if c.childOverride != nil {
		internalProps := &nn.Element{}
		if el, ok := defaultRoot.(*nn.Element); ok {
			internalProps = el.Clone()
		}

		internalProps.Merge(customerProps)
		if receiver, ok := c.childOverride.(interface{ MergeProps(*nn.Element) }); ok {
			receiver.MergeProps(internalProps)
		} else if el, ok := c.childOverride.(*nn.Element); ok {
			el.Merge(internalProps)
		}

		return c.childOverride.AsNode()
	}

	if el, ok := defaultRoot.(*nn.Element); ok {
		el.Merge(customerProps)
		return el
	}

	return defaultRoot
}

func (c *uiComponent[T]) ApplyPropsWithChildren(defaultRoot nn.Node) nn.Node {
	dest := c.ApplyProps(defaultRoot)

	if el, ok := defaultRoot.(*nn.Element); ok {
		el.Children(c.children...)
	}

	return dest
}

type simpleBuilder struct {
	baseBuilder[*simpleBuilder]
}

func (s *simpleBuilder) build(*buildContext) {}

func simple(tag string) *simpleBuilder {
	b := &simpleBuilder{}
	b.baseBuilder = baseBuilder[*simpleBuilder]{
		tag: tag,
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
