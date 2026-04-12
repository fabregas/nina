package ui

import (
	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

func InputGroup() *simpleBuilder {
	return simple(
		nn.Div().
			Attr("data-slot", "input-group").
			Attr("role", "group").
			Class("group/input-group relative flex h-9 w-full min-w-0 items-center rounded-4xl border border-transparent bg-input/50 transition-[color,box-shadow,background-color] outline-none in-data-[slot=combobox-content]:focus-within:border-inherit in-data-[slot=combobox-content]:focus-within:ring-0 has-data-[align=block-end]:rounded-3xl has-data-[align=block-start]:rounded-3xl has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-3 has-[[data-slot=input-group-control]:focus-visible]:ring-ring/30 has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-3 has-[[data-slot][aria-invalid=true]]:ring-destructive/20 has-[textarea]:rounded-2xl has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>textarea]:h-auto dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40 has-[>[data-align=block-end]]:[&>input]:pt-3 has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=inline-end]]:[&>input]:pr-1.5 has-[>[data-align=inline-start]]:[&>input]:pl-1.5"),
	)
}

//--------------------------------

type inputGroupAddonBuilder struct {
	baseBuilder[*inputGroupAddonBuilder]
}

func InputGroupAddon() *inputGroupAddonBuilder {
	el := nn.Div().
		Attr("role", "group").
		Attr("data-slot", "input-group-addon").
		Class("flex items-center justify-center [&_svg:not([class*='size-'])]:size-4").
		OnClick(func(e nn.Event) {
			target := e.Target()
			closestBtn := target.Call("closest", "button")
			if !closestBtn.IsNull() && !closestBtn.IsUndefined() {
				return
			}

			group := target.Call("closest", "[data-slot='input-group']")

			if !group.IsNull() && !group.IsUndefined() {
				input := group.Call("querySelector", "input")

				if !input.IsNull() && !input.IsUndefined() {
					input.Call("focus")
				}
			}
		})

	g := &inputGroupAddonBuilder{}
	g.baseBuilder = base(g, el)
	g.AlignInlineStart()

	return g
}

func (g *inputGroupAddonBuilder) AlignInlineStart() *inputGroupAddonBuilder {
	g.el.Attr("data-align", "inline-start")
	g.el.Class("order-first pl-3 has-[>button]:-ml-1 has-[>kbd]:-ml-1")
	return g
}

func (g *inputGroupAddonBuilder) AlignInlineEnd() *inputGroupAddonBuilder {
	g.el.Attr("data-align", "inline-end")
	g.el.Class("order-last pr-3 has-[>button]:-mr-1 has-[>kbd]:-mr-1")
	return g
}

func (g *inputGroupAddonBuilder) AlignBlockStart() *inputGroupAddonBuilder {
	g.el.Attr("data-align", "block-start")
	g.el.Class("order-first w-full justify-start px-3 pt-3 group-has-[>input]/input-group:pt-3.5 [.border-b]:pb-3.5")
	return g
}

func (g *inputGroupAddonBuilder) AlignBlockEnd() *inputGroupAddonBuilder {
	g.el.Attr("data-align", "block-end")
	g.el.Class("order-last w-full justify-start px-3 pb-3 group-has-[>input]/input-group:pb-3.5 [.border-t]:pt-3.5")
	return g
}

func (g *inputGroupAddonBuilder) build() *nn.Element {
	return g.el
}

// ==========================================
// INPUT GROUP BUTTON
// ==========================================

func InputGroupButton() *buttonBuilder {
	return Button().
		Ghost().
		SizeXs().
		Class("flex items-center gap-2 rounded-4xl text-sm shadow-none").
		Attr("data-slot", "input-group-button").
		overrideSize("xs", "h-6 gap-1 rounded-xl px-1.5 [&>svg:not([class*='size-'])]:size-3.5").
		overrideSize("sm", "").
		overrideSize("icon-xs", "size-6 rounded-xl p-0 has-[>svg]:p-0").
		overrideSize("icon-sm", "size-8 p-0 has-[>svg]:p-0")
}

// ==========================================
// INPUT GROUP INPUT
// ==========================================

func InputGroupInput() *inputBuilder {
	return Input().
		Class("flex-1 rounded-none border-0 bg-transparent shadow-none ring-0 focus-visible:ring-0 aria-invalid:ring-0 dark:bg-transparent").
		Attr("data-slot", "input-group-control")
}

// ==========================================
// INPUT GROUP TEXT
// ==========================================

func InputGroupText() *simpleBuilder {
	baseClass := "flex items-center gap-2 text-sm text-muted-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"

	return simple(
		nn.Span().
			Attr("data-slot", "input-group-text").
			Class(baseClass),
	)
}

// ---------- TODO move to some helpers components ------

type pwdInputState struct {
	isShowing bool
}

type passwordInput struct {
	nn.State[pwdInputState]

	id    string
	label string
	val   *string
}

func PasswordInputWrapper(id, label string, val *string) *passwordInput {
	return &passwordInput{
		id:    id,
		label: label,
		val:   val,
	}
}

func (i *passwordInput) toggle() {
	i.S.isShowing = !i.S.isShowing
}

func (i *passwordInput) View() *nn.Element {
	eyeBtn := Button().
		Ghost().
		SizeIconSm().
		OnClick(func() { i.toggle() })

	input := InputGroupInput().ID(i.id).Bind(i.val)

	if i.S.isShowing {
		input.TypeText()
		eyeBtn.Children(icons.Eye())

	} else {
		input.TypePassword()
		eyeBtn.Children(icons.EyeOff())
	}

	return InputGroup().Children(
		input,
		InputGroupAddon().AlignInlineEnd().Children(
			eyeBtn,
		),
	).El()
}
