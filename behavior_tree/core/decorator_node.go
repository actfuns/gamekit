package core

// DecoratorNode is the base class for all decorator nodes
type DecoratorNode struct {
	TreeNode
}

// NewDecoratorNode creates a new decorator node
func NewDecoratorNode(name string, config NodeConfig) DecoratorNode {
	return DecoratorNode{
		TreeNode: NewTreeNode(name, config),
	}
}

// Type returns the node type
func (dn *DecoratorNode) Type() NodeType {
	return NodeTypeDecorator
}
