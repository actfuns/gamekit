package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// ReactiveSequence 反应式序列节点
type ReactiveSequence struct {
	core.ControlNode
}

// NewReactiveSequence 创建新的ReactiveSequence实例
func NewReactiveSequence(name string, config core.NodeConfig) *ReactiveSequence {
	node := &ReactiveSequence{}
	node.ControlNode = core.NewControlNode(name, config)
	return node
}

// Tick 执行节点逻辑
func (rs *ReactiveSequence) Tick() core.NodeStatus {
	children := rs.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	for _, child := range children {
		status := child.Tick()
		if status != core.NodeStatusSuccess {
			return status
		}
	}

	return core.NodeStatusSuccess
}
