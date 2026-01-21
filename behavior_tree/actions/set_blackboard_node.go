package actions

import "github.com/actfuns/gamekit/behavior_tree/core"

// SetBlackboardNode 动作节点 - 设置黑板值
type SetBlackboardNode struct {
	core.ActionNodeBase
	key   string
	value string
}

// NewSetBlackboardNode 创建新的SetBlackboardNode实例
func NewSetBlackboardNode(name string, config core.NodeConfig) *SetBlackboardNode {
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
		ActionNodeBase: core.NewActionNodeBase(name, config),
		key:            key,
		value:          value,
	}
	return node
}

// Tick 执行动作节点逻辑
func (sbn *SetBlackboardNode) Tick() core.NodeStatus {
	blackboard := sbn.Config().Blackboard
	if blackboard == nil {
		return core.NodeStatusFailure
	}

	// 使用预先存储的 key 和 value
	if sbn.key != "" {
		// TODO: 需要实现将 value 转换为适当类型并存储到黑板
		// 目前作为占位实现，直接返回SUCCESS
		return core.NodeStatusSuccess
	}

	return core.NodeStatusFailure
}
