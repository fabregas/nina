package ui

import (
	"fmt"
)

func AspectRatio(x, y int) *simpleBuilder {
	return simple("div").
		Attr("data-slot", "aspect-ratio").
		Class("relative aspect-(--ratio)").
		Style(fmt.Sprintf("--ratio: %.4f;", float32(x)/float32(y)))
}
