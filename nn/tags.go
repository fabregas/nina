package nn

// Tag створює новий HTML елемент із вказаним ім'ям
func Tag(name string) *Element {
	return &Element{
		tag: name,
	}
}

// ==========================================
// (Layout & Semantics)
// ==========================================
func Div() *Element     { return Tag("div") }
func Header() *Element  { return Tag("header") }
func Footer() *Element  { return Tag("footer") }
func Main() *Element    { return Tag("main") }
func Section() *Element { return Tag("section") }
func Article() *Element { return Tag("article") }
func Nav() *Element     { return Tag("nav") }
func Aside() *Element   { return Tag("aside") }

// ==========================================
// (Text Content)
// ==========================================
func P() *Element          { return Tag("p") }
func Span() *Element       { return Tag("span") }
func H1() *Element         { return Tag("h1") }
func H2() *Element         { return Tag("h2") }
func H3() *Element         { return Tag("h3") }
func H4() *Element         { return Tag("h4") }
func H5() *Element         { return Tag("h5") }
func H6() *Element         { return Tag("h6") }
func Strong() *Element     { return Tag("strong") }
func Em() *Element         { return Tag("em") }
func Br() *Element         { return Tag("br") }
func Hr() *Element         { return Tag("hr") }
func Pre() *Element        { return Tag("pre") }
func Code() *Element       { return Tag("code") }
func Blockquote() *Element { return Tag("blockquote") }

// ==========================================
// (Lists)
// ==========================================
func Ul() *Element { return Tag("ul") }
func Ol() *Element { return Tag("ol") }
func Li() *Element { return Tag("li") }

// ==========================================
// (Links & Media)
// ==========================================
func A() *Element      { return Tag("a") }
func Img() *Element    { return Tag("img") }
func Audio() *Element  { return Tag("audio") }
func Video() *Element  { return Tag("video") }
func Source() *Element { return Tag("source") }
func Iframe() *Element { return Tag("iframe") }

// ==========================================
// (Forms)
// ==========================================
func Form() *Element     { return Tag("form") }
func Input() *Element    { return Tag("input") }
func Textarea() *Element { return Tag("textarea") }
func Button() *Element   { return Tag("button") }
func Select() *Element   { return Tag("select") }
func Option() *Element   { return Tag("option") }
func Label() *Element    { return Tag("label") }
func Fieldset() *Element { return Tag("fieldset") }
func Legend() *Element   { return Tag("legend") }

// ==========================================
// (Tables)
// ==========================================
func Table() *Element { return Tag("table") }
func Thead() *Element { return Tag("thead") }
func Tbody() *Element { return Tag("tbody") }
func Tfoot() *Element { return Tag("tfoot") }
func Tr() *Element    { return Tag("tr") }
func Th() *Element    { return Tag("th") }
func Td() *Element    { return Tag("td") }

func Text(v string) *TextNode {
	return &TextNode{value: v}
}
