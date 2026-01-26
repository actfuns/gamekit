package fsm

import (
	"fmt"
	"maps"
)

// BaseFSM is the finite state machine implementation
type BaseFSM struct {
	current  State
	initial  State
	previous State
	states   map[StateId]State
	name     string
}

// NewBaseFSM creates a new finite state machine instance
func NewBaseFSM(name string) *BaseFSM {
	return &BaseFSM{
		states: make(map[StateId]State),
		name:   name,
	}
}

// Name returns the machine name
func (m *BaseFSM) Name() string {
	return m.name
}

// Current returns the current state
func (m *BaseFSM) Current() State {
	return m.current
}

// Previous returns the previous state
func (m *BaseFSM) Previous() State {
	return m.previous
}

// States returns all states in the machine
func (m *BaseFSM) States() map[StateId]State {
	// Return a shallow copy to avoid exposing internal map
	copyMap := make(map[StateId]State, len(m.states))
	maps.Copy(copyMap, m.states)
	return copyMap
}

// Add adds a state to the machine
func (m *BaseFSM) Add(state State) {
	m.states[state.Id()] = state

	if setter, ok := state.(MachineSetter); ok {
		setter.SetMachine(m)
		setter.SetSuperMachine(m)
	}
}

// Change changes to a specific state by ID
func (m *BaseFSM) Change(id StateId) State {
	if nextState, exists := m.states[id]; exists {
		// Exit current state if it exists
		if m.current != nil {
			m.current.Exit()
		}

		// Update previous state
		m.previous = m.current

		// Set new current state
		m.current = nextState

		// Enter new state
		if m.current != nil {
			m.current.Enter()
		}
		return nextState
	}

	return nil
}

// Contains checks if the machine contains a specific state
func (m *BaseFSM) Contains(id StateId) bool {
	_, exists := m.states[id]
	return exists
}

// Get retrieves a specific state by ID
func (m *BaseFSM) Get(id StateId) State {
	return m.states[id]
}

// IsCurrent checks if the given state ID is the current state
func (m *BaseFSM) IsCurrent(id StateId) bool {
	if m.current != nil {
		return m.current.Id() == id
	}
	return false
}

// IsPrevious checks if the given state ID is the previous state
func (m *BaseFSM) IsPrevious(id StateId) bool {
	if m.previous != nil {
		return m.previous.Id() == id
	}
	return false
}

// Remove removes a state from the machine
func (m *BaseFSM) Remove(id StateId) {
	// Remove the state from the registry
	if s, exists := m.states[id]; exists {
		// If the removed state is current/previous/initial, clear those references
		if m.current != nil && m.current == s {
			m.current = nil
		}
		if m.previous != nil && m.previous == s {
			m.previous = nil
		}
		if m.initial != nil && m.initial == s {
			m.initial = nil
		}
	}

	delete(m.states, id)
}

// Back transitions to the previous state
func (m *BaseFSM) Back() {
	if m.previous != nil {
		m.Change(m.previous.Id())
	}
}

// HandleEvent handles events for the entire machine
func (m *BaseFSM) HandleEvent(eventType EventType) {
	if m.current != nil {
		m.current.HandleEvent(eventType)
	}
}

// Clear removes all states from the machine
func (m *BaseFSM) Clear() {
	m.states = make(map[StateId]State)
	m.current = nil
	m.previous = nil
	m.initial = nil
}

// Init sets the initial state as the current state
func (m *BaseFSM) Init() {
	for _, state := range m.states {
		state.Init()
	}

	if m.initial != nil {
		// Exit current state if it exists and is different from the initial state
		if m.current != nil && m.current != m.initial {
			m.current.Exit()
		}

		// Set current state to initial state
		m.current = m.initial

		// Enter the initial state
		if m.current != nil {
			m.current.Enter()
		}
	}
}

// SetInitial sets the initial state
func (m *BaseFSM) SetInitial(id StateId) {
	if state, exists := m.states[id]; exists {
		m.initial = state
	}
}

// String returns a string representation of the machine
func (m *BaseFSM) String() string {
	return fmt.Sprintf("FSM(%s, current=%v, initial=%v)",
		m.name,
		m.current,
		m.initial)
}
