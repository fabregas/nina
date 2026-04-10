package ui

import "github.com/fabregas/nina/nn"

type builder interface {
	build() *nn.Element
}

type wrapperBuilder interface {
	wrap(target *nn.Element) *nn.Element
}

type baseBuilder[T builder] struct {
	el            *nn.Element
	self          T
	customClasses []string
}

func base[T builder](self T, el *nn.Element) baseBuilder[T] {
	return baseBuilder[T]{
		el:   el,
		self: self,
	}
}

func (b *baseBuilder[T]) Text(text string) T {
	b.el.Text(text)
	return b.self
}

func (b *baseBuilder[T]) Class(class string) T {
	if class != "" {
		b.customClasses = append(b.customClasses, class)
	}
	return b.self
}

func (b *baseBuilder[T]) Children(items ...nn.IntoNode) T {
	b.el.Children(items...)
	return b.self
}

func (b *baseBuilder[T]) Attr(key, value string) T {
	b.el.Attr(key, value)
	return b.self
}

func (b *baseBuilder[T]) ID(id string) T {
	return b.Attr("id", id)
}

func (b *baseBuilder[T]) For(id string) T {
	return b.Attr("for", id)
}

func (b *baseBuilder[T]) ToNode() nn.Node { return b.El() }

func (b *baseBuilder[T]) El() *nn.Element {
	finalEl := b.self.build()

	for _, cls := range b.customClasses {
		finalEl.Class(cls)
	}

	if wrapper, ok := any(b.self).(wrapperBuilder); ok {
		return wrapper.wrap(finalEl)
	}

	return finalEl
}

type simpleBuilder struct {
	baseBuilder[*simpleBuilder]
}

func (s *simpleBuilder) build() *nn.Element {
	return s.el
}

func simple(el *nn.Element) *simpleBuilder {
	b := &simpleBuilder{}
	b.baseBuilder = baseBuilder[*simpleBuilder]{
		el:   el,
		self: b,
	}

	return b
}
