package actions

import "github.com/actfuns/gamekit/behavior_tree/core"

// PopFromQueue 动作节点 - 从队列中弹出元素
type PopFromQueue struct {
	core.ActionNodeBase
	queueKey  string
	outputKey string
}

// NewPopFromQueue 创建新的PopFromQueue实例
func NewPopFromQueue(name string, config core.NodeConfig) *PopFromQueue {
	// 从Manifest的Ports中获取端口信息
	queueKey := ""
	outputKey := ""

	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["queue"]; exists {
			queueKey = portInfo.TypeName
		}
		if portInfo, exists := ports["output_key"]; exists {
			outputKey = portInfo.TypeName
		}
	}

	node := &PopFromQueue{
		ActionNodeBase: core.NewActionNodeBase(name, config),
		queueKey:       queueKey,
		outputKey:      outputKey,
	}
	return node
}

// Tick 执行动作节点逻辑
func (p *PopFromQueue) Tick() core.NodeStatus {
	blackboard := p.Config().Blackboard
	if blackboard == nil {
		return core.NodeStatusFailure
	}

	// TODO: 实现从队列弹出值的逻辑
	// 目前返回SUCCESS作为占位实现
	return core.NodeStatusSuccess
}
