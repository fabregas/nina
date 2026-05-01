package ui

// ==========================================
// BUTTON GROUP
// ==========================================

type buttonGroupBuilder struct {
	baseBuilder[*buttonGroupBuilder]
	orientationAttr string
}

func ButtonGroup() *buttonGroupBuilder {
	baseClass := "flex w-fit items-stretch *:focus-visible:relative *:focus-visible:z-10 has-[>[data-slot=button-group]]:gap-2 has-[>[data-variant=outline]]:*:data-[slot=input-group]:border-border has-[>[data-variant=outline]]:*:data-[slot=select-trigger]:border-border has-[>[data-variant=outline]]:[&>[data-slot=input-group]:has(:focus-visible)]:border-ring has-[>[data-variant=outline]]:[&>[data-slot=select-trigger]:focus-visible]:border-ring has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-4xl [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1 has-[>[data-variant=outline]]:[&>input]:border-border has-[>[data-variant=outline]]:[&>input:focus-visible]:border-ring"

	b := &buttonGroupBuilder{}
	b.baseBuilder = base(b, "div")
	b.Attr("role", "group").
		Attr("data-slot", "button-group").
		Class(baseClass).
		Horizontal()

	return b
}

func (g *buttonGroupBuilder) Horizontal() *buttonGroupBuilder {
	g.orientationAttr = "horizontal"

	return g
}

func (g *buttonGroupBuilder) Vertical() *buttonGroupBuilder {
	g.orientationAttr = "vertical"

	return g
}

func (g *buttonGroupBuilder) build(ctx *buildContext) {
	ctx.Props.Attr("data-orientation", g.orientationAttr)

	switch g.orientationAttr {
	case "horizontal":
		ctx.Props.Class("*:data-slot:rounded-r-none [&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-4xl! [&>[data-slot]~[data-slot]]:rounded-l-none [&>[data-slot]~[data-slot]]:border-l-0")
	case "vertical":
		ctx.Props.Class("flex-col *:data-slot:rounded-b-none [&>[data-slot]:not(:has(~[data-slot]))]:rounded-b-4xl! [&>[data-slot]~[data-slot]]:rounded-t-none [&>[data-slot]~[data-slot]]:border-t-0")
	}
}

// ==========================================
// BUTTON GROUP TEXT
// ==========================================

func ButtonGroupText() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "button-group-text").
		Class("flex items-center gap-2 rounded-4xl border bg-muted px-2.5 text-sm font-medium [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4")
}

// ==========================================
// BUTTON GROUP SEPARATOR
// ==========================================

type buttonGroupSeparatorBuilder struct {
	baseBuilder[*buttonGroupSeparatorBuilder]
}

func ButtonGroupSeparator() *buttonGroupSeparatorBuilder {
	b := &buttonGroupSeparatorBuilder{}
	b.baseBuilder = base(b, "div")

	sepConfig := Separator().
		Attr("data-slot", "button-group-separator").
		Class("relative self-stretch bg-input data-horizontal:mx-px data-horizontal:w-auto data-vertical:my-px data-vertical:h-auto")

	b.MergeProps(&sepConfig.props)

	b.Vertical()

	return b
}

func (s *buttonGroupSeparatorBuilder) Vertical() *buttonGroupSeparatorBuilder {
	s.Attr("orientation", "vertical")
	return s
}

func (s *buttonGroupSeparatorBuilder) Horizontal() *buttonGroupSeparatorBuilder {
	s.Attr("orientation", "horizontal")
	return s
}
func (s *buttonGroupSeparatorBuilder) build(_ *buildContext) {}
