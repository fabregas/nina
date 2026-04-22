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

func If(condition bool, node IntoNode) IntoNode {
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

func List[T any](lst []T, conv func(T, int) IntoNode, emptyNode IntoNode) IntoNode {
	if len(lst) == 0 {
		return emptyNode
	}

	nodes := make([]IntoNode, 0, len(lst))

	for i, v := range lst {
		nodes = append(nodes, conv(v, i))
	}

	g := &groupNode{
		children: intoNodesList(nodes),
	}

	if len(g.children) == 0 {
		return emptyNode
	}

	return g
}
