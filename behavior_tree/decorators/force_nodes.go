package decorators

import "github.com/actfuns/gamekit/behavior_tree/core"

// ForceSuccessNode always returns SUCCESS when the child completes,
// regardless of whether the child returned SUCCESS or FAILURE.
// - If the child returns RUNNING, this node returns RUNNING.
// - If the child returns SUCCESS or FAILURE, this node returns SUCCESS.
type ForceSuccessNode struct {
	core.DecoratorNode
}

// NewForceSuccessNode creates a new ForceSuccessNode
func NewForceSuccessNode(name string, config core.NodeConfig) *ForceSuccessNode {
	return &ForceSuccessNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
	}
}

// Tick executes the force success logic
func (fsn *ForceSuccessNode) Tick() core.NodeStatus {
	children := fsn.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	if core.IsStatusCompleted(status) {
		child.HaltAndReset()
		return core.NodeStatusSuccess
	}

	// RUNNING or skipping
	return status
}

// ForceFailureNode always returns FAILURE when the child completes,
// regardless of whether the child returned SUCCESS or FAILURE.
// - If the child returns RUNNING, this node returns RUNNING.
// - If the child returns SUCCESS or FAILURE, this node returns FAILURE.
type ForceFailureNode struct {
	core.DecoratorNode
}

// NewForceFailureNode creates a new ForceFailureNode
func NewForceFailureNode(name string, config core.NodeConfig) *ForceFailureNode {
	return &ForceFailureNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
	}
}

// Tick executes the force failure logic
func (ffn *ForceFailureNode) Tick() core.NodeStatus {
	children := ffn.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	if core.IsStatusCompleted(status) {
		child.HaltAndReset()
		return core.NodeStatusFailure
	}

	// RUNNING or skipping
	return status
}
