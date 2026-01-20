package controls

import "github.com/actfuns/gamekit/behavior_tree"

// FallbackNode executes children in order until one succeeds or all fail
type FallbackNode struct {
	behavior_tree.ControlNode
	currentChildIdx int
	asynch          bool
	skippedCount    int
}

// NewFallbackNode creates a new FallbackNode
func NewFallbackNode(name string, config behavior_tree.NodeConfig, asynch bool) *FallbackNode {
	return &FallbackNode{
		ControlNode:     *behavior_tree.NewControlNode(name, config),
		currentChildIdx: 0,
		asynch:          asynch,
		skippedCount:    0,
	}
}

// Tick executes the fallback logic
func (fn *FallbackNode) Tick() behavior_tree.NodeStatus {
	children := fn.Children()
	childrenCount := len(children)

	if !behavior_tree.IsStatusActive(fn.Status()) {
		fn.skippedCount = 0
	}

	fn.SetStatus(behavior_tree.NodeStatusRunning)

	for fn.currentChildIdx < childrenCount {
		currentChild := children[fn.currentChildIdx]
		prevStatus := currentChild.Status()
		childStatus := currentChild.Tick()

		switch childStatus {
		case behavior_tree.NodeStatusRunning:
			return childStatus
		case behavior_tree.NodeStatusSuccess:
			fn.ResetChildren()
			fn.currentChildIdx = 0
			return childStatus
		case behavior_tree.NodeStatusFailure:
			fn.currentChildIdx++
			// For async mode, return running to make it interruptible
			if fn.asynch && fn.RequiresWakeUp() && prevStatus == behavior_tree.NodeStatusIdle &&
				fn.currentChildIdx < childrenCount {
				// In Go implementation, we don't have emitWakeUpSignal, so we just return running
				return behavior_tree.NodeStatusRunning
			}
		case behavior_tree.NodeStatusSkipped:
			fn.currentChildIdx++
			fn.skippedCount++
		case behavior_tree.NodeStatusIdle:
			// This should not happen in normal operation
			return behavior_tree.NodeStatusFailure
		}
	}

	// All children failed
	allChildrenSkipped := (fn.skippedCount == childrenCount)
	if fn.currentChildIdx == childrenCount {
		fn.ResetChildren()
		fn.currentChildIdx = 0
		fn.skippedCount = 0
	}

	if allChildrenSkipped {
		return behavior_tree.NodeStatusSkipped
	}
	return behavior_tree.NodeStatusFailure
}

// Halt handles halting the fallback node
func (fn *FallbackNode) Halt() {
	fn.currentChildIdx = 0
	fn.skippedCount = 0
	fn.ControlNode.Halt()
}

// ReactiveFallbackNode executes children reactively (resets previous children on each tick)
type ReactiveFallbackNode struct {
	behavior_tree.ControlNode
	runningChild int
}

// NewReactiveFallbackNode creates a new ReactiveFallbackNode
func NewReactiveFallbackNode(name string, config behavior_tree.NodeConfig) *ReactiveFallbackNode {
	return &ReactiveFallbackNode{
		ControlNode:  *behavior_tree.NewControlNode(name, config),
		runningChild: -1,
	}
}

// Tick executes the reactive fallback logic
func (rfn *ReactiveFallbackNode) Tick() behavior_tree.NodeStatus {
	children := rfn.Children()
	allSkipped := true

	if rfn.Status() == behavior_tree.NodeStatusIdle {
		rfn.runningChild = -1
	}

	rfn.SetStatus(behavior_tree.NodeStatusRunning)

	for index := 0; index < len(children); index++ {
		currentChild := children[index]
		childStatus := currentChild.Tick()

		allSkipped = allSkipped && (childStatus == behavior_tree.NodeStatusSkipped)

		switch childStatus {
		case behavior_tree.NodeStatusRunning:
			// Reset previous children to ensure they are in IDLE state next time
			for i := 0; i < len(children); i++ {
				if i != index {
					children[i].HaltAndReset()
				}
			}
			if rfn.runningChild == -1 {
				rfn.runningChild = index
			}
			// In Go implementation, we don't throw exceptions for multiple running children
			return behavior_tree.NodeStatusRunning
		case behavior_tree.NodeStatusFailure:
			// Continue to next child
		case behavior_tree.NodeStatusSuccess:
			rfn.ResetChildren()
			return behavior_tree.NodeStatusSuccess
		case behavior_tree.NodeStatusSkipped:
			// Reset the node to allow it to be skipped again
			currentChild.HaltAndReset()
		case behavior_tree.NodeStatusIdle:
			// This should not happen in normal operation
			return behavior_tree.NodeStatusFailure
		}
	}

	rfn.ResetChildren()

	if allSkipped {
		return behavior_tree.NodeStatusSkipped
	}
	return behavior_tree.NodeStatusFailure
}

// Halt handles halting the reactive fallback node
func (rfn *ReactiveFallbackNode) Halt() {
	rfn.runningChild = -1
	rfn.ControlNode.Halt()
}
