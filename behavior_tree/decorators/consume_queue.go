package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// ConsumeQueue 装饰器节点 - 消费队列
type ConsumeQueue struct {
	behavior_tree.DecoratorNode
	queueKey string
}

// NewConsumeQueue 创建新的ConsumeQueue实例
func NewConsumeQueue(name string, config behavior_tree.NodeConfig) *ConsumeQueue {
	queueKey := ""
	
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["queue"]; exists {
			queueKey = portInfo.TypeName
		}
	}
	
	node := &ConsumeQueue{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
		queueKey:      queueKey,
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (cq *ConsumeQueue) Tick() behavior_tree.NodeStatus {
	children := cq.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}
	
	child := children[0]
	
	blackboard := cq.Config().Blackboard
	if blackboard == nil {
		return behavior_tree.NodeStatusFailure
	}
	
	// TODO: 实现队列消费逻辑
	// 目前直接执行子节点作为占位实现
	return child.Tick()
}