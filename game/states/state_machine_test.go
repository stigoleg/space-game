package states

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

// mockState is a mock implementation of the State interface for testing
type mockState struct {
	stateType   StateType
	enterCount  int
	exitCount   int
	updateCount int
}

func (m *mockState) GetType() StateType {
	return m.stateType
}

func (m *mockState) Enter(prevState StateType) {
	m.enterCount++
}

func (m *mockState) Exit(nextState StateType) {
	m.exitCount++
}

func (m *mockState) Update() error {
	m.updateCount++
	return nil
}

func (m *mockState) Draw(screen *ebiten.Image) {
	// No-op for testing
}

func TestStateMachineCreation(t *testing.T) {
	sm := NewStateMachine()
	if sm == nil {
		t.Fatal("NewStateMachine() returned nil")
	}
	if sm.states == nil {
		t.Error("StateMachine.states map is nil")
	}
	if sm.transitions == nil {
		t.Error("StateMachine.transitions map is nil")
	}
}

func TestRegisterState(t *testing.T) {
	sm := NewStateMachine()
	mockMenu := &mockState{stateType: TypeMenu}

	sm.RegisterState(mockMenu)

	if len(sm.states) != 1 {
		t.Errorf("Expected 1 registered state, got %d", len(sm.states))
	}

	state, ok := sm.states[TypeMenu]
	if !ok {
		t.Error("Menu state not found in registry")
	}
	if state != mockMenu {
		t.Error("Registered state doesn't match original")
	}
}

func TestAllowTransition(t *testing.T) {
	sm := NewStateMachine()

	sm.AllowTransition(TypeMenu, TypePlaying)

	transition := Transition{From: TypeMenu, To: TypePlaying}
	if !sm.transitions[transition] {
		t.Error("Transition from Menu to Playing not registered")
	}
}

func TestSetInitialState(t *testing.T) {
	sm := NewStateMachine()
	mockMenu := &mockState{stateType: TypeMenu}
	sm.RegisterState(mockMenu)

	err := sm.SetInitialState(TypeMenu)
	if err != nil {
		t.Errorf("SetInitialState failed: %v", err)
	}

	if sm.GetCurrentStateType() != TypeMenu {
		t.Errorf("Expected current state TypeMenu, got %v", sm.GetCurrentStateType())
	}

	// Enter should not be called for initial state
	if mockMenu.enterCount != 0 {
		t.Errorf("Enter called %d times, expected 0", mockMenu.enterCount)
	}
}

func TestTransitionTo(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	sm.AllowTransition(TypeMenu, TypePlaying)

	// Set initial state
	sm.SetInitialState(TypeMenu)

	// Transition to Playing
	err := sm.TransitionTo(TypePlaying)
	if err != nil {
		t.Errorf("TransitionTo failed: %v", err)
	}

	// Verify current state changed
	if sm.GetCurrentStateType() != TypePlaying {
		t.Errorf("Expected current state TypePlaying, got %v", sm.GetCurrentStateType())
	}

	// Verify exit and enter were called
	if mockMenu.exitCount != 1 {
		t.Errorf("Menu exit called %d times, expected 1", mockMenu.exitCount)
	}
	if mockPlaying.enterCount != 1 {
		t.Errorf("Playing enter called %d times, expected 1", mockPlaying.enterCount)
	}
}

func TestInvalidTransition(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	// Don't allow transition

	sm.SetInitialState(TypeMenu)

	// Try invalid transition
	err := sm.TransitionTo(TypePlaying)
	if err == nil {
		t.Error("Expected error for invalid transition, got nil")
	}

	// State should not change
	if sm.GetCurrentStateType() != TypeMenu {
		t.Error("State changed despite invalid transition")
	}
}

func TestConfigureDefaultTransitions(t *testing.T) {
	sm := NewStateMachine()
	sm.ConfigureDefaultTransitions()

	// Test key transitions
	testCases := []struct {
		from StateType
		to   StateType
		desc string
	}{
		{TypeMenu, TypePlaying, "Menu to Playing"},
		{TypePlaying, TypePaused, "Playing to Paused"},
		{TypePlaying, TypeGameOver, "Playing to GameOver"},
		{TypePaused, TypePlaying, "Paused to Playing"},
		{TypePaused, TypeMenu, "Paused to Menu"},
		{TypeGameOver, TypeMenu, "GameOver to Menu"},
		{TypeGameOver, TypePlaying, "GameOver to Playing (retry)"},
	}

	for _, tc := range testCases {
		transition := Transition{From: tc.from, To: tc.to}
		if !sm.transitions[transition] {
			t.Errorf("Transition %s not configured", tc.desc)
		}
	}
}

func TestTransitionHooks(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	sm.AllowTransition(TypeMenu, TypePlaying)

	beforeCalled := false
	afterCalled := false

	sm.SetOnBeforeTransition(func(from, to StateType) {
		beforeCalled = true
		if from != TypeMenu || to != TypePlaying {
			t.Error("Before hook received wrong state types")
		}
	})

	sm.SetOnAfterTransition(func(from, to StateType) {
		afterCalled = true
		if from != TypeMenu || to != TypePlaying {
			t.Error("After hook received wrong state types")
		}
	})

	sm.SetInitialState(TypeMenu)
	sm.TransitionTo(TypePlaying)

	if !beforeCalled {
		t.Error("Before transition hook not called")
	}
	if !afterCalled {
		t.Error("After transition hook not called")
	}
}

func TestGetHistory(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	sm.AllowTransition(TypeMenu, TypePlaying)

	sm.SetInitialState(TypeMenu)
	sm.TransitionTo(TypePlaying)

	history := sm.GetHistory()
	if len(history) != 2 {
		t.Errorf("Expected history length 2, got %d", len(history))
	}
	if history[0] != TypeMenu {
		t.Error("First history entry should be Menu")
	}
	if history[1] != TypePlaying {
		t.Error("Second history entry should be Playing")
	}
}

func TestGetPreviousStateType(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	sm.AllowTransition(TypeMenu, TypePlaying)

	sm.SetInitialState(TypeMenu)

	// No previous state initially
	prev := sm.GetPreviousStateType()
	if prev != -1 {
		t.Errorf("Expected -1 for no previous state, got %v", prev)
	}

	sm.TransitionTo(TypePlaying)

	// Previous state should be Menu
	prev = sm.GetPreviousStateType()
	if prev != TypeMenu {
		t.Errorf("Expected previous state Menu, got %v", prev)
	}
}

func TestCanTransitionTo(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	mockPlaying := &mockState{stateType: TypePlaying}

	sm.RegisterState(mockMenu)
	sm.RegisterState(mockPlaying)
	sm.AllowTransition(TypeMenu, TypePlaying)

	sm.SetInitialState(TypeMenu)

	// Valid transition
	if !sm.CanTransitionTo(TypePlaying) {
		t.Error("CanTransitionTo returned false for valid transition")
	}

	// Invalid transition
	if sm.CanTransitionTo(TypeGameOver) {
		t.Error("CanTransitionTo returned true for invalid transition")
	}
}

func TestUpdate(t *testing.T) {
	sm := NewStateMachine()

	mockMenu := &mockState{stateType: TypeMenu}
	sm.RegisterState(mockMenu)
	sm.SetInitialState(TypeMenu)

	err := sm.Update()
	if err != nil {
		t.Errorf("Update failed: %v", err)
	}

	if mockMenu.updateCount != 1 {
		t.Errorf("State Update called %d times, expected 1", mockMenu.updateCount)
	}
}
