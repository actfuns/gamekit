package actions

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// UnsetBlackboardNode 动作节点 - 删除黑板值
type UnsetBlackboardNode struct {
	*behavior_tree.ActionNodeBase
	key string
}

// NewUnsetBlackboardNode 创建新的UnsetBlackboardNode实例
func NewUnsetBlackboardNode(name string, config behavior_tree.NodeConfig) *UnsetBlackboardNode {
	key := ""
	
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["key"]; exists {
			key = portInfo.TypeName
		}
	}
	
	node := &UnsetBlackboardNode{
		ActionNodeBase: behavior_tree.NewActionNodeBase(name, config),
		key:            key,
	}
	return node
}

// Tick 执行动作节点逻辑
func (ubn *UnsetBlackboardNode) Tick() behavior_tree.NodeStatus {
	blackboard := ubn.ActionNodeBase.Config().Blackboard
	if blackboard == nil {
		return behavior_tree.NodeStatusFailure
	}
	
	// TODO: 实现删除黑板值的逻辑
	// 目前返回SUCCESS作为占位实现
	return behavior_tree.NodeStatusSuccess
}