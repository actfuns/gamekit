package fsm

import "testing"

////////////////////////////////////////////////////////////
// Test Helpers
////////////////////////////////////////////////////////////

// MockState is a simple test double for State
type MockState struct {
	*BaseState
	entered bool
	exited  bool
}

func NewMockState(id StateId, name string) *MockState {
	return &MockState{
		BaseState: NewBaseState(id, name),
	}
}

func (m *MockState) Enter() {
	m.entered = true
}

func (m *MockState) Exit() {
	m.exited = true
}

func (m *MockState) HandleEvent(eventType EventType) {}

////////////////////////////////////////////////////////////
// BaseFSM – 基础能力测试
////////////////////////////////////////////////////////////

func TestBaseFSM_StateIDOperations(t *testing.T) {
	fsm := NewBaseFSM("TestFSM")

	s1 := NewMockState(1, "State1")
	s2 := NewMockState(2, "State2")
	s3 := NewMockState(3, "State3")

	fsm.Add(s1)
	fsm.Add(s2)
	fsm.Add(s3)

	// Initial
	fsm.SetInitial(1)
	fsm.Init()

	if !fsm.IsCurrent(1) {
		t.Fatalf("expected initial state 1, got %v", fsm.Current())
	}

	// Change
	fsm.Change(2)
	if !fsm.IsCurrent(2) {
		t.Fatalf("expected current state 2, got %v", fsm.Current())
	}

	// Back
	fsm.Change(3)
	fsm.Back()
	if !fsm.IsCurrent(2) {
		t.Fatalf("expected back to state 2, got %v", fsm.Current())
	}
}

func TestBaseFSM_StateLifecycle(t *testing.T) {
	fsm := NewBaseFSM("LifecycleFSM")

	idle := NewMockState(1, "Idle")
	walk := NewMockState(2, "Walk")

	fsm.Add(idle)
	fsm.Add(walk)

	fsm.SetInitial(1)
	fsm.Init()

	if !idle.entered {
		t.Fatalf("Idle.Enter was not called")
	}

	fsm.Change(2)

	if !idle.exited {
		t.Fatalf("Idle.Exit was not called")
	}
	if !walk.entered {
		t.Fatalf("Walk.Enter was not called")
	}
}

////////////////////////////////////////////////////////////
// Event Handling
////////////////////////////////////////////////////////////

func TestFSM_EventHandling(t *testing.T) {
	fsm := NewBaseFSM("EventFSM")

	s1 := NewMockState(1, "State1")
	fsm.Add(s1)

	fsm.SetInitial(1)
	fsm.Init()

	// Should not panic
	fsm.HandleEvent(EventType(100))
}

////////////////////////////////////////////////////////////
// HState – 层级状态测试
////////////////////////////////////////////////////////////

func TestHState_Basic(t *testing.T) {
	h := NewBaseHState(10, "RootHState")

	sub1 := NewMockState(101, "Sub1")
	sub2 := NewMockState(102, "Sub2")

	h.Add(sub1)
	h.Add(sub2)
	h.SetInitial(101)

	h.Init()
	h.Enter()

	if h.Current() == nil {
		t.Fatalf("expected current sub-state, got nil")
	}

	if !h.IsCurrent(101) {
		t.Fatalf("expected current sub-state 101, got %v", h.Current())
	}
}

func TestHState_SubStateTransition(t *testing.T) {
	h := NewBaseHState(20, "Combat")

	a := NewMockState(201, "Attack")
	d := NewMockState(202, "Defend")

	h.Add(a)
	h.Add(d)
	h.SetInitial(201)

	h.Init()
	h.Enter()

	h.Change(202)

	if !h.IsCurrent(202) {
		t.Fatalf("expected current sub-state 202, got %v", h.Current())
	}
}

////////////////////////////////////////////////////////////
// WorkerFSM – 业务 FSM
////////////////////////////////////////////////////////////

type WorkState int

const (
	WorkIdle WorkState = iota
	WorkWalk
	WorkBuild
	WorkCarryGo
	WorkCarryBack
)

type WorkerState struct {
	*BaseState
	work WorkState
}

func NewWorkerState(id StateId, name string, work WorkState) *WorkerState {
	return &WorkerState{
		BaseState: NewBaseState(id, name),
		work:      work,
	}
}

func (w *WorkerState) Enter()                          {}
func (w *WorkerState) Exit()                           {}
func (w *WorkerState) HandleEvent(eventType EventType) {}

type WorkerFSM struct {
	*BaseFSM
	currentWork WorkState
}

func NewWorkerFSM(name string) *WorkerFSM {
	return &WorkerFSM{
		BaseFSM: NewBaseFSM(name),
	}
}

func (w *WorkerFSM) SetWork(work WorkState) {
	w.currentWork = work
}

////////////////////////////////////////////////////////////
// WorkerFSM 测试
////////////////////////////////////////////////////////////

func TestWorkerFSM_BasicFlow(t *testing.T) {
	fsm := NewWorkerFSM("Worker")

	idle := NewWorkerState(0, "Idle", WorkIdle)
	walk := NewWorkerState(1, "Walk", WorkWalk)
	build := NewWorkerState(2, "Build", WorkBuild)

	fsm.Add(idle)
	fsm.Add(walk)
	fsm.Add(build)

	fsm.SetInitial(1)
	fsm.Init()

	if !fsm.IsCurrent(1) {
		t.Fatalf("expected initial Walk state")
	}

	fsm.Change(2)
	if !fsm.IsCurrent(2) {
		t.Fatalf("expected Build state")
	}

	fsm.Back()
	if !fsm.IsCurrent(1) {
		t.Fatalf("expected back to Walk state")
	}
}

////////////////////////////////////////////////////////////
// WorkerFSM + HState（完整业务示例）
////////////////////////////////////////////////////////////

func TestWorkerFSM_WithHState(t *testing.T) {
	fsm := NewWorkerFSM("WorkerMain")

	idle := NewWorkerState(0, "Idle", WorkIdle)
	walk := NewWorkerState(1, "Walk", WorkWalk)

	buildH := NewBaseHState(2, "BuildH")
	start := NewWorkerState(201, "Start", WorkBuild)
	progress := NewWorkerState(202, "Progress", WorkBuild)
	done := NewWorkerState(203, "Done", WorkBuild)

	buildH.Add(start)
	buildH.Add(progress)
	buildH.Add(done)
	buildH.SetInitial(201)

	fsm.Add(idle)
	fsm.Add(walk)
	fsm.Add(buildH)

	fsm.SetInitial(1)
	fsm.Init()

	fsm.Change(2)
	if !fsm.IsCurrent(2) {
		t.Fatalf("expected BuildH as current state")
	}

	build := fsm.Get(2).(*BaseHState)
	build.Change(202)

	if !build.IsCurrent(202) {
		t.Fatalf("expected Build Progress state")
	}

	fsm.Back()
	if !fsm.IsCurrent(1) {
		t.Fatalf("expected back to Walk state")
	}
}
