package behavior_tree

import (
	"fmt"
	"sync"
)

// TreeNodeCreator is a function that creates a tree node
type TreeNodeCreator func(string, NodeConfig) (TreeNode, error)

// BehaviorTreeFactory is used to register and create tree nodes
type BehaviorTreeFactory struct {
	manifests    map[string]TreeNodeManifest
	constructors map[string]TreeNodeCreator
	mutex        sync.RWMutex
}

// NewBehaviorTreeFactory creates a new behavior tree factory
func NewBehaviorTreeFactory() *BehaviorTreeFactory {
	return &BehaviorTreeFactory{
		manifests:    make(map[string]TreeNodeManifest),
		constructors: make(map[string]TreeNodeCreator),
	}
}

// RegisterBuilder registers a node builder with the factory
func (f *BehaviorTreeFactory) RegisterBuilder(registrationID string, manifest TreeNodeManifest, creator TreeNodeCreator) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, exists := f.constructors[registrationID]; exists {
		return fmt.Errorf("registration ID '%s' already exists", registrationID)
	}

	f.manifests[registrationID] = manifest
	f.constructors[registrationID] = creator
	return nil
}

// UnregisterBuilder removes a registered builder
func (f *BehaviorTreeFactory) UnregisterBuilder(registrationID string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	delete(f.manifests, registrationID)
	delete(f.constructors, registrationID)
}

// CreateNode creates a node using the registered constructor
func (f *BehaviorTreeFactory) CreateNode(registrationID string, name string, config NodeConfig) (TreeNode, error) {
	f.mutex.RLock()
	creator, exists := f.constructors[registrationID]
	f.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("registration ID '%s' not found", registrationID)
	}

	return creator(name, config)
}

// GetManifest returns the manifest for a registration ID
func (f *BehaviorTreeFactory) GetManifest(registrationID string) (TreeNodeManifest, bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	manifest, exists := f.manifests[registrationID]
	return manifest, exists
}

// RegisteredNodes returns all registered node IDs
func (f *BehaviorTreeFactory) RegisteredNodes() []string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	ids := make([]string, 0, len(f.constructors))
	for id := range f.constructors {
		ids = append(ids, id)
	}
	return ids
}

// Clear removes all registered builders
func (f *BehaviorTreeFactory) Clear() {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.manifests = make(map[string]TreeNodeManifest)
	f.constructors = make(map[string]TreeNodeCreator)
}
