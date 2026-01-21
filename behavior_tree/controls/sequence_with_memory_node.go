package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// SequenceWithMemoryNode executes children in sequence, remembering the last running child
type SequenceWithMemoryNode struct {
	core.ControlNode
	currentChild int
}

// NewSequenceWithMemoryNode creates a new sequence with memory node
func NewSequenceWithMemoryNode(name string, config core.NodeConfig) *SequenceWithMemoryNode {
	node := &SequenceWithMemoryNode{
		ControlNode:  core.NewControlNode(name, config),
		currentChild: 0,
	}
	return node
}

// Tick executes the sequence with memory logic
func (node *SequenceWithMemoryNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	// Start from the current child
	for i := node.currentChild; i < len(children); i++ {
		child := children[i]
		status := child.ExecuteTick()

		switch status {
		case core.NodeStatusRunning:
			// Remember the current child for next tick
			node.currentChild = i
			return core.NodeStatusRunning
		case core.NodeStatusFailure:
			// Sequence failed, reset current child and halt all children
			node.currentChild = 0
			node.ResetChildren()
			return core.NodeStatusFailure
		case core.NodeStatusSuccess:
			// Continue to next child
			continue
		}
	}

	// All children succeeded
	node.currentChild = 0
	return core.NodeStatusSuccess
}

// Halt stops execution and resets the node
func (node *SequenceWithMemoryNode) Halt() {
	node.currentChild = 0
	node.ResetChildren()
	node.ControlNode.Halt()
}
