package decorators

import (
	"fmt"
	"strconv"

	"github.com/actfuns/gamekit/behavior_tree"
)

const (
	// RetryNumAttempts is the port name for number of attempts in RetryNode
	RetryNumAttempts = "num_attempts"
)

// RetryNode retries its child until it succeeds or reaches the maximum number of attempts.
// If max_attempts is -1, it retries indefinitely until success.
type RetryNode struct {
	behavior_tree.DecoratorNode
	maxAttempts   int
	tryCount      int
	readFromPorts bool
}

// NewRetryNode creates a new RetryNode with either direct parameter or from config
// If maxAttempts is provided, use it directly (code usage)
// If no maxAttempts provided, read from ports (XML usage)
func NewRetryNode(name string, config behavior_tree.NodeConfig, maxAttempts ...int) *RetryNode {
	retryNode := &RetryNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
		tryCount:      0,
	}
	
	if len(maxAttempts) > 0 {
		// Direct parameter provided (for code usage)
		retryNode.maxAttempts = maxAttempts[0]
		retryNode.readFromPorts = false
	} else {
		// Read from ports (for XML usage)
		retryNode.maxAttempts = -1 // default to infinite
		retryNode.readFromPorts = true
		
		if numAttemptsStr, ok := config.InputPorts[RetryNumAttempts]; ok {
			if numAttempts, err := strconv.Atoi(numAttemptsStr); err == nil {
				retryNode.maxAttempts = numAttempts
			}
		}
	}
	
	return retryNode
}

// Tick executes the retry logic
func (rn *RetryNode) Tick() behavior_tree.NodeStatus {
	if rn.readFromPorts {
		if value, ok := rn.GetInput(RetryNumAttempts); ok {
			if intValue, err := strconv.Atoi(value); err == nil {
				rn.maxAttempts = intValue
			} else {
				panic(fmt.Sprintf("Invalid parameter [%s] in RetryNode: %v", RetryNumAttempts, err))
			}
		} else {
			panic("Missing parameter [" + RetryNumAttempts + "] in RetryNode")
		}
	}

	children := rn.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	doLoop := rn.tryCount < rn.maxAttempts || rn.maxAttempts == -1

	for doLoop {
		prevStatus := child.Status()
		status := child.Tick()

		switch status {
		case behavior_tree.NodeStatusSuccess:
			rn.tryCount = 0
			child.HaltAndReset()
			return behavior_tree.NodeStatusSuccess

		case behavior_tree.NodeStatusFailure:
			rn.tryCount++
			// Refresh maxAttempts in case it changed in one of the child nodes
			if rn.readFromPorts {
				if value, ok := rn.GetInput(RetryNumAttempts); ok {
					if intValue, err := strconv.Atoi(value); err == nil {
						rn.maxAttempts = intValue
					}
				}
			}
			doLoop = rn.tryCount < rn.maxAttempts || rn.maxAttempts == -1
			child.HaltAndReset()

			// For async children, return RUNNING to make it interruptible
			if child.RequiresWakeUp() && prevStatus == behavior_tree.NodeStatusIdle && doLoop {
				rn.EmitWakeUpSignal()
				return behavior_tree.NodeStatusRunning
			}

		case behavior_tree.NodeStatusRunning:
			return behavior_tree.NodeStatusRunning

		case behavior_tree.NodeStatusSkipped:
			child.HaltAndReset()
			return behavior_tree.NodeStatusSkipped

		case behavior_tree.NodeStatusIdle:
			panic("[" + rn.Name() + "]: A child should not return IDLE")
		}
	}

	rn.tryCount = 0
	return behavior_tree.NodeStatusFailure
}

// Halt handles halting the retry node
func (rn *RetryNode) Halt() {
	rn.tryCount = 0
	children := rn.Children()
	if len(children) > 0 {
		children[0].Halt()
	}
	rn.DecoratorNode.Halt()
}