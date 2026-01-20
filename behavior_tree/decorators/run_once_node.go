package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// RunOnceNode 装饰器节点 - 只运行一次
type RunOnceNode struct {
	behavior_tree.DecoratorNode
	hasRun bool
}

// NewRunOnceNode 创建新的RunOnceNode实例
func NewRunOnceNode(name string, config behavior_tree.NodeConfig) *RunOnceNode {
	node := &RunOnceNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
		hasRun:        false,
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (ron *RunOnceNode) Tick() behavior_tree.NodeStatus {
	if ron.hasRun {
		return behavior_tree.NodeStatusSuccess
	}
	
	children := ron.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}
	
	child := children[0]
	status := child.Tick()
	
	if status != behavior_tree.NodeStatusRunning {
		ron.hasRun = true
	}
	
	return status
}

// Halt 重置运行状态
func (ron *RunOnceNode) Halt() {
	ron.hasRun = false
	children := ron.Children()
	if len(children) > 0 {
		children[0].HaltAndReset()
	}
}