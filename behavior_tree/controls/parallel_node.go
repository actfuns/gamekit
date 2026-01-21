package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// ParallelNode executes all children in parallel until thresholds are reached
type ParallelNode struct {
	core.ControlNode
	successThreshold int
	failureThreshold int
	activeChildren   []bool
	completedList    map[int]bool
}

// NewParallelNode creates a new parallel node with thresholds
func NewParallelNode(name string, config core.NodeConfig, successThreshold, failureThreshold int) *ParallelNode {
	return &ParallelNode{
		ControlNode:      core.NewControlNode(name, config),
		successThreshold: successThreshold,
		failureThreshold: failureThreshold,
		activeChildren:   make([]bool, 0),
		completedList:    make(map[int]bool),
	}
}

// Tick executes the parallel logic
func (node *ParallelNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	// Initialize active children slice if needed
	if len(node.activeChildren) != len(children) {
		node.activeChildren = make([]bool, len(children))
		for i := range node.activeChildren {
			node.activeChildren[i] = true
		}
	}

	// Initialize completed list if needed
	if node.completedList == nil {
		node.completedList = make(map[int]bool)
	}

	successCount := 0
	failureCount := 0
	runningCount := 0

	// Execute all active children
	for i, child := range children {
		if !node.activeChildren[i] {
			continue
		}

		status := child.ExecuteTick()
		switch status {
		case core.NodeStatusSuccess:
			node.completedList[i] = true
			node.activeChildren[i] = false
			successCount++
		case core.NodeStatusFailure:
			node.completedList[i] = true
			node.activeChildren[i] = false
			failureCount++
		case core.NodeStatusRunning:
			runningCount++
		}
	}

	// Check termination conditions
	if successCount >= node.successThreshold {
		// Success threshold reached
		for _, child := range children {
			child.HaltAndReset()
		}
		node.resetState()
		return core.NodeStatusSuccess
	}

	if failureCount >= node.failureThreshold {
		// Failure threshold reached
		for _, child := range children {
			child.HaltAndReset()
		}
		node.resetState()
		return core.NodeStatusFailure
	}

	// Still running
	return core.NodeStatusRunning
}

// resetState resets the internal state for next execution
func (node *ParallelNode) resetState() {
	node.activeChildren = make([]bool, 0)
	node.completedList = make(map[int]bool)
}

// Halt stops execution and resets the node
func (node *ParallelNode) Halt() {
	for _, child := range node.Children() {
		child.HaltAndReset()
	}
	node.resetState()
	node.ControlNode.Halt()
}
