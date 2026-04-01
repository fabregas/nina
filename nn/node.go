package nn

// general interface for VirtualDOM
type Node interface {
	isNode()
	getKey() string
	isNil() bool
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

func IfElse(condition bool, trueNode, falseNode Node) Node {
	if condition {
		return trueNode
	}
	return falseNode
}
