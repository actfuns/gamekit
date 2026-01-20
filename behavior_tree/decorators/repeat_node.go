package decorators

import (
	"fmt"
	"strconv"

	"github.com/actfuns/gamekit/behavior_tree"
)

const (
	// RepeatNumCycles is the port name for number of cycles in RepeatNode
	RepeatNumCycles = "num_cycles"
)

// RepeatNode repeats its child until it fails or reaches the maximum number of repetitions.
// If num_cycles is -1, it repeats indefinitely until failure.
type RepeatNode struct {
	behavior_tree.DecoratorNode
	numCycles     int
	repeatCount   int
	readFromPorts bool
}

// NewRepeatNode creates a new RepeatNode
func NewRepeatNode(name string, config behavior_tree.NodeConfig) *RepeatNode {
	repeatNode := &RepeatNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
		numCycles:     -1, // default to infinite
		repeatCount:   0,
		readFromPorts: true,
	}
	
	// Try to get num_cycles from input ports
	if numCyclesStr, ok := config.InputPorts[RepeatNumCycles]; ok {
		if numCycles, err := strconv.Atoi(numCyclesStr); err == nil {
			repeatNode.numCycles = numCycles
		}
	}
	
	return repeatNode
}

// NewRepeatNodeFromConfig creates a new RepeatNode that reads num_cycles from ports
func NewRepeatNodeFromConfig(name string, config behavior_tree.NodeConfig) *RepeatNode {
	return &RepeatNode{
		DecoratorNode: *behavior_tree.NewDecoratorNode(name, config),
		numCycles:     0,
		repeatCount:   0,
		readFromPorts: true,
	}
}

// Tick executes the repeat logic
func (rn *RepeatNode) Tick() behavior_tree.NodeStatus {
	if rn.readFromPorts {
		if value, ok := rn.GetInput(RepeatNumCycles); ok {
			if intValue, err := strconv.Atoi(value); err == nil {
				rn.numCycles = intValue
			} else {
				panic(fmt.Sprintf("Invalid parameter [%s] in RepeatNode: %v", RepeatNumCycles, err))
			}
		} else {
			panic("Missing parameter [" + RepeatNumCycles + "] in RepeatNode")
		}
	}

	children := rn.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	doLoop := rn.repeatCount < rn.numCycles || rn.numCycles == -1

	for doLoop {
		prevStatus := child.Status()
		status := child.Tick()

		switch status {
		case behavior_tree.NodeStatusSuccess:
			rn.repeatCount++
			doLoop = rn.repeatCount < rn.numCycles || rn.numCycles == -1
			child.HaltAndReset()

			// For async children, return RUNNING to make it interruptible
			if child.RequiresWakeUp() && prevStatus == behavior_tree.NodeStatusIdle && doLoop {
				rn.EmitWakeUpSignal()
				return behavior_tree.NodeStatusRunning
			}

		case behavior_tree.NodeStatusFailure:
			rn.repeatCount = 0
			child.HaltAndReset()
			return behavior_tree.NodeStatusFailure

		case behavior_tree.NodeStatusRunning:
			return behavior_tree.NodeStatusRunning

		case behavior_tree.NodeStatusSkipped:
			child.HaltAndReset()
			return behavior_tree.NodeStatusSkipped

		case behavior_tree.NodeStatusIdle:
			panic("[" + rn.Name() + "]: A child should not return IDLE")
		}
	}

	rn.repeatCount = 0
	return behavior_tree.NodeStatusSuccess
}

// Halt handles halting the repeat node
func (rn *RepeatNode) Halt() {
	rn.repeatCount = 0
	children := rn.Children()
	if len(children) > 0 {
		children[0].Halt()
	}
	rn.DecoratorNode.Halt()
}