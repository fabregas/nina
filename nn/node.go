package nn

// general interface for VirtualDOM
type Node interface {
	isNode()
	getKey() string
	isNil() bool
}

type AsNode interface {
	AsNode() Node
}

func isNilNode(n Node) bool {
	return n == nil || n.isNil()
}

func If(condition bool, node AsNode) AsNode {
	if condition {
		return node
	}
	return nil
}

func IfElse(condition bool, trueNode, falseNode AsNode) AsNode {
	if condition {
		return trueNode
	}
	return falseNode
}

func List[T any](lst []T, conv func(T, int) AsNode, emptyNode AsNode) AsNode {
	if len(lst) == 0 {
		return emptyNode
	}

	nodes := make([]AsNode, 0, len(lst))

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
