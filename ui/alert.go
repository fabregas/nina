package ui

// ==========================================
// ALERT
// ==========================================

type alertBuilder struct {
	baseBuilder[*alertBuilder]
	variantAttr string
}

func Alert() *alertBuilder {
	b := &alertBuilder{}

	baseClass := "group/alert relative grid w-full gap-0.5 rounded-2xl border px-4 py-3 text-left text-sm has-data-[slot=alert-action]:relative has-data-[slot=alert-action]:pr-18 has-[>svg]:grid-cols-[auto_1fr] has-[>svg]:gap-x-2.5 *:[svg]:row-span-2 *:[svg]:translate-y-0.5 *:[svg]:text-current *:[svg:not([class*='size-'])]:size-4"

	b.baseBuilder = base(b, "div")
	b.Attr("data-slot", "alert").
		Attr("role", "alert").
		Class(baseClass)
	b.Default()

	return b
}

func (a *alertBuilder) Default() *alertBuilder {
	a.variantAttr = "default"
	return a
}

func (a *alertBuilder) Destructive() *alertBuilder {
	a.variantAttr = "destructive"
	return a
}

func (a *alertBuilder) build(ctx *buildContext) {
	var variantClass string

	switch a.variantAttr {
	case "destructive":
		variantClass = "bg-card text-destructive border-destructive/50 *:data-[slot=alert-description]:text-destructive/90 *:[svg]:text-current"
	default: // "default"
		variantClass = "bg-card text-card-foreground"
	}

	ctx.Props.Class(variantClass)
}

// ==========================================
// ALERT TITLE
// ==========================================

func AlertTitle() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "alert-title").
		Class("cn-font-heading font-medium group-has-[>svg]/alert:col-start-2 [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground")
}

// ==========================================
// ALERT DESCRIPTION
// ==========================================

func AlertDescription() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "alert-description").
		Class("text-sm text-balance text-muted-foreground md:text-pretty [&_a]:underline [&_a]:underline-offset-3 [&_a]:hover:text-foreground [&_p:not(:last-child)]:mb-4")

}

// ==========================================
// ALERT ACTION
// ==========================================

func AlertAction() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "alert-action").
		Class("absolute top-2.5 right-3")
}
