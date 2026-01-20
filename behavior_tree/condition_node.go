package behavior_tree

// ConditionNode is the base class for all condition nodes
type ConditionNode struct {
	TreeNodeBase
	tickFunc func() NodeStatus
}

// NewConditionNode creates a new condition node
func NewConditionNode(name string, config NodeConfig, tickFunc func() NodeStatus) *ConditionNode {
	return &ConditionNode{
		TreeNodeBase: *NewTreeNode(name, config),
		tickFunc:     tickFunc,
	}
}

// Type returns the node type
func (cn *ConditionNode) Type() NodeType {
	return NodeTypeCondition
}

// Tick executes the condition
func (cn *ConditionNode) Tick() NodeStatus {
	if cn.tickFunc != nil {
		return cn.tickFunc()
	}
	return NodeStatusSuccess
}

// StatefulConditionNode is a condition node that can maintain state
type StatefulConditionNode struct {
	ConditionNode
	onStartFunc   func() NodeStatus
	onRunningFunc func() NodeStatus
	onHaltedFunc  func()
}

// NewStatefulConditionNode creates a new stateful condition node
func NewStatefulConditionNode(name string, config NodeConfig,
	startFunc func() NodeStatus,
	runningFunc func() NodeStatus,
	haltedFunc func()) *StatefulConditionNode {
	return &StatefulConditionNode{
		ConditionNode: *NewConditionNode(name, config, nil),
		onStartFunc:   startFunc,
		onRunningFunc: runningFunc,
		onHaltedFunc:  haltedFunc,
	}
}

// Tick executes the stateful condition
func (scn *StatefulConditionNode) Tick() NodeStatus {
	switch scn.Status() {
	case NodeStatusIdle:
		if scn.onStartFunc != nil {
			status := scn.onStartFunc()
			if status == NodeStatusRunning {
				scn.SetStatus(NodeStatusRunning)
			}
			return status
		}
		return NodeStatusSuccess

	case NodeStatusRunning:
		if scn.onRunningFunc != nil {
			return scn.onRunningFunc()
		}
		return NodeStatusSuccess

	default:
		return NodeStatusFailure
	}
}

// Halt handles halting the stateful condition
func (scn *StatefulConditionNode) Halt() {
	if scn.onHaltedFunc != nil {
		scn.onHaltedFunc()
	}
}
