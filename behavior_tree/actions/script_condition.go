package actions

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// ScriptCondition 动作节点 - 脚本条件
type ScriptCondition struct {
	behavior_tree.ActionNodeBase
	script string
}

// NewScriptCondition 创建新的ScriptCondition实例
func NewScriptCondition(name string, config behavior_tree.NodeConfig) *ScriptCondition {
	script := ""
	
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["script"]; exists {
			script = portInfo.TypeName
		}
	}
	
	node := &ScriptCondition{
		ActionNodeBase: *behavior_tree.NewActionNodeBase(name, config),
		script:         script,
	}
	return node
}

// Tick 执行动作节点逻辑
func (sc *ScriptCondition) Tick() behavior_tree.NodeStatus {
	// TODO: 实现脚本执行逻辑
	// 暂时返回成功
	return behavior_tree.NodeStatusSuccess
}