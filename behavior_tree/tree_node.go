package behavior_tree

import (
	"sync"
)

// TreeNode represents a node in the behavior tree
type TreeNode interface {
	Name() string
	Status() NodeStatus
	SetStatus(NodeStatus)
	Config() NodeConfig
	Blackboard() *Blackboard
	AddChild(TreeNode)
	Children() []TreeNode
	SetParent(*TreeNode)
	Parent() *TreeNode
	Tick() NodeStatus
	ExecuteTick() NodeStatus
	Halt()
	HaltAndReset()
	RequiresWakeUp() bool
	EmitWakeUpSignal()
}

// TreeNodeBase represents a tree node implementing TreeNodeInterface
type TreeNodeBase struct {
	name     string
	config   NodeConfig
	status   NodeStatus
	mutex    sync.RWMutex
	children []TreeNode
	parent   *TreeNode
	factory  *BehaviorTreeFactory
}

// NewTreeNode creates a new tree node
func NewTreeNode(name string, config NodeConfig) *TreeNodeBase {
	return &TreeNodeBase{
		name:   name,
		config: config,
		status: NodeStatusIdle,
	}
}

// Name returns the name of the node
func (tn *TreeNodeBase) Name() string {
	return tn.name
}

// Status returns the current status of the node
func (tn *TreeNodeBase) Status() NodeStatus {
	tn.mutex.RLock()
	defer tn.mutex.RUnlock()
	return tn.status
}

// SetStatus sets the status of the node
func (tn *TreeNodeBase) SetStatus(status NodeStatus) {
	tn.mutex.Lock()
	defer tn.mutex.Unlock()
	tn.status = status
}

// Config returns the node configuration
func (tn *TreeNodeBase) Config() NodeConfig {
	return tn.config
}

// Blackboard returns the blackboard associated with this node
func (tn *TreeNodeBase) Blackboard() *Blackboard {
	return tn.config.Blackboard
}

// AddChild adds a child node
func (tn *TreeNodeBase) AddChild(child TreeNode) {
	tn.children = append(tn.children, child)
}

// Children returns the child nodes
func (tn *TreeNodeBase) Children() []TreeNode {
	return tn.children
}

// SetParent sets the parent node
func (tn *TreeNodeBase) SetParent(parent *TreeNode) {
	tn.parent = parent
}

// Parent returns the parent node
func (tn *TreeNodeBase) Parent() *TreeNode {
	return tn.parent
}

// Tick is the main execution method that must be implemented by derived classes
// For the base TreeNode, we return SUCCESS as a default behavior
func (tn *TreeNodeBase) Tick() NodeStatus {
	return NodeStatusSuccess
}

// Halt is called when the node is halted
func (tn *TreeNodeBase) Halt() {
	// Default implementation does nothing
}

// Type returns the node type
func (tn *TreeNodeBase) Type() NodeType {
	return NodeTypeUndefined
}

// UID returns the unique identifier of the node
func (tn *TreeNodeBase) UID() uint16 {
	return tn.config.UID
}

// Manifest returns the node manifest
func (tn *TreeNodeBase) Manifest() TreeNodeManifest {
	return tn.config.Manifest
}

// GetInput retrieves an input port value
func (tn *TreeNodeBase) GetInput(key string) (string, bool) {
	// First check if the key exists in the input ports remapping
	if value, exists := tn.config.InputPorts[key]; exists {
		return value, true
	}
	
	// If not found in input ports, check if it's defined in manifest as input/inout port
	if portInfo, exists := tn.config.Manifest.Ports[key]; exists {
		if portInfo.Direction == PortDirectionInput || portInfo.Direction == PortDirectionInOut {
			// Return empty string with true to indicate the port exists but has no value
			return "", true
		}
	}
	return "", false
}

// ExecuteTick executes a tick and handles status changes
func (tn *TreeNodeBase) ExecuteTick() NodeStatus {
	if tn.Status() == NodeStatusRunning {
		newStatus := tn.Tick()
		tn.SetStatus(newStatus)
		return newStatus
	}

	// If not running, start fresh
	tn.SetStatus(NodeStatusIdle)
	newStatus := tn.Tick()
	tn.SetStatus(newStatus)
	return newStatus
}

// HaltAndReset halts the node and resets its status to Idle
func (tnb *TreeNodeBase) HaltAndReset() {
	tnb.Halt()
	tnb.SetStatus(NodeStatusIdle)
}

// RequiresWakeUp returns whether the node requires wake up signal
func (tnb *TreeNodeBase) RequiresWakeUp() bool {
	return false
}

// EmitWakeUpSignal emits wake up signal (placeholder implementation)
func (tnb *TreeNodeBase) EmitWakeUpSignal() {
	// Placeholder implementation
}
