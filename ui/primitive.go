package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

type UIProps[T any] struct {
	instance T

	props         nn.Props
	children      []nn.AsNode
	childOverride nn.AsNode // optional for AsChild()

	isInit bool
}

func (p *UIProps[T]) init(instance T) {
	p.instance = instance
	p.isInit = true
}

func (p *UIProps[T]) MergeProps(other *nn.Props) {
	p.props.Merge(other)
}

func (p *UIProps[T]) Class(c string) T {
	p.assertInit()
	p.props.Class(c)

	return p.instance
}

func (p *UIProps[T]) Style(style string) T {
	p.assertInit()
	p.props.Style(style)

	return p.instance
}

func (p *UIProps[T]) Attr(k, v string) T {
	p.assertInit()
	p.props.Attr(k, v)

	return p.instance
}

func (p *UIProps[T]) Key(key string) T {
	p.assertInit()
	p.props.Key(key)

	return p.instance
}

func (p *UIProps[T]) Ref(ref *nn.Ref) T {
	p.assertInit()
	p.props.Ref(ref)

	return p.instance
}

func (p *UIProps[T]) On(event string, handler func(nn.Event)) T {
	p.assertInit()
	p.props.On(event, handler)

	return p.instance
}

func (p *UIProps[T]) OnGlobal(event string, handler func(nn.Event)) T {
	p.assertInit()
	p.props.OnGlobal(event, handler)

	return p.instance
}

func (p *UIProps[T]) OnClick(handler func(nn.Event)) T {
	return p.On("click", handler)
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

func (p *UIProps[T]) Text(t string) T {
	return p.Children(nn.Text(t))
}

func (p *UIProps[T]) AsChild(child nn.AsNode) T {
	p.assertInit()
	p.childOverride = child

	return p.instance
}

// ==========================================
// baseBuilder
// ==========================================

type buildContext struct {
	Props    *nn.Props
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
	ctx := &buildContext{
		Props:    b.props.Clone(),
		Children: append([]nn.AsNode(nil), b.children...),
	}

	b.instance.build(ctx)

	var coreNode nn.Node

	if b.childOverride != nil {
		if receiver, ok := b.childOverride.(interface{ MergeProps(*nn.Props) }); ok {
			receiver.MergeProps(ctx.Props)
			coreNode = b.childOverride.AsNode()
		} else {
			node := b.childOverride.AsNode()
			if el, ok := node.(*nn.Element); ok {
				ctx.Props.ApplyTo(el)
				coreNode = el
			}
		}
	} else {
		freshEl := nn.Tag(b.tag)
		ctx.Props.ApplyTo(freshEl)
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
		var internalProps nn.Props
		if el, ok := defaultRoot.(*nn.Element); ok {
			internalProps.MergeFromElement(el)
		}

		internalProps.Merge(customerProps)
		if receiver, ok := c.childOverride.(interface{ MergeProps(*nn.Props) }); ok {
			receiver.MergeProps(&internalProps)
		}

		return c.childOverride.AsNode()
	}

	if el, ok := defaultRoot.(*nn.Element); ok {
		customerProps.ApplyTo(el)
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
