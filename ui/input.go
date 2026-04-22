package ui

import "github.com/fabregas/nina/nn"

type inputBuilder struct {
	baseBuilder[*inputBuilder]

	value   *string
	onInput func(string)
	onEvent func(nn.Event)
}

func Input() *inputBuilder {
	b := &inputBuilder{}
	inputEl := nn.Input().
		Attr("data-slot", "input").
		Class("h-8 w-full min-w-0 rounded-3xl border border-transparent bg-input/50 px-3 py-1 text-base transition-[color,box-shadow,background-color] outline-none file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/30 disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40")

	inputEl.OnInput(func(e nn.Event) {
		if b.value != nil {
			*b.value = e.TargetValue()
		}

		if b.onInput != nil {
			b.onInput(e.TargetValue())
		}
		if b.onEvent != nil {
			b.onEvent(e)
		}
	})

	b.baseBuilder = base(b, inputEl)
	b.TypeText()

	return b
}

func (i *inputBuilder) TypeText() *inputBuilder     { i.el.Type("text"); return i }
func (i *inputBuilder) TypePassword() *inputBuilder { i.el.Type("password"); return i }
func (i *inputBuilder) TypeEmail() *inputBuilder    { i.el.Type("email"); return i }
func (i *inputBuilder) TypeNumber() *inputBuilder   { i.el.Type("number"); return i }

func (i *inputBuilder) NoArrows() *inputBuilder {
	i.el.Class("[appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none")
	return i
}
func (i *inputBuilder) Placeholder(p string) *inputBuilder {
	i.el.Attr("placeholder", p)
	return i
}

func (i *inputBuilder) Bind(v *string) *inputBuilder { i.value = v; return i }

func (i *inputBuilder) Value(v string) *inputBuilder {
	if i.value == nil {
		i.value = new(string)
	}
	*i.value = v
	return i
}

func (i *inputBuilder) OnInput(fn func(string)) *inputBuilder {
	i.onInput = fn
	return i
}

func (i *inputBuilder) OnEvent(fn func(nn.Event)) *inputBuilder {
	i.onEvent = fn
	return i
}

func (i *inputBuilder) build() *nn.Element {
	if i.value != nil {
		i.el.Value(*i.value)
	}

	return i.el
}
