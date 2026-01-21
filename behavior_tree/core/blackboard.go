package core

import (
	"reflect"
	"sync"
)

// Any represents a type-safe container for any value
type Any struct {
	value interface{}
}

// NewAny creates a new Any from a value
func NewAny(value interface{}) Any {
	return Any{value: value}
}

// Value returns the contained value
func (a Any) Value() interface{} {
	return a.value
}

// Type returns the type of the contained value
func (a Any) Type() reflect.Type {
	if a.value == nil {
		return nil
	}
	return reflect.TypeOf(a.value)
}

// Entry represents a blackboard entry with value and type info
type Entry struct {
	Value      Any
	Info       TypeInfo
	SequenceID uint64
}

// TypeInfo contains type information
type TypeInfo struct {
	TypeName string
}

// Blackboard is used by BehaviorTrees to exchange typed data
type Blackboard struct {
	entries  map[string]Entry
	parent   *Blackboard
	mutex    sync.RWMutex
	portInfo map[string]PortInfo
}

// NewBlackboard creates a new blackboard
func NewBlackboard() *Blackboard {
	return &Blackboard{
		entries: make(map[string]Entry),
	}
}

// NewBlackboardWithParent creates a new blackboard with a parent
func NewBlackboardWithParent(parent *Blackboard) *Blackboard {
	return &Blackboard{
		entries: make(map[string]Entry),
		parent:  parent,
	}
}

// Set sets a value in the blackboard
func (bb *Blackboard) Set(key string, value interface{}) error {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()

	bb.entries[key] = Entry{
		Value: NewAny(value),
		Info:  TypeInfo{TypeName: reflect.TypeOf(value).String()},
	}
	return nil
}

// Get retrieves a value from the blackboard
func (bb *Blackboard) Get(key string) (interface{}, bool) {
	bb.mutex.RLock()
	defer bb.mutex.RUnlock()

	if entry, exists := bb.entries[key]; exists {
		return entry.Value.Value(), true
	}

	if bb.parent != nil {
		return bb.parent.Get(key)
	}

	return nil, false
}

// HasKey checks if a key exists in the blackboard
func (bb *Blackboard) HasKey(key string) bool {
	bb.mutex.RLock()
	defer bb.mutex.RUnlock()

	if _, exists := bb.entries[key]; exists {
		return true
	}

	if bb.parent != nil {
		return bb.parent.HasKey(key)
	}

	return false
}

// Clear removes all entries from the blackboard
func (bb *Blackboard) Clear() {
	bb.mutex.Lock()
	defer bb.mutex.Unlock()
	bb.entries = make(map[string]Entry)
}

// RegisterPort registers a port with the blackboard
func (bb *Blackboard) RegisterPort(name string, direction PortDirection, typeName string, description string) {
	if bb.portInfo == nil {
		bb.portInfo = make(map[string]PortInfo)
	}
	bb.portInfo[name] = PortInfo{
		Direction:   direction,
		TypeName:    typeName,
		Description: description,
	}
}

// GetPortInfo retrieves port information
func (bb *Blackboard) GetPortInfo(name string) (PortInfo, bool) {
	if bb.portInfo == nil {
		return PortInfo{}, false
	}
	info, exists := bb.portInfo[name]
	return info, exists
}

// GetEntry retrieves a blackboard entry
func (bb *Blackboard) GetEntry(key string) *Entry {
	bb.mutex.RLock()
	defer bb.mutex.RUnlock()

	if entry, exists := bb.entries[key]; exists {
		return &entry
	}

	if bb.parent != nil {
		return bb.parent.GetEntry(key)
	}

	return nil
}
