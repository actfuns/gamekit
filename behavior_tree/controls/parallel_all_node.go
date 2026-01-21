package controls

import (
	"fmt"
	"strconv"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// ParallelAllNode executes all children in parallel and succeeds only when all succeed
type ParallelAllNode struct {
	core.ControlNode
	failureThreshold int
	completedList    map[int]bool
}

// NewParallelAllNode creates a new parallel all node
func NewParallelAllNode(name string, config core.NodeConfig) *ParallelAllNode {
	node := &ParallelAllNode{
		ControlNode:      core.NewControlNode(name, config),
		failureThreshold: 1,
		completedList:    make(map[int]bool),
	}
	return node
}

// SetFailureThreshold sets the failure threshold
func (node *ParallelAllNode) SetFailureThreshold(threshold int) {
	node.failureThreshold = threshold
}

// Tick executes the parallel all logic
func (node *ParallelAllNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	// Get max_failures parameter
	maxFailures := 0
	if value, ok := node.GetInput("max_failures"); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			maxFailures = intValue
		} else {
			panic(fmt.Sprintf("Invalid parameter [max_failures] in ParallelAllNode: %v", err))
		}
	} else {
		panic("Missing parameter [max_failures] in ParallelAllNode")
	}
	node.SetFailureThreshold(maxFailures)

	if len(children) < node.failureThreshold {
		panic("Number of children is less than threshold. Can never fail.")
	}

	// Initialize completed list if empty
	if node.completedList == nil {
		node.completedList = make(map[int]bool)
	}

	// Execute all children that haven't completed yet
	failureCount := 0
	successCount := 0
	runningCount := 0

	for i, child := range children {
		if node.completedList[i] {
			// Child already completed, check its status
			childStatus := child.Status()
			switch childStatus {
			case core.NodeStatusSuccess:
				successCount++
			case core.NodeStatusFailure:
				failureCount++
			}
			continue
		}

		// Execute child
		status := child.ExecuteTick()
		switch status {
		case core.NodeStatusSuccess:
			node.completedList[i] = true
			successCount++
		case core.NodeStatusFailure:
			node.completedList[i] = true
			failureCount++
		case core.NodeStatusRunning:
			runningCount++
		}
	}

	// Check termination conditions
	if failureCount >= node.failureThreshold {
		// Too many failures, halt all children and return failure
		for _, child := range children {
			child.HaltAndReset()
		}
		node.completedList = make(map[int]bool) // Reset for next execution
		return core.NodeStatusFailure
	}

	if successCount == len(children) {
		// All children succeeded
		node.completedList = make(map[int]bool) // Reset for next execution
		return core.NodeStatusSuccess
	}

	// Still running
	return core.NodeStatusRunning
}

// Halt stops execution and resets the node
func (node *ParallelAllNode) Halt() {
	for _, child := range node.Children() {
		child.HaltAndReset()
	}
	node.completedList = make(map[int]bool)
	node.ControlNode.Halt()
}
