package ui

import (
	"github.com/fabregas/nina/nn"
)

type buttonBuilder struct {
	baseBuilder[*buttonBuilder]

	variantClass      string
	sizeAttr          string
	customSizeClasses map[string]string

	onClick      func()
	onEventClick func(nn.Event)
}

func Button() *buttonBuilder {
	b := &buttonBuilder{}

	baseClass := "group/button inline-flex shrink-0 items-center justify-center rounded-4xl border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30 active:not-aria-[haspopup]:translate-y-px disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"

	btn := nn.Button().
		Attr("data-slot", "button").
		Class(baseClass)

	btn.OnClick(func(e nn.Event) {
		e.PreventDefault()

		if b.onEventClick != nil {
			b.onEventClick(e)
		} else if b.onClick != nil {
			b.onClick()
		}
	})

	b.baseBuilder = base(b, btn)

	return b.Primary().SizeMd()
}

func (b *buttonBuilder) overrideSize(sizeName string, tailwindClasses string) *buttonBuilder {
	if b.customSizeClasses == nil {
		b.customSizeClasses = make(map[string]string)
	}
	b.customSizeClasses[sizeName] = tailwindClasses
	return b
}

func (b *buttonBuilder) Primary() *buttonBuilder {
	b.variantClass = "bg-primary text-primary-foreground hover:bg-primary/80"
	return b
}

func (b *buttonBuilder) Outline() *buttonBuilder {
	b.variantClass = "border-border bg-background hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:bg-transparent dark:hover:bg-input/30"
	return b
}

func (b *buttonBuilder) Secondary() *buttonBuilder {
	b.variantClass = "bg-secondary text-secondary-foreground hover:bg-secondary/80 aria-expanded:bg-secondary aria-expanded:text-secondary-foreground"
	return b
}

func (b *buttonBuilder) Ghost() *buttonBuilder {
	b.variantClass = "hover:bg-muted hover:text-foreground aria-expanded:bg-muted aria-expanded:text-foreground dark:hover:bg-muted/50"
	return b
}

func (b *buttonBuilder) Destructive() *buttonBuilder {
	b.variantClass = "bg-destructive/10 text-destructive hover:bg-destructive/20 focus-visible:border-destructive/40 focus-visible:ring-destructive/20 dark:bg-destructive/20 dark:hover:bg-destructive/30 dark:focus-visible:ring-destructive/40"
	return b
}

func (b *buttonBuilder) Link() *buttonBuilder {
	b.variantClass = "text-primary underline-offset-4 hover:underline"
	return b
}

func (b *buttonBuilder) SizeXs() *buttonBuilder {
	b.sizeAttr = "xs"
	return b
}

func (b *buttonBuilder) SizeSm() *buttonBuilder {
	b.sizeAttr = "sm"
	return b
}

func (b *buttonBuilder) SizeMd() *buttonBuilder {
	b.sizeAttr = "md"
	return b
}

func (b *buttonBuilder) SizeLg() *buttonBuilder {
	b.sizeAttr = "lg"
	return b
}

func (b *buttonBuilder) SizeIcon() *buttonBuilder {
	b.sizeAttr = "icon"
	return b
}

func (b *buttonBuilder) SizeIconXs() *buttonBuilder {
	b.sizeAttr = "icon-xs"
	return b
}

func (b *buttonBuilder) SizeIconSm() *buttonBuilder {
	b.sizeAttr = "icon-sm"
	return b
}
func (b *buttonBuilder) SizeIconLg() *buttonBuilder {
	b.sizeAttr = "icon-lg"
	return b
}

func (b *buttonBuilder) OnClick(fn func()) *buttonBuilder {
	b.onClick = fn
	return b
}

func (b *buttonBuilder) OnEvent(fn func(nn.Event)) *buttonBuilder {
	b.onEventClick = fn
	return b
}

func (b *buttonBuilder) build() *nn.Element {
	var sizeClass string
	if custom, ok := b.customSizeClasses[b.sizeAttr]; ok {
		sizeClass = custom
	} else {
		switch b.sizeAttr {
		case "xs":
			sizeClass = "h-6 gap-1 px-2.5 text-xs has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2 [&_svg:not([class*='size-'])]:size-3"
		case "sm":
			sizeClass = "h-8 gap-1 px-3 has-data-[icon=inline-end]:pr-2 has-data-[icon=inline-start]:pl-2"
		case "lg":
			sizeClass = "h-10 gap-1.5 px-4 has-data-[icon=inline-end]:pr-3 has-data-[icon=inline-start]:pl-3"
		case "icon":
			sizeClass = "size-9"
		case "icon-xs":
			sizeClass = "size-6 [&_svg:not([class*='size-'])]:size-3"
		case "icon-sm":
			sizeClass = "size-8"
		case "icon-lg":
			sizeClass = "size-10"
		default: // default md
			sizeClass = "h-9 gap-1.5 px-3 has-data-[icon=inline-end]:pr-2.5 has-data-[icon=inline-start]:pl-2.5"
		}
	}
	b.el.Class(b.variantClass, sizeClass)

	return b.el
}
