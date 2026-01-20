package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// KeepRunningUntilFailureNode 装饰器节点 - 持续运行直到失败
type KeepRunningUntilFailureNode struct {
	behavior_tree.DecoratorNode
}

// NewKeepRunningUntilFailureNode 创建新的KeepRunningUntilFailureNode实例
func NewKeepRunningUntilFailureNode(name string, config behavior_tree.NodeConfig) *KeepRunningUntilFailureNode {
	node := &KeepRunningUntilFailureNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (kruf *KeepRunningUntilFailureNode) Tick() behavior_tree.NodeStatus {
	children := kruf.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}
	
	child := children[0]
	status := child.Tick()
	
	if status == behavior_tree.NodeStatusFailure {
		return behavior_tree.NodeStatusSuccess
	}
	
	return behavior_tree.NodeStatusRunning
}