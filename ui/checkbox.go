package ui

import (
	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

type checkboxBuilder struct {
	baseBuilder[*checkboxBuilder]

	isChecked *bool
	onChange  func(bool)
}

func Checkbox() *checkboxBuilder {
	b := &checkboxBuilder{}

	baseClass := "peer relative flex size-4 shrink-0 items-center justify-center rounded-[5px] border border-transparent bg-input/90 transition-shadow outline-none group-has-disabled/field:opacity-50 after:absolute after:-inset-x-3 after:-inset-y-2 focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 aria-invalid:aria-checked:border-primary dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40 data-checked:border-primary data-checked:bg-primary data-checked:text-primary-foreground dark:data-checked:bg-primary"

	btn := nn.Button().
		Attr("type", "button").
		Attr("role", "checkbox").
		Attr("data-slot", "checkbox").
		Attr("aria-checked", "false").
		Class(baseClass)

	btn.OnClick(func(e nn.Event) {
		if b.isChecked != nil {
			*b.isChecked = !(*b.isChecked)
		}

		if b.onChange != nil {
			b.onChange(*b.isChecked)
		}
	})

	b.baseBuilder = base(b, btn)

	return b
}
func (c *checkboxBuilder) build() *nn.Element { return c.el }

func (c *checkboxBuilder) Bind(v *bool) *checkboxBuilder {
	c.isChecked = v

	return c.Checked(*v)
}

func (c *checkboxBuilder) Checked(checked bool) *checkboxBuilder {
	if c.isChecked == nil {
		c.isChecked = new(bool)
	}
	*c.isChecked = checked

	if c.isChecked != nil && *c.isChecked {
		c.el.Attr("data-checked", "true")
		c.el.Attr("aria-checked", "true")

		indicator := nn.Div().
			Attr("data-slot", "checkbox-indicator").
			Class("grid place-content-center text-current transition-none [&>svg]:size-3.5").
			Children(icons.Check())

		c.el.Children(indicator)
	}

	return c
}

func (c *checkboxBuilder) OnChange(fn func(bool)) *checkboxBuilder {
	c.onChange = fn
	return c
}
