package decorators

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// TimeoutNode adds a timeout to its child execution.
// If the child doesn't complete within the specified time, it returns FAILURE.
type TimeoutNode struct {
	core.DecoratorNode
	msec           uint
	timeoutStarted bool
	childHalted    bool
	timeoutMutex   sync.Mutex
	readFromPorts  bool
}

// NewTimeoutNode creates a new TimeoutNode
func NewTimeoutNode(name string, config core.NodeConfig) *TimeoutNode {
	timeoutNode := &TimeoutNode{
		DecoratorNode:  core.NewDecoratorNode(name, config),
		msec:           0,
		timeoutStarted: false,
		childHalted:    false,
		readFromPorts:  true,
	}

	// Try to get msec from input ports
	if msecStr, ok := config.InputPorts["msec"]; ok {
		if msec, err := strconv.Atoi(msecStr); err == nil && msec >= 0 {
			timeoutNode.msec = uint(msec)
		}
	}

	return timeoutNode
}

// NewTimeoutNodeFromConfig creates a new TimeoutNode that reads msec from ports
func NewTimeoutNodeFromConfig(name string, config core.NodeConfig) *TimeoutNode {
	return &TimeoutNode{
		DecoratorNode:  core.NewDecoratorNode(name, config),
		msec:           0,
		timeoutStarted: false,
		childHalted:    false,
		readFromPorts:  true,
	}
}

// Tick executes the timeout logic
func (tn *TimeoutNode) Tick() core.NodeStatus {
	if tn.readFromPorts {
		var msec uint
		if value, ok := tn.GetInput("msec"); ok {
			if intValue, err := strconv.ParseUint(value, 10, 32); err == nil {
				msec = uint(intValue)
				tn.msec = msec
			} else {
				panic(fmt.Sprintf("Invalid parameter [msec] in TimeoutNode: %v", err))
			}
		} else {
			panic("Missing parameter [msec] in TimeoutNode")
		}

	}

	if !tn.timeoutStarted {
		tn.timeoutStarted = true
		tn.SetStatus(core.NodeStatusRunning)
		tn.childHalted = false

		if tn.msec > 0 {
			// Start the timeout timer
			go func() {
				time.Sleep(time.Duration(tn.msec) * time.Millisecond)
				tn.timeoutMutex.Lock()
				defer tn.timeoutMutex.Unlock()

				children := tn.Children()
				if len(children) > 0 && children[0].Status() == core.NodeStatusRunning {
					tn.childHalted = true
					children[0].Halt()
					tn.EmitWakeUpSignal()
				}
			}()
		}
	}

	tn.timeoutMutex.Lock()
	defer tn.timeoutMutex.Unlock()

	if tn.childHalted {
		tn.timeoutStarted = false
		return core.NodeStatusFailure
	}

	children := tn.Children()
	if len(children) == 0 {
		tn.timeoutStarted = false
		return core.NodeStatusFailure
	}

	child := children[0]
	childStatus := child.Tick()
	if core.IsStatusCompleted(childStatus) {
		tn.timeoutStarted = false
		child.HaltAndReset()
	}

	return childStatus
}

// Halt handles halting the timeout node
func (tn *TimeoutNode) Halt() {
	tn.timeoutStarted = false
	// Note: In Go, we can't easily cancel the goroutine, but the flag will prevent action
	children := tn.Children()
	if len(children) > 0 {
		children[0].Halt()
	}
	tn.DecoratorNode.Halt()
}
