package decorators

import "github.com/actfuns/gamekit/behavior_tree"

// ForceSuccessNode always returns SUCCESS when the child completes,
// regardless of whether the child returned SUCCESS or FAILURE.
// - If the child returns RUNNING, this node returns RUNNING.
// - If the child returns SUCCESS or FAILURE, this node returns SUCCESS.
type ForceSuccessNode struct {
	behavior_tree.DecoratorNode
}

// NewForceSuccessNode creates a new ForceSuccessNode
func NewForceSuccessNode(name string, config behavior_tree.NodeConfig) *ForceSuccessNode {
	return &ForceSuccessNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
	}
}

// Tick executes the force success logic
func (fsn *ForceSuccessNode) Tick() behavior_tree.NodeStatus {
	children := fsn.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	if behavior_tree.IsStatusCompleted(status) {
		child.HaltAndReset()
		return behavior_tree.NodeStatusSuccess
	}

	// RUNNING or skipping
	return status
}

// ForceFailureNode always returns FAILURE when the child completes,
// regardless of whether the child returned SUCCESS or FAILURE.
// - If the child returns RUNNING, this node returns RUNNING.
// - If the child returns SUCCESS or FAILURE, this node returns FAILURE.
type ForceFailureNode struct {
	behavior_tree.DecoratorNode
}

// NewForceFailureNode creates a new ForceFailureNode
func NewForceFailureNode(name string, config behavior_tree.NodeConfig) *ForceFailureNode {
	return &ForceFailureNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
	}
}

// Tick executes the force failure logic
func (ffn *ForceFailureNode) Tick() behavior_tree.NodeStatus {
	children := ffn.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	if behavior_tree.IsStatusCompleted(status) {
		child.HaltAndReset()
		return behavior_tree.NodeStatusFailure
	}

	// RUNNING or skipping
	return status
}
