package fsm

import "fmt"

// BaseHState is a state that can contain other states (composite state)
type BaseHState struct {
	*BaseState
	subMachine *BaseFSM
}

// NewBaseHState creates a new hierarchical state instance
func NewBaseHState(id StateId, name string) *BaseHState {
	return &BaseHState{
		BaseState:  NewBaseState(id, name),
		subMachine: NewBaseFSM(name),
	}
}

// SetInitial sets the initial state of the sub-machine
func (hs *BaseHState) SetInitial(id StateId) {
	hs.subMachine.SetInitial(id)
}

// Current returns the current state of the sub-machine
func (hs *BaseHState) Current() State {
	return hs.subMachine.Current()
}

// Previous returns the previous state of the sub-machine
func (hs *BaseHState) Previous() State {
	return hs.subMachine.Previous()
}

// States returns all states in the sub-machine
func (hs *BaseHState) States() map[StateId]State {
	return hs.subMachine.States()
}

// Add adds a state to the machine
func (hs *BaseHState) Add(state State) {
	hs.subMachine.Add(state)
}

// Change changes to a specific state by ID in the sub-machine
func (hs *BaseHState) Change(id StateId) State {
	return hs.subMachine.Change(id)
}

// Contains checks if the sub-machine contains a specific state
func (hs *BaseHState) Contains(id StateId) bool {
	return hs.subMachine.Contains(id)
}

// Get retrieves a specific state by ID from the sub-machine
func (hs *BaseHState) Get(id StateId) State {
	return hs.subMachine.Get(id)
}

// IsCurrent checks if the given state ID is the current state in the sub-machine
func (hs *BaseHState) IsCurrent(id StateId) bool {
	return hs.subMachine.IsCurrent(id)
}

// IsPrevious checks if the given state ID is the previous state in the sub-machine
func (hs *BaseHState) IsPrevious(id StateId) bool {
	return hs.subMachine.IsPrevious(id)
}

// Remove removes a state from the sub-machine
func (hs *BaseHState) Remove(id StateId) {
	hs.subMachine.Remove(id)
}

// Back transitions to the previous state in the sub-machine
func (hs *BaseHState) Back() {
	hs.subMachine.Back()
}

// HandleEvent handles events for the entire sub-machine
func (hs *BaseHState) HandleEvent(eventType EventType) {
	// Handle the event in the current sub-state
	if hs.subMachine.Current() != nil {
		hs.subMachine.Current().HandleEvent(eventType)
	}
}

// Clear removes all states from the sub-machine
func (hs *BaseHState) Clear() {
	hs.subMachine.Clear()
}

// Init initializes this hierarchical state and its sub-states
func (hs *BaseHState) Init() {
	for _, state := range hs.subMachine.states {
		if setter, ok := state.(MachineSetter); ok {
			setter.SetSuperMachine(hs.superMachine)
		}

		state.Init()
	}
}

// Enter is called when entering this hstate state
func (hs *BaseHState) Enter() {
	if hs.subMachine.initial != nil {
		hs.subMachine.Change(hs.subMachine.initial.Id())
	}
}

// Exit is called when exiting this hstate state
func (hs *BaseHState) Exit() {
	if hs.subMachine.current != nil {
		hs.subMachine.current.Exit()
		hs.subMachine.current = nil
	}
}

// String returns a string representation of the hstate state
func (hs *BaseHState) String() string {
	return fmt.Sprintf("HState(%d, %s, sub=%v)",
		hs.Id(), hs.Name(), hs.subMachine)
}
