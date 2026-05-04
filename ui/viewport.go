package ui

import (
	"github.com/fabregas/nina/nn"
)

type Breakpoint string

const (
	BreakpointXS  Breakpoint = "xs"  // < 640px (mobile)
	BreakpointSM  Breakpoint = "sm"  // >= 640px (big mobiles)
	BreakpointMD  Breakpoint = "md"  // >= 768px (tablets)
	BreakpointLG  Breakpoint = "lg"  // >= 1024px (laptops)
	BreakpointXL  Breakpoint = "xl"  // >= 1280px (desktops)
	BreakpointXXL Breakpoint = "2xl" // >= 1536px (large displays)
)

var Viewport = nn.NewSignal[Breakpoint](BreakpointXS)

//var resizeCallback js.Func

func IsMobile(c nn.Component) bool {
	bp := Viewport.Get(c)
	return bp == BreakpointXS || bp == BreakpointSM
}

func getBreakpoint(width int) Breakpoint {
	switch {
	case width >= 1536:
		return BreakpointXXL
	case width >= 1280:
		return BreakpointXL
	case width >= 1024:
		return BreakpointLG
	case width >= 768:
		return BreakpointMD
	case width >= 640:
		return BreakpointSM
	default:
		return BreakpointXS
	}
}

func init() {
	nn.RegisterInitHook(func(r nn.Renderer) {
		vp := r.GetViewport()
		Viewport.Set(getBreakpoint(int(vp.Width)))

		r.AddEventListener(r.Window(), "resize", func(nn.Event) {
			vp := r.GetViewport()
			newBp := getBreakpoint(int(vp.Width))

			if newBp != Viewport.Get(nil) {
				Viewport.Set(newBp)
			}
		})
	})
}
