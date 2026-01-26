package fsm

import "fmt"

// BaseState is the base state implementation
type BaseState struct {
	id           StateId
	name         string
	machine      FSM
	superMachine FSM
}

// NewBaseState creates a new base state instance
func NewBaseState(id StateId, name string) *BaseState {
	return &BaseState{
		id:   id,
		name: name,
	}
}

// Id returns the unique identifier for this state
func (s *BaseState) Id() StateId {
	return s.id
}

// Name returns the name of this state
func (s *BaseState) Name() string {
	return s.name
}

// SetMachine sets the state machine this state belongs to
func (s *BaseState) SetMachine(m FSM) {
	s.machine = m
}

// SetSuperMachine sets the super machine this state belongs to
func (s *BaseState) SetSuperMachine(m FSM) {
	s.superMachine = m
}

// Machine returns the state machine this state belongs to
func (s *BaseState) Machine() FSM {
	return s.machine
}

// SuperMachine returns the super machine this state belongs to
func (s *BaseState) SuperMachine() FSM {
	return s.superMachine
}

// Enter is called every time this state is activated
func (s *BaseState) Enter() {
	// Default implementation does nothing
}

// Exit is called every time this state is deactivated
func (s *BaseState) Exit() {
	// Default implementation does nothing
}

// Init initializes this state instance
func (s *BaseState) Init() {
	// Default implementation does nothing
}

// HandleEvent handles different types of events
func (s *BaseState) HandleEvent(eventType EventType) {
	// Default implementation does nothing
}

// String returns a string representation of the state
func (s *BaseState) String() string {
	return fmt.Sprintf("State(%d, %s)", s.id, s.name)
}
