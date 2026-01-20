package decorators

import "github.com/actfuns/gamekit/behavior_tree"

// InverterNode inverts the result of its child:
// - SUCCESS becomes FAILURE
// - FAILURE becomes SUCCESS  
// - RUNNING and SKIPPED remain unchanged
type InverterNode struct {
	behavior_tree.DecoratorNode
}

// NewInverterNode creates a new InverterNode
func NewInverterNode(name string, config behavior_tree.NodeConfig) *InverterNode {
	node := &InverterNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
	}
	return node
}

// Tick executes the inverter logic
func (in *InverterNode) Tick() behavior_tree.NodeStatus {
	children := in.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	switch status {
	case behavior_tree.NodeStatusSuccess:
		return behavior_tree.NodeStatusFailure
	case behavior_tree.NodeStatusFailure:
		return behavior_tree.NodeStatusSuccess
	default:
		// RUNNING, SKIPPED remain unchanged
		return status
	}
}