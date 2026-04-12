package ui

import (
	"fmt"
	"syscall/js"

	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

// ==========================================
// SELECT GROUP
// ==========================================

func SelectGroup() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "select-group").
			Attr("role", "group").
			Class("scroll-my-1.5 p-1.5"),
	)
}

// ==========================================
// SELECT VALUE
// ==========================================

func SelectValue() *simpleBuilder {
	return simple(
		nn.Span().
			Attr("data-slot", "select-value").
			Class("flex flex-1 text-left"),
	)
}

// ==========================================
// SELECT TRIGGER
// ==========================================

type selectTriggerBuilder struct {
	baseBuilder[*selectTriggerBuilder]

	onClick func(nn.Event)
}

func SelectTrigger() *selectTriggerBuilder {
	baseClass := "flex w-fit items-center justify-between gap-1.5 rounded-3xl border border-transparent bg-input/50 px-3 py-2 text-sm whitespace-nowrap transition-[color,box-shadow,background-color] outline-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 data-placeholder:text-muted-foreground data-[size=default]:h-9 data-[size=sm]:h-8 *:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-1.5 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

	btn := nn.Button().
		Attr("type", "button").
		Attr("role", "combobox").
		Attr("data-slot", "select-trigger").
		Attr("data-size", "default").
		Class(baseClass)

	b := &selectTriggerBuilder{}
	b.baseBuilder = base(b, btn)

	return b
}

func (t *selectTriggerBuilder) build() *nn.Element {
	// icon must be last child
	return t.el.Children(icons.ChevronDown())
}

func (t *selectTriggerBuilder) OnClick(fn func(nn.Event)) *selectTriggerBuilder {
	t.el.OnClick(fn)
	return t
}

func (t *selectTriggerBuilder) SizeDefault() *selectTriggerBuilder {
	t.el.Attr("data-size", "default")
	return t
}

func (t *selectTriggerBuilder) SizeSm() *selectTriggerBuilder {
	t.el.Attr("data-size", "sm")
	return t
}

// ==========================================
// SELECT CONTENT
// ==========================================

type selectContentBuilder struct {
	baseBuilder[*selectContentBuilder]
}

func SelectContent() *selectContentBuilder {
	baseClass := "cn-menu-target cn-menu-translucent relative isolate z-50 max-h-(--available-height) w-(--anchor-width) min-w-36 origin-(--transform-origin) overflow-x-hidden overflow-y-auto rounded-3xl bg-popover text-popover-foreground shadow-lg ring-1 ring-foreground/5 duration-100 data-[align-trigger=true]:animate-none data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 dark:ring-foreground/10 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"

	el := nn.Div().
		Attr("data-slot", "select-content").
		Class(baseClass)

	b := &selectContentBuilder{}
	b.baseBuilder = base(b, el)

	b.SideLeft().AlignTrigger(true)

	return b
}

func (c *selectContentBuilder) SideTop() *selectContentBuilder {
	c.el.Attr("data-side", "top")
	return c
}
func (c *selectContentBuilder) SideBottom() *selectContentBuilder {
	c.el.Attr("data-side", "bottom")
	return c
}
func (c *selectContentBuilder) SideLeft() *selectContentBuilder {
	c.el.Attr("data-side", "left")
	return c
}
func (c *selectContentBuilder) SideRight() *selectContentBuilder {
	c.el.Attr("data-side", "right")
	return c
}

func (c *selectContentBuilder) AlignTrigger(align bool) *selectContentBuilder {
	if align {
		c.el.Attr("data-align-trigger", "true")
	} else {
		c.el.Attr("data-align-trigger", "false")
	}
	return c
}

func (c *selectContentBuilder) build() *nn.Element {
	return c.el
}

// --- SELECT LABEL ---

func SelectLabel() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "select-label").
			Class("px-3 py-2.5 text-xs text-muted-foreground"),
	)
}

// ==========================================
// SELECT ITEM
// ==========================================

type selectItemBuilder struct {
	baseBuilder[*selectItemBuilder]

	onClick    func(nn.Event)
	isSelected bool
	isDisabled bool
}

func SelectItem(value string) *selectItemBuilder {
	item := nn.Span().
		Class("flex flex-1 shrink-0 gap-2 whitespace-nowrap")

	b := &selectItemBuilder{}
	b.baseBuilder = base(b, item)

	return b
}

func (b *selectItemBuilder) OnClick(fn func(nn.Event)) *selectItemBuilder {
	b.onClick = fn
	return b
}
func (i *selectItemBuilder) Selected(selected bool) *selectItemBuilder {
	i.isSelected = selected
	return i
}

func (i *selectItemBuilder) Disabled(disabled bool) *selectItemBuilder {
	i.isDisabled = disabled
	return i
}

func (i *selectItemBuilder) build() *nn.Element {
	return i.el
}

func (i *selectItemBuilder) wrap(target *nn.Element) *nn.Element {
	baseClass := "relative flex w-full cursor-default items-center gap-2.5 rounded-2xl py-2 pr-8 pl-3 text-sm font-medium outline-hidden select-none hover:bg-accent hover:text-accent-foreground not-data-[variant=destructive]:hover:**:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 *:[span]:last:flex *:[span]:last:items-center *:[span]:last:gap-2"

	wrapper := nn.Div().
		Attr("role", "option").
		Attr("data-slot", "select-item").
		Class(baseClass).
		Children(target)

	if i.isDisabled {
		wrapper.Attr("data-disabled", "true")
		wrapper.Attr("aria-disabled", "true")
	}

	if i.isSelected {
		wrapper.Attr("data-state", "checked")
		wrapper.Attr("aria-selected", "true")

		indicator := nn.Span().
			Class("pointer-events-none absolute right-2 flex size-4 items-center justify-center").
			Children(icons.Check())

		wrapper.Children(indicator)
	} else {
		wrapper.Attr("data-state", "unchecked")
		wrapper.Attr("aria-selected", "false")
	}

	if i.onClick != nil {
		wrapper.OnClick(i.onClick)
	}

	return wrapper
}

// ==========================================
// SELECT SEPARATOR
// ==========================================

func SelectSeparator() *simpleBuilder {
	return simple(nn.Div().
		Attr("data-slot", "select-separator").
		Attr("role", "separator").
		Class("pointer-events-none -mx-1.5 my-1.5 h-px bg-border"),
	)
}

// ==========================================
// SELECT SCROLL UP BUTTON
// ==========================================

func SelectScrollUpButton() *simpleBuilder {
	baseClass := "top-0 z-10 flex w-full cursor-default items-center justify-center bg-popover py-1 [&_svg:not([class*='size-'])]:size-4"

	return simple(
		nn.Div().
			Attr("data-slot", "select-scroll-up-button").
			Attr("aria-hidden", "true").
			Children(icons.ChevronUp()).
			Class(baseClass),
	)
}

// ==========================================
// SELECT SCROLL DOWN BUTTON
// ==========================================

func SelectScrollDownButton() *simpleBuilder {
	baseClass := "bottom-0 z-10 flex w-full cursor-default items-center justify-center bg-popover py-1 [&_svg:not([class*='size-'])]:size-4"

	return simple(
		nn.Div().
			Attr("data-slot", "select-scroll-down-button").
			Attr("aria-hidden", "true").
			Class(baseClass).
			Children(icons.ChevronDown()),
	)
}

// ==========================================
// SELECT CONTROLLER
// ==========================================

type SelectOption struct {
	Label string
	Value string
}

type SelectState struct {
	IsOpen   bool
	Value    string
	OnChange func(string)

	top   float64
	left  float64
	width float64
}

func SelectController(id string, state *SelectState, options []SelectOption, placeholder string) nn.IntoNode {
	if state == nil {
		state = &SelectState{}
	}

	triggerID := id
	menuID := id + "-menu"

	displayValue := placeholder
	for _, opt := range options {
		if opt.Value == state.Value {
			displayValue = opt.Label
			break
		}
	}

	globalClick := func(e nn.Event) {
		target := e.Target()

		doc := js.Global().Get("document")
		triggerEl := doc.Call("getElementById", triggerID)
		menuEl := doc.Call("getElementById", menuID)

		clickedInTrigger := !triggerEl.IsNull() && triggerEl.Call("contains", target).Bool()
		clickedInMenu := !menuEl.IsNull() && menuEl.Call("contains", target).Bool()

		if !clickedInTrigger && !clickedInMenu {
			state.IsOpen = false
		}
	}

	toggleOpen := func(e nn.Event) {
		state.IsOpen = !state.IsOpen

		if state.IsOpen {
			btn := e.CurrentTarget()
			if !btn.IsNull() {
				rect := btn.Call("getBoundingClientRect")

				window := js.Global().Get("window")
				scrollY := window.Get("scrollY").Float()
				scrollX := window.Get("scrollX").Float()

				state.top = rect.Get("bottom").Float() + scrollY + 4
				state.left = rect.Get("left").Float() + scrollX
				state.width = rect.Get("width").Float()
			}
		}
	}

	makeSelectHandler := func(val string) func(nn.Event) {
		return func(e nn.Event) {
			state.Value = val
			state.IsOpen = false
			if state.OnChange != nil {
				state.OnChange(val)
			}
		}
	}

	trigger := SelectTrigger().
		ID(triggerID).
		OnClick(toggleOpen).
		Children(
			SelectValue().Children(nn.Text(displayValue)),
		).El().OnGlobal("click", globalClick)

	var content nn.IntoNode
	if state.IsOpen {
		var items []nn.IntoNode
		for _, opt := range options {
			item := SelectItem(opt.Value).
				Selected(state.Value == opt.Value).
				OnClick(makeSelectHandler(opt.Value)).
				Children(nn.Text(opt.Label))

			items = append(items, item)
		}

		styleStr := fmt.Sprintf("position: absolute; top: %.1fpx; left: %.1fpx; width: %.1fpx; z-index: 50;",
			state.top, state.left, state.width)
		content = nn.Portal(
			SelectContent().
				AlignTrigger(true).
				Class("min-w-[8rem]").
				Children(
					SelectGroup().Children(items...),
				).El().Attr("style", styleStr).Attr("id", menuID),
		)
	}

	return nn.Div().
		Class("relative inline-block w-full").
		Children(
			trigger.Class("w-full"),
			content,
		)
}
