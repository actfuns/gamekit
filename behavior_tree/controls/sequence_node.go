package controls

import "github.com/actfuns/gamekit/behavior_tree"

// SequenceNode executes children in order until one fails or all succeed
type SequenceNode struct {
	behavior_tree.ControlNode
	currentChildIdx int
	skippedCount    int
}

// NewSequenceNode creates a new sequence node
func NewSequenceNode(name string, config behavior_tree.NodeConfig) *SequenceNode {
	return &SequenceNode{
		ControlNode:     *behavior_tree.NewControlNode(name, config),
		currentChildIdx: 0,
		skippedCount:    0,
	}
}

// Tick executes the sequence logic
func (sn *SequenceNode) Tick() behavior_tree.NodeStatus {
	children := sn.Children()
	childrenCount := len(children)

	if !behavior_tree.IsStatusActive(sn.Status()) {
		sn.skippedCount = 0
	}

	sn.SetStatus(behavior_tree.NodeStatusRunning)

	for sn.currentChildIdx < childrenCount {
		currentChild := children[sn.currentChildIdx]
		childStatus := currentChild.Tick()

		switch childStatus {
		case behavior_tree.NodeStatusRunning:
			return childStatus
		case behavior_tree.NodeStatusFailure:
			sn.ResetChildren()
			sn.currentChildIdx = 0
			return childStatus
		case behavior_tree.NodeStatusSuccess:
			sn.currentChildIdx++
		case behavior_tree.NodeStatusSkipped:
			sn.currentChildIdx++
			sn.skippedCount++
		case behavior_tree.NodeStatusIdle:
			// This should not happen in normal operation
			return behavior_tree.NodeStatusFailure
		}
	}

	// All children succeeded
	allChildrenSkipped := (sn.skippedCount == childrenCount)
	if sn.currentChildIdx == childrenCount {
		sn.ResetChildren()
		sn.currentChildIdx = 0
		sn.skippedCount = 0
	}

	if allChildrenSkipped {
		return behavior_tree.NodeStatusSkipped
	}
	return behavior_tree.NodeStatusSuccess
}

// Halt handles halting the sequence node
func (sn *SequenceNode) Halt() {
	sn.currentChildIdx = 0
	sn.skippedCount = 0
	sn.ControlNode.Halt()
}
