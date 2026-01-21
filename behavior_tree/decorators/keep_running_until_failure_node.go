package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree/core"
)

// KeepRunningUntilFailureNode 装饰器节点 - 持续运行直到失败
type KeepRunningUntilFailureNode struct {
	core.DecoratorNode
}

// NewKeepRunningUntilFailureNode 创建新的KeepRunningUntilFailureNode实例
func NewKeepRunningUntilFailureNode(name string, config core.NodeConfig) *KeepRunningUntilFailureNode {
	node := &KeepRunningUntilFailureNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (kruf *KeepRunningUntilFailureNode) Tick() core.NodeStatus {
	children := kruf.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	if status == core.NodeStatusFailure {
		return core.NodeStatusSuccess
	}

	return core.NodeStatusRunning
}
