package actions

import (
	"github.com/actfuns/gamekit/behavior_tree"
)

// ScriptNode 动作节点 - 执行脚本
type ScriptNode struct {
	behavior_tree.ActionNodeBase
	script string
}

// NewScriptNode 创建新的ScriptNode实例
func NewScriptNode(name string, config behavior_tree.NodeConfig) *ScriptNode {
	script := ""
	
	if ports := config.Manifest.Ports; ports != nil {
		if portInfo, exists := ports["script"]; exists {
			script = portInfo.TypeName
		}
	}
	
	node := &ScriptNode{
		ActionNodeBase: *behavior_tree.NewActionNodeBase(name, config),
		script:         script,
	}
	return node
}

// Tick 执行动作节点逻辑
func (sn *ScriptNode) Tick() behavior_tree.NodeStatus {
	// TODO: 实现脚本执行逻辑
	// 暂时返回成功
	return behavior_tree.NodeStatusSuccess
}