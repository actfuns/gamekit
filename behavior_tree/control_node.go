package behavior_tree

// ControlNode is the base class for all control nodes
type ControlNode struct {
	TreeNodeBase
}

// NewControlNode creates a new control node
func NewControlNode(name string, config NodeConfig) *ControlNode {
	node := &ControlNode{
		TreeNodeBase: *NewTreeNode(name, config),
	}
	// Override the tick method to prevent panic
	node.TreeNodeBase = *NewTreeNode(name, config)
	return node
}

// Type returns the node type
func (cn *ControlNode) Type() NodeType {
	return NodeTypeControl
}

// ResetChildren resets all children nodes
func (cn *ControlNode) ResetChildren() {
	for _, child := range cn.Children() {
		child.HaltAndReset()
	}
}

// SequenceNode executes children in sequence until one fails
type SequenceNode struct {
	ControlNode
	currentChild int
}

// NewSequenceNode creates a new sequence node
func NewSequenceNode(name string, config NodeConfig) *SequenceNode {
	node := &SequenceNode{
		ControlNode:  *NewControlNode(name, config),
		currentChild: 0,
	}
	// Ensure the embedded TreeNode is properly initialized
	node.ControlNode.TreeNodeBase = *NewTreeNode(name, config)
	return node
}

// Tick executes the sequence logic
func (sn *SequenceNode) Tick() NodeStatus {
	children := sn.Children()
	if len(children) == 0 {
		return NodeStatusSuccess
	}

	for i := sn.currentChild; i < len(children); i++ {
		child := children[i]
		status := child.ExecuteTick()

		switch status {
		case NodeStatusRunning:
			sn.currentChild = i
			return NodeStatusRunning
		case NodeStatusFailure:
			sn.currentChild = 0
			return NodeStatusFailure
		case NodeStatusSuccess:
			// Continue to next child
			continue
		default:
			sn.currentChild = 0
			return status
		}
	}

	// All children succeeded
	sn.currentChild = 0
	return NodeStatusSuccess
}

// Halt handles halting the sequence
func (sn *SequenceNode) Halt() {
	for i := range sn.Children() {
		sn.Children()[i].HaltAndReset()
	}
	sn.currentChild = 0
}

// ReactiveSequenceNode executes all children every tick (reactive)
type ReactiveSequenceNode struct {
	ControlNode
}

// NewReactiveSequenceNode creates a new reactive sequence node
func NewReactiveSequenceNode(name string, config NodeConfig) *ReactiveSequenceNode {
	node := &ReactiveSequenceNode{
		ControlNode: *NewControlNode(name, config),
	}
	// Ensure the embedded TreeNode is properly initialized
	node.ControlNode.TreeNodeBase = *NewTreeNode(name, config)
	return node
}

// Tick executes the reactive sequence logic
func (rsn *ReactiveSequenceNode) Tick() NodeStatus {
	children := rsn.Children()
	if len(children) == 0 {
		return NodeStatusSuccess
	}

	for i := range children {
		child := children[i]
		status := child.ExecuteTick()

		if status != NodeStatusSuccess {
			return status
		}
	}

	return NodeStatusSuccess
}

// FallbackNode executes children in sequence until one succeeds
type FallbackNode struct {
	ControlNode
	currentChild int
}

// NewFallbackNode creates a new fallback node
func NewFallbackNode(name string, config NodeConfig) *FallbackNode {
	node := &FallbackNode{
		ControlNode:  *NewControlNode(name, config),
		currentChild: 0,
	}
	// Ensure the embedded TreeNode is properly initialized
	node.ControlNode.TreeNodeBase = *NewTreeNode(name, config)
	return node
}

// Tick executes the fallback logic
func (fn *FallbackNode) Tick() NodeStatus {
	children := fn.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	for i := fn.currentChild; i < len(children); i++ {
		child := children[i]
		status := child.ExecuteTick()

		switch status {
		case NodeStatusRunning:
			fn.currentChild = i
			return NodeStatusRunning
		case NodeStatusSuccess:
			fn.currentChild = 0
			return NodeStatusSuccess
		case NodeStatusFailure:
			// Continue to next child
			continue
		default:
			fn.currentChild = 0
			return status
		}
	}

	// All children failed
	fn.currentChild = 0
	return NodeStatusFailure
}

// Halt handles halting the fallback
func (fn *FallbackNode) Halt() {
	for i := range fn.Children() {
		fn.Children()[i].HaltAndReset()
	}
	fn.currentChild = 0
}

// ReactiveFallbackNode executes all children every tick (reactive)
type ReactiveFallbackNode struct {
	ControlNode
}

// NewReactiveFallbackNode creates a new reactive fallback node
func NewReactiveFallbackNode(name string, config NodeConfig) *ReactiveFallbackNode {
	node := &ReactiveFallbackNode{
		ControlNode: *NewControlNode(name, config),
	}
	// Ensure the embedded TreeNode is properly initialized
	node.ControlNode.TreeNodeBase = *NewTreeNode(name, config)
	return node
}

// Tick executes the reactive fallback logic
func (rfn *ReactiveFallbackNode) Tick() NodeStatus {
	children := rfn.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	for i := range children {
		child := children[i]
		status := child.ExecuteTick()

		if status != NodeStatusFailure {
			return status
		}
	}

	return NodeStatusFailure
}
