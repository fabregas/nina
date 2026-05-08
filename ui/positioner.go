package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

type positionerState struct {
	top             float64
	left            float64
	availableHeight float64
	availableWidth  float64
	anchorWidth     float64
	anchorHeight    float64
	isReady         bool
}

type positioner struct {
	uiComponent[*positioner]
	nn.State[positionerState]

	side      string
	align     string
	flip      bool
	offset    float64
	anchorRef *nn.Ref
	children  []nn.AsNode
	onClose   func()

	contentRef *nn.Ref
}

func Positioner(anchorRef *nn.Ref) *positioner {
	p := &positioner{
		anchorRef:  anchorRef,
		contentRef: nn.NewRef(),
	}
	p.init(p)

	return p
}

func (c *positioner) Children(nodes ...nn.AsNode) *positioner {
	c.children = nodes
	return c
}

func (c *positioner) OnClose(fn func()) *positioner {
	c.onClose = fn
	return c
}

func (c *positioner) Close() {
	content := c.contentRef.Current
	if content == nil {
		return
	}

	var cleaner func()

	onAnimationEnd := func(e nn.Event) {
		if cleaner != nil {
			cleaner()
		}

		e.Renderer().RemoveAttribute(content, "data-closed")

		c.Data.isReady = false

		c.Update()

		if c.onClose != nil {
			c.onClose()
		}
	}

	r := c.contentRef.Renderer
	cleaner = r.AddEventListener(content, "animationend", onAnimationEnd)
	r.RemoveAttribute(content, "data-open")
	r.SetAttribute(content, "data-closed", "true")
}

func (c *positioner) View() nn.Node {
	visibility := "hidden"
	if c.Data.isReady {
		visibility = "visible"
	}

	if c.Data.availableHeight == 0 {
		c.Data.availableHeight = 1000
	}

	style := fmt.Sprintf("top: 0px; left: 0px; transform: translate(%.1fpx, %.1fpx); will-change: transform; --available-width: %.1fpx; --available-height: %.1fpx; --anchor-width: %.1fpx; --anchor-height: %.1fpx; --transform-origin: 146px -6px; visibility: %s", c.Data.left, c.Data.top, c.Data.availableWidth, c.Data.availableHeight, c.Data.anchorWidth, c.Data.anchorHeight, visibility)

	class := "duration-300 data-[side=bottom]:slide-in-from-top-2 data-[side=inline-end]:slide-in-from-left-2 data-[side=inline-start]:slide-in-from-right-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95"

	content := nn.Div().Children(c.children...).
		Class(class).
		Attr("data-side", c.side).
		Attr("data-align", c.align).
		Ref(c.contentRef)

	if c.Data.isReady {
		content.Attr("data-open", "true")
	}

	return nn.Div().Children(
		nn.Portal(
			nn.Div().
				Class("absolute z-50").
				Style(style).
				Children(content).
				OnGlobal("keydown", c.onKeyDown).
				OnGlobal("scroll", c.onScroll).
				OnGlobal("click", c.onGlobalClick),
		),
	)
}

func (c *positioner) OnMount() {
	c.recalculatePosition()
	c.Data.isReady = true

	c.Update()
}

func (c *positioner) onScroll(e nn.Event) {
	c.recalculatePosition()
}

func (c *positioner) onKeyDown(e nn.Event) {
	key := e.Key()
	if key == "Escape" {
		c.Close()
	}
}

func (c *positioner) onGlobalClick(e nn.Event) {
	target := e.Target()

	r := e.Renderer()
	trigger := c.anchorRef.Current
	clickedInTrigger := trigger != nil && r.Contains(trigger, target)
	content := c.contentRef.Current
	clickedInContent := content != nil && r.Contains(content, target)

	if !clickedInTrigger && !clickedInContent {
		c.Close()
	} else {
		e.PreventUpdate()
	}
}

func (c *positioner) recalculatePosition() {
	anchorDOM := c.anchorRef.Current
	floatingDOM := c.contentRef.Current

	if anchorDOM == nil || floatingDOM == nil {
		return
	}

	r := c.anchorRef.Renderer
	aRect := r.GetBoundingClientRect(anchorDOM)
	fRect := r.GetBoundingClientRect(floatingDOM)

	aTop := aRect.Top
	aBottom := aRect.Bottom
	aLeft := aRect.Left
	aWidth := aRect.Width
	aHeight := aRect.Height

	c.Data.anchorWidth = aWidth
	c.Data.anchorHeight = aHeight

	fWidth := fRect.Width
	fHeight := fRect.Height

	offset := c.offset

	var x, y float64

	// ==========================================
	// base positioning
	// ==========================================
	switch c.side {
	case "inline-start":
		x = aLeft - offset - fWidth // align by left side
	case "inline-end":
		x = aLeft + offset + aWidth
	case "top":
		y = aTop - fHeight - offset
	default: // bottom
		y = aBottom + offset

	}
	switch c.align {
	case "start":
		switch c.side {
		case "inline-start", "inline-end":
			y = aTop
		default:
			x = aLeft
		}
	case "end":
		switch c.side {
		case "inline-start", "inline-end":
			y = aTop + aHeight - fHeight
		default:
			x = aLeft + aWidth - fWidth
		}
	default: // center
		switch c.side {
		case "inline-start", "inline-end":
			y = aTop + (aHeight / 2) - (fHeight / 2) //center vertically
		default:
			x = aLeft + (aWidth / 2) - (fWidth / 2) // center horizontally

		}
	}

	viewport := r.GetViewport()

	vh := viewport.Height
	vw := viewport.Width
	spaceAbove := aTop - offset
	spaceBelow := vh - (aBottom + offset)
	c.Data.availableWidth = vw

	// ==========================================
	// calc available height
	// ==========================================
	switch c.side {
	case "bottom":
		c.Data.availableHeight = spaceBelow

	case "top":
		c.Data.availableHeight = spaceAbove
	}

	// ==========================================
	// smart Flip
	// ==========================================
	if c.flip {
		if c.side == "bottom" && (y+fHeight > vh) {
			if spaceAbove >= fHeight {
				c.Data.availableHeight = spaceAbove
				y = aTop - fHeight - offset
			}
		}

		if c.side == "top" && (y < 0) {
			if spaceBelow >= fHeight {
				c.Data.availableHeight = spaceBelow
				y = aBottom + offset
			}
		}

		if x+fWidth > vw {
			x = vw - fWidth - 8
		}
		if x < 0 {
			x = 8
		}
	}

	// ==========================================
	// convert into abs coordinates
	// ==========================================

	c.Data.left = x + viewport.ScrollX
	c.Data.top = y + viewport.ScrollY
	c.Data.availableHeight += viewport.ScrollY
	//if c.Data.top < 0 {
	//	c.Data.top = 0
	//}
}

type positionerContext[T any] struct {
	instance T
	pos      *positioner
}

func (c *positionerContext[T]) SideTop() T {
	c.pos.side = "top"
	return c.instance
}

func (c *positionerContext[T]) SideBottom() T {
	c.pos.side = "bottom"
	return c.instance
}

func (c *positionerContext[T]) SideLeft() T {
	c.pos.side = "inline-start"
	return c.instance
}

func (c *positionerContext[T]) SideRight() T {
	c.pos.side = "inline-end"
	return c.instance
}

func (c *positionerContext[T]) AlignStart() T {
	c.pos.align = "start"
	return c.instance
}

func (c *positionerContext[T]) AlignCenter() T {
	c.pos.align = "center"
	return c.instance
}

func (c *positionerContext[T]) AlignEnd() T {
	c.pos.align = "end"
	return c.instance
}

func (c *positionerContext[T]) Flip() T {
	c.pos.flip = true
	return c.instance
}

func (c *positionerContext[T]) Offset(o float64) T {
	c.pos.offset = o
	return c.instance
}

func getPositionerContext[T any](instance T, pos *positioner) *positionerContext[T] {
	return &positionerContext[T]{instance: instance, pos: pos}
}
