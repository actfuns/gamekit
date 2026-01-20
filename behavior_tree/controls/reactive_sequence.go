package controls

import "github.com/actfuns/gamekit/behavior_tree"

// ReactiveSequence 反应式序列节点
type ReactiveSequence struct {
	behavior_tree.ControlNode
}

// NewReactiveSequence 创建新的ReactiveSequence实例
func NewReactiveSequence(name string, config behavior_tree.NodeConfig) *ReactiveSequence {
	node := &ReactiveSequence{}
	node.ControlNode = *behavior_tree.NewControlNode(name, config)
	return node
}

// Tick 执行节点逻辑
func (rs *ReactiveSequence) Tick() behavior_tree.NodeStatus {
	children := rs.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusSuccess
	}

	for _, child := range children {
		status := child.Tick()
		if status != behavior_tree.NodeStatusSuccess {
			return status
		}
	}

	return behavior_tree.NodeStatusSuccess
}