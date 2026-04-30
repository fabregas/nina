package ui

import (
	"fmt"
	"syscall/js"

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
	nn.BaseComponent
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
	if content.IsUndefined() || content.IsNull() {
		return
	}

	var onAnimationEnd js.Func
	onAnimationEnd = js.FuncOf(func(this js.Value, args []js.Value) any {
		content.Call("removeEventListener", "animationend", onAnimationEnd)
		onAnimationEnd.Release()

		content.Call("removeAttribute", "data-closed")
		c.Data.isReady = false

		c.Update()

		if c.onClose != nil {
			c.onClose()
		}

		return nil
	})

	content.Call("addEventListener", "animationend", onAnimationEnd)
	content.Call("removeAttribute", "data-open")
	content.Call("setAttribute", "data-closed", "true")
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

	trigger := c.anchorRef.Current
	clickedInTrigger := !trigger.IsUndefined() && !trigger.IsNull() && trigger.Call("contains", target).Bool()
	content := c.contentRef.Current
	clickedInContent := !content.IsUndefined() && !content.IsNull() && content.Call("contains", target).Bool()

	if !clickedInTrigger && !clickedInContent {
		c.Close()
	} else {
		e.PreventUpdate()
	}
}

func (c *positioner) recalculatePosition() {
	anchorDOM := c.anchorRef.Current
	floatingDOM := c.contentRef.Current

	window := js.Global().Get("window")

	if anchorDOM.IsNull() || anchorDOM.IsUndefined() || floatingDOM.IsNull() || floatingDOM.IsUndefined() {
		return
	}

	aRect := anchorDOM.Call("getBoundingClientRect")
	fRect := floatingDOM.Call("getBoundingClientRect")

	aTop := aRect.Get("top").Float()
	aBottom := aRect.Get("bottom").Float()
	aLeft := aRect.Get("left").Float()
	aWidth := aRect.Get("width").Float()
	aHeight := aRect.Get("height").Float()
	c.Data.anchorWidth = aWidth
	c.Data.anchorHeight = aHeight

	fWidth := fRect.Get("width").Float()
	fHeight := fRect.Get("height").Float()

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

	vh := window.Get("innerHeight").Float()
	vw := window.Get("innerWidth").Float()
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

	scrollX := window.Get("scrollX").Float()
	scrollY := window.Get("scrollY").Float()

	c.Data.left = x + scrollX
	c.Data.top = y + scrollY
	c.Data.availableHeight += scrollY
	//if c.Data.top < 0 {
	//	c.Data.top = 0
	//}
}

type positionerContext struct {
	pos *positioner
}

func (c *positionerContext) SideTop() *positionerContext {
	c.pos.side = "top"
	return c
}

func (c *positionerContext) SideBottom() *positionerContext {
	c.pos.side = "bottom"
	return c
}

func (c *positionerContext) SideLeft() *positionerContext {
	c.pos.side = "inline-start"
	return c
}

func (c *positionerContext) SideRight() *positionerContext {
	c.pos.side = "inline-end"
	return c
}

func (c *positionerContext) AlignStart() *positionerContext {
	c.pos.align = "start"
	return c
}

func (c *positionerContext) AlignCenter() *positionerContext {
	c.pos.align = "center"
	return c
}

func (c *positionerContext) AlignEnd() *positionerContext {
	c.pos.align = "end"
	return c
}

func (c *positionerContext) Flip() *positionerContext {
	c.pos.flip = true
	return c
}

func (c *positionerContext) Offset(o float64) *positionerContext {
	c.pos.offset = o
	return c
}

func (c *positioner) context() *positionerContext {
	return &positionerContext{pos: c}
}
