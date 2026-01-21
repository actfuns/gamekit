package decorators

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// DelayNode delays the execution of its child by a specified number of milliseconds.
type DelayNode struct {
	core.DecoratorNode
	msec          uint
	delayStarted  bool
	delayComplete bool
	delayAborted  bool
	delayMutex    sync.Mutex
	readFromPorts bool
}

// NewDelayNode creates a new DelayNode
func NewDelayNode(name string, config core.NodeConfig) *DelayNode {
	delayNode := &DelayNode{
		DecoratorNode: core.NewDecoratorNode(name, config),
		msec:          0,
		delayStarted:  false,
		delayComplete: false,
		delayAborted:  false,
		readFromPorts: true,
	}

	// Try to get delay_msec from input ports
	if msecStr, ok := config.InputPorts["delay_msec"]; ok {
		if msec, err := strconv.Atoi(msecStr); err == nil && msec >= 0 {
			delayNode.msec = uint(msec)
		}
	}

	return delayNode
}

// Tick executes the delay logic
func (dn *DelayNode) Tick() core.NodeStatus {
	if dn.readFromPorts {
		var delayMsec uint
		if value, ok := dn.GetInput("delay_msec"); ok {
			if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
				delayMsec = uint(intValue)
				dn.msec = delayMsec
			} else {
				panic(fmt.Sprintf("Invalid parameter [delay_msec] in DelayNode: %v", err))
			}
		} else {
			panic("Missing parameter [delay_msec] in DelayNode")
		}
	}

	if !dn.delayStarted {
		dn.delayComplete = false
		dn.delayAborted = false
		dn.delayStarted = true
		dn.SetStatus(core.NodeStatusRunning)

		// Start the delay timer
		go func() {
			time.Sleep(time.Duration(dn.msec) * time.Millisecond)
			dn.delayMutex.Lock()
			dn.delayComplete = true
			dn.delayMutex.Unlock()
			dn.EmitWakeUpSignal()
		}()
	}

	dn.delayMutex.Lock()
	defer dn.delayMutex.Unlock()

	if dn.delayAborted {
		dn.delayAborted = false
		dn.delayStarted = false
		return core.NodeStatusFailure
	}

	if dn.delayComplete {
		children := dn.Children()
		if len(children) == 0 {
			dn.delayStarted = false
			dn.delayAborted = false
			return core.NodeStatusFailure
		}

		child := children[0]
		childStatus := child.Tick()
		if core.IsStatusCompleted(childStatus) {
			dn.delayStarted = false
			dn.delayAborted = false
			child.HaltAndReset()
		}
		return childStatus
	}

	return core.NodeStatusRunning
}

// Halt handles halting the delay node
func (dn *DelayNode) Halt() {
	dn.delayStarted = false
	dn.delayAborted = true
	children := dn.Children()
	if len(children) > 0 {
		children[0].Halt()
	}
	dn.DecoratorNode.Halt()
}
