package ui

type separatorBuilder struct {
	baseBuilder[*separatorBuilder]
}

func Separator() *separatorBuilder {
	b := &separatorBuilder{}
	b.baseBuilder = base(b, "div")
	b.Attr("data-slot", "separator").
		Class("shrink-0 bg-border data-horizontal:h-px data-horizontal:w-full data-vertical:w-px data-vertical:self-stretch")

	return b
}

func (s *separatorBuilder) Vertical() *separatorBuilder {
	s.Attr("data-vertical", "true")
	return s
}

func (s *separatorBuilder) Horizontal() *separatorBuilder {
	s.Attr("data-horizontal", "true")
	return s
}

func (s *separatorBuilder) build(_ *buildContext) {}
