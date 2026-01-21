package controls

import (
	"fmt"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// SwitchNode executes a specific child based on a switch value
type SwitchNode struct {
	core.ControlNode
	runningChildIdx int
}

// NewSwitchNode creates a new switch node
func NewSwitchNode(name string, config core.NodeConfig) *SwitchNode {
	node := &SwitchNode{
		ControlNode:     core.NewControlNode(name, config),
		runningChildIdx: -1,
	}
	return node
}

// Tick executes the switch logic
func (node *SwitchNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	// Get the switch value from input
	switchValue := ""
	if value, ok := node.GetInput("switch"); ok {
		switchValue = value
	} else {
		panic("Missing required input [switch] in SwitchNode")
	}

	// Find the child with matching name
	selectedIndex := -1
	for i, child := range children {
		if child.Name() == switchValue {
			selectedIndex = i
			break
		}
	}

	if selectedIndex == -1 {
		panic(fmt.Sprintf("Can't find requested child [%s]", switchValue))
	}

	// If we have a running child, check if it's the same as selected
	if node.runningChildIdx >= 0 {
		if node.runningChildIdx == selectedIndex {
			// Continue executing the same child
			status := children[selectedIndex].ExecuteTick()
			if status != core.NodeStatusRunning {
				node.runningChildIdx = -1
			}
			return status
		} else {
			// Different child selected, halt the running child
			children[node.runningChildIdx].HaltAndReset()
			node.runningChildIdx = -1
		}
	}

	// Execute the selected child
	status := children[selectedIndex].ExecuteTick()
	if status == core.NodeStatusRunning {
		node.runningChildIdx = selectedIndex
	}
	return status
}

// Halt stops execution and resets the node
func (node *SwitchNode) Halt() {
	if node.runningChildIdx >= 0 {
		children := node.Children()
		if node.runningChildIdx < len(children) {
			children[node.runningChildIdx].HaltAndReset()
		}
		node.runningChildIdx = -1
	}
	node.ControlNode.Halt()
}
