package ui

type separatorBuilder struct {
	baseBuilder[*separatorBuilder]
}

func Separator() *separatorBuilder {
	b := &separatorBuilder{}
	b.baseBuilder = base(b, "div")
	b.Attr("data-slot", "separator").
		Class("shrink-0 bg-border data-horizontal:h-px data-horizontal:w-full data-vertical:w-px data-vertical:self-stretch").
		Horizontal()

	return b
}

func (s *separatorBuilder) Vertical() *separatorBuilder {
	s.Attr("orientation", "vertical")
	return s
}

func (s *separatorBuilder) Horizontal() *separatorBuilder {
	s.Attr("orientation", "horizontal")
	return s
}

func (s *separatorBuilder) build(_ *buildContext) {}
