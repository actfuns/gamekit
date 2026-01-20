package behavior_tree

import "time"

// DecoratorNode is the base class for all decorator nodes
type DecoratorNode struct {
	TreeNodeBase
}

// NewDecoratorNode creates a new decorator node
func NewDecoratorNode(name string, config NodeConfig) *DecoratorNode {
	return &DecoratorNode{
		TreeNodeBase: *NewTreeNode(name, config),
	}
}

// Type returns the node type
func (dn *DecoratorNode) Type() NodeType {
	return NodeTypeDecorator
}

// InverterNode inverts the result of its child
type InverterNode struct {
	DecoratorNode
}

// NewInverterNode creates a new inverter node
func NewInverterNode(name string, config NodeConfig) *InverterNode {
	return &InverterNode{
		DecoratorNode: *NewDecoratorNode(name, config),
	}
}

// Tick executes the inverter logic
func (in *InverterNode) Tick() NodeStatus {
	children := in.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	child := children[0]
	// 直接调用child的Tick，而不是ExecuteTick
	// 这样装饰器可以控制子节点的状态变化
	status := child.Tick()

	switch status {
	case NodeStatusSuccess:
		return NodeStatusFailure
	case NodeStatusFailure:
		return NodeStatusSuccess
	default:
		return status
	}
}

// RetryNode retries its child until success or max attempts
type RetryNode struct {
	DecoratorNode
	maxAttempts int
	attempts    int
}

// NewRetryNode creates a new retry node
func NewRetryNode(name string, config NodeConfig, maxAttempts int) *RetryNode {
	return &RetryNode{
		DecoratorNode: *NewDecoratorNode(name, config),
		maxAttempts:   maxAttempts,
		attempts:      0,
	}
}

// Tick executes the retry logic
func (rn *RetryNode) Tick() NodeStatus {
	children := rn.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	switch status {
	case NodeStatusSuccess:
		rn.attempts = 0
		return NodeStatusSuccess
	case NodeStatusFailure:
		rn.attempts++
		if rn.attempts >= rn.maxAttempts {
			rn.attempts = 0
			return NodeStatusFailure
		}
		// Reset child and retry
		child.HaltAndReset()
		return NodeStatusRunning
	default:
		return status
	}
}

// Halt handles halting the retry
func (rn *RetryNode) Halt() {
	children := rn.Children()
	if len(children) > 0 {
		children[0].HaltAndReset()
	}
	rn.attempts = 0
}

// RepeatNode repeats its child until failure or max repetitions
type RepeatNode struct {
	DecoratorNode
	maxRepetitions int
	repetitions    int
}

// NewRepeatNode creates a new repeat node
func NewRepeatNode(name string, config NodeConfig, maxRepetitions int) *RepeatNode {
	return &RepeatNode{
		DecoratorNode:  *NewDecoratorNode(name, config),
		maxRepetitions: maxRepetitions,
		repetitions:    0,
	}
}

// Tick executes the repeat logic
func (rn *RepeatNode) Tick() NodeStatus {
	children := rn.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	child := children[0]
	status := child.Tick()

	switch status {
	case NodeStatusFailure:
		rn.repetitions = 0
		return NodeStatusFailure
	case NodeStatusSuccess:
		rn.repetitions++
		if rn.repetitions >= rn.maxRepetitions {
			rn.repetitions = 0
			return NodeStatusSuccess
		}
		// Reset child and repeat
		child.HaltAndReset()
		return NodeStatusRunning
	default:
		return status
	}
}

// Halt handles halting the repeat
func (rn *RepeatNode) Halt() {
	children := rn.Children()
	if len(children) > 0 {
		children[0].HaltAndReset()
	}
	rn.repetitions = 0
}

// TimeoutNode adds a timeout to its child
type TimeoutNode struct {
	DecoratorNode
	timeoutMs int64
	startTime time.Time
}

// NewTimeoutNode creates a new timeout node
func NewTimeoutNode(name string, config NodeConfig, timeoutMs int64) *TimeoutNode {
	return &TimeoutNode{
		DecoratorNode: *NewDecoratorNode(name, config),
		timeoutMs:     timeoutMs,
	}
}

// Tick executes the timeout logic
func (tn *TimeoutNode) Tick() NodeStatus {
	children := tn.Children()
	if len(children) == 0 {
		return NodeStatusFailure
	}

	currentTime := time.Now()

	if tn.Status() == NodeStatusIdle {
		tn.startTime = currentTime
	}

	child := children[0]
	status := child.Tick()

	if status == NodeStatusRunning {
		elapsed := currentTime.Sub(tn.startTime)
		if int64(elapsed/time.Millisecond) >= tn.timeoutMs {
			child.Halt()
			return NodeStatusFailure
		}
		return NodeStatusRunning
	}

	return status
}
