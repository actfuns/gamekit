package actions

import "github.com/actfuns/gamekit/behavior_tree/core"

// ScriptNode 动作节点 - 执行脚本
type ScriptNode struct {
	core.ActionNodeBase
	script string
}

// NewScriptNode 创建新的ScriptNode实例
func NewScriptNode(name string, config core.NodeConfig) *ScriptNode {
	script := ""

	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["script"]; exists {
			script = portInfo.TypeName
		}
	}

	node := &ScriptNode{
		ActionNodeBase: core.NewActionNodeBase(name, config),
		script:         script,
	}
	return node
}

// Tick 执行动作节点逻辑
func (sn *ScriptNode) Tick() core.NodeStatus {
	// TODO: 实现脚本执行逻辑
	// 暂时返回成功
	return core.NodeStatusSuccess
}
