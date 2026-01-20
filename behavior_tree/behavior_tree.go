package behavior_tree

import (
	"fmt"
	"sync"
)

// BehaviorTree represents a complete behavior tree
type BehaviorTree struct {
	rootNode   TreeNode
	blackboard *Blackboard
	factory    *BehaviorTreeFactory
	mutex      sync.RWMutex
}

// NewBehaviorTree creates a new behavior tree
func NewBehaviorTree(rootNode TreeNode, blackboard *Blackboard, factory *BehaviorTreeFactory) *BehaviorTree {
	return &BehaviorTree{
		rootNode:   rootNode,
		blackboard: blackboard,
		factory:    factory,
	}
}

// RootNode returns the root node of the tree
func (bt *BehaviorTree) RootNode() TreeNode {
	return bt.rootNode
}

// Blackboard returns the blackboard associated with the tree
func (bt *BehaviorTree) Blackboard() *Blackboard {
	return bt.blackboard
}

// Factory returns the factory used to create nodes
func (bt *BehaviorTree) Factory() *BehaviorTreeFactory {
	return bt.factory
}

// Tick executes one tick of the behavior tree
func (bt *BehaviorTree) Tick() NodeStatus {
	if bt.rootNode == nil {
		return NodeStatusFailure
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
func (bt *BehaviorTree) ApplyVisitor(visitor func(TreeNode)) {
	if bt.rootNode != nil {
		ApplyRecursiveVisitor(bt.rootNode, visitor)
	}
}

// Create creates a behavior tree from a root node
func Create(rootNode TreeNode) (*BehaviorTree, error) {
	if rootNode == nil {
		return nil, fmt.Errorf("root node cannot be nil")
	}

	blackboard := NewBlackboard()
	factory := NewBehaviorTreeFactory()

	return NewBehaviorTree(rootNode, blackboard, factory), nil
}

// CreateWithBlackboard creates a behavior tree with a specific blackboard
func CreateWithBlackboard(rootNode TreeNode, blackboard *Blackboard) (*BehaviorTree, error) {
	if rootNode == nil {
		return nil, fmt.Errorf("root node cannot be nil")
	}

	if blackboard == nil {
		blackboard = NewBlackboard()
	}

	factory := NewBehaviorTreeFactory()

	return NewBehaviorTree(rootNode, blackboard, factory), nil
}

// ApplyRecursiveVisitor applies a visitor function to all nodes in the tree
func ApplyRecursiveVisitor(rootNode TreeNode, visitor func(TreeNode)) {
	if rootNode == nil {
		return
	}

	visitor(rootNode)
	for _, child := range rootNode.Children() {
		ApplyRecursiveVisitor(child, visitor)
	}
}

// PrintTreeRecursively prints the tree hierarchy recursively
func PrintTreeRecursively(rootNode TreeNode, indent string) {
	if rootNode == nil {
		return
	}

	fmt.Printf("%s%s (%s)\n", indent, rootNode.Name(), rootNode.Status().String())
	for _, child := range rootNode.Children() {
		PrintTreeRecursively(child, indent+"  ")
	}
}
