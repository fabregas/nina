package nn

type groupNode struct {
	children []Node
}

func (g *groupNode) isNode() {}

func (g *groupNode) getKey() string {
	return ""
}

func (g *groupNode) isNil() bool {
	return g == nil
}

func (g *groupNode) AsNode() Node {
	return g
}

func Group(nodes ...AsNode) AsNode {
	return &groupNode{
		children: intoNodesList(nodes),
	}
}
