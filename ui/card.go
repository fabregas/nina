package ui

import "github.com/fabregas/nina/nn"

// --- CARD ---

type cardBuilder struct {
	baseBuilder[*cardBuilder]
}

func Card() *cardBuilder {
	baseClass := "group/card flex flex-col gap-6 overflow-hidden rounded-4xl bg-card py-6 text-sm text-card-foreground shadow-md ring-1 ring-foreground/5 has-[>img:first-child]:pt-0 data-[size=sm]:gap-4 data-[size=sm]:py-4 dark:ring-foreground/10 *:[img:first-child]:rounded-t-4xl *:[img:last-child]:rounded-b-4xl"

	el := nn.Div().
		Attr("data-slot", "card").
		Attr("data-size", "default").
		Class(baseClass)

	c := &cardBuilder{}

	c.baseBuilder = base(c, el)

	return c
}

func (c *cardBuilder) SizeDefault() *cardBuilder {
	c.el.Attr("data-size", "default")
	return c
}

func (c *cardBuilder) SizeSm() *cardBuilder {
	c.el.Attr("data-size", "sm")
	return c
}
func (c *cardBuilder) build() *nn.Element { return c.el }

// --- CARD HEADER ---

func CardHeader() *simpleBuilder {
	baseClass := "group/card-header @container/card-header grid auto-rows-min items-start gap-1.5 rounded-t-4xl px-6 group-data-[size=sm]/card:px-4 has-data-[slot=card-action]:grid-cols-[1fr_auto] has-data-[slot=card-description]:grid-rows-[auto_auto] [.border-b]:pb-6 group-data-[size=sm]/card:[.border-b]:pb-4"

	return simple(
		nn.Div().
			Attr("data-slot", "card-header").
			Class(baseClass),
	)
}

// --- CARD TITLE ---

func CardTitle() *simpleBuilder {
	baseClass := "cn-font-heading text-base font-medium"
	return simple(
		nn.Div().
			Attr("data-slot", "card-title").
			Class(baseClass),
	)

}

// --- CARD DESCRIPTION ---

func CardDescription() *simpleBuilder {
	baseClass := "text-sm text-muted-foreground"
	return simple(
		nn.Div().
			Attr("data-slot", "card-description").
			Class(baseClass),
	)

}

// --- CARD ACTION ---

func CardAction() *simpleBuilder {
	baseClass := "col-start-2 row-span-2 row-start-1 self-start justify-self-end"
	return simple(
		nn.Div().
			Attr("data-slot", "card-action").
			Class(baseClass),
	)
}

// --- CARD CONTENT ---

func CardContent() *simpleBuilder {
	baseClass := "px-6 group-data-[size=sm]/card:px-4"
	return simple(
		nn.Div().
			Attr("data-slot", "card-content").
			Class(baseClass),
	)
}

// --- CARD FOOTER ---

func CardFooter() *simpleBuilder {
	baseClass := "flex items-center rounded-b-4xl px-6 group-data-[size=sm]/card:px-4 [.border-t]:pt-6 group-data-[size=sm]/card:[.border-t]:pt-4"
	return simple(
		nn.Div().
			Attr("data-slot", "card-footer").
			Class(baseClass),
	)
}
