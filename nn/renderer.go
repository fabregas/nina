package nn

type NativeNode interface {
	isNative()
	Equal(NativeNode) bool

	Raw() any
}

type Renderer interface {
	CreateElement(tag string) NativeNode
	CreateElementNS(ns, tag string) NativeNode
	CreateTextNode(text string) NativeNode
	CreateComment(comment string) NativeNode
	CreateDocumentFragment() NativeNode
	SetAttribute(node NativeNode, key, value string)
	RemoveAttribute(node NativeNode, key string)
	AppendChild(parent, child NativeNode)
	InsertBefore(parent, child, anchor NativeNode)
	//RemoveChild(parent, child NativeNode)
	//SetText(node NativeNode, text string)
	Remove(node NativeNode)
	AddEventListener(node NativeNode, event string, handler func(Event)) func()
	AddEventListenerWithCapture(node NativeNode, event string, handler func(Event)) func()
	AddResizeObserver(node NativeNode, handler func(Event)) func()
	SetInnerHTML(node NativeNode, html string)
	GetElementById(id string) NativeNode
	SetNodeValue(node NativeNode, val string)
	FirstChild(node NativeNode) NativeNode
	NextSibling(node NativeNode) NativeNode
	Contains(n1, n2 NativeNode) bool
	Closest(node NativeNode, selector string) NativeNode
	QuerySelector(node NativeNode, selector string) NativeNode
	QuerySelectorAll(node NativeNode, selector string) []NativeNode
	HasAttribute(node NativeNode, attr string) bool
	GetAttribute(node NativeNode, attr string) string

	GetBoundingClientRect(node NativeNode) NativeNodeRect
	GetViewport() Viewport
	ScrollIntoView(node NativeNode, options map[string]any)
	Focus(node NativeNode)

	RootNode() NativeNode
	Window() NativeNode

	PushState(path string)
	OnPopState(handler func(path string)) func()
	GetCurrentPath() string

	initRequestAnimationFrame(cb func()) (reqNext, cleaner func())
	waitNextFrame() <-chan struct{}
}

type NativeNodeRect struct {
	X, Y                     float64
	Left, Top, Right, Bottom float64
	Width, Height            float64
}

type Viewport struct {
	ScrollX float64
	ScrollY float64
	Width   float64
	Height  float64
}
