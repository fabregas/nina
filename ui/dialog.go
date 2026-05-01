package ui

import (
	"github.com/fabregas/nina/nn"
	"github.com/fabregas/nina/ui/icons"
)

// ==========================================
// DIALOG CONTENT
// ==========================================

func dialogOverlay() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "dialog-overlay").
		Class("fixed inset-0 isolate z-50 bg-black/30 duration-300 supports-backdrop-filter:backdrop-blur-sm data-open:animate-in data-open:fade-in-0 data-closed:animate-out data-closed:fade-out-0")
}

type dialogContent struct {
	uiComponent[*dialogContent]

	hideCloseBnt bool
}

func DialogContent() *dialogContent {
	b := &dialogContent{}
	b.init(b)

	return b
}

func (c *dialogContent) ShowCloseButton(t bool) *dialogContent {
	c.hideCloseBnt = !t
	return c
}

func (c *dialogContent) View() nn.Node {
	ctx := nn.GetContext[*dialogInternalCtx](c)

	el := nn.Div().
		Attr("data-slot", "dialog-content").
		Class("fixed top-1/2 left-1/2 z-50 grid w-full max-w-[calc(100%-2rem)] -translate-x-1/2 -translate-y-1/2 gap-6 rounded-4xl bg-popover p-6 text-sm text-popover-foreground shadow-xl ring-1 ring-foreground/5 duration-300 outline-none sm:max-w-md dark:ring-foreground/10 data-open:animate-in data-open:fade-in-0 data-open:zoom-in-95 data-closed:animate-out data-closed:fade-out-0 data-closed:zoom-out-95")

	if ctx.isClosing() {
		el.Attr("data-closed", "true")
	} else {
		el.Attr("data-open", "true")
	}

	el.Children(c.children...)

	closeFn := func(e nn.Event) {
		e.PreventUpdate()
		if ctx.closeDialog != nil {
			ctx.closeDialog()
		}
	}

	if !c.hideCloseBnt {
		el.Children(
			Button().
				Ghost().
				SizeSm().
				Class("absolute top-4 right-4").
				Attr("data-slot", "dialog-close").
				Children(icons.X()).
				OnClick(closeFn),
		)
	}

	return c.ApplyProps(el)
}

// ==========================================
// DIALOG HEADER
// ==========================================

func DialogHeader() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "dialog-header").
		Class("flex flex-col gap-1.5")
}

// ==========================================
// DIALOG FOOTER
// ==========================================

func DialogFooter() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "dialog-footer").
		Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end")
}

// ==========================================
// DIALOG TITLE
// ==========================================

func DialogTitle() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "dialog-title").
		Class("cn-font-heading text-base leading-none font-medium")
}

// ==========================================
// DIALOG DESCRIPTION
// ==========================================

func DialogDescription() *simpleBuilder {
	return simple("div").
		Attr("data-slot", "dialog-description").
		Class("text-sm text-muted-foreground *:[a]:underline *:[a]:underline-offset-3 *:[a]:hover:text-foreground")
}

// ==========================================
// DIALOG COMPONENT
// ==========================================

type dialogInternalCtx struct {
	closeDialog func()
	isClosing   func() bool
}

type dialogState struct {
	isMounted bool
	isClosing bool
}

type dialog struct {
	uiComponent[*dialog]
	nn.State[dialogState]

	onClose             func()
	isOpen              bool
	closeOnOutsideClick bool
}

func Dialog() *dialog {
	d := &dialog{}
	d.init(d)

	nn.ProvideContext(d,
		&dialogInternalCtx{
			closeDialog: func() {
				if d.onClose != nil {
					d.onClose()
				}
			},
			isClosing: func() bool { return d.Data.isClosing },
		},
	)

	return d
}

func (d *dialog) Open(b bool) *dialog {
	d.isOpen = b

	return d
}

func (d *dialog) CloseOnOutsideClick(b bool) *dialog {
	d.closeOnOutsideClick = b

	return d
}

func (d *dialog) OnClose(cb func()) *dialog {
	d.onClose = cb
	return d
}

func (d *dialog) View() nn.Node {
	if d.isOpen && !d.Data.isMounted {
		d.Data.isMounted = true
		d.Data.isClosing = false
	} else if !d.isOpen && d.Data.isMounted && !d.Data.isClosing {
		d.Data.isClosing = true
	}

	if !d.Data.isMounted {
		return nil
	}

	stateStr := "open"
	if d.Data.isClosing {
		stateStr = "closed"
	}

	overlay := dialogOverlay().
		Attr("data-"+stateStr, "true").
		On("animationend", func(nn.Event) {
			if d.Data.isClosing {
				d.Data.isClosing = false
				d.Data.isMounted = false
				d.Update()
			}
		})

	if d.closeOnOutsideClick {
		overlay.OnClick(func(nn.Event) {
			if d.onClose != nil {
				d.onClose()
			}
		})
	}

	return nn.Portal(
		nn.Div().
			Children(overlay).
			Children(d.children...),
	)

}
