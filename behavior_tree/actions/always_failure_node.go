package actions

import "github.com/actfuns/gamekit/behavior_tree/core"

// AlwaysFailureNode 总是返回失败的动作节点
type AlwaysFailureNode struct {
	core.ActionNodeBase
}

// NewAlwaysFailureNode 创建新的AlwaysFailureNode实例
func NewAlwaysFailureNode(name string, config core.NodeConfig) *AlwaysFailureNode {
	node := &AlwaysFailureNode{}
	node.ActionNodeBase = core.NewActionNodeBase(name, config)
	return node
}

// Tick 执行节点逻辑，总是返回失败
func (afn *AlwaysFailureNode) Tick() core.NodeStatus {
	return core.NodeStatusFailure
}
