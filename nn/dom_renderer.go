//go:build js && wasm

package nn

import (
	"fmt"
	"syscall/js"
)

type domNode struct {
	js.Value
}

func (b domNode) isNative() {}

func (d domNode) Equal(node NativeNode) bool {
	if node == nil {
		return false
	}
	n := node.(domNode).Value

	return d.Value.Equal(n)
}

func (d domNode) Raw() any {
	return d.Value
}

type domRenderer struct {
	doc    js.Value
	global js.Value
}

func newDomRenderer() *domRenderer {
	return &domRenderer{
		doc:    js.Global().Get("document"),
		global: js.Global(),
	}
}
func (d *domRenderer) RootNode() NativeNode {
	return domNode{d.doc}
}

func (d *domRenderer) Window() NativeNode {
	return domNode{d.global.Get("window")}
}

func (d *domRenderer) CreateElement(tag string) NativeNode {
	val := d.doc.Call("createElement", tag)
	return domNode{val}
}
func (d *domRenderer) CreateElementNS(ns, tag string) NativeNode {
	val := d.doc.Call("createElementNS", ns, tag)
	return domNode{val}
}
func (d *domRenderer) CreateTextNode(text string) NativeNode {
	val := d.doc.Call("createTextNode", text)
	return domNode{val}
}
func (d *domRenderer) CreateComment(comment string) NativeNode {
	val := d.doc.Call("createComment", comment)
	return domNode{val}
}
func (d *domRenderer) CreateDocumentFragment() NativeNode {
	val := d.doc.Call("createDocumentFragment")
	return domNode{val}
}

func (d *domRenderer) SetAttribute(node NativeNode, key, val string) {
	domNode := node.(domNode).Value

	if domNode.IsNull() || domNode.IsUndefined() {
		fmt.Println("[DOMRenderer] Warning: attempt to setAttribute to null node")
		return
	}

	switch key {
	case "value":
		domNode.Set("value", val)
		return
	case "checked", "disabled", "readonly", "hidden", "required", "autofocus":
		isTrue := val == "true" || val == "" || val == key
		domNode.Set(key, isTrue)

		if isTrue {
			domNode.Call("setAttribute", key, "")
		} else {
			domNode.Call("removeAttribute", key)
		}
		return
	}

	domNode.Call("setAttribute", key, val)
}

func (d *domRenderer) RemoveAttribute(node NativeNode, key string) {
	domNode := node.(domNode).Value

	if domNode.IsNull() || domNode.IsUndefined() {
		fmt.Println("[DOMRenderer] Warning: attempt to removeAttribute from null node")
		return
	}

	switch key {
	case "value":
		domNode.Set("value", "")
	case "checked", "disabled", "readonly", "hidden", "required":
		domNode.Set(key, false)
	}

	domNode.Call("removeAttribute", key)
}

func (d *domRenderer) HasAttribute(node NativeNode, key string) bool {
	domNode := node.(domNode).Value
	if domNode.IsNull() || domNode.IsUndefined() {
		return false
	}

	return domNode.Call("hasAttribute", key).Bool()
}

func (d *domRenderer) GetAttribute(node NativeNode, key string) string {
	domNode := node.(domNode).Value
	if domNode.IsNull() || domNode.IsUndefined() {
		return ""
	}

	return domNode.Call("getAttribute", key).String()
}

func (d *domRenderer) AppendChild(parent, child NativeNode) {
	p := parent.(domNode).Value
	c := child.(domNode).Value

	if p.IsNull() || p.IsUndefined() || c.IsNull() || c.IsUndefined() {
		fmt.Println("[DOMRenderer] Warning: attempt to append to null node")
		return
	}

	p.Call("appendChild", c)
}

func (d *domRenderer) InsertBefore(parent, child, anchor NativeNode) {
	p := parent.(domNode).Value
	c := child.(domNode).Value

	if p.IsNull() || p.IsUndefined() || c.IsNull() || c.IsUndefined() {
		fmt.Println("[DOMRenderer] Warning: attempt to insert before to null node")
		return
	}

	if anchor == nil {
		p.Call("appendChild", c)
	} else {
		a := anchor.(domNode).Value
		p.Call("insertBefore", c, a)
	}
}

func (d *domRenderer) Remove(node NativeNode) {
	n := node.(domNode).Value

	removeFunc := n.Get("remove")

	if removeFunc.Type() == js.TypeFunction {
		n.Call("remove")
	} else {
		n.Set("textContent", "")
	}
}

func (d *domRenderer) AddEventListener(node NativeNode, event string, handler func(Event)) func() {
	var jsNode js.Value
	if node == nil {
		jsNode = d.doc
	} else {
		jsNode = node.(domNode).Value
	}

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		syntheticEvent := newDOMEvent(args[0], d)
		handler(syntheticEvent)
		return nil
	})

	jsNode.Call("addEventListener", event, cb)

	return func() {
		jsNode.Call("removeEventListener", event, cb)
		cb.Release()
	}
}

func (d *domRenderer) AddEventListenerWithCapture(node NativeNode, event string, handler func(Event)) func() {
	var jsNode js.Value
	if node == nil {
		jsNode = d.doc
	} else {
		jsNode = node.(domNode).Value
	}

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		syntheticEvent := newDOMEvent(args[0], d)
		handler(syntheticEvent)
		return nil
	})

	jsNode.Call("addEventListener", event, cb, js.ValueOf(true))

	return func() {
		jsNode.Call("removeEventListener", event, cb, js.ValueOf(true))
		cb.Release()
	}
}

// TODO fix this
// TODO global observer singleton
func (d *domRenderer) AddResizeObserver(node NativeNode, handler func(Event)) func() {
	jsNode := node.(domNode).Value

	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		syntheticEvent := newDOMEvent(args[0], d)
		handler(syntheticEvent)
		return nil
	})

	observer := d.global.Get("ResizeObserver").New(cb)
	observer.Call("observe", jsNode)

	return func() {
		observer.Call("disconnect")
		cb.Release()
	}
}

func (d *domRenderer) NextSibling(node NativeNode) NativeNode {
	jsNode := node.(domNode).Value
	next := jsNode.Get("nextSibling")

	if next.IsNull() || next.IsUndefined() {
		return nil
	}

	return domNode{next}
}

func (d *domRenderer) FirstChild(node NativeNode) NativeNode {
	jsNode := node.(domNode).Value
	first := jsNode.Get("firstChild")

	if first.IsNull() || first.IsUndefined() {
		return nil
	}

	return domNode{first}
}

func (d *domRenderer) SetInnerHTML(node NativeNode, html string) {
	jsNode := node.(domNode).Value
	jsNode.Set("innerHTML", html)
}

func (d *domRenderer) SetNodeValue(node NativeNode, val string) {
	jsNode := node.(domNode).Value
	jsNode.Set("nodeValue", val)
}

func (d *domRenderer) GetElementById(id string) NativeNode {
	node := d.doc.Call("getElementById", id)

	return d.toNativeNode(node)
}

func (d *domRenderer) toNativeNode(v js.Value) NativeNode {
	if v.IsNull() || v.IsUndefined() {
		return nil
	}
	return domNode{v}
}

func (d *domRenderer) Contains(node1, node2 NativeNode) bool {
	jsNode1 := node1.(domNode).Value
	if jsNode1.IsNull() || jsNode1.IsUndefined() {
		return false
	}
	jsNode2 := node2.(domNode).Value
	if jsNode2.IsNull() || jsNode2.IsUndefined() {
		return false
	}

	return jsNode1.Call("contains", jsNode2).Bool()
}

func (d *domRenderer) Closest(node NativeNode, selector string) NativeNode {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return nil
	}

	ret := jsNode.Call("closest", selector)

	return d.toNativeNode(ret)
}

func (d *domRenderer) QuerySelector(node NativeNode, selector string) NativeNode {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return nil
	}

	ret := jsNode.Call("querySelector", selector)

	return d.toNativeNode(ret)
}

func (d *domRenderer) QuerySelectorAll(node NativeNode, selector string) []NativeNode {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return nil
	}

	items := jsNode.Call("querySelectorAll", selector)
	length := items.Get("length").Int()
	if length == 0 {
		return nil
	}

	var ret []NativeNode
	for i := 0; i < length; i++ {
		ret = append(ret, d.toNativeNode(items.Index(i)))
	}

	return ret
}

func (d *domRenderer) ScrollIntoView(node NativeNode, options map[string]any) {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return
	}

	jsNode.Call("scrollIntoView", options)
}

func (d *domRenderer) Focus(node NativeNode) {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return
	}

	jsNode.Call("focus")
}

func (d *domRenderer) GetBoundingClientRect(node NativeNode) NativeNodeRect {
	jsNode := node.(domNode).Value
	if jsNode.IsNull() || jsNode.IsUndefined() {
		return NativeNodeRect{}
	}

	rect := jsNode.Call("getBoundingClientRect")

	return NativeNodeRect{
		X:      rect.Get("x").Float(),
		Y:      rect.Get("y").Float(),
		Left:   rect.Get("left").Float(),
		Top:    rect.Get("top").Float(),
		Right:  rect.Get("right").Float(),
		Bottom: rect.Get("bottom").Float(),
		Width:  rect.Get("width").Float(),
		Height: rect.Get("height").Float(),
	}
}

func (d *domRenderer) GetViewport() Viewport {
	window := d.global.Get("window")

	return Viewport{
		ScrollX: window.Get("scrollX").Float(),
		ScrollY: window.Get("scrollY").Float(),
		Width:   window.Get("innerWidth").Float(),
		Height:  window.Get("innerHeight").Float(),
	}
}

func (d *domRenderer) ToggleHTMLClass(class string) {
	classList := d.doc.Get("documentElement").Get("classList")
	classList.Call("toggle", class)
}

func (d *domRenderer) initRequestAnimationFrame(cb func()) (reqNext func(), cleaner func()) {
	cbFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		cb()
		return nil
	})

	reqNext = func() {
		d.global.Call("requestAnimationFrame", cbFunc)
	}

	cleaner = func() {
		cbFunc.Release()
	}

	return
}

func (d *domRenderer) waitNextFrame() <-chan struct{} {
	ch := make(chan struct{})

	var cb js.Func

	cb = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer cb.Release()

		close(ch)
		return nil
	})

	d.global.Call("requestAnimationFrame", cb)

	return ch
}

func (d *domRenderer) PushState(path string) {
	d.global.Get("history").Call("pushState", nil, "", path)
}

func (d *domRenderer) OnPopState(handler func(path string)) func() {
	cb := js.FuncOf(func(this js.Value, args []js.Value) any {
		loc := js.Global().Get("window").Get("location")
		fullPath := loc.Get("pathname").String() + loc.Get("search").String()

		handler(fullPath)
		return nil
	})

	window := d.global.Get("window")
	window.Call("addEventListener", "popstate", cb)

	return func() {
		window.Call("removeEventListener", "popstate", cb)
		cb.Release()
	}
}

func (d *domRenderer) GetCurrentPath() string {
	loc := d.global.Get("window").Get("location")
	return loc.Get("pathname").String() + loc.Get("search").String()
}
