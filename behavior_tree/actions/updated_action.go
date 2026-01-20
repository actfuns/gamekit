package actions

import (
	"fmt"
	"github.com/actfuns/gamekit/behavior_tree"
)

// EntryUpdatedAction 检查条目是否自上次检查以来已更新
type EntryUpdatedAction struct {
	behavior_tree.ActionNodeBase
	sequenceID uint64
	entryKey   string
}

// NewEntryUpdatedAction 创建新的EntryUpdatedAction
func NewEntryUpdatedAction(name string, config behavior_tree.NodeConfig) *EntryUpdatedAction {
	// 检查必需的输入端口 "entry"
	entryPort, exists := config.Manifest.Ports["entry"]
	if !exists {
		panic(fmt.Sprintf("Missing port 'entry' in %s", name))
	}
	
	// 处理黑板指针
	var entryKey string
	if isBlackboardPointer(entryPort.TypeName) {
		entryKey = stripBlackboardPointer(entryPort.TypeName)
	} else {
		entryKey = entryPort.TypeName
	}
	
	node := &EntryUpdatedAction{
		ActionNodeBase: *behavior_tree.NewActionNodeBase(name, config),
		sequenceID:     0,
		entryKey:       entryKey,
	}
	return node
}

// Tick 执行节点逻辑
func (node *EntryUpdatedAction) Tick() behavior_tree.NodeStatus {
	blackboard := node.Config().Blackboard
	if blackboard == nil {
		return behavior_tree.NodeStatusFailure
	}
	
	// 获取条目
	entry := blackboard.GetEntry(node.entryKey)
	if entry == nil {
		return behavior_tree.NodeStatusFailure
	}
	
	// 检查序列ID是否变化
	currentID := entry.SequenceID
	previousID := node.sequenceID
	node.sequenceID = currentID
	
	if previousID != currentID {
		return behavior_tree.NodeStatusSuccess
	}
	return behavior_tree.NodeStatusFailure
}

// isBlackboardPointer 检查是否为黑板指针（以{开头，以}结尾）
func isBlackboardPointer(value string) bool {
	return len(value) > 2 && value[0] == '{' && value[len(value)-1] == '}'
}

// stripBlackboardPointer 去除黑板指针的花括号
func stripBlackboardPointer(value string) string {
	if len(value) > 2 && value[0] == '{' && value[len(value)-1] == '}' {
		return value[1 : len(value)-1]
	}
	return value
}