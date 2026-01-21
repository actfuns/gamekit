package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree/core"
)

// LoopNode 装饰器循环执行子节点指定次数
type LoopNode struct {
	core.DecoratorNode
	numIterations    int
	currentIteration int
}

// NewLoopNode 创建新的LoopNode实例
func NewLoopNode(name string, config core.NodeConfig) *LoopNode {
	numIterations := -1 // -1 表示无限循环

	if ports := config.Manifest.Ports; len(ports) > 0 {
		if _, exists := ports["num_cycles"]; exists {
			// TODO: 需要解析端口值，这里暂时使用默认值
			numIterations = -1
		}
	}

	node := &LoopNode{
		DecoratorNode:    core.NewDecoratorNode(name, config),
		numIterations:    numIterations,
		currentIteration: 0,
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (ln *LoopNode) Tick() core.NodeStatus {
	children := ln.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]

	status := child.Tick()

	if status == core.NodeStatusRunning {
		return core.NodeStatusRunning
	}

	// 子节点已完成（成功或失败），开始下一次循环
	ln.currentIteration++

	// 检查是否达到指定的循环次数
	if ln.numIterations > 0 && ln.currentIteration >= ln.numIterations {
		return core.NodeStatusSuccess
	}

	// 由于Node接口没有Reset方法，我们直接返回Running状态
	// 下次Tick会自动重新执行子节点
	return core.NodeStatusRunning
}

// Halt 重置循环计数器
func (ln *LoopNode) Halt() {
	ln.currentIteration = 0
	ln.DecoratorNode.Halt()
}
