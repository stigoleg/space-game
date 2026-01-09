package entities

// AbilityType represents different ability types
type AbilityType string

const (
	AbilityTypeDash          AbilityType = "dash"
	AbilityTypeSlowTime      AbilityType = "slow_time"
	AbilityTypeBarrier       AbilityType = "barrier"
	AbilityTypeWeaponBoost   AbilityType = "weapon_boost"
	AbilityTypeEMPPulse      AbilityType = "emp_pulse"
	AbilityTypeOrbitalShield AbilityType = "orbital_shield"
)

// Ability represents a player ability
type Ability struct {
	Type          AbilityType
	Name          string
	Description   string
	KeyBinding    string
	Cooldown      float64 // Seconds
	CooldownTimer float64 // Current countdown
	Available     bool
	ShieldCost    float64 // Shield points required to use
	IconEmoji     string
	Level         int // 1-3, affects power
}

// AbilityManager manages all player abilities
type AbilityManager struct {
	Abilities           map[AbilityType]*Ability
	ActiveAbilities     map[AbilityType]bool    // Currently active effects
	ActiveAbilityTimers map[AbilityType]float64 // Duration timers
}

// NewAbilityManager creates a new ability manager
func NewAbilityManager() *AbilityManager {
	am := &AbilityManager{
		Abilities:           make(map[AbilityType]*Ability),
		ActiveAbilities:     make(map[AbilityType]bool),
		ActiveAbilityTimers: make(map[AbilityType]float64),
	}

	// Initialize base abilities
	am.Abilities[AbilityTypeDash] = &Ability{
		Type:        AbilityTypeDash,
		Name:        "Dash",
		Description: "Quick dodge in any direction",
		KeyBinding:  "Q",
		Cooldown:    1.0,
		ShieldCost:  25,
		IconEmoji:   "💨",
		Level:       1,
	}

	am.Abilities[AbilityTypeSlowTime] = &Ability{
		Type:        AbilityTypeSlowTime,
		Name:        "Bullet Time",
		Description: "Slow time to 50% for 2 seconds",
		KeyBinding:  "E",
		Cooldown:    5.0,
		ShieldCost:  20,
		IconEmoji:   "⏱️",
		Level:       1,
	}

	am.Abilities[AbilityTypeBarrier] = &Ability{
		Type:        AbilityTypeBarrier,
		Name:        "Barrier",
		Description: "Create protective shield",
		KeyBinding:  "R",
		Cooldown:    3.0,
		ShieldCost:  30,
		IconEmoji:   "🛡️",
		Level:       1,
	}

	return am
}

// AddAbility adds a new ability to the manager
func (am *AbilityManager) AddAbility(abilityType AbilityType) bool {
	if _, exists := am.Abilities[abilityType]; exists {
		return false // Already have this ability
	}

	var ability *Ability
	switch abilityType {
	case AbilityTypeWeaponBoost:
		ability = &Ability{
			Type:        AbilityTypeWeaponBoost,
			Name:        "Weapon Overcharge",
			Description: "Double fire rate for 3 seconds",
			KeyBinding:  "F",
			Cooldown:    8.0,
			ShieldCost:  25,
			IconEmoji:   "⚡",
			Level:       1,
		}
	case AbilityTypeEMPPulse:
		ability = &Ability{
			Type:        AbilityTypeEMPPulse,
			Name:        "EMP Pulse",
			Description: "Stun all enemies for 1 second",
			KeyBinding:  "G",
			Cooldown:    10.0,
			ShieldCost:  40,
			IconEmoji:   "💥",
			Level:       1,
		}
	case AbilityTypeOrbitalShield:
		ability = &Ability{
			Type:        AbilityTypeOrbitalShield,
			Name:        "Orbital Defense",
			Description: "Projectiles orbit and protect",
			KeyBinding:  "H",
			Cooldown:    12.0,
			ShieldCost:  35,
			IconEmoji:   "🔵",
			Level:       1,
		}
	default:
		return false
	}

	am.Abilities[abilityType] = ability
	return true
}

// CanUseAbility checks if an ability can be used
func (am *AbilityManager) CanUseAbility(abilityType AbilityType) bool {
	ability, exists := am.Abilities[abilityType]
	if !exists {
		return false
	}
	return ability.CooldownTimer <= 0
}

// UseAbility activates an ability (caller must check CanUseAbility first)
func (am *AbilityManager) UseAbility(abilityType AbilityType) bool {
	if !am.CanUseAbility(abilityType) {
		return false
	}

	ability := am.Abilities[abilityType]
	ability.CooldownTimer = ability.Cooldown
	am.ActiveAbilities[abilityType] = true

	// Set duration for effects that have duration
	switch abilityType {
	case AbilityTypeSlowTime:
		am.ActiveAbilityTimers[abilityType] = 2.0 // 2 seconds
	case AbilityTypeBarrier:
		am.ActiveAbilityTimers[abilityType] = 5.0 // 5 seconds or until hit
	case AbilityTypeWeaponBoost:
		am.ActiveAbilityTimers[abilityType] = 3.0 // 3 seconds
	case AbilityTypeEMPPulse:
		am.ActiveAbilityTimers[abilityType] = 0.2 // Quick pulse
	case AbilityTypeOrbitalShield:
		am.ActiveAbilityTimers[abilityType] = 8.0 // 8 seconds
	}

	return true
}

// IsAbilityActive checks if an ability effect is currently active
func (am *AbilityManager) IsAbilityActive(abilityType AbilityType) bool {
	return am.ActiveAbilities[abilityType]
}

// Update updates all ability cooldowns and timers
func (am *AbilityManager) Update() {
	for _, ability := range am.Abilities {
		if ability.CooldownTimer > 0 {
			ability.CooldownTimer -= 1.0 / 60.0 // 60 FPS
		}
	}

	// Update active ability timers
	for abilityType, timer := range am.ActiveAbilityTimers {
		if timer > 0 {
			am.ActiveAbilityTimers[abilityType] = timer - 1.0/60.0
		} else if timer <= 0 && am.ActiveAbilities[abilityType] {
			am.ActiveAbilities[abilityType] = false
		}
	}
}

// GetAbilityByKey returns ability by key binding
func (am *AbilityManager) GetAbilityByKey(key string) *Ability {
	for _, ability := range am.Abilities {
		if ability.KeyBinding == key {
			return ability
		}
	}
	return nil
}

// GetAllAbilities returns all learned abilities
func (am *AbilityManager) GetAllAbilities() []*Ability {
	var abilities []*Ability
	for _, ability := range am.Abilities {
		abilities = append(abilities, ability)
	}
	return abilities
}

// GetActiveAbilities returns currently active ability effects
func (am *AbilityManager) GetActiveAbilities() []AbilityType {
	var active []AbilityType
	for abilityType := range am.ActiveAbilities {
		if am.ActiveAbilities[abilityType] {
			active = append(active, abilityType)
		}
	}
	return active
}

// HasAbility checks if player has learned an ability
func (am *AbilityManager) HasAbility(abilityType AbilityType) bool {
	_, exists := am.Abilities[abilityType]
	return exists
}

// GetAbilityCooldownPercent returns cooldown progress (0.0-1.0)
func (am *AbilityManager) GetAbilityCooldownPercent(abilityType AbilityType) float64 {
	if ability, exists := am.Abilities[abilityType]; exists {
		if ability.Cooldown == 0 {
			return 0
		}
		return ability.CooldownTimer / ability.Cooldown
	}
	return 0
}

// DashEffect contains parameters for dash ability
type DashEffect struct {
	Direction [2]float64
	Speed     float64
	Duration  float64
}

// SlowTimeEffect contains parameters for slow time ability
type SlowTimeEffect struct {
	TimeScale float64 // 0.5 = 50% speed
	Duration  float64
}

// BarrierEffect contains parameters for barrier ability
type BarrierEffect struct {
	Radius   float64
	Health   float64
	Active   bool
	Duration float64
}

// WeaponBoostEffect contains parameters for weapon boost
type WeaponBoostEffect struct {
	FireRateMultiplier float64
	Duration           float64
}

// EMPPulseEffect contains parameters for EMP pulse
type EMPPulseEffect struct {
	Radius   float64
	StunTime float64
}

// OrbitalShieldEffect contains parameters for orbital shield
type OrbitalShieldEffect struct {
	ProjectileCount int
	OrbitalRadius   float64
	RotationSpeed   float64
	Duration        float64
}
