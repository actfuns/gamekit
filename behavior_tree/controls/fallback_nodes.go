package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// FallbackNode executes children in order until one succeeds or all fail
type FallbackNode struct {
	core.ControlNode
	currentChildIdx int
	asynch          bool
	skippedCount    int
}

// NewFallbackNode creates a new FallbackNode
func NewFallbackNode(name string, config core.NodeConfig, asynch bool) *FallbackNode {
	return &FallbackNode{
		ControlNode:     core.NewControlNode(name, config),
		currentChildIdx: 0,
		asynch:          asynch,
		skippedCount:    0,
	}
}

// Tick executes the fallback logic
func (fn *FallbackNode) Tick() core.NodeStatus {
	children := fn.Children()
	childrenCount := len(children)

	if !core.IsStatusActive(fn.Status()) {
		fn.skippedCount = 0
	}

	fn.SetStatus(core.NodeStatusRunning)

	for fn.currentChildIdx < childrenCount {
		currentChild := children[fn.currentChildIdx]
		prevStatus := currentChild.Status()
		childStatus := currentChild.Tick()

		switch childStatus {
		case core.NodeStatusRunning:
			return childStatus
		case core.NodeStatusSuccess:
			fn.ResetChildren()
			fn.currentChildIdx = 0
			return childStatus
		case core.NodeStatusFailure:
			fn.currentChildIdx++
			// For async mode, return running to make it interruptible
			if fn.asynch && fn.RequiresWakeUp() && prevStatus == core.NodeStatusIdle &&
				fn.currentChildIdx < childrenCount {
				// In Go implementation, we don't have emitWakeUpSignal, so we just return running
				return core.NodeStatusRunning
			}
		case core.NodeStatusSkipped:
			fn.currentChildIdx++
			fn.skippedCount++
		case core.NodeStatusIdle:
			// This should not happen in normal operation
			return core.NodeStatusFailure
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
		return core.NodeStatusSkipped
	}
	return core.NodeStatusFailure
}

// Halt handles halting the fallback node
func (fn *FallbackNode) Halt() {
	fn.currentChildIdx = 0
	fn.skippedCount = 0
	fn.ControlNode.Halt()
}

// ReactiveFallbackNode executes children reactively (resets previous children on each tick)
type ReactiveFallbackNode struct {
	core.ControlNode
	runningChild int
}

// NewReactiveFallbackNode creates a new ReactiveFallbackNode
func NewReactiveFallbackNode(name string, config core.NodeConfig) *ReactiveFallbackNode {
	return &ReactiveFallbackNode{
		ControlNode:  core.NewControlNode(name, config),
		runningChild: -1,
	}
}

// Tick executes the reactive fallback logic
func (rfn *ReactiveFallbackNode) Tick() core.NodeStatus {
	children := rfn.Children()
	allSkipped := true

	if rfn.Status() == core.NodeStatusIdle {
		rfn.runningChild = -1
	}

	rfn.SetStatus(core.NodeStatusRunning)

	for index := 0; index < len(children); index++ {
		currentChild := children[index]
		childStatus := currentChild.Tick()

		allSkipped = allSkipped && (childStatus == core.NodeStatusSkipped)

		switch childStatus {
		case core.NodeStatusRunning:
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
			return core.NodeStatusRunning
		case core.NodeStatusFailure:
			// Continue to next child
		case core.NodeStatusSuccess:
			rfn.ResetChildren()
			return core.NodeStatusSuccess
		case core.NodeStatusSkipped:
			// Reset the node to allow it to be skipped again
			currentChild.HaltAndReset()
		case core.NodeStatusIdle:
			// This should not happen in normal operation
			return core.NodeStatusFailure
		}
	}

	rfn.ResetChildren()

	if allSkipped {
		return core.NodeStatusSkipped
	}
	return core.NodeStatusFailure
}

// Halt handles halting the reactive fallback node
func (rfn *ReactiveFallbackNode) Halt() {
	rfn.runningChild = -1
	rfn.ControlNode.Halt()
}
