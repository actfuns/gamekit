package behavior_tree

import (
	"fmt"
	"time"
)

// NodeType enumerates the possible types of nodes
type NodeType int

const (
	NodeTypeUndefined NodeType = iota
	NodeTypeAction
	NodeTypeCondition
	NodeTypeControl
	NodeTypeDecorator
	NodeTypeSubtree
)

func (nt NodeType) String() string {
	switch nt {
	case NodeTypeAction:
		return "ACTION"
	case NodeTypeCondition:
		return "CONDITION"
	case NodeTypeControl:
		return "CONTROL"
	case NodeTypeDecorator:
		return "DECORATOR"
	case NodeTypeSubtree:
		return "SUBTREE"
	default:
		return "UNDEFINED"
	}
}

// NodeStatus enumerates the states every node can be in after execution
type NodeStatus int

const (
	NodeStatusIdle NodeStatus = iota
	NodeStatusRunning
	NodeStatusSuccess
	NodeStatusFailure
	NodeStatusSkipped
)

func (ns NodeStatus) String() string {
	switch ns {
	case NodeStatusIdle:
		return "IDLE"
	case NodeStatusRunning:
		return "RUNNING"
	case NodeStatusSuccess:
		return "SUCCESS"
	case NodeStatusFailure:
		return "FAILURE"
	case NodeStatusSkipped:
		return "SKIPPED"
	default:
		return "UNKNOWN"
	}
}

// IsStatusActive returns true if status is not IDLE or SKIPPED
func IsStatusActive(status NodeStatus) bool {
	return status != NodeStatusIdle && status != NodeStatusSkipped
}

// IsStatusCompleted returns true if status is SUCCESS or FAILURE
func IsStatusCompleted(status NodeStatus) bool {
	return status == NodeStatusSuccess || status == NodeStatusFailure
}

// PortDirection defines the direction of a port
type PortDirection int

const (
	PortDirectionInput PortDirection = iota
	PortDirectionOutput
	PortDirectionInOut
)

// Timestamp represents a timestamp for blackboard entries
type Timestamp time.Time

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value string
}

// KeyValueVector is a vector of key-value pairs
type KeyValueVector []KeyValue

// Expected represents a result that may contain a value or an error
type Expected[T any] struct {
	value T
	err   error
	valid bool
}

// NewExpected creates a new Expected with a value
func NewExpected[T any](value T) Expected[T] {
	return Expected[T]{value: value, valid: true}
}

// NewExpectedError creates a new Expected with an error
func NewExpectedError[T any](err error) Expected[T] {
	return Expected[T]{err: err, valid: false}
}

// Value returns the contained value
func (e Expected[T]) Value() T {
	return e.value
}

// Error returns the contained error
func (e Expected[T]) Error() error {
	return e.err
}

// Valid returns true if the Expected contains a valid value
func (e Expected[T]) Valid() bool {
	return e.valid
}

// AnyTypeAllowed is a marker type for any type allowed
type AnyTypeAllowed struct{}

// PortInfo represents information about a port
type PortInfo struct {
	Direction   PortDirection
	TypeName    string
	Description string
	DefaultValue string
}

// SetDefaultValue sets the default value for the port
func (p *PortInfo) SetDefaultValue(value interface{}) {
	// For simplicity, convert to string representation
	p.DefaultValue = fmt.Sprintf("%v", value)
}

// SetDescription sets the description for the port
func (p *PortInfo) SetDescription(description string) {
	p.Description = description
}

// PreCond defines pre-condition types
type PreCond int

const (
	PreCondFailureIf PreCond = iota
	PreCondSuccessIf
	PreCondSkipIf
	PreCondWhileTrue
	PreCondCount
)

// PostCond defines post-condition types  
type PostCond int

const (
	PostCondOnHalted PostCond = iota
	PostCondOnFailure
	PostCondOnSuccess
	PostCondAlways
	PostCondCount
)

// NodeConfig contains configuration for a tree node
type NodeConfig struct {
	Blackboard      *Blackboard
	Enums           map[string]int // ScriptingEnumsRegistry equivalent
	InputPorts      PortsRemapping
	OutputPorts     PortsRemapping
	OtherAttributes NonPortAttributes
	Manifest        TreeNodeManifest
	UID             uint16
	Path            string
	PreConditions   map[PreCond]string
	PostConditions  map[PostCond]string
}

// TreeNodeManifest contains information about a tree node
type TreeNodeManifest struct {
	Type           NodeType
	RegistrationID string
	Ports          PortsList
	Metadata       KeyValueVector
}

// PortsList defines a list of ports
type PortsList map[string]PortInfo

// PortsRemapping maps port names to their remapped values
type PortsRemapping map[string]string

// NonPortAttributes maps attribute names to their values
type NonPortAttributes map[string]string

