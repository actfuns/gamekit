package controls

import (
	"strconv"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// ManualSelectorNode allows manual selection of which child to execute
type ManualSelectorNode struct {
	core.ControlNode
	runningChildIdx       int
	previouslyExecutedIdx int
	repeatLastSelection   bool
	selectedChildIndex    int
}

// NewManualSelectorNode creates a new manual selector node
func NewManualSelectorNode(name string, config core.NodeConfig) *ManualSelectorNode {
	node := &ManualSelectorNode{
		ControlNode:           core.NewControlNode(name, config),
		runningChildIdx:       -1,
		previouslyExecutedIdx: -1,
		repeatLastSelection:   false,
		selectedChildIndex:    -1,
	}
	return node
}

// Tick executes the manual selection logic
func (node *ManualSelectorNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return node.selectStatus()
	}

	// Get input for repeat last selection
	repeatLast := false
	if value, ok := node.GetInput("REPEAT_LAST_SELECTION"); ok {
		// 尝试转换为bool值
		if boolValue, err := strconv.ParseBool(value); err == nil {
			repeatLast = boolValue
		}
	}
	node.repeatLastSelection = repeatLast

	// Get selected child index from input
	selectedIndex := -1
	if value, ok := node.GetInput("SELECTED_CHILD_INDEX"); ok {
		// 尝试转换为int值
		if intValue, err := strconv.Atoi(value); err == nil {
			selectedIndex = intValue
		}
	}
	if selectedIndex >= 0 && selectedIndex < len(children) {
		node.selectedChildIndex = selectedIndex
	}

	// If no child is selected and we should repeat last selection
	if node.selectedChildIndex == -1 && node.repeatLastSelection && node.previouslyExecutedIdx >= 0 {
		node.selectedChildIndex = node.previouslyExecutedIdx
	}

	// If still no child selected, return running
	if node.selectedChildIndex == -1 {
		return core.NodeStatusRunning
	}

	// Execute the selected child
	selectedChild := children[node.selectedChildIndex]
	status := selectedChild.ExecuteTick()

	if status == core.NodeStatusRunning {
		node.runningChildIdx = node.selectedChildIndex
		return core.NodeStatusRunning
	}

	// Child completed
	node.runningChildIdx = -1
	node.previouslyExecutedIdx = node.selectedChildIndex
	node.selectedChildIndex = -1 // Reset for next tick

	return status
}

// Halt stops execution and resets the node
func (node *ManualSelectorNode) Halt() {
	if node.runningChildIdx >= 0 {
		children := node.Children()
		if node.runningChildIdx < len(children) {
			children[node.runningChildIdx].HaltAndReset()
		}
		node.runningChildIdx = -1
	}
	node.ControlNode.Halt()
}

// selectStatus handles the case when there are no children
func (node *ManualSelectorNode) selectStatus() core.NodeStatus {
	// In Go implementation, we'll return success by default when no children
	// The actual manual selection would be handled through inputs
	return core.NodeStatusSuccess
}
