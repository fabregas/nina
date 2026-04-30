package ui

import (
	"github.com/fabregas/nina/nn"
)

type collapsibleState struct {
	isOpen *nn.Signal[bool]
}

// ==========================================
// COLLAPSIBLE TRIGGER
// ==========================================

type collapsibleTrigger struct {
	uiComponent[*collapsibleTrigger]
}

func CollapsibleTrigger() *collapsibleTrigger {
	t := &collapsibleTrigger{}
	t.init(t)

	return t
}

func (t *collapsibleTrigger) View() nn.Node {
	state := nn.GetContext[*collapsibleState](t)

	el := nn.Div().
		Attr("data-slot", "collapsible-trigger").
		OnClick(func(e nn.Event) {
			cur := state.isOpen.Get(nil)
			state.isOpen.Set(!cur)
			e.PreventUpdate()
		})

	return t.ApplyProps(el)
}

// ==========================================
// COLLAPSIBLE CONTENT
// ==========================================

type collapsibleContent struct {
	uiComponent[*collapsibleContent]
}

func CollapsibleContent() *collapsibleContent {
	c := &collapsibleContent{}
	c.init(c)

	return c
}

func (c *collapsibleContent) View() nn.Node {
	state := nn.GetContext[*collapsibleState](c)

	isOpen := state.isOpen.Get(c)

	el := nn.Div().Attr("data-slot", "collapsible-content")
	c.props.ApplyTo(el)
	if isOpen {
		el.Children(c.children...)
	}

	return el
}

// ==========================================
// COLLAPSIBLE
// ==========================================

type collapsible struct {
	uiComponent[*collapsible]
	nn.State[collapsibleState]
}

func Collapsible() *collapsible {
	c := &collapsible{}
	c.init(c)

	state := &collapsibleState{
		isOpen: nn.NewSignal(false),
	}
	c.Data = state

	nn.ProvideContext(c, state)

	return c
}

func (c *collapsible) View() nn.Node {
	return c.ApplyProps(
		nn.Div().
			Attr("data-slot", "collapsible"),
	)
}
