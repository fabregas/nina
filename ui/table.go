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

	table := nn.Table().
		Attr("data-slot", "table").
		Class(baseClass)

	b := &tableBuilder{}
	b.baseBuilder = base(b, table)

	return b
}

func (t *tableBuilder) build() *nn.Element {
	return t.el
}

func (t *tableBuilder) wrap(target *nn.Element) *nn.Element {
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
	return simple(
		nn.Thead().
			Attr("data-slot", "table-header").
			Class(baseClass),
	)
}

// ==========================================
// TABLE BODY
// ==========================================

func TableBody() *simpleBuilder {
	baseClass := "[&_tr:last-child]:border-0"
	return simple(
		nn.Tbody().
			Attr("data-slot", "table-body").
			Class(baseClass),
	)
}

// ==========================================
// TABLE FOOTER
// ==========================================

func TableFooter() *simpleBuilder {
	baseClass := "border-t bg-muted/50 font-medium [&>tr]:last:border-b-0"
	return simple(
		nn.Tfoot().
			Attr("data-slot", "table-footer").
			Class(baseClass),
	)
}

// ==========================================
// TABLE ROW
// ==========================================

type tableRowBuilder struct {
	baseBuilder[*tableRowBuilder]
}

func TableRow() *tableRowBuilder {
	baseClass := "border-b transition-colors hover:bg-muted/50 has-aria-expanded:bg-muted/50 data-[state=selected]:bg-muted"

	el := nn.Tr().Attr("data-slot", "table-row").Class(baseClass)

	b := &tableRowBuilder{}
	b.baseBuilder = base(b, el)

	return b
}

func (t *tableRowBuilder) build() *nn.Element { return t.el }

func (t *tableRowBuilder) Selected(selected bool) *tableRowBuilder {
	t.el.Attr("data-state", "selected")
	return t
}

// ==========================================
// TABLE HEAD (TH)
// ==========================================

func TableHead() *simpleBuilder {
	baseClass := "h-12 px-3 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0"
	return simple(
		nn.Th().
			Attr("data-slot", "table-head").
			Class(baseClass),
	)
}

// ==========================================
// TABLE CELL (TD)
// ==========================================

type tableCellBuilder struct {
	baseBuilder[*tableCellBuilder]
}

func TableCell() *tableCellBuilder {
	baseClass := "p-3 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0"
	td := nn.Td().Attr("data-slot", "table-cell").Class(baseClass)

	b := &tableCellBuilder{}
	b.baseBuilder = base(b, td)

	return b
}

func (c *tableCellBuilder) build() *nn.Element { return c.el }

func (c *tableCellBuilder) ColSpan(colSpan int) *tableCellBuilder {
	c.el.Attr("colspan", fmt.Sprintf("%d", colSpan))
	return c
}

// ==========================================
// TABLE CAPTION
// ==========================================

func TableCaption() *simpleBuilder {
	baseClass := "mt-4 text-sm text-muted-foreground"
	return simple(
		nn.Caption().
			Attr("data-slot", "table-caption").
			Class(baseClass),
	)
}
