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

	return t.ApplyPropsWithChildren(el)
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

	if isOpen {
		el.Children(c.children...)
	}

	return c.ApplyProps(el)
}

// ==========================================
// COLLAPSIBLE
// ==========================================

type collapsible struct {
	uiComponent[*collapsible]
	nn.State[collapsibleState]

	initOpen bool
}

func Collapsible() *collapsible {
	c := &collapsible{}
	c.init(c)

	c.InitState(func() *collapsibleState {
		return &collapsibleState{
			isOpen: nn.NewSignal(false),
		}
	})

	nn.ProvideContextDefer(c, func() *collapsibleState {
		return c.Data
	})

	return c
}

func (c *collapsible) IsOpen(isOpen bool) *collapsible {
	c.initOpen = isOpen
	return c
}

func (c *collapsible) OnMount() {
	if c.initOpen {
		c.Data.isOpen.Set(true)
	}
}

func (c *collapsible) View() nn.Node {
	isOpen := c.Data.isOpen.Get(c)
	stateStr := "closed"
	if isOpen {
		stateStr = "open"
	}

	return c.ApplyProps(
		nn.Div().
			Attr("data-slot", "collapsible").
			Attr("data-state", stateStr).
			Class("group/collapsible"),
	)
}
