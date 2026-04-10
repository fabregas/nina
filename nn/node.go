package nn

// general interface for VirtualDOM
type Node interface {
	isNode()
	getKey() string
	isNil() bool
}

type IntoNode interface {
	ToNode() Node
}

func isNilNode(n Node) bool {
	return n == nil || n.isNil()
}

func If(condition bool, node Node) Node {
	if condition {
		return node
	}
	return nil
}

func IfElse(condition bool, trueNode, falseNode IntoNode) IntoNode {
	if condition {
		return trueNode
	}
	return falseNode
}
