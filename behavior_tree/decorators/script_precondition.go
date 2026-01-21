package decorators

import (
	"github.com/actfuns/gamekit/behavior_tree/core"
)

// ScriptPrecondition 装饰器节点 - 脚本前置条件
type ScriptPrecondition struct {
	core.DecoratorNode
	script string
}

// NewScriptPrecondition 创建新的ScriptPrecondition实例
func NewScriptPrecondition(name string, config core.NodeConfig) *ScriptPrecondition {
	script := ""

	// 使用 config.Manifest.Ports 而不是 config.Props
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["script"]; exists {
			script = portInfo.TypeName
		}
	}

	node := &ScriptPrecondition{
		DecoratorNode: core.NewDecoratorNode(name, config),
		script:        script,
	}
	return node
}

// Tick 执行装饰器节点逻辑
func (sp *ScriptPrecondition) Tick() core.NodeStatus {
	children := sp.Children()
	if len(children) == 0 {
		return core.NodeStatusFailure
	}

	child := children[0]

	blackboard := sp.Config().Blackboard
	if blackboard == nil {
		return core.NodeStatusFailure
	}

	// TODO: 实现脚本执行和条件检查逻辑
	// 目前直接执行子节点作为占位实现
	return child.Tick()
}
