package ui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

type comboboxInternalCtx struct {
	anchorRef      *nn.Ref
	onInput        func(nn.Event)
	onKeyDown      func(nn.Event)
	onToggle       func(nn.Event)
	onAnchorResize func(nn.Event)
	openContent    func()
	closeContent   func()
	deactivate     func(string)
}

// ==========================================
// COMBOBOX TRIGGER
// ==========================================

type comboboxTriggerBuilder struct {
	baseBuilder[*comboboxTriggerBuilder]
}

func ComboboxTrigger() *comboboxTriggerBuilder {
	b := &comboboxTriggerBuilder{}
	b.baseBuilder = base(b, "button")
	b.Attr("data-slot", "combobox-trigger").
		Class("[&_svg:not([class*='size-'])]:size-4")

	return b
}

func (t *comboboxTriggerBuilder) build(ctx *buildContext) {
	ctx.Children = append(
		ctx.Children,
		icons.ChevronDown().Class("pointer-events-none size-4 text-muted-foreground"),
	)
}

// ==========================================
// COMBOBOX CLEAR
// ==========================================

func ComboboxClear() *buttonBuilder {
	return InputGroupButton().Ghost().SizeIconXs().
		Attr("data-slot", "combobox-clear").
		Class("inline-flex items-center justify-center rounded-md hover:bg-accent hover:text-accent-foreground size-6").
		Children(
			icons.X().Class("pointer-events-none size-4"),
		)
}

// ==========================================
// COMBOBOX INPUT
// ==========================================

type comboboxInput struct {
	uiComponent[*comboboxInput]

	disabled    bool
	invalid     bool
	hideTrigger bool
	showClear   bool
	placeholder string
}

func ComboboxInput() *comboboxInput {
	b := &comboboxInput{}
	b.init(b)

	return b
}

func (b *comboboxInput) View() nn.Node {
	ctx := nn.GetContext[*comboboxInternalCtx](b)

	el := InputGroup().
		Class("w-auto").
		Attr("data-slot", "combobox-input-group").
		Ref(ctx.anchorRef).
		On("keydown", ctx.onKeyDown)

	input := InputGroupInput().
		Attr("data-slot", "combobox-input").
		Disabled(b.disabled).
		Placeholder(b.placeholder).
		On("input", ctx.onInput).
		On("focus", ctx.onToggle).
		On("focusout", func(nn.Event) {
			ctx.closeContent()
		})

	if b.invalid {
		input.Attr("aria-invalid", "true")
	}

	addon := InputGroupAddon().AlignInlineEnd()

	if !b.hideTrigger {
		addon.Children(
			InputGroupButton().
				SizeIconXs().
				Ghost().
				Attr("tabindex", "-1").
				Attr("data-slot", "combobox-trigger").
				Class("group-has-data-[slot=combobox-clear]/input-group:hidden data-pressed:bg-transparent").
				Disabled(b.disabled).
				RenderAs(
					ComboboxTrigger(),
				).OnClick(ctx.onToggle),
		)
	}

	if b.showClear {
		addon.Children(
			ComboboxClear().Disabled(b.disabled),
		)
	}

	el.Children(input, addon)
	el.Children(b.children...)

	return b.ApplyProps(el.AsNode())
}

func (b *comboboxInput) Placeholder(p string) *comboboxInput {
	b.placeholder = p
	return b
}
func (b *comboboxInput) Disabled(d bool) *comboboxInput {
	b.disabled = d
	return b
}
func (b *comboboxInput) Invalid(d bool) *comboboxInput {
	b.invalid = d
	return b
}

func (b *comboboxInput) ShowClear(s bool) *comboboxInput {
	b.showClear = s
	return b
}

func (b *comboboxInput) HideTrigger(h bool) *comboboxInput {
	b.hideTrigger = h
	return b
}

// ==========================================
// COMBOBOX GROUP
// ==========================================

func ComboboxGroup(id string) *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-group").
		Key("group-" + id)
}

// ==========================================
// COMBOBOX LABEL
// ==========================================

func ComboboxLabel() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-label").
		Class("px-3 py-2.5 text-xs text-muted-foreground")
}

// ==========================================
// COMBOBOX COLLECTION
// ==========================================

func ComboboxCollection() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-collection")
}

// ==========================================
// COMBOBOX SEPARATOR
// ==========================================

func ComboboxSeparator() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-separator").
		Class("-mx-1.5 my-1.5 h-px bg-border")
}

// ==========================================
// COMBOBOX CONTENT
// ==========================================

func ComboboxContent() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-content").
		Class("cn-menu-target cn-menu-translucent group/combobox-content relative max-h-(--available-height) w-(--anchor-width) max-w-(--available-width) min-w-(--anchor-width) origin-(--transform-origin) overflow-hidden rounded-3xl bg-popover text-popover-foreground shadow-lg ring-1 ring-foreground/5 duration-300 data-[chips=true]:min-w-(--anchor-width) data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 *:data-[slot=input-group]:m-1.5 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:border-input/30 *:data-[slot=input-group]:bg-input/50 *:data-[slot=input-group]:shadow-none dark:ring-foreground/10 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95")
}

// ==========================================
// COMBOBOX LIST
// ==========================================

func ComboboxList() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-list").
		Class("no-scrollbar max-h-[min(calc(--spacing(72)---spacing(9)),calc(var(--available-height)---spacing(9)))] scroll-py-1.5 overflow-y-auto overscroll-contain p-1.5 data-empty:p-0")
}

// ==========================================
// COMBOBOX ITEM
// ==========================================

type comboboxItemBuilder struct {
	baseBuilder[*comboboxItemBuilder]

	selected bool
	value    string
}

func ComboboxItem(value string) *comboboxItemBuilder {
	b := &comboboxItemBuilder{value: value}
	b.baseBuilder = base(b, "div")
	b.Attr("role", "option").
		Attr("data-slot", "combobox-item").
		Attr("data-value", value).
		Class("relative flex w-full cursor-default items-center gap-2.5 rounded-2xl py-2 pr-8 pl-3 text-sm font-medium outline-hidden select-none data-highlighted:bg-accent data-highlighted:text-accent-foreground not-data-[variant=destructive]:data-highlighted:**:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4")

	return b
}

func (b *comboboxItemBuilder) Selected(s bool) *comboboxItemBuilder {
	b.selected = s
	return b
}

func (b *comboboxItemBuilder) build(ctx *buildContext) {
	if b.selected {
		indicator := nn.Span().
			Attr("data-slot", "combobox-item-indicator").
			Class("absolute right-2 flex size-4 items-center justify-center pointer-events-none").
			Children(
				icons.Check().Class("pointer-events-none"),
			)

		ctx.Children = append(ctx.Children, indicator)
	}

	ctx.Props.Key("item-" + b.value)
}

// ==========================================
// COMBOBOX EMPTY
// ==========================================

func ComboboxEmpty() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "combobox-empty").
		Class("w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/combobox-content:flex")
}

// ==========================================
// COMBOBOX CHIPS
// ==========================================

type comboboxChips struct {
	uiComponent[*comboboxChips]
}

func ComboboxChips() *comboboxChips {
	c := &comboboxChips{}
	c.init(c)

	return c
}

func (c *comboboxChips) View() nn.Node {
	ctx := nn.GetContext[*comboboxInternalCtx](c)

	return c.ApplyPropsWithChildren(
		nn.Div().
			Attr("data-slot", "combobox-chips").
			Class("flex min-h-9 flex-wrap items-center gap-1.5 rounded-3xl border border-transparent bg-input/50 bg-clip-padding px-3 py-1.5 text-sm transition-[color,box-shadow,background-color] focus-within:border-ring focus-within:ring-3 focus-within:ring-ring/30 has-aria-invalid:border-destructive has-aria-invalid:ring-3 has-aria-invalid:ring-destructive/20 has-data-[slot=combobox-chip]:px-1.5 dark:has-aria-invalid:border-destructive/50 dark:has-aria-invalid:ring-destructive/40").
			Ref(ctx.anchorRef).
			OnKeyDown(noUpdate(ctx.onKeyDown)).
			OnResize(noUpdate(ctx.onAnchorResize)),
	)
}

// ==========================================
// COMBOBOX CHIPS INPUT
// ==========================================

type comboboxChipInput struct {
	uiComponent[*comboboxChipInput]

	invalid     bool
	placeholder string
}

func ComboboxChipInput() *comboboxChipInput {
	b := &comboboxChipInput{}
	b.init(b)

	return b
}

func (b *comboboxChipInput) Placeholder(p string) *comboboxChipInput {
	b.placeholder = p
	return b
}

func (b *comboboxChipInput) Invalid(d bool) *comboboxChipInput {
	b.invalid = d
	return b
}

func (b *comboboxChipInput) View() nn.Node {
	ctx := nn.GetContext[*comboboxInternalCtx](b)

	el := nn.Input().
		Attr("data-slot", "combobox-chip-input").
		Key("combobox-chip-input").
		Class("min-w-16 flex-1 outline-none").
		Placeholder(b.placeholder).
		OnInput(noUpdate(ctx.onInput)).
		On("focus", noUpdate(ctx.onToggle)).
		On("focusout", noUpdate(func(nn.Event) {
			ctx.closeContent()
		}))

	if b.invalid {
		el.Attr("aria-invalid", "true")
	}

	return b.ApplyPropsWithChildren(el)
}

// ==========================================
// COMBOBOX CHIP
// ==========================================

type comboboxChip struct {
	uiComponent[*comboboxChip]

	value      string
	showRemove bool
}

func ComboboxChip(value string) *comboboxChip {
	b := &comboboxChip{value: value}
	b.init(b)

	return b
}

func (b *comboboxChip) ShowRemove(sr bool) *comboboxChip {
	b.showRemove = sr
	return b
}

func (b *comboboxChip) View() nn.Node {
	ctx := nn.GetContext[*comboboxInternalCtx](b)

	el := nn.Div().
		Attr("data-slot", "combobox-chip").
		Attr("tabindex", "-1").
		Class("flex h-[calc(--spacing(5.5))] w-fit items-center justify-center gap-1 rounded-3xl bg-input px-2 text-xs font-medium whitespace-nowrap text-foreground has-disabled:pointer-events-none has-disabled:cursor-not-allowed has-disabled:opacity-50 has-data-[slot=combobox-chip-remove]:pr-0 dark:bg-input/60").
		Key("chip-"+b.value).
		On("mousedown", func(e nn.Event) {
			e.PreventDefault()
			e.PreventUpdate()
		})

	el.Children(b.children...)

	if b.showRemove {
		btn := Button().
			Ghost().
			SizeIconXs().
			Class("-ml-1 opacity-50 hover:opacity-100 hover:bg-transparent dark:hover:bg-transparent").
			Attr("data-slot", "combobox-chip-close").
			Attr("bind-val", b.value).
			Attr("tabindex", "-1").
			Children(
				icons.X().Class("pointer-events-none"),
			).
			OnClick(func(e nn.Event) {
				e.PreventDefault()
				e.PreventUpdate()
				val := e.Renderer().GetAttribute(e.Target(), "bind-val")

				ctx.deactivate(val)
			})

		el.Children(btn)
	}

	return b.ApplyProps(el)
}

// ==========================================
// COMBOBOX COMPONENT
// ==========================================

type ComboboxConfig[T any] struct {
	Items []T

	MultipleSelect bool

	ItemToString func(item T) string

	ItemToValue func(item T) string
}

type ComboboxContext[T any] struct {
	IsOpen        bool
	VisibleItems  []T
	SelectedItems []T
	IsActive      func(item T) bool
	MenuPosition  *positionerContext[*combobox[T]]
}

type comboboxState[T any] struct {
	searchQuery string
	isOpen      bool

	activeOpts []T
}

type combobox[T any] struct {
	uiComponent[*combobox[T]]
	nn.State[comboboxState[T]]

	config ComboboxConfig[T]

	anchorRef *nn.Ref

	pos *positioner

	renderTrigger func(ctx ComboboxContext[T]) nn.AsNode
	renderContent func(ctx ComboboxContext[T]) nn.AsNode
}

func Combobox[T any](cfg ComboboxConfig[T]) *combobox[T] {
	c := &combobox[T]{
		config:    cfg,
		anchorRef: nn.NewRef(),
	}
	c.init(c)

	c.pos = Positioner(c.anchorRef).
		OnClose(func() {
			c.Data.isOpen = false

			trigger := c.anchorRef.Current
			r := c.anchorRef.Renderer
			if trigger != nil {
				input := r.QuerySelector(trigger, "input")
				if !c.config.MultipleSelect && len(c.Data.activeOpts) > 0 {
					r.SetAttribute(input, "value", c.config.ItemToString(c.Data.activeOpts[0]))
				} else {
					r.SetAttribute(input, "value", "")
				}
			}

			c.Update()
		})

	posCtx := getPositionerContext(c, c.pos)
	posCtx.Flip()
	posCtx.AlignStart()
	posCtx.Offset(8)

	intCtx := &comboboxInternalCtx{
		anchorRef: c.anchorRef,
		onInput: func(e nn.Event) {
			c.Data.searchQuery = e.TargetValue()
			if !c.Data.isOpen {
				c.openContent()
			}
			c.Update()
		},
		deactivate: func(val string) {
			c.Data.activeOpts = slices.DeleteFunc(
				c.Data.activeOpts,
				func(o T) bool { return c.config.ItemToValue(o) == val },
			)
			c.Update()
		},
		onAnchorResize: func(e nn.Event) { c.pos.recalculatePosition(); c.pos.Update() },
		onKeyDown:      c.onKeyDown,
		onToggle:       c.onToggle,
		openContent:    c.openContent,
		closeContent:   c.closeContent,
	}

	nn.ProvideContext(c, intCtx)

	return c
}

func (c *combobox[T]) visibleItems() []T {
	var result []T
	query := strings.ToLower(c.Data.searchQuery)

	for _, item := range c.config.Items {
		itemText := strings.ToLower(c.config.ItemToString(item))
		if strings.Contains(itemText, query) {
			result = append(result, item)
		}
	}

	return result
}

func (c *combobox[T]) Trigger(cb func(ctx ComboboxContext[T]) nn.AsNode) *combobox[T] {
	c.renderTrigger = cb
	return c
}

func (c *combobox[T]) Content(cb func(ctx ComboboxContext[T]) nn.AsNode) *combobox[T] {
	c.renderContent = cb
	return c
}

func (c *combobox[T]) wrapContent(content nn.AsNode) nn.AsNode {
	n := content.AsNode()
	el, ok := n.(*nn.Element)
	if !ok {
		fmt.Println("[ERR] expected *Element as combobox content")
		return nil
	}

	menu := WrapMenu(el, func(val string) { c.activateSelected(val); c.Update() })

	return c.pos.Children(menu)
}

func (c *combobox[T]) View() nn.Node {
	ctx := ComboboxContext[T]{
		IsOpen:        c.Data.isOpen,
		VisibleItems:  c.visibleItems(),
		SelectedItems: c.Data.activeOpts,
		IsActive:      c.isActive,
		MenuPosition:  getPositionerContext(c, c.pos),
	}

	var children []nn.AsNode

	if c.renderTrigger != nil {
		triggerNode := c.renderTrigger(ctx)
		children = append(children, triggerNode)
	}

	if c.Data.isOpen && c.renderContent != nil {
		contentNode := c.renderContent(ctx)
		children = append(children, c.wrapContent(contentNode))
	}

	return nn.Div().Children(children...)
}

func (c *combobox[T]) isActive(v T) bool {
	for _, i := range c.Data.activeOpts {
		if c.config.ItemToValue(v) == c.config.ItemToValue(i) {
			return true
		}
	}
	return false
}

func (c *combobox[T]) onToggle(e nn.Event) {
	if c.Data.isOpen {
		c.closeContent()
	} else {
		c.openContent()
	}
}

func (c *combobox[T]) openContent() {
	c.Data.isOpen = true

	c.Update()
}

func (c *combobox[T]) activateSelected(val string) {
	var (
		newOpt T
		found  bool
	)
	for _, opt := range c.config.Items {
		if c.config.ItemToValue(opt) == val {
			newOpt = opt
			found = true
			break
		}
	}

	if found {
		if c.config.MultipleSelect {
			if c.isOptionActive(newOpt) {
				c.Data.activeOpts = slices.DeleteFunc(
					c.Data.activeOpts,
					func(o T) bool { return c.config.ItemToValue(o) == c.config.ItemToValue(newOpt) },
				)
			} else {
				c.Data.activeOpts = append(c.Data.activeOpts, newOpt)
			}
		} else {
			c.Data.activeOpts = []T{newOpt}
		}
	}

	if !c.config.MultipleSelect {
		c.closeContent()
	}
}

func (c *combobox[T]) isOptionActive(o T) bool {
	return slices.ContainsFunc(
		c.Data.activeOpts,
		func(v T) bool { return c.config.ItemToValue(v) == c.config.ItemToValue(o) },
	)
}

func (c *combobox[T]) onKeyDown(e nn.Event) {
	key := e.Key()

	if c.Data.isOpen == false && key == "ArrowDown" {
		c.openContent()
	}

}

func (c *combobox[T]) closeContent() {
	if !c.Data.isOpen {
		return
	}

	c.pos.Close()
	c.Data.searchQuery = ""
}
