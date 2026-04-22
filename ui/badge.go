package ui

import "github.com/fabregas/nina/nn"

type badgeBuilder struct {
	baseBuilder[*badgeBuilder]

	variant string
}

func Badge() *badgeBuilder {
	baseClass := "group/badge inline-flex h-5 w-fit shrink-0 items-center justify-center gap-1 overflow-hidden rounded-3xl border border-transparent px-2 py-0.5 text-xs font-medium whitespace-nowrap transition-all focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 has-data-[icon=inline-end]:pr-1.5 has-data-[icon=inline-start]:pl-1.5 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&>svg]:pointer-events-none [&>svg]:size-3!"

	el := nn.Span().
		Attr("data-slot", "badge").
		Class(baseClass)

	b := &badgeBuilder{}

	b.baseBuilder = base(b, el)
	b.Default()

	return b
}

func (b *badgeBuilder) Default() *badgeBuilder {
	b.variant = "default"
	return b
}

func (b *badgeBuilder) Secondary() *badgeBuilder {
	b.variant = "secondary"
	return b
}

func (b *badgeBuilder) Destructive() *badgeBuilder {
	b.variant = "destructive"
	return b
}

func (b *badgeBuilder) Outline() *badgeBuilder {
	b.variant = "outline"
	return b
}

func (b *badgeBuilder) Ghost() *badgeBuilder {
	b.variant = "ghost"
	return b
}

func (b *badgeBuilder) Link() *badgeBuilder {
	b.variant = "link"
	return b
}

func (b *badgeBuilder) build() *nn.Element {
	switch b.variant {
	case "secondary":
		b.Class("bg-secondary text-secondary-foreground [a]:hover:bg-secondary/80")
	case "destructive":
		b.Class("bg-destructive/10 text-destructive focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:focus-visible:ring-destructive/40 [a]:hover:bg-destructive/20")
	case "outline":
		b.Class("border-border text-foreground [a]:hover:bg-muted [a]:hover:text-muted-foreground")
	case "ghost":
		b.Class("hover:bg-muted hover:text-muted-foreground dark:hover:bg-muted/50")
	case "link":
		b.Class("text-primary underline-offset-4 hover:underline")
	default:
		b.Class("bg-primary text-primary-foreground [a]:hover:bg-primary/80")
	}

	return b.el
}
