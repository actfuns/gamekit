package behavior_tree

import (
	"fmt"
	"sync"

	"github.com/actfuns/gamekit/behavior_tree/core"
)

// BehaviorTree represents a complete behavior tree
type BehaviorTree struct {
	rootNode   core.Node
	blackboard *core.Blackboard
	mutex      sync.RWMutex
}

// NewBehaviorTree creates a new behavior tree
func NewBehaviorTree(rootNode core.Node, blackboard *core.Blackboard) *BehaviorTree {
	return &BehaviorTree{
		rootNode:   rootNode,
		blackboard: blackboard,
	}
}

// RootNode returns the root node of the tree
func (bt *BehaviorTree) RootNode() core.Node {
	return bt.rootNode
}

// Blackboard returns the blackboard associated with the tree
func (bt *BehaviorTree) Blackboard() *core.Blackboard {
	return bt.blackboard
}

// Tick executes one tick of the behavior tree
func (bt *BehaviorTree) Tick() core.NodeStatus {
	if bt.rootNode == nil {
		return core.NodeStatusFailure
	}

	bt.mutex.Lock()
	defer bt.mutex.Unlock()

	status := bt.rootNode.ExecuteTick()
	return status
}

// Halt halts the entire behavior tree
func (bt *BehaviorTree) Halt() {
	if bt.rootNode != nil {
		bt.rootNode.HaltAndReset()
	}
}

// PrintTree prints the tree structure
func (bt *BehaviorTree) PrintTree() {
	if bt.rootNode != nil {
		PrintTreeRecursively(bt.rootNode, "")
	}
}

// ApplyVisitor applies a visitor function to all nodes in the tree
func (bt *BehaviorTree) ApplyVisitor(visitor func(core.Node)) {
	if bt.rootNode != nil {
		ApplyRecursiveVisitor(bt.rootNode, visitor)
	}
}

// Create creates a behavior tree from a root node
func Create(rootNode core.Node) (*BehaviorTree, error) {
	if rootNode == nil {
		return nil, fmt.Errorf("root node cannot be nil")
	}

	blackboard := core.NewBlackboard()
	return NewBehaviorTree(rootNode, blackboard), nil
}

// CreateWithBlackboard creates a behavior tree with a specific blackboard
func CreateWithBlackboard(rootNode core.Node, blackboard *core.Blackboard) (*BehaviorTree, error) {
	if rootNode == nil {
		return nil, fmt.Errorf("root node cannot be nil")
	}

	if blackboard == nil {
		blackboard = core.NewBlackboard()
	}
	return NewBehaviorTree(rootNode, blackboard), nil
}

// ApplyRecursiveVisitor applies a visitor function to all nodes in the tree
func ApplyRecursiveVisitor(rootNode core.Node, visitor func(core.Node)) {
	if rootNode == nil {
		return
	}

	visitor(rootNode)
	for _, child := range rootNode.Children() {
		ApplyRecursiveVisitor(child, visitor)
	}
}

// PrintTreeRecursively prints the tree hierarchy recursively
func PrintTreeRecursively(rootNode core.Node, indent string) {
	if rootNode == nil {
		return
	}

	fmt.Printf("%s%s (%s)\n", indent, rootNode.Name(), rootNode.Status().String())
	for _, child := range rootNode.Children() {
		PrintTreeRecursively(child, indent+"  ")
	}
}
