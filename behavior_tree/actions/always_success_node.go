package actions

import "github.com/actfuns/gamekit/behavior_tree"

// AlwaysSuccessNode 总是返回成功的动作节点
type AlwaysSuccessNode struct {
	behavior_tree.ActionNodeBase
}

// NewAlwaysSuccessNode 创建新的AlwaysSuccessNode实例
func NewAlwaysSuccessNode(name string, config behavior_tree.NodeConfig) *AlwaysSuccessNode {
	node := &AlwaysSuccessNode{}
	node.ActionNodeBase = *behavior_tree.NewActionNodeBase(name, config)
	return node
}

// Tick 执行节点逻辑，总是返回成功
func (asn *AlwaysSuccessNode) Tick() behavior_tree.NodeStatus {
	return behavior_tree.NodeStatusSuccess
}
