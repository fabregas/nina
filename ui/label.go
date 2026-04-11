package ui

import "github.com/fabregas/nina/nn"

func Label() *simpleBuilder {
	baseClass := "flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50"

	return simple(
		nn.Label().
			Attr("data-slot", "label").
			Class(baseClass),
	)
}
