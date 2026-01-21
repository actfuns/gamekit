package actions

import (
	"strconv"
	"time"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// SleepNode sleeps for a specified duration (in milliseconds)
type SleepNode struct {
	core.StatefulActionNode
	duration   time.Duration
	startTime  time.Time
	isSleeping bool
}

// NewSleepNode creates a new sleep node
func NewSleepNode(name string, config core.NodeConfig) *SleepNode {
	node := &SleepNode{}
	statefulNode := core.NewStatefulActionNode(name, config,
		node.onStart,
		node.onRunning,
		node.onHalted)
	node.StatefulActionNode = *statefulNode
	return node
}

// onStart is called when the node starts
func (sn *SleepNode) onStart() core.NodeStatus {
	// Get duration from input port "msec"
	msecStr, ok := sn.Config().InputPorts["msec"]
	if !ok {
		// Try to get from other attributes if not in input ports
		msecStr, ok = sn.Config().OtherAttributes["msec"]
		if !ok {
			return core.NodeStatusFailure
		}
	}

	msec, err := strconv.Atoi(msecStr)
	if err != nil {
		return core.NodeStatusFailure
	}

	if msec <= 0 {
		return core.NodeStatusSuccess
	}

	sn.duration = time.Duration(msec) * time.Millisecond
	sn.startTime = time.Now()
	sn.isSleeping = true

	return core.NodeStatusRunning
}

// onRunning is called while the node is running
func (sn *SleepNode) onRunning() core.NodeStatus {
	if !sn.isSleeping {
		return core.NodeStatusSuccess
	}

	elapsed := time.Since(sn.startTime)
	if elapsed >= sn.duration {
		sn.isSleeping = false
		return core.NodeStatusSuccess
	}

	return core.NodeStatusRunning
}

// onHalted is called when the node is halted
func (sn *SleepNode) onHalted() {
	sn.isSleeping = false
}
