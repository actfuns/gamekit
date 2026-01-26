package fsm

// StateId represents a unique identifier for a state in the finite state machine
type StateId int

// EventType represents different types of events that can be handled by states
type EventType int

// MachineSetter interface for states that can have their machine reference set
type MachineSetter interface {
	SetMachine(m FSM)
	SetSuperMachine(m FSM)
}

// State represents a state in the finite state machine
type State interface {
	// Id returns the unique identifier for this state
	Id() StateId

	// Name returns the name of this state
	Name() string

	// Machine returns the machine that this state belongs to
	Machine() FSM

	// SuperMachine returns the super machine that this state belongs to
	SuperMachine() FSM

	// Enter is called every time this state is activated
	Enter()

	// Exit is called every time this state is deactivated
	Exit()

	// HandleEvent handles different types of events
	HandleEvent(eventType EventType)

	// Init initializes this state instance
	Init()
}

// FSM represents a finite state machine
type FSM interface {
	// Current returns the current state
	Current() State

	// Previous returns the previous state
	Previous() State

	// States returns a copy of all states in the machine
	States() map[StateId]State

	// Add adds a state to the machine
	Add(state State)

	// Change changes to a specific state by ID
	Change(id StateId) State

	// Contains checks if the machine contains a specific state
	Contains(id StateId) bool

	// Get retrieves a specific state by ID
	Get(id StateId) State

	// IsCurrent checks if the given state ID is the current state
	IsCurrent(id StateId) bool

	// IsPrevious checks if the given state ID is the previous state
	IsPrevious(id StateId) bool

	// Remove removes a state from the machine
	Remove(id StateId)

	// Back transitions to the previous state
	Back()

	// HandleEvent handles events for the entire machine
	HandleEvent(eventType EventType)

	// Clear removes all states from the machine
	Clear()

	// Init sets the initial state as the current state
	Init()
}

// HState represents a hierarchical state that can contain other states
type HState interface {
	State
	FSM
}
