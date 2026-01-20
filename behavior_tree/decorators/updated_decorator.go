package decorators

import "github.com/actfuns/gamekit/behavior_tree"

// EntryUpdatedDecorator executes its child only when the specified blackboard entry is updated.
type EntryUpdatedDecorator struct {
	behavior_tree.DecoratorNode
	entryKey           string
	sequenceId         uint64
	stillExecutingChild bool
	ifNotUpdated       behavior_tree.NodeStatus
}

// NewEntryUpdatedDecorator creates a new EntryUpdatedDecorator
func NewEntryUpdatedDecorator(name string, config behavior_tree.NodeConfig, ifNotUpdated behavior_tree.NodeStatus) *EntryUpdatedDecorator {
	// Get the entry port
	entryPort, exists := config.InputPorts["entry"]
	if !exists || entryPort == "" {
		panic("Missing port 'entry' in " + name)
	}

	// Extract the entry key (handle blackboard pointer syntax)
	entryKey := entryPort
	if len(entryKey) > 0 && entryKey[0] == '{' && entryKey[len(entryKey)-1] == '}' {
		entryKey = entryKey[1 : len(entryKey)-1]
	}

	return &EntryUpdatedDecorator{
		DecoratorNode:       *behavior_tree.NewDecoratorNode(name, config),
		entryKey:            entryKey,
		sequenceId:          0,
		stillExecutingChild: false,
		ifNotUpdated:        ifNotUpdated,
	}
}

// Tick executes the updated decorator logic
func (eud *EntryUpdatedDecorator) Tick() behavior_tree.NodeStatus {
	// Continue executing an asynchronous child
	if eud.stillExecutingChild {
		children := eud.Children()
		if len(children) == 0 {
			eud.stillExecutingChild = false
			return behavior_tree.NodeStatusFailure
		}

		child := children[0]
		status := child.Tick()
		eud.stillExecutingChild = (status == behavior_tree.NodeStatusRunning)
		return status
	}

	// Check if the blackboard entry has been updated
	blackboard := eud.Config().Blackboard
	if blackboard != nil {
		entry := blackboard.GetEntry(eud.entryKey)
		if entry != nil {
			currentId := entry.SequenceID
			previousId := eud.sequenceId
			eud.sequenceId = currentId

			if previousId == currentId {
				// Entry not updated
				return eud.ifNotUpdated
			}
		} else {
			// Entry doesn't exist
			return eud.ifNotUpdated
		}
	} else {
		// No blackboard available
		return eud.ifNotUpdated
	}

	// Execute child since entry was updated
	children := eud.Children()
	if len(children) == 0 {
		return behavior_tree.NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()
	eud.stillExecutingChild = (status == behavior_tree.NodeStatusRunning)
	return status
}

// Halt handles halting the updated decorator
func (eud *EntryUpdatedDecorator) Halt() {
	eud.stillExecutingChild = false
	children := eud.Children()
	if len(children) > 0 {
		children[0].Halt()
	}
	eud.DecoratorNode.Halt()
}