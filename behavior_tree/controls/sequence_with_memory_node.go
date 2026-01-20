package controls

import "github.com/actfuns/gamekit/behavior_tree"

// SequenceWithMemoryNode executes children in sequence, remembering the last running child
type SequenceWithMemoryNode struct {
	behavior_tree.ControlNode
	currentChild int
}

// NewSequenceWithMemoryNode creates a new sequence with memory node
func NewSequenceWithMemoryNode(name string, config behavior_tree.NodeConfig) *SequenceWithMemoryNode {
	node := &SequenceWithMemoryNode{
		ControlNode:  *behavior_tree.NewControlNode(name, config),
		currentChild: 0,
	}
	return node
}

// Tick executes the sequence with memory logic
func (node *SequenceWithMemoryNode) Tick() behavior_tree.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusSuccess
	}

	// Start from the current child
	for i := node.currentChild; i < len(children); i++ {
		child := children[i]
		status := child.ExecuteTick()

		switch status {
		case behavior_tree.NodeStatusRunning:
			// Remember the current child for next tick
			node.currentChild = i
			return behavior_tree.NodeStatusRunning
		case behavior_tree.NodeStatusFailure:
			// Sequence failed, reset current child and halt all children
			node.currentChild = 0
			node.ResetChildren()
			return behavior_tree.NodeStatusFailure
		case behavior_tree.NodeStatusSuccess:
			// Continue to next child
			continue
		}
	}

	// All children succeeded
	node.currentChild = 0
	return behavior_tree.NodeStatusSuccess
}

// Halt stops execution and resets the node
func (node *SequenceWithMemoryNode) Halt() {
	node.currentChild = 0
	node.ResetChildren()
	node.ControlNode.Halt()
}