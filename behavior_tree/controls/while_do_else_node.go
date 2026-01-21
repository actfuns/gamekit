package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// WhileDoElseNode executes a condition and loops while it's true, with an else branch
type WhileDoElseNode struct {
	core.ControlNode
	whileStatus core.NodeStatus
}

// NewWhileDoElseNode creates a new while-do-else node
func NewWhileDoElseNode(name string, config core.NodeConfig) *WhileDoElseNode {
	node := &WhileDoElseNode{
		ControlNode: core.NewControlNode(name, config),
		whileStatus: core.NodeStatusIdle,
	}
	return node
}

// Tick executes the while-do-else logic
func (node *WhileDoElseNode) Tick() core.NodeStatus {
	children := node.Children()
	if len(children) < 2 || len(children) > 3 {
		return core.NodeStatusFailure
	}

	condition := children[0]
	thenBranch := children[1]
	var elseBranch core.TreeNode = nil
	if len(children) == 3 {
		elseBranch = children[2]
	}

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
		// Halt else branch if it exists
		if elseBranch != nil {
			elseBranch.HaltAndReset()
		}
		return thenStatus
	} else {
		// Condition failed, execute else branch if it exists
		if elseBranch != nil {
			elseStatus := elseBranch.ExecuteTick()
			if elseStatus == core.NodeStatusRunning {
				return core.NodeStatusRunning
			}
			// Halt then branch
			thenBranch.HaltAndReset()
			return elseStatus
		} else {
			// No else branch, just return failure
			thenBranch.HaltAndReset()
			return core.NodeStatusFailure
		}
	}
}

// Halt stops execution and resets the node
func (node *WhileDoElseNode) Halt() {
	for _, child := range node.Children() {
		child.HaltAndReset()
	}
	node.whileStatus = core.NodeStatusIdle
	node.ControlNode.Halt()
}
