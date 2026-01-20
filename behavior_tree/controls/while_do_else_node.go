package controls

import "github.com/actfuns/gamekit/behavior_tree"

// WhileDoElseNode executes a condition and loops while it's true, with an else branch
type WhileDoElseNode struct {
	behavior_tree.ControlNode
	whileStatus behavior_tree.NodeStatus
}

// NewWhileDoElseNode creates a new while-do-else node
func NewWhileDoElseNode(name string, config behavior_tree.NodeConfig) *WhileDoElseNode {
	node := &WhileDoElseNode{
		ControlNode: *behavior_tree.NewControlNode(name, config),
		whileStatus: behavior_tree.NodeStatusIdle,
	}
	return node
}

// Tick executes the while-do-else logic
func (node *WhileDoElseNode) Tick() behavior_tree.NodeStatus {
	children := node.Children()
	if len(children) < 2 || len(children) > 3 {
		return behavior_tree.NodeStatusFailure
	}

	condition := children[0]
	thenBranch := children[1]
	var elseBranch behavior_tree.TreeNode = nil
	if len(children) == 3 {
		elseBranch = children[2]
	}

	// Execute condition
	conditionStatus := condition.ExecuteTick()
	if conditionStatus == behavior_tree.NodeStatusRunning {
		return behavior_tree.NodeStatusRunning
	}

	if conditionStatus == behavior_tree.NodeStatusSuccess {
		// Condition succeeded, execute then branch
		thenStatus := thenBranch.ExecuteTick()
		if thenStatus == behavior_tree.NodeStatusRunning {
			return behavior_tree.NodeStatusRunning
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
			if elseStatus == behavior_tree.NodeStatusRunning {
				return behavior_tree.NodeStatusRunning
			}
			// Halt then branch
			thenBranch.HaltAndReset()
			return elseStatus
		} else {
			// No else branch, just return failure
			thenBranch.HaltAndReset()
			return behavior_tree.NodeStatusFailure
		}
	}
}

// Halt stops execution and resets the node
func (node *WhileDoElseNode) Halt() {
	for _, child := range node.Children() {
		child.HaltAndReset()
	}
	node.whileStatus = behavior_tree.NodeStatusIdle
	node.ControlNode.Halt()
}
