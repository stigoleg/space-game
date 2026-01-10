package states

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// StateType represents the type of game state
type StateType int

const (
	TypeMenu StateType = iota
	TypePlaying
	TypePaused
	TypeGameOver
)

// maxHistorySize limits the state transition history to prevent unbounded growth
const maxHistorySize = 100

// String returns the string representation of the state type
func (s StateType) String() string {
	switch s {
	case TypeMenu:
		return "Menu"
	case TypePlaying:
		return "Playing"
	case TypePaused:
		return "Paused"
	case TypeGameOver:
		return "GameOver"
	default:
		return "Unknown"
	}
}

// State represents a game state with update and draw logic
type State interface {
	// GetType returns the type of this state
	GetType() StateType

	// Enter is called when transitioning into this state
	Enter(prevState StateType)

	// Exit is called when transitioning out of this state
	Exit(nextState StateType)

	// Update updates the state logic
	Update() error

	// Draw renders the state
	Draw(screen *ebiten.Image)
}

// Transition represents a valid state transition
type Transition struct {
	From StateType
	To   StateType
}

// StateMachine manages game state transitions
type StateMachine struct {
	// Current state
	currentState State

	// State registry
	states map[StateType]State

	// Valid transitions
	transitions map[Transition]bool

	// Transition hooks
	onBeforeTransition func(from, to StateType)
	onAfterTransition  func(from, to StateType)

	// History for debugging
	history []StateType
}

// NewStateMachine creates a new state machine
func NewStateMachine() *StateMachine {
	return &StateMachine{
		states:      make(map[StateType]State),
		transitions: make(map[Transition]bool),
		history:     make([]StateType, 0, 10),
	}
}

// RegisterState registers a state with the state machine
func (sm *StateMachine) RegisterState(state State) {
	sm.states[state.GetType()] = state
}

// AllowTransition defines a valid state transition
func (sm *StateMachine) AllowTransition(from, to StateType) {
	sm.transitions[Transition{From: from, To: to}] = true
}

// SetOnBeforeTransition sets a callback to run before state transitions
func (sm *StateMachine) SetOnBeforeTransition(fn func(from, to StateType)) {
	sm.onBeforeTransition = fn
}

// SetOnAfterTransition sets a callback to run after state transitions
func (sm *StateMachine) SetOnAfterTransition(fn func(from, to StateType)) {
	sm.onAfterTransition = fn
}

// SetInitialState sets the initial state without calling Enter
func (sm *StateMachine) SetInitialState(stateType StateType) error {
	state, ok := sm.states[stateType]
	if !ok {
		return fmt.Errorf("state %v not registered", stateType)
	}

	sm.currentState = state
	sm.appendToHistory(stateType)
	return nil
}

// TransitionTo transitions to a new state
func (sm *StateMachine) TransitionTo(stateType StateType) error {
	// Get the target state
	targetState, ok := sm.states[stateType]
	if !ok {
		return fmt.Errorf("state %v not registered", stateType)
	}

	// Check if transition is allowed
	currentType := sm.GetCurrentStateType()
	transition := Transition{From: currentType, To: stateType}

	if !sm.transitions[transition] {
		return fmt.Errorf("transition from %v to %v not allowed", currentType, stateType)
	}

	// Run before transition hook
	if sm.onBeforeTransition != nil {
		sm.onBeforeTransition(currentType, stateType)
	}

	// Exit current state
	if sm.currentState != nil {
		sm.currentState.Exit(stateType)
	}

	// Enter new state
	targetState.Enter(currentType)

	// Update current state
	sm.currentState = targetState
	sm.appendToHistory(stateType)

	// Run after transition hook
	if sm.onAfterTransition != nil {
		sm.onAfterTransition(currentType, stateType)
	}

	return nil
}

// GetCurrentStateType returns the type of the current state
func (sm *StateMachine) GetCurrentStateType() StateType {
	if sm.currentState == nil {
		return -1 // Invalid state
	}
	return sm.currentState.GetType()
}

// GetCurrentState returns the current state
func (sm *StateMachine) GetCurrentState() State {
	return sm.currentState
}

// Update updates the current state
func (sm *StateMachine) Update() error {
	if sm.currentState == nil {
		return fmt.Errorf("no current state")
	}
	return sm.currentState.Update()
}

// Draw draws the current state
func (sm *StateMachine) Draw(screen *ebiten.Image) {
	if sm.currentState != nil {
		sm.currentState.Draw(screen)
	}
}

// GetHistory returns the state transition history
func (sm *StateMachine) GetHistory() []StateType {
	return sm.history
}

// GetPreviousStateType returns the previous state type (if any)
func (sm *StateMachine) GetPreviousStateType() StateType {
	if len(sm.history) < 2 {
		return -1 // No previous state
	}
	return sm.history[len(sm.history)-2]
}

// CanTransitionTo checks if transitioning to a state is valid
func (sm *StateMachine) CanTransitionTo(stateType StateType) bool {
	currentType := sm.GetCurrentStateType()
	transition := Transition{From: currentType, To: stateType}
	return sm.transitions[transition]
}

// ConfigureDefaultTransitions sets up the standard game state transitions
func (sm *StateMachine) ConfigureDefaultTransitions() {
	// From Menu
	sm.AllowTransition(TypeMenu, TypePlaying)

	// From Playing
	sm.AllowTransition(TypePlaying, TypePaused)
	sm.AllowTransition(TypePlaying, TypeGameOver)

	// From Paused
	sm.AllowTransition(TypePaused, TypePlaying)
	sm.AllowTransition(TypePaused, TypeMenu)

	// From GameOver
	sm.AllowTransition(TypeGameOver, TypeMenu)
	sm.AllowTransition(TypeGameOver, TypePlaying) // For retry
}

// appendToHistory adds a state to history, keeping it bounded to maxHistorySize
func (sm *StateMachine) appendToHistory(stateType StateType) {
	sm.history = append(sm.history, stateType)
	if len(sm.history) > maxHistorySize {
		// Keep the last maxHistorySize entries by shifting
		copy(sm.history, sm.history[len(sm.history)-maxHistorySize:])
		sm.history = sm.history[:maxHistorySize]
	}
}

// ClearHistory resets the state transition history
func (sm *StateMachine) ClearHistory() {
	sm.history = sm.history[:0]
}
