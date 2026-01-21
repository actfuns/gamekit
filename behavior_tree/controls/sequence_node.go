package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// SequenceNode executes children in order until one fails or all succeed
type SequenceNode struct {
	core.ControlNode
	currentChildIdx int
	skippedCount    int
}

// NewSequenceNode creates a new sequence node
func NewSequenceNode(name string, config core.NodeConfig) *SequenceNode {
	return &SequenceNode{
		ControlNode:     core.NewControlNode(name, config),
		currentChildIdx: 0,
		skippedCount:    0,
	}
}

// Tick executes the sequence logic
func (sn *SequenceNode) Tick() core.NodeStatus {
	children := sn.Children()
	childrenCount := len(children)

	if !core.IsStatusActive(sn.Status()) {
		sn.skippedCount = 0
	}

	sn.SetStatus(core.NodeStatusRunning)

	for sn.currentChildIdx < childrenCount {
		currentChild := children[sn.currentChildIdx]
		childStatus := currentChild.Tick()

		switch childStatus {
		case core.NodeStatusRunning:
			return childStatus
		case core.NodeStatusFailure:
			sn.ResetChildren()
			sn.currentChildIdx = 0
			return childStatus
		case core.NodeStatusSuccess:
			sn.currentChildIdx++
		case core.NodeStatusSkipped:
			sn.currentChildIdx++
			sn.skippedCount++
		case core.NodeStatusIdle:
			// This should not happen in normal operation
			return core.NodeStatusFailure
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
		return core.NodeStatusSkipped
	}
	return core.NodeStatusSuccess
}

// Halt handles halting the sequence node
func (sn *SequenceNode) Halt() {
	sn.currentChildIdx = 0
	sn.skippedCount = 0
	sn.ControlNode.Halt()
}
