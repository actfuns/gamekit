package behavior_tree

// ActionNodeBase is the base class for all action nodes
type ActionNodeBase struct {
	TreeNodeBase
}

// NewActionNodeBase creates a new action node base
func NewActionNodeBase(name string, config NodeConfig) *ActionNodeBase {
	return &ActionNodeBase{
		TreeNodeBase: *NewTreeNode(name, config),
	}
}

// Type returns the node type
func (an *ActionNodeBase) Type() NodeType {
	return NodeTypeAction
}

// ActionNode is a simple action node that executes once
type ActionNode struct {
	ActionNodeBase
	tickFunc func() NodeStatus
}

// NewActionNode creates a new action node with a tick function
func NewActionNode(name string, config NodeConfig, tickFunc func() NodeStatus) *ActionNode {
	return &ActionNode{
		ActionNodeBase: *NewActionNodeBase(name, config),
		tickFunc:       tickFunc,
	}
}

// Tick executes the action
func (an *ActionNode) Tick() NodeStatus {
	if an.tickFunc != nil {
		return an.tickFunc()
	}
	return NodeStatusSuccess
}

// AsyncActionNode is an action node that can run asynchronously
type AsyncActionNode struct {
	ActionNodeBase
	tickFunc func() NodeStatus
}

// NewAsyncActionNode creates a new async action node
func NewAsyncActionNode(name string, config NodeConfig, tickFunc func() NodeStatus) *AsyncActionNode {
	return &AsyncActionNode{
		ActionNodeBase: *NewActionNodeBase(name, config),
		tickFunc:       tickFunc,
	}
}

// Tick executes the async action
func (aan *AsyncActionNode) Tick() NodeStatus {
	if aan.tickFunc != nil {
		return aan.tickFunc()
	}
	return NodeStatusSuccess
}

// StatefulActionNode is an action node that maintains state between ticks
type StatefulActionNode struct {
	ActionNodeBase
	onStartFunc   func() NodeStatus
	onRunningFunc func() NodeStatus
	onHaltedFunc  func()
}

// NewStatefulActionNode creates a new stateful action node
func NewStatefulActionNode(name string, config NodeConfig,
	startFunc func() NodeStatus,
	runningFunc func() NodeStatus,
	haltedFunc func()) *StatefulActionNode {
	return &StatefulActionNode{
		ActionNodeBase: *NewActionNodeBase(name, config),
		onStartFunc:    startFunc,
		onRunningFunc:  runningFunc,
		onHaltedFunc:   haltedFunc,
	}
}

// Tick executes the stateful action
func (san *StatefulActionNode) Tick() NodeStatus {
	switch san.Status() {
	case NodeStatusIdle:
		if san.onStartFunc != nil {
			status := san.onStartFunc()
			if status == NodeStatusRunning {
				san.SetStatus(NodeStatusRunning)
			}
			return status
		}
		return NodeStatusSuccess

	case NodeStatusRunning:
		if san.onRunningFunc != nil {
			return san.onRunningFunc()
		}
		return NodeStatusSuccess

	default:
		return NodeStatusFailure
	}
}

// Halt handles halting the stateful action
func (san *StatefulActionNode) Halt() {
	if san.onHaltedFunc != nil {
		san.onHaltedFunc()
	}
}
