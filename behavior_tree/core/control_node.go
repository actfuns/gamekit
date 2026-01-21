package core

// ControlNode is the base class for all control nodes
type ControlNode struct {
	TreeNode
}

// NewControlNode creates a new control node
func NewControlNode(name string, config NodeConfig) ControlNode {
	return ControlNode{
		TreeNode: NewTreeNode(name, config),
	}
}

// Type returns the node type
func (cn *ControlNode) Type() NodeType {
	return NodeTypeControl
}

// ResetChildren resets all children nodes
func (cn *ControlNode) ResetChildren() {
	for _, child := range cn.Children() {
		child.HaltAndReset()
	}
}
