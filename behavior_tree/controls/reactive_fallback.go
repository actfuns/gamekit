package controls

import "github.com/actfuns/gamekit/behavior_tree"

// ReactiveFallback 反应式回退节点
type ReactiveFallback struct {
	behavior_tree.ControlNode
}

// NewReactiveFallback 创建新的ReactiveFallback实例
func NewReactiveFallback(name string, config behavior_tree.NodeConfig) *ReactiveFallback {
	node := &ReactiveFallback{}
	node.ControlNode = *behavior_tree.NewControlNode(name, config)
	return node
}

// Tick 执行节点逻辑
func (rf *ReactiveFallback) Tick() behavior_tree.NodeStatus {
	children := rf.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusSuccess
	}

	for _, child := range children {
		status := child.Tick()
		if status != behavior_tree.NodeStatusFailure {
			return status
		}
	}

	return behavior_tree.NodeStatusFailure
}