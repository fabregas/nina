package ui

import "github.com/fabregas/nina/nn"

func Icon(svgContent string) *nn.Element {
	return nn.Raw("svg", svgContent).
		Attr("xmlns", "http://www.w3.org/2000/svg").
		Attr("viewBox", "0 0 24 24").
		Attr("fill", "none").
		Attr("stroke", "currentColor").
		Attr("stroke-width", "2").
		Attr("stroke-linecap", "round").
		Attr("stroke-linejoin", "round")
}

func IconCheck() *nn.Element {
	return Icon(`<path d="M20 6 9 17l-5-5"/>`)
}

func IconClose() *nn.Element {
	return Icon(`<path d="M18 6 6 18"/><path d="m6 6 12 12"/>`)
}

func IconEyeOff() *nn.Element {
	return Icon(`<path d="M10.733 5.076a10.744 10.744 0 0 1 11.205 6.575 1 1 0 0 1 0 .696 10.747 10.747 0 0 1-1.444 2.49" /> <path d="M14.084 14.158a3 3 0 0 1-4.242-4.242" /> <path d="M17.479 17.499a10.75 10.75 0 0 1-15.417-5.151 1 1 0 0 1 0-.696 10.75 10.75 0 0 1 4.446-5.143" /> <path d="m2 2 20 20"/>`)
}

func IconEye() *nn.Element {
	return Icon(`<path d="M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0" /> <circle cx="12" cy="12" r="3" />`)
}

func IconChevronDown() *nn.Element {
	return Icon(`<path d="m6 9 6 6 6-6" />`)
}

func IconChevronUp() *nn.Element {
	return Icon(`<circle cx="12" cy="12" r="10" /> <path d="m8 14 4-4 4 4" />`)
}

func IconX() *nn.Element {
	return Icon(`<path d="M18 6 6 18" /> <path d="m6 6 12 12" />`)
}
