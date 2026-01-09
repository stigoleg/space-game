package core

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"stellar-siege/game/entities"
)

// InputAction represents an action triggered by user input
type InputAction int

const (
	InputActionNone InputAction = iota
	InputActionPause
	InputActionCycleWeapon
	InputActionSwitchWeapon
	InputActionActivateAbility
	InputActionShoot
)

// InputEvent represents a single input event with associated data
type InputEvent struct {
	Action      InputAction
	WeaponType  entities.WeaponType  // For InputActionSwitchWeapon
	AbilityType entities.AbilityType // For InputActionActivateAbility
}

// InputHandler handles all game input and converts it to actions
type InputHandler struct {
	// Configuration
	weaponTypes []entities.WeaponType
	abilityKeys map[ebiten.Key]entities.AbilityType
}

// NewInputHandler creates a new input handler
func NewInputHandler() *InputHandler {
	return &InputHandler{
		weaponTypes: []entities.WeaponType{
			entities.WeaponTypeSpread,
			entities.WeaponTypeBlaster,
			entities.WeaponTypeFollowingRocket,
			entities.WeaponTypeChainLightning,
			entities.WeaponTypeFlamethrower,
			entities.WeaponTypeIonBeam,
			entities.WeaponTypeLaser,
			entities.WeaponTypeShotgun,
			entities.WeaponTypePlasma,
		},
		abilityKeys: map[ebiten.Key]entities.AbilityType{
			ebiten.KeyQ: entities.AbilityTypeDash,
			ebiten.KeyE: entities.AbilityTypeSlowTime,
			ebiten.KeyR: entities.AbilityTypeBarrier,
			ebiten.KeyF: entities.AbilityTypeWeaponBoost,
			ebiten.KeyG: entities.AbilityTypeEMPPulse,
			ebiten.KeyH: entities.AbilityTypeOrbitalShield,
		},
	}
}

// PollGameplayInput polls all gameplay input and returns a list of input events
// This is called during the game's playing state
func (h *InputHandler) PollGameplayInput() []InputEvent {
	var events []InputEvent

	// Check for pause
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		events = append(events, InputEvent{Action: InputActionPause})
		return events // Return immediately on pause
	}

	// Check for weapon cycle (Tab key)
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		events = append(events, InputEvent{Action: InputActionCycleWeapon})
	}

	// Check for direct weapon selection (1-9 keys)
	for i, wt := range h.weaponTypes {
		if i >= 9 {
			break // Only 1-9 keys
		}
		key := ebiten.Key(int(ebiten.Key1) + i)
		if inpututil.IsKeyJustPressed(key) {
			events = append(events, InputEvent{
				Action:     InputActionSwitchWeapon,
				WeaponType: wt,
			})
		}
	}

	// Check for ability activation (Q, E, R, F, G, H keys)
	for key, abilityType := range h.abilityKeys {
		if inpututil.IsKeyJustPressed(key) {
			events = append(events, InputEvent{
				Action:      InputActionActivateAbility,
				AbilityType: abilityType,
			})
		}
	}

	// Check for shooting (Space or left mouse button)
	if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		events = append(events, InputEvent{Action: InputActionShoot})
	}

	return events
}

// IsPausePressed returns true if the pause key is pressed (for menu navigation)
func (h *InputHandler) IsPausePressed() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP)
}

// IsShootPressed returns true if the shoot input is pressed
func (h *InputHandler) IsShootPressed() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

// GetWeaponCycleList returns the full list of weapon types for cycling
func (h *InputHandler) GetWeaponCycleList() []entities.WeaponType {
	// Return a list with all weapons including those not mapped to number keys
	return []entities.WeaponType{
		entities.WeaponTypeSpread,
		entities.WeaponTypeBlaster,
		entities.WeaponTypeFollowingRocket,
		entities.WeaponTypeChainLightning,
		entities.WeaponTypeFlamethrower,
		entities.WeaponTypeIonBeam,
		entities.WeaponTypeLaser,
		entities.WeaponTypeShotgun,
		entities.WeaponTypePlasma,
		entities.WeaponTypeHoming,
		entities.WeaponTypeRailgun,
	}
}

// GetDashDirection returns the direction for a dash ability based on WASD/Arrow keys
// Returns (dx, dy) normalized direction. If no input, returns (0, -1) for upward dash.
func (h *InputHandler) GetDashDirection() (float64, float64) {
	dx, dy := 0.0, 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx = 1
	}

	// If no movement input, dash upward
	if dx == 0 && dy == 0 {
		dy = -1
	}

	return dx, dy
}
