package nn

import "syscall/js"

// PortalNode — special node that tells engine that Child should be rendered inside TargetSelector
type PortalNode struct {
	targetSelector string // (for example "body" or "#portal-root")
	child          Node

	domNode         js.Value
	placeholderNode js.Value
}

func (e *PortalNode) isNode() {}

func (e *PortalNode) getKey() string {
	return ""
}

func (e *PortalNode) isNil() bool {
	return e == nil
}

func (p *PortalNode) ToNode() Node {
	return p
}

func Portal(child IntoNode) *PortalNode {
	return PortalTo("body", child)
}

func PortalTo(targetSelector string, child IntoNode) *PortalNode {
	var n Node
	if child != nil {
		n = child.ToNode()
	}

	return &PortalNode{
		targetSelector: targetSelector,
		child:          n,
	}
}
