package decorators

import "github.com/actfuns/gamekit/behavior_tree/core"

// InverterNode inverts the result of its child:
// - SUCCESS becomes FAILURE
// - FAILURE becomes SUCCESS
// - RUNNING and SKIPPED remain unchanged
type InverterNode struct {
	core.DecoratorNode
}

// NewInverterNode creates a new InverterNode
func NewInverterNode(name string, config core.NodeConfig) *InverterNode {
	node := &InverterNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
	}
	return node
}

// Tick executes the inverter logic
func (in *InverterNode) Tick() core.NodeStatus {
	children := in.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	switch status {
	case core.NodeStatusSuccess:
		return core.NodeStatusFailure
	case core.NodeStatusFailure:
		return core.NodeStatusSuccess
	default:
		// RUNNING, SKIPPED remain unchanged
		return status
	}
}
