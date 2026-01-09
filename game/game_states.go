package game

import (
	"stellar-siege/game/states"

	"github.com/hajimehoshi/ebiten/v2"
)

// GameStateHandler wraps game state logic to work with the state machine
type GameStateHandler struct {
	stateType  states.StateType
	game       *Game
	updateFunc func(*Game) error
	drawFunc   func(*Game, *ebiten.Image)
}

// NewGameStateHandler creates a new game state handler
func NewGameStateHandler(stateType states.StateType, game *Game, updateFunc func(*Game) error, drawFunc func(*Game, *ebiten.Image)) *GameStateHandler {
	return &GameStateHandler{
		stateType:  stateType,
		game:       game,
		updateFunc: updateFunc,
		drawFunc:   drawFunc,
	}
}

// GetType returns the state type
func (h *GameStateHandler) GetType() states.StateType {
	return h.stateType
}

// Enter is called when entering this state
func (h *GameStateHandler) Enter(prevState states.StateType) {
	// State-specific entry logic can be added here
	switch h.stateType {
	case states.TypeMenu:
		// Menu entry
	case states.TypePlaying:
		// Playing entry
	case states.TypePaused:
		// Paused entry
	case states.TypeGameOver:
		// Game over entry
	}
}

// Exit is called when leaving this state
func (h *GameStateHandler) Exit(nextState states.StateType) {
	// State-specific exit logic can be added here
}

// Update updates the state
func (h *GameStateHandler) Update() error {
	return h.updateFunc(h.game)
}

// Draw renders the state
func (h *GameStateHandler) Draw(screen *ebiten.Image) {
	h.drawFunc(h.game, screen)
}

// initializeStateMachine sets up the state machine with all game states
func (g *Game) initializeStateMachine() {
	g.stateMachine = states.NewStateMachine()

	// Register all states
	g.stateMachine.RegisterState(NewGameStateHandler(
		states.TypeMenu,
		g,
		func(game *Game) error {
			game.updateMenu()
			return nil
		},
		func(game *Game, screen *ebiten.Image) {
			// Drawing is handled by main Draw() method
		},
	))

	g.stateMachine.RegisterState(NewGameStateHandler(
		states.TypePlaying,
		g,
		func(game *Game) error {
			game.updatePlaying()
			return nil
		},
		func(game *Game, screen *ebiten.Image) {
			// Drawing is handled by main Draw() method
		},
	))

	g.stateMachine.RegisterState(NewGameStateHandler(
		states.TypePaused,
		g,
		func(game *Game) error {
			game.updatePaused()
			return nil
		},
		func(game *Game, screen *ebiten.Image) {
			// Drawing is handled by main Draw() method
		},
	))

	g.stateMachine.RegisterState(NewGameStateHandler(
		states.TypeGameOver,
		g,
		func(game *Game) error {
			game.updateGameOver()
			return nil
		},
		func(game *Game, screen *ebiten.Image) {
			// Drawing is handled by main Draw() method
		},
	))

	// Configure valid transitions
	g.stateMachine.ConfigureDefaultTransitions()

	// Set transition hooks
	g.stateMachine.SetOnAfterTransition(func(from, to states.StateType) {
		// Play transition sound or effects
		if to == states.TypePlaying && from == states.TypeMenu {
			// Starting new game
		} else if to == states.TypePaused {
			// Game paused
		} else if to == states.TypeGameOver {
			// Game over
		}
	})

	// Set initial state
	g.stateMachine.SetInitialState(states.TypeMenu)
}

// transitionToState transitions to a new state using the state machine
func (g *Game) transitionToState(newState GameState) {
	// Map old GameState to new StateType
	var stateType states.StateType
	switch newState {
	case StateMenu:
		stateType = states.TypeMenu
	case StatePlaying:
		stateType = states.TypePlaying
	case StatePaused:
		stateType = states.TypePaused
	case StateGameOver:
		stateType = states.TypeGameOver
	default:
		return
	}

	// Use state machine if available
	if g.stateMachine != nil {
		if err := g.stateMachine.TransitionTo(stateType); err != nil {
			// Fallback to manual state change if transition fails
			g.state = newState
		} else {
			// Update the legacy state field to keep it in sync
			g.state = newState
		}
	} else {
		// Fallback to manual state change
		g.state = newState
	}
}
