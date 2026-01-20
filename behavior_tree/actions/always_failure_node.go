package actions

import "github.com/actfuns/gamekit/behavior_tree"

// AlwaysFailureNode 总是返回失败的动作节点
type AlwaysFailureNode struct {
	behavior_tree.ActionNodeBase
}

// NewAlwaysFailureNode 创建新的AlwaysFailureNode实例
func NewAlwaysFailureNode(name string, config behavior_tree.NodeConfig) *AlwaysFailureNode {
	node := &AlwaysFailureNode{}
	node.ActionNodeBase = *behavior_tree.NewActionNodeBase(name, config)
	return node
}

// Tick 执行节点逻辑，总是返回失败
func (afn *AlwaysFailureNode) Tick() behavior_tree.NodeStatus {
	return behavior_tree.NodeStatusFailure
}
