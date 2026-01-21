package actions

import "github.com/actfuns/gamekit/behavior_tree/core"

// AlwaysSuccessNode 总是返回成功的动作节点
type AlwaysSuccessNode struct {
	core.ActionNodeBase
}

// NewAlwaysSuccessNode 创建新的AlwaysSuccessNode实例
func NewAlwaysSuccessNode(name string, config core.NodeConfig) *AlwaysSuccessNode {
	node := &AlwaysSuccessNode{}
	node.ActionNodeBase = core.NewActionNodeBase(name, config)
	return node
}

// Tick 执行节点逻辑，总是返回成功
func (asn *AlwaysSuccessNode) Tick() core.NodeStatus {
	return core.NodeStatusSuccess
}
