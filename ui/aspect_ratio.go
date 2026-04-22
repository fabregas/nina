package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

func AspectRatio(x, y int) *simpleBuilder {
	el := nn.Div().
		Attr("data-slot", "aspect-ratio").
		Class("relative aspect-(--ratio)").
		Style(fmt.Sprintf("--ratio: %.4f;", float32(x)/float32(y)))

	return simple(el)
}
