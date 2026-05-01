package ui

import (
	"fmt"

	"github.com/fabregas/nina/nn"
)

// ==========================================
// TABLE
// ==========================================

type tableBuilder struct {
	baseBuilder[*tableBuilder]
}

func Table() *tableBuilder {
	baseClass := "w-full caption-bottom text-sm"

	b := &tableBuilder{}
	b.baseBuilder = base(b, "table")
	b.Attr("data-slot", "table").
		Class(baseClass)

	return b
}

func (t *tableBuilder) build(_ *buildContext) {}

func (t *tableBuilder) wrap(target nn.Node) nn.Node {
	return nn.Div().
		Attr("data-slot", "table-container").
		Class("relative w-full overflow-x-auto").
		Children(target)
}

// ==========================================
// TABLE HEADER
// ==========================================

func TableHeader() *simpleBuilder {
	baseClass := "[&_tr]:border-b"
	return simple("thead").
		Attr("data-slot", "table-header").
		Class(baseClass)
}

// ==========================================
// TABLE BODY
// ==========================================

func TableBody() *simpleBuilder {
	baseClass := "[&_tr:last-child]:border-0"
	return simple("tbody").
		Attr("data-slot", "table-body").
		Class(baseClass)
}

// ==========================================
// TABLE FOOTER
// ==========================================

func TableFooter() *simpleBuilder {
	baseClass := "border-t bg-muted/50 font-medium [&>tr]:last:border-b-0"
	return simple("tfoot").
		Attr("data-slot", "table-footer").
		Class(baseClass)
}

// ==========================================
// TABLE ROW
// ==========================================

type tableRowBuilder struct {
	baseBuilder[*tableRowBuilder]
}

func TableRow() *tableRowBuilder {
	baseClass := "border-b transition-colors hover:bg-muted/50 has-aria-expanded:bg-muted/50 data-[state=selected]:bg-muted"

	b := &tableRowBuilder{}
	b.baseBuilder = base(b, "tr")
	b.Attr("data-slot", "table-row").Class(baseClass)

	return b
}

func (t *tableRowBuilder) build(_ *buildContext) {}

func (t *tableRowBuilder) Selected(selected bool) *tableRowBuilder {
	t.Attr("data-state", "selected")
	return t
}

// ==========================================
// TABLE HEAD (TH)
// ==========================================

func TableHead() *simpleBuilder {
	baseClass := "h-12 px-3 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0"
	return simple("th").
		Attr("data-slot", "table-head").
		Class(baseClass)
}

// ==========================================
// TABLE CELL (TD)
// ==========================================

type tableCellBuilder struct {
	baseBuilder[*tableCellBuilder]
}

func TableCell() *tableCellBuilder {
	baseClass := "p-3 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0"

	b := &tableCellBuilder{}
	b.baseBuilder = base(b, "td")
	b.Attr("data-slot", "table-cell").Class(baseClass)

	return b
}

func (c *tableCellBuilder) build(_ *buildContext) {}

func (c *tableCellBuilder) ColSpan(colSpan int) *tableCellBuilder {
	c.Attr("colspan", fmt.Sprintf("%d", colSpan))
	return c
}

// ==========================================
// TABLE CAPTION
// ==========================================

func TableCaption() *simpleBuilder {
	baseClass := "mt-4 text-sm text-muted-foreground"
	return simple("caption").
		Attr("data-slot", "table-caption").
		Class(baseClass)
}
