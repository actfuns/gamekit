package actions

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// SetBlackboardNode 动作节点 - 设置黑板值
type SetBlackboardNode struct {
	behavior_tree.ActionNodeBase
	key   string
	value string
}

// NewSetBlackboardNode 创建新的SetBlackboardNode实例
func NewSetBlackboardNode(name string, config behavior_tree.NodeConfig) *SetBlackboardNode {
	key := ""
	value := ""
	
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["key"]; exists {
			key = portInfo.TypeName
		}
		if portInfo, exists := ports["value"]; exists {
			value = portInfo.TypeName
		}
	}
	
	node := &SetBlackboardNode{
		ActionNodeBase: *behavior_tree.NewActionNodeBase(name, config),
		key:            key,
		value:          value,
	}
	return node
}

// Tick 执行动作节点逻辑
func (sbn *SetBlackboardNode) Tick() behavior_tree.NodeStatus {
	blackboard := sbn.Config().Blackboard
	if blackboard == nil {
		return behavior_tree.NodeStatusFailure
	}
	
	// 使用预先存储的 key 和 value
	if sbn.key != "" {
		// TODO: 需要实现将 value 转换为适当类型并存储到黑板
		// 目前作为占位实现，直接返回SUCCESS
		return behavior_tree.NodeStatusSuccess
	}
	
	return behavior_tree.NodeStatusFailure
}