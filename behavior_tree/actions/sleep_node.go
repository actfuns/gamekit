package actions

import (
	"strconv"
	"time"
	"github.com/actfuns/gamekit/behavior_tree"
)

// SleepNode sleeps for a specified duration (in milliseconds)
type SleepNode struct {
	behavior_tree.StatefulActionNode
	duration    time.Duration
	startTime   time.Time
	isSleeping  bool
}

// NewSleepNode creates a new sleep node
func NewSleepNode(name string, config behavior_tree.NodeConfig) *SleepNode {
	node := &SleepNode{}
	statefulNode := behavior_tree.NewStatefulActionNode(name, config,
		node.onStart,
		node.onRunning,
		node.onHalted)
	node.StatefulActionNode = *statefulNode
	return node
}

// onStart is called when the node starts
func (sn *SleepNode) onStart() behavior_tree.NodeStatus {
	// Get duration from input port "msec"
	msecStr, ok := sn.Config().InputPorts["msec"]
	if !ok {
		// Try to get from other attributes if not in input ports
		msecStr, ok = sn.Config().OtherAttributes["msec"]
		if !ok {
			return behavior_tree.NodeStatusFailure
		}
	}
	
	msec, err := strconv.Atoi(msecStr)
	if err != nil {
		return behavior_tree.NodeStatusFailure
	}
	
	if msec <= 0 {
		return behavior_tree.NodeStatusSuccess
	}
	
	sn.duration = time.Duration(msec) * time.Millisecond
	sn.startTime = time.Now()
	sn.isSleeping = true
	
	return behavior_tree.NodeStatusRunning
}

// onRunning is called while the node is running
func (sn *SleepNode) onRunning() behavior_tree.NodeStatus {
	if !sn.isSleeping {
		return behavior_tree.NodeStatusSuccess
	}
	
	elapsed := time.Since(sn.startTime)
	if elapsed >= sn.duration {
		sn.isSleeping = false
		return behavior_tree.NodeStatusSuccess
	}
	
	return behavior_tree.NodeStatusRunning
}

// onHalted is called when the node is halted
func (sn *SleepNode) onHalted() {
	sn.isSleeping = false
}