package decorators

import "github.com/actfuns/gamekit/behavior_tree"

// SubtreeNode is a decorator that wraps a subtree.
// It simply executes its child and returns the same status.
type SubtreeNode struct {
	behavior_tree.DecoratorNode
}

// NewSubtreeNode creates a new SubtreeNode
func NewSubtreeNode(name string, config behavior_tree.NodeConfig) *SubtreeNode {
	return &SubtreeNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
	}
}

// ProvidedPorts returns the ports provided by SubtreeNode
func (stn *SubtreeNode) ProvidedPorts() map[string]behavior_tree.PortInfo {
	port := behavior_tree.PortInfo{
		Direction: behavior_tree.PortDirectionInput,
		TypeName:  "bool",
	}
	port.SetDefaultValue(false)
	port.SetDescription("If true, all the ports with the same name will be remapped")

	return map[string]behavior_tree.PortInfo{
		"_autoremap": port,
	}
}

// Tick executes the subtree logic
func (stn *SubtreeNode) Tick() behavior_tree.NodeStatus {
	prevStatus := stn.Status()
	if prevStatus == behavior_tree.NodeStatusIdle {
		stn.SetStatus(behavior_tree.NodeStatusRunning)
	}

	children := stn.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	childStatus := child.Tick()
	if behavior_tree.IsStatusCompleted(childStatus) {
		child.HaltAndReset()
	}

	return childStatus
}