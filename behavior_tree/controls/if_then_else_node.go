package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// IfThenElseNode is a control node that executes the first child as condition,
// then executes the second child if condition succeeds, or the third child if condition fails
type IfThenElseNode struct {
	core.ControlNode
}

// NewIfThenElseNode creates a new if-then-else node
func NewIfThenElseNode(name string, config core.NodeConfig) *IfThenElseNode {
	node := &IfThenElseNode{
		ControlNode: core.NewControlNode(name, config),
	}
	return node
}

// Tick executes the if-then-else logic
func (node *IfThenElseNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) != 3 {
		return core.NodeStatusFailure
	}

	condition := children[0]
	thenBranch := children[1]
	elseBranch := children[2]

	// Execute condition
	conditionStatus := condition.ExecuteTick()
	if conditionStatus == core.NodeStatusRunning {
		return core.NodeStatusRunning
	}

	if conditionStatus == core.NodeStatusSuccess {
		// Condition succeeded, execute then branch
		thenStatus := thenBranch.ExecuteTick()
		if thenStatus == core.NodeStatusRunning {
			return core.NodeStatusRunning
		}
		// Halt else branch since we're not using it
		elseBranch.HaltAndReset()
		return thenStatus
	} else {
		// Condition failed, execute else branch
		elseStatus := elseBranch.ExecuteTick()
		if elseStatus == core.NodeStatusRunning {
			return core.NodeStatusRunning
		}
		// Halt then branch since we're not using it
		thenBranch.HaltAndReset()
		return elseStatus
	}
}

// Halt stops execution and resets the node
func (node *IfThenElseNode) Halt() {
	for _, child := range node.Children() {
		child.HaltAndReset()
	}
	node.ControlNode.Halt()
}
