package controls

import "github.com/actfuns/gamekit/behavior_tree/core"

// ReactiveFallback 反应式回退节点
type ReactiveFallback struct {
	core.ControlNode
}

// NewReactiveFallback 创建新的ReactiveFallback实例
func NewReactiveFallback(name string, config core.NodeConfig) *ReactiveFallback {
	node := &ReactiveFallback{}
	node.ControlNode = core.NewControlNode(name, config)
	return node
}

// Tick 执行节点逻辑
func (rf *ReactiveFallback) Tick() core.NodeStatus {
	children := rf.Children()
	if len(children) == 0 {
		return core.NodeStatusSuccess
	}

	for _, child := range children {
		status := child.Tick()
		if status != core.NodeStatusFailure {
			return status
		}
	}

	return core.NodeStatusFailure
}
