package entities

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	X, Y         float64
	VelX, VelY   float64
	Radius       float64
	Speed        float64
	Health       int
	MaxHealth    int
	Shield       int
	MaxShield    int
	WeaponLevel  int             // Deprecated: kept for compatibility
	WeaponMgr    *WeaponManager  // New weapon system
	AbilityMgr   *AbilityManager // Ability system
	FireRate     float64         // Deprecated: now in WeaponManager
	FireCooldown float64         // Deprecated: now in WeaponManager
	InvincTimer  float64
	Active       bool
	EngineGlow   float64

	// Difficulty-dependent settings
	ShieldRegenRate   float64 // HP per frame
	InvincibilityTime float64 // seconds
	ShieldRegenDelay  float64 // seconds before regen starts
	LastDamageTime    float64 // when damage was last taken
	PrevShield        int     // Track previous shield value for sound effects
	ShieldRegenAccum  float64 // Accumulator for fractional shield regeneration

	// Special attack mechanics
	ChargeLevel       float64 // 0 to 1.0 (charge for special attack)
	UltimateCharge    float64 // 0 to 1.0 (builds from combat)
	MaxUltimateCharge float64 // 1.0
	UltimateActive    bool    // Ultimate ability activated
	UltimateTimer     float64 // Duration of ultimate effect

	// Thruster trail system
	ThrusterTrail []struct{ X, Y, Life float64 }

	// Mystery Power-Up temporary effects
	SpeedBoostTimer      float64 // Speed boost duration
	SpeedBoostMultiplier float64 // Speed multiplier (1.5 = +50%)
	RapidFireTimer       float64 // Rapid fire duration
	RapidFireMultiplier  float64 // Fire rate multiplier
	ScoreMultiplierTimer float64 // Score multiplier duration
	ScoreMultiplier      float64 // Score multiplier value
	ControlReversed      bool    // Controls are reversed
	ControlReversalTimer float64 // Duration of control reversal
	SlowFireTimer        float64 // Slow fire duration
	SlowFireMultiplier   float64 // Fire rate reduction
	InvincibilityTimer   float64 // Invincibility from power-up
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:                 x,
		Y:                 y,
		Radius:            20,
		Speed:             6,
		Health:            100,
		MaxHealth:         100,
		Shield:            50,
		MaxShield:         50,
		WeaponLevel:       1,
		WeaponMgr:         NewWeaponManager(),  // Initialize weapon system
		AbilityMgr:        NewAbilityManager(), // Initialize ability system
		FireRate:          0.12,
		Active:            true,
		EngineGlow:        0,
		ShieldRegenRate:   0.5,  // Default (Normal difficulty)
		InvincibilityTime: 0.25, // Default (Normal difficulty)
		ShieldRegenDelay:  3.5,  // Default (Normal difficulty)
		LastDamageTime:    -999, // Initialize to long ago so regen starts immediately
		PrevShield:        50,   // Initialize to starting shield
		ChargeLevel:       0,    // No charge initially
		UltimateCharge:    0,    // No ultimate initially
		MaxUltimateCharge: 1.0,  // Max ultimate charge
		UltimateActive:    false,
		UltimateTimer:     0,
		ThrusterTrail:     make([]struct{ X, Y, Life float64 }, 0, 10),

		// Mystery power-up effects
		SpeedBoostTimer:      0,
		SpeedBoostMultiplier: 1.0,
		RapidFireTimer:       0,
		RapidFireMultiplier:  1.0,
		ScoreMultiplierTimer: 0,
		ScoreMultiplier:      1.0,
		ControlReversed:      false,
		ControlReversalTimer: 0,
		SlowFireTimer:        0,
		SlowFireMultiplier:   1.0,
		InvincibilityTimer:   0,
	}
}

func (p *Player) Update(screenWidth, screenHeight int, gameTime float64) {
	// Update weapon manager
	p.WeaponMgr.Update()

	// Apply fire rate modifiers from mystery power-ups to weapon manager
	if p.RapidFireTimer > 0 {
		// Temporarily boost current weapon fire rate
		weapon := p.WeaponMgr.GetCurrentWeapon()
		if weapon != nil && p.RapidFireTimer > 0 && p.RapidFireMultiplier > 1.0 {
			// Apply multiplier (will be reset when timer expires)
			// Note: This is applied dynamically in Shoot() function
		}
	}

	// Update mystery power-up effect timers
	if p.SpeedBoostTimer > 0 {
		p.SpeedBoostTimer -= 1.0 / 60.0
		if p.SpeedBoostTimer <= 0 {
			p.SpeedBoostMultiplier = 1.0
		}
	}
	if p.RapidFireTimer > 0 {
		p.RapidFireTimer -= 1.0 / 60.0
		if p.RapidFireTimer <= 0 {
			p.RapidFireMultiplier = 1.0
		}
	}
	if p.ScoreMultiplierTimer > 0 {
		p.ScoreMultiplierTimer -= 1.0 / 60.0
		if p.ScoreMultiplierTimer <= 0 {
			p.ScoreMultiplier = 1.0
		}
	}
	if p.ControlReversalTimer > 0 {
		p.ControlReversalTimer -= 1.0 / 60.0
		if p.ControlReversalTimer <= 0 {
			p.ControlReversed = false
		}
	}
	if p.SlowFireTimer > 0 {
		p.SlowFireTimer -= 1.0 / 60.0
		if p.SlowFireTimer <= 0 {
			p.SlowFireMultiplier = 1.0
		}
	}
	if p.InvincibilityTimer > 0 {
		p.InvincibilityTimer -= 1.0 / 60.0
	}

	// Handle movement
	p.VelX = 0
	p.VelY = 0

	// Apply control reversal if active
	controlMult := 1.0
	if p.ControlReversed {
		controlMult = -1.0
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.VelY = -p.Speed * controlMult
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.VelY = p.Speed * controlMult
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.VelX = -p.Speed * controlMult
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.VelX = p.Speed * controlMult
	}

	// Normalize diagonal movement
	if p.VelX != 0 && p.VelY != 0 {
		p.VelX *= 0.707
		p.VelY *= 0.707
	}

	// Apply speed boost/malfunction multiplier
	p.VelX *= p.SpeedBoostMultiplier
	p.VelY *= p.SpeedBoostMultiplier

	p.X += p.VelX
	p.Y += p.VelY

	// Add thruster trail particles when moving
	if (p.VelX != 0 || p.VelY != 0) && rand.Float64() < 0.6 {
		p.ThrusterTrail = append(p.ThrusterTrail, struct{ X, Y, Life float64 }{
			X:    p.X + (rand.Float64()-0.5)*8,
			Y:    p.Y + p.Radius + (rand.Float64()-0.5)*4,
			Life: 0.5,
		})
		// Keep trail list bounded
		if len(p.ThrusterTrail) > 15 {
			p.ThrusterTrail = p.ThrusterTrail[1:]
		}
	}

	// Update thruster trail (iterate backwards to avoid skipping elements on removal)
	for i := len(p.ThrusterTrail) - 1; i >= 0; i-- {
		p.ThrusterTrail[i].Life -= 1.0 / 60.0
		p.ThrusterTrail[i].Y += 1.5 // Trail drifts down slightly
		if p.ThrusterTrail[i].Life <= 0 {
			p.ThrusterTrail = append(p.ThrusterTrail[:i], p.ThrusterTrail[i+1:]...)
		}
	}

	// Bounds check
	margin := p.Radius
	if p.X < margin {
		p.X = margin
	}
	if p.X > float64(screenWidth)-margin {
		p.X = float64(screenWidth) - margin
	}
	if p.Y < margin {
		p.Y = margin
	}
	if p.Y > float64(screenHeight)-margin {
		p.Y = float64(screenHeight) - margin
	}

	// Update cooldowns
	if p.FireCooldown > 0 {
		p.FireCooldown -= 1.0 / 60.0
	}
	if p.InvincTimer > 0 {
		p.InvincTimer -= 1.0 / 60.0
	}

	// Regenerate shield slowly - only after delay since last damage
	timeSinceDamage := gameTime - p.LastDamageTime
	if p.Shield < p.MaxShield && p.InvincTimer <= 0 && timeSinceDamage >= p.ShieldRegenDelay {
		// Use accumulator for fractional regeneration
		p.ShieldRegenAccum += p.ShieldRegenRate
		if p.ShieldRegenAccum >= 1.0 {
			regenAmount := int(p.ShieldRegenAccum)
			p.Shield += regenAmount
			p.ShieldRegenAccum -= float64(regenAmount)

			if p.Shield > p.MaxShield {
				p.Shield = p.MaxShield
				p.ShieldRegenAccum = 0 // Reset accumulator when maxed
			}
		}
	}

	// Engine glow animation
	p.EngineGlow += 0.2

	// Charge mechanics
	// Handle charge attack (hold space to charge)
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// Charging shot (slower than normal shooting)
		if p.ChargeLevel < 1.0 {
			p.ChargeLevel += 0.02 // Charge over ~3 seconds
		}
	} else {
		p.ChargeLevel = 0 // Reset when not charging
	}

	// Ultimate ability mechanics - builds from combat
	if p.UltimateActive {
		// Ultimate is active
		p.UltimateTimer -= 1.0 / 60.0
		if p.UltimateTimer <= 0 {
			p.UltimateActive = false
			p.UltimateTimer = 0
		}
	}

	// Slow ultimate charge regen (1% per second during gameplay)
	if !p.UltimateActive && p.UltimateCharge < p.MaxUltimateCharge {
		p.UltimateCharge += 0.01 / 60.0
		if p.UltimateCharge > p.MaxUltimateCharge {
			p.UltimateCharge = p.MaxUltimateCharge
		}
	}
}

func (p *Player) TakeDamage(damage int, gameTime float64) {
	// Check both regular invincibility and mystery power-up invincibility
	if p.InvincTimer > 0 || p.InvincibilityTimer > 0 {
		return
	}

	// Shield absorbs damage first
	if p.Shield > 0 {
		absorbed := min(p.Shield, damage)
		p.Shield -= absorbed
		damage -= absorbed
	}

	p.Health -= damage
	p.LastDamageTime = gameTime
	p.InvincTimer = p.InvincibilityTime // Use difficulty-dependent invincibility time
}

// Interface implementation methods

// IsActive returns whether the player is active
func (p *Player) IsActive() bool {
	return p.Active
}

// GetPosition returns the player's position
func (p *Player) GetPosition() (x, y float64) {
	return p.X, p.Y
}

// GetCollisionBounds returns the player's collision bounds
func (p *Player) GetCollisionBounds() (x, y, radius float64) {
	return p.X, p.Y, p.Radius
}

// GetHealth returns the player's current health
func (p *Player) GetHealth() int {
	return p.Health
}

// GetMaxHealth returns the player's maximum health
func (p *Player) GetMaxHealth() int {
	return p.MaxHealth
}

// IsDead returns whether the player is dead
func (p *Player) IsDead() bool {
	return p.Health <= 0
}

// GetVelocity returns the player's velocity
func (p *Player) GetVelocity() (vx, vy float64) {
	return p.VelX, p.VelY
}

// SetVelocity sets the player's velocity
func (p *Player) SetVelocity(vx, vy float64) {
	p.VelX = vx
	p.VelY = vy
}

// GetSpeed returns the player's speed
func (p *Player) GetSpeed() float64 {
	return p.Speed
}

// SetSpeed sets the player's speed
func (p *Player) SetSpeed(speed float64) {
	p.Speed = speed
}
