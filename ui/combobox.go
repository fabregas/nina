package ui

import (
	"fmt"
	"slices"
	"strings"
	"syscall/js"

	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

// ==========================================
// COMBOBOX TRIGGER
// ==========================================

type comboboxTriggerBuilder struct {
	baseBuilder[*comboboxTriggerBuilder]
}

func ComboboxTrigger() *comboboxTriggerBuilder {
	el := nn.Button().
		Attr("data-slot", "combobox-trigger").
		Class("[&_svg:not([class*='size-'])]:size-4")

	b := &comboboxTriggerBuilder{}
	b.baseBuilder = base(b, el)

	return b
}

func (t *comboboxTriggerBuilder) build() *nn.Element {
	t.Children(icons.ChevronDown().Class("pointer-events-none size-4 text-muted-foreground"))

	return t.el
}

// ==========================================
// COMBOBOX CLEAR
// ==========================================

func ComboboxClear() *simpleBuilder {
	return simple(
		InputGroupButton().Ghost().SizeIconXs().
			Attr("data-slot", "combobox-clear").
			Class("inline-flex items-center justify-center rounded-md hover:bg-accent hover:text-accent-foreground size-6").
			Children(
				icons.X().Class("pointer-events-none size-4"),
			).El(),
	)
}

// ==========================================
// COMBOBOX INPUT
// ==========================================

type comboboxInputBuilder struct {
	baseBuilder[*comboboxInputBuilder]

	disabled    bool
	invalid     bool
	hideTrigger bool
	showClear   bool
	placeholder string
}

func ComboboxInput() *comboboxInputBuilder {
	el := InputGroup().
		Class("w-auto").
		Attr("data-slot", "combobox-input-group").El()

	b := &comboboxInputBuilder{}
	b.baseBuilder = base(b, el)

	return b
}

func (b *comboboxInputBuilder) Placeholder(p string) *comboboxInputBuilder {
	b.placeholder = p
	return b
}
func (b *comboboxInputBuilder) Disabled(d bool) *comboboxInputBuilder {
	b.disabled = d
	return b
}
func (b *comboboxInputBuilder) Invalid(d bool) *comboboxInputBuilder {
	b.invalid = d
	return b
}

func (b *comboboxInputBuilder) ShowClear(s bool) *comboboxInputBuilder {
	b.showClear = s
	return b
}

func (b *comboboxInputBuilder) HideTrigger(h bool) *comboboxInputBuilder {
	b.hideTrigger = h
	return b
}

func (b *comboboxInputBuilder) build() *nn.Element {
	input := InputGroupInput().
		Attr("data-slot", "combobox-input").
		Disabled(b.disabled).
		Placeholder(b.placeholder)

	if b.invalid {
		input.Attr("aria-invalid", "true")
	}

	addon := InputGroupAddon().AlignInlineEnd()

	if !b.hideTrigger {
		addon.Children(
			InputGroupButton().
				SizeIconXs().
				Ghost().
				Attr("data-slot", "combobox-trigger").
				Attr("tabindex", "-1").
				Class("group-has-data-[slot=combobox-clear]/input-group:hidden data-pressed:bg-transparent").
				Disabled(b.disabled).
				AsChild(
					ComboboxTrigger().El(),
				),
		)
	}

	if b.showClear {
		addon.Children(
			ComboboxClear().Disabled(b.disabled),
		)
	}

	return b.el.Children(input, addon)
}

// ==========================================
// COMBOBOX GROUP
// ==========================================

func ComboboxGroup(id string) *simpleBuilder {
	return simple(
		nn.Div().Attr("data-slot", "combobox-group").Key("group-" + id),
	)
}

// ==========================================
// COMBOBOX LABEL
// ==========================================

func ComboboxLabel() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-label").
			Class("px-3 py-2.5 text-xs text-muted-foreground"),
	)
}

// ==========================================
// COMBOBOX COLLECTION
// ==========================================

func ComboboxCollection() *simpleBuilder {
	return simple(
		nn.Div().Attr("data-slot", "combobox-collection"),
	)
}

// ==========================================
// COMBOBOX SEPARATOR
// ==========================================

func ComboboxSeparator() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-separator").
			Class("-mx-1.5 my-1.5 h-px bg-border"),
	)
}

// ==========================================
// COMBOBOX CONTENT
// ==========================================

type comboboxContentPosition struct {
	availableHeight float64
	availableWidth  float64
	anchorWidth     float64
	anchorHeight    float64
	top             float64
	left            float64
}

type comboboxContentPositioner struct {
	el  *nn.Element
	ref *nn.Ref
	pos comboboxContentPosition
}

func (b *comboboxContentPositioner) build() *nn.Element {
	onContentHover := func(e nn.Event) {
		e.PreventUpdate()

		target := e.Target()

		item := target.Call("closest", "[data-slot='combobox-item']")
		if item.IsNull() || item.IsUndefined() {
			return
		}

		container := e.CurrentTarget()

		current := container.Call("querySelector", "[data-highlighted]")
		if !current.IsNull() && !current.Equal(item) {
			current.Call("removeAttribute", "data-highlighted")
		}

		if !item.Call("hasAttribute", "data-highlighted").Bool() {
			item.Call("setAttribute", "data-highlighted", "")
		}
	}

	onMouseDown := func(e nn.Event) {
		e.PreventDefault()
		e.PreventUpdate()
	}

	onMouseLeave := func(e nn.Event) {
		e.PreventUpdate()

		container := e.CurrentTarget()

		current := container.Call("querySelector", "[data-highlighted]")
		if !current.IsNull() {
			current.Call("removeAttribute", "data-highlighted")
		}
	}

	el := b.el.Ref(b.ref).
		On("pointermove", onContentHover).
		On("mousedown", onMouseDown).
		On("mouseleave", onMouseLeave)

	positioner := nn.Div().
		Class("isolate z-50 fixed").
		Attr("style", fmt.Sprintf("top: 0px; left: 0px; transform: translate(%.1fpx, %.1fpx); will-change: transform; --available-width: %.1fpx; --available-height: %.1fpx; --anchor-width: %.1fpx; --anchor-height: %.1fpx; --transform-origin: 146px -6px;", b.pos.left, b.pos.top, b.pos.availableWidth, b.pos.availableHeight, b.pos.anchorWidth, b.pos.anchorHeight)).
		Children(el)

	portal := nn.Portal(positioner)

	return nn.Span().Children(portal)
}

func ComboboxContent() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-content").
			Class("cn-menu-target cn-menu-translucent group/combobox-content relative max-h-(--available-height) w-(--anchor-width) max-w-(--available-width) min-w-(--anchor-width) origin-(--transform-origin) overflow-hidden rounded-3xl bg-popover text-popover-foreground shadow-lg ring-1 ring-foreground/5 duration-300 data-[chips=true]:min-w-(--anchor-width) data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 *:data-[slot=input-group]:m-1.5 *:data-[slot=input-group]:mb-0 *:data-[slot=input-group]:h-8 *:data-[slot=input-group]:border-input/30 *:data-[slot=input-group]:bg-input/50 *:data-[slot=input-group]:shadow-none dark:ring-foreground/10 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"),
	)
}

// ==========================================
// COMBOBOX LIST
// ==========================================

func ComboboxList() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-list").
			Class("no-scrollbar max-h-[min(calc(--spacing(72)---spacing(9)),calc(var(--available-height)---spacing(9)))] scroll-py-1.5 overflow-y-auto overscroll-contain p-1.5 data-empty:p-0"),
	)
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
	el := nn.Div().
		Attr("role", "option").
		Attr("data-slot", "combobox-item").
		Attr("value", value).
		Class("relative flex w-full cursor-default items-center gap-2.5 rounded-2xl py-2 pr-8 pl-3 text-sm font-medium outline-hidden select-none data-highlighted:bg-accent data-highlighted:text-accent-foreground not-data-[variant=destructive]:data-highlighted:**:text-accent-foreground data-disabled:pointer-events-none data-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4")

	b := &comboboxItemBuilder{value: value}
	b.baseBuilder = base(b, el)

	return b
}

func (b *comboboxItemBuilder) Selected(s bool) *comboboxItemBuilder {
	b.selected = s
	return b
}

func (b *comboboxItemBuilder) build() *nn.Element {
	if b.selected {
		indicator := nn.Span().
			Attr("data-slot", "combobox-item-indicator").
			Class("absolute right-2 flex size-4 items-center justify-center pointer-events-none").
			Children(
				icons.Check().Class("pointer-events-none"),
			)

		b.el.Children(indicator)
	}

	return b.el.Key("item-" + b.value)
}

// ==========================================
// COMBOBOX EMPTY
// ==========================================

func ComboboxEmpty() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-empty").
			Class("w-full justify-center py-2 text-center text-sm text-muted-foreground group-data-empty/combobox-content:flex"),
	)
}

// ==========================================
// COMBOBOX CHIPS
// ==========================================

func ComboboxChips() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "combobox-chips").
			Class("flex min-h-9 flex-wrap items-center gap-1.5 rounded-3xl border border-transparent bg-input/50 bg-clip-padding px-3 py-1.5 text-sm transition-[color,box-shadow,background-color] focus-within:border-ring focus-within:ring-3 focus-within:ring-ring/30 has-aria-invalid:border-destructive has-aria-invalid:ring-3 has-aria-invalid:ring-destructive/20 has-data-[slot=combobox-chip]:px-1.5 dark:has-aria-invalid:border-destructive/50 dark:has-aria-invalid:ring-destructive/40"),
	)
}

// ==========================================
// COMBOBOX CHIPS INPUT
// ==========================================

type comboboxChipInputBuilder struct {
	baseBuilder[*comboboxChipInputBuilder]

	invalid     bool
	placeholder string
}

func ComboboxChipInput() *comboboxChipInputBuilder {
	el := nn.Input().
		Attr("data-slot", "combobox-chip-input").
		Key("combobox-chip-input").
		Class("min-w-16 flex-1 outline-none")

	b := &comboboxChipInputBuilder{}
	b.baseBuilder = base(b, el)

	return b
}

func (b *comboboxChipInputBuilder) Placeholder(p string) *comboboxChipInputBuilder {
	b.placeholder = p
	return b
}

func (b *comboboxChipInputBuilder) Invalid(d bool) *comboboxChipInputBuilder {
	b.invalid = d
	return b
}

func (b *comboboxChipInputBuilder) build() *nn.Element {
	if b.invalid {
		b.el.Attr("aria-invalid", "true")
	}
	b.el.Placeholder(b.placeholder)

	return b.el
}

// ==========================================
// COMBOBOX CHIP
// ==========================================

type comboboxChipBuilder struct {
	baseBuilder[*comboboxChipBuilder]

	value      string
	showRemove bool
}

func ComboboxChip(value string) *comboboxChipBuilder {
	el := nn.Div().
		Attr("data-slot", "combobox-chip").
		Attr("tabindex", "-1").
		Class("flex h-[calc(--spacing(5.5))] w-fit items-center justify-center gap-1 rounded-3xl bg-input px-2 text-xs font-medium whitespace-nowrap text-foreground has-disabled:pointer-events-none has-disabled:cursor-not-allowed has-disabled:opacity-50 has-data-[slot=combobox-chip-remove]:pr-0 dark:bg-input/60")

	b := &comboboxChipBuilder{value: value}
	b.baseBuilder = base(b, el)

	return b
}

func (b *comboboxChipBuilder) ShowRemove(sr bool) *comboboxChipBuilder {
	b.showRemove = sr
	return b
}

func (b *comboboxChipBuilder) build() *nn.Element {
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
			)

		b.el.Children(
			btn,
		).On("mousedown", func(e nn.Event) { e.PreventDefault(); e.PreventUpdate() })
	}

	return b.el.Key("chip-" + b.value)
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
}

type comboboxState[T any] struct {
	searchQuery string
	isOpen      bool
	isMenuReady bool

	activeOpts []T
	contentPos comboboxContentPosition
}

type combobox[T any] struct {
	nn.BaseComponent
	nn.State[comboboxState[T]]

	config ComboboxConfig[T]

	anchorRef  *nn.Ref
	contentRef *nn.Ref

	renderTrigger func(ctx ComboboxContext[T]) nn.IntoNode
	renderContent func(ctx ComboboxContext[T]) nn.IntoNode
}

func Combobox[T any](cfg ComboboxConfig[T]) *combobox[T] {
	return &combobox[T]{
		config:     cfg,
		anchorRef:  nn.NewRef(),
		contentRef: nn.NewRef(),
	}
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

func (c *combobox[T]) BindTrigger(cb func(ctx ComboboxContext[T]) nn.IntoNode) *combobox[T] {
	c.renderTrigger = cb
	return c
}

func (c *combobox[T]) BindContent(cb func(ctx ComboboxContext[T]) nn.IntoNode) *combobox[T] {
	c.renderContent = cb
	return c
}

func (c *combobox[T]) bindTriggerEvents(trigger nn.IntoNode) *nn.Element {
	n := trigger.ToNode()
	el, ok := n.(*nn.Element)
	if !ok {
		fmt.Println("[ERR] expected *Element as combobox trigger")
		return nil
	}

	procEl := func(el *nn.Element) {
		switch el.GetAttr("data-slot") {
		case "combobox-input-group":
			el.Ref(c.anchorRef).
				OnKeyDown(c.onKeyDown)

		case "combobox-input":
			el.OnInput(
				func(e nn.Event) {
					c.Data.searchQuery = e.TargetValue()
					if !c.Data.isOpen {
						c.openContent()
					}
				}).
				On("focus", c.onToggle).
				On("focusout", func(nn.Event) {
					if c.Data.isOpen {
						c.closeContent()
					}
				})

			if len(c.Data.activeOpts) > 0 {
				el.Value(c.config.ItemToString(c.Data.activeOpts[0]))
			}
		case "combobox-trigger":
			el.OnClick(c.onToggle)

			/*CHIPS*/
		case "combobox-chips":
			el.Ref(c.anchorRef).
				OnKeyDown(c.onKeyDown).
				OnResize(
					func(nn.Event) {
						c.Data.contentPos = c.calcContentPosition(c.anchorRef.Current)
					},
				)

		case "combobox-chip-input":
			el.OnInput(
				func(e nn.Event) {
					v := e.TargetValue()
					c.Data.searchQuery = strings.ToLower(v)
					if !c.Data.isOpen {
						c.openContent()
					}
				}).
				On("focus", c.onToggle).
				On("focusout", func(nn.Event) {
					if c.Data.isOpen {
						c.closeContent()
					}
				})
		case "combobox-chip-close":
			el.OnClick(func(e nn.Event) {
				e.PreventDefault()

				val := el.GetAttr("bind-val")
				c.Data.activeOpts = slices.DeleteFunc(
					c.Data.activeOpts,
					func(o T) bool { return c.config.ItemToValue(o) == val },
				)
			})

		}
	}

	if el != nil {
		el.Walk(procEl)
	}

	return el
}

func (c *combobox[T]) bindContentEvents(content nn.IntoNode) *nn.Element {
	n := content.ToNode()
	el, ok := n.(*nn.Element)
	if !ok {
		fmt.Println("[ERR] expected *Element as combobox content")
		return nil
	}

	if el != nil {
		el.Style("visibility: hidden;").
			OnGlobal("click", c.onGlobalClick).
			OnGlobal("scroll", c.onScroll)

		if c.Data.isMenuReady {
			el.Style("visibility: visible;")
		} else {
			go func() {
				<-nn.WaitForPaint()

				c.Data.isMenuReady = true
				c.Update()
			}()
		}
	}

	c.Data.contentPos = c.calcContentPosition(c.anchorRef.Current)
	positioner := &comboboxContentPositioner{
		el:  el,
		ref: c.contentRef,
		pos: c.Data.contentPos,
	}

	return positioner.build()
}

func (c *combobox[T]) View() *nn.Element {
	ctx := ComboboxContext[T]{
		IsOpen:        c.Data.isOpen,
		VisibleItems:  c.visibleItems(),
		SelectedItems: c.Data.activeOpts,
		IsActive:      c.isActive,
	}

	var children []nn.IntoNode

	if c.renderTrigger != nil {
		triggerNode := c.renderTrigger(ctx)
		children = append(children, c.bindTriggerEvents(triggerNode))
	}

	if c.Data.isOpen && c.renderContent != nil {
		contentNode := c.renderContent(ctx)
		children = append(children, c.bindContentEvents(contentNode))
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

func (c *combobox[T]) calcContentPosition(parent js.Value) comboboxContentPosition {
	var ret comboboxContentPosition
	cur := parent

	var minHeight float64

	if !cur.IsNull() && !cur.IsUndefined() {
		rect := cur.Call("getBoundingClientRect")

		window := js.Global().Get("window")
		scrollY := window.Get("scrollY").Float()
		scrollX := window.Get("scrollX").Float()

		ret.anchorWidth = rect.Get("width").Float()
		ret.anchorHeight = rect.Get("height").Float()
		ret.availableWidth = window.Get("innerWidth").Float()
		ret.availableHeight = window.Get("innerHeight").Float() - rect.Get("bottom").Float()

		content := c.contentRef.Current
		if !content.IsUndefined() && !content.IsNull() {
			cRect := content.Call("getBoundingClientRect")
			minHeight = cRect.Get("height").Float()
		}
		if minHeight < 30 {
			minHeight = 200
		}

		if ret.availableHeight >= minHeight+36 {
			ret.top = rect.Get("bottom").Float() + scrollY + 4
			ret.left = rect.Get("left").Float() + scrollX
		} else {
			ret.top = rect.Get("top").Float() + scrollY - minHeight - 4
			ret.left = rect.Get("left").Float() + scrollX
			ret.availableHeight = minHeight + 36
		}
	}

	return ret
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
	c.Data.isMenuReady = false
}

func (c *combobox[T]) activateSelected() {
	content := c.contentRef.Current
	current := content.Call("querySelector", "[data-highlighted]")
	if current.IsNull() {
		return
	}

	val := current.Call("getAttribute", "value").String()
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

func (c *combobox[T]) onGlobalClick(e nn.Event) {
	target := e.Target()

	trigger := c.anchorRef.Current
	clickedInTrigger := !trigger.IsUndefined() && !trigger.IsNull() && trigger.Call("contains", target).Bool()
	content := c.contentRef.Current
	clickedInContent := !content.IsUndefined() && !content.IsNull() && content.Call("contains", target).Bool()

	if clickedInContent {
		c.activateSelected()
		return
	}

	if !clickedInTrigger && !clickedInContent {
		c.closeContent()
	} else {
		e.PreventUpdate()
	}
}

func (c *combobox[T]) onScroll(e nn.Event) {
	if c.Data.isOpen {
		c.Data.contentPos = c.calcContentPosition(c.anchorRef.Current)
	} else {
		e.PreventUpdate()
	}
}

func (c *combobox[T]) onKeyDown(e nn.Event) {
	key := e.Key()

	switch key {
	case "ArrowDown", "ArrowUp":
	case "Enter":
	default:
		return
	}

	e.PreventDefault()

	if c.Data.isOpen == false {
		if key == "ArrowDown" {
			c.openContent()
		}
		return
	}

	content := c.contentRef.Current
	if content.IsUndefined() || content.IsNull() {
		return
	}

	items := content.Call("querySelectorAll", "[data-slot='combobox-item']")
	length := items.Get("length").Int()
	if length == 0 {
		return
	}

	currentIndex := -1
	for i := 0; i < length; i++ {
		if items.Index(i).Call("hasAttribute", "data-highlighted").Bool() {
			currentIndex = i
			break
		}
	}

	if currentIndex >= 0 && key == "Enter" {
		c.activateSelected()
		return
	}

	e.PreventUpdate()
	newIndex := currentIndex
	if key == "ArrowDown" {
		newIndex++
		if newIndex >= length {
			newIndex = 0
		}
	} else if key == "ArrowUp" {
		newIndex--
		if newIndex < 0 {
			newIndex = length - 1
		}
	}

	if currentIndex != -1 {
		items.Index(currentIndex).Call("removeAttribute", "data-highlighted")
	}

	newItem := items.Index(newIndex)
	newItem.Call("setAttribute", "data-highlighted", "")

	scrollOptions := map[string]any{"block": "nearest"}
	newItem.Call("scrollIntoView", scrollOptions)
}

func (c *combobox[T]) closeContent() {
	content := c.contentRef.Current
	if content.IsUndefined() || content.IsNull() {
		return
	}

	trigger := c.anchorRef.Current
	if !trigger.IsUndefined() && !trigger.IsNull() {
		input := trigger.Call("querySelector", "input")
		if !c.config.MultipleSelect && len(c.Data.activeOpts) > 0 {
			input.Set("value", c.config.ItemToString(c.Data.activeOpts[0]))
		} else {
			input.Set("value", "")
		}
	}

	var onAnimationEnd js.Func
	onAnimationEnd = js.FuncOf(func(this js.Value, args []js.Value) any {
		c.Data.isOpen = false
		c.Data.searchQuery = ""

		content.Call("removeEventListener", "animationend", onAnimationEnd)
		onAnimationEnd.Release()

		content.Call("removeAttribute", "data-closed")
		content.Call("setAttribute", "data-open", "true")
		c.Update()

		return nil
	})

	content.Call("addEventListener", "animationend", onAnimationEnd)
	content.Call("removeAttribute", "data-open")
	content.Call("setAttribute", "data-closed", "true")
}
