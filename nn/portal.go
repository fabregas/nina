package nn

import "syscall/js"

// portalNode — special node that tells engine that Child should be rendered inside TargetSelector
type portalNode struct {
	targetSelector string // (for example "body" or "#portal-root")
	child          Node

	domNode         js.Value
	placeholderNode js.Value
}

func (e *portalNode) isNode() {}

func (e *portalNode) getKey() string {
	return ""
}

func (e *portalNode) isNil() bool {
	return e == nil
}

func (p *portalNode) ToNode() Node {
	return p
}

func Portal(child IntoNode) *portalNode {
	return PortalTo("body", child)
}

func PortalTo(targetSelector string, child IntoNode) *portalNode {
	var n Node
	if child != nil {
		n = child.ToNode()
	}

	return &portalNode{
		targetSelector: targetSelector,
		child:          n,
	}
}
