package nn

type Props struct {
	classes      string
	key          string
	attrs        map[string]string
	events       map[string]func(Event)
	globalEvents map[string]func(Event)
	ref          *Ref
}

func (p *Props) Clone() *Props {
	clone := &Props{
		classes: p.classes,
		key:     p.key,
		ref:     p.ref,
	}
	if p.attrs != nil {
		clone.attrs = make(map[string]string, len(p.attrs))
		for k, v := range p.attrs {
			clone.attrs[k] = v
		}
	}
	if p.events != nil {
		clone.events = make(map[string]func(Event), len(p.events))
		for k, v := range p.events {
			clone.events[k] = v
		}
	}
	if p.globalEvents != nil {
		clone.globalEvents = make(map[string]func(Event), len(p.globalEvents))
		for k, v := range p.globalEvents {
			clone.globalEvents[k] = v
		}
	}

	return clone
}

func (p *Props) Class(class string) {
	if class == "" {
		return
	}
	if p.classes != "" {
		p.classes += " "
	}
	p.classes += class
}

func (p *Props) Key(key string) {
	p.key = key
}

func (p *Props) Ref(r *Ref) {
	p.ref = r
}

func (p *Props) Style(style string) {
	s, ok := p.attrs["style"]
	if ok {
		s += ";" + style
	} else {
		s = style
	}
	p.Attr("style", style)
}

func (p *Props) Attr(k, v string) {
	if p.attrs == nil {
		p.attrs = make(map[string]string)
	}
	p.attrs[k] = v
}

func (p *Props) On(e string, h func(Event)) {
	if p.events == nil {
		p.events = make(map[string]func(Event))
	}
	p.events[e] = h
}

func (p *Props) OnGlobal(e string, h func(Event)) {
	if p.globalEvents == nil {
		p.globalEvents = make(map[string]func(Event))
	}
	p.globalEvents[e] = h
}

func (p *Props) Merge(other *Props) {
	p.Class(other.classes)

	for k, v := range other.attrs {
		p.Attr(k, v)
	}
	for evt, handler := range other.events {
		p.On(evt, handler)
	}
	for evt, handler := range other.globalEvents {
		p.OnGlobal(evt, handler)
	}
	if other.key != "" {
		p.key = other.key
	}
	if other.ref != nil {
		p.ref = other.ref

	}
}

func (p *Props) MergeFromElement(el *Element) {
	p.Class(el.classes)

	for k, v := range el.attrs {
		p.Attr(k, v)
	}
	if el.listeners != nil {
		for einfo, handler := range el.listeners.events {
			if einfo.isGlobal {
				p.OnGlobal(einfo.name, handler)
			} else {
				p.On(einfo.name, handler)
			}
		}
	}

	if el.key != "" {
		p.key = el.key
	}

	if el.ref != nil {
		p.ref = el.ref
	}
}

func (p *Props) ApplyTo(el *Element) *Element {
	if p.classes != "" {
		el.Class(p.classes)
	}
	for k, v := range p.attrs {
		el.Attr(k, v)
	}
	for evt, handler := range p.events {
		el.On(evt, handler)
	}
	for evt, handler := range p.globalEvents {
		el.OnGlobal(evt, handler)
	}

	if p.key != "" {
		el.Key(p.key)
	}
	if p.ref != nil {
		el.Ref(p.ref)
	}

	return el
}
