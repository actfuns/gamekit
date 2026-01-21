package decorators

import "github.com/actfuns/gamekit/behavior_tree/core"

// SubtreeNode is a decorator that wraps a subtree.
// It simply executes its child and returns the same status.
type SubtreeNode struct {
	core.DecoratorNode
}

// NewSubtreeNode creates a new SubtreeNode
func NewSubtreeNode(name string, config core.NodeConfig) *SubtreeNode {
	return &SubtreeNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
	}
}

// ProvidedPorts returns the ports provided by SubtreeNode
func (stn *SubtreeNode) ProvidedPorts() map[string]core.PortInfo {
	port := core.PortInfo{
		Direction: core.PortDirectionInput,
		TypeName:  "bool",
	}
	port.SetDefaultValue(false)
	port.SetDescription("If true, all the ports with the same name will be remapped")

	return map[string]core.PortInfo{
		"_autoremap": port,
	}
}

// Tick executes the subtree logic
func (stn *SubtreeNode) Tick() core.NodeStatus {
	prevStatus := stn.Status()
	if prevStatus == core.NodeStatusIdle {
		stn.SetStatus(core.NodeStatusRunning)
	}

	children := stn.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]
	childStatus := child.Tick()
	if core.IsStatusCompleted(childStatus) {
		child.HaltAndReset()
	}

	return childStatus
}
