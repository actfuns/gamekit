package actions

import (
	"time"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// TestNodeConfig 配置TestNode的行为
type TestNodeConfig struct {
	// 节点完成时返回的状态
	ReturnStatus core.NodeStatus
	// 成功时执行的脚本（暂不实现）
	SuccessScript string
	// 失败时执行的脚本（暂不实现）
	FailureScript string
	// 完成后执行的脚本（暂不实现）
	PostScript string
	// 异步延迟时间，如果大于0则变为异步动作
	AsyncDelay time.Duration
	// 完成时调用的函数（暂不实现）
	CompleteFunc func() core.NodeStatus
}

// TestNode 是一个可配置的测试节点
type TestNode struct {
	core.StatefulActionNode
	config    *TestNodeConfig
	completed bool
	timer     *time.Timer
}

// NewTestNode 创建新的TestNode
func NewTestNode(name string, config core.NodeConfig, testConfig *TestNodeConfig) *TestNode {
	if testConfig.ReturnStatus == core.NodeStatusIdle {
		panic("TestNode can not return IDLE")
	}

	node := &TestNode{
		config:    testConfig,
		completed: false,
	}

	// 初始化StatefulActionNode
	statefulNode := core.NewStatefulActionNode(name, config,
		node.OnStart,
		node.OnRunning,
		node.OnHalted)
	node.StatefulActionNode = *statefulNode

	return node
}

// OnStart 开始执行
func (node *TestNode) OnStart() core.NodeStatus {
	if node.config.AsyncDelay <= 0 {
		return node.onCompleted()
	}

	// 异步操作，启动定时器
	node.completed = false
	node.timer = time.AfterFunc(node.config.AsyncDelay, func() {
		node.completed = true
		// 在实际实现中需要唤醒信号，这里简化处理
	})

	return core.NodeStatusRunning
}

// OnRunning 持续执行
func (node *TestNode) OnRunning() core.NodeStatus {
	if node.completed {
		return node.onCompleted()
	}
	return core.NodeStatusRunning
}

// OnHalted 被中断时调用
func (node *TestNode) OnHalted() {
	if node.timer != nil {
		node.timer.Stop()
		node.timer = nil
	}
	node.completed = false
}

// onCompleted 完成时的处理
func (node *TestNode) onCompleted() core.NodeStatus {
	status := node.config.ReturnStatus

	// 如果有CompleteFunc，使用它返回的状态
	if node.config.CompleteFunc != nil {
		status = node.config.CompleteFunc()
	}

	// TODO: 执行脚本逻辑（简化版暂不实现）

	return status
}
