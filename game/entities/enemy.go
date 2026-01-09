package entities

import (
	"math"
)

// Enemy represents an enemy entity in the game
type Enemy struct {
	X, Y       float64
	VelX, VelY float64
	Radius     float64
	Speed      float64
	Health     int
	MaxHealth  int
	Points     int
	Type       EnemyType
	Active     bool
	ShootTimer float64
	ShootRate  float64
	AnimTimer  float64
	Phase      float64 // For wave movement

	// Burning DoT system
	Burning       bool
	BurnDuration  float64
	BurnDamage    int
	BurnTickTimer float64

	// Formation system
	FormationType     FormationType
	FormationID       int     // Groups enemies in same formation
	IsFormationLeader bool    // Is this the formation leader?
	FormationTargetX  float64 // Target position for formation
	FormationTargetY  float64
	FormationIndex    int      // Position in formation (0 = leader)
	NearbyAllies      []*Enemy // References to nearby allies in formation
	LastShootTime     float64
	CoorditatedShoot  bool // Should coordinate fire with formation

	// New enemy-specific abilities
	ShieldPoints     int     // For ShieldBearer - regenerating shield
	MaxShieldPoints  int     // Maximum shield capacity
	ShieldRegenTimer float64 // Timer for shield regeneration
	HasSplit         bool    // For Splitter - tracks if already split
	SniperLockTimer  float64 // For Sniper - time to lock onto target
	SniperLocked     bool    // For Sniper - is currently locked on player
	SniperTargetX    float64 // For Sniper - locked target position
	SniperTargetY    float64
}

// TryShoot attempts to shoot a projectile if the enemy's shoot timer is ready
func (e *Enemy) TryShoot() *Projectile {
	if e.ShootRate <= 0 || e.ShootTimer < e.ShootRate {
		return nil
	}
	e.ShootTimer = 0

	switch e.Type {
	case EnemyDrone, EnemyHunter:
		return NewProjectile(e.X, e.Y+e.Radius, 0, 6, false, 10)
	case EnemyTank:
		return NewProjectile(e.X, e.Y+e.Radius, 0, 5, false, 20)
	case EnemySniper:
		// Shoots precise fast projectiles at locked position
		if e.SniperLocked {
			dx := e.SniperTargetX - e.X
			dy := e.SniperTargetY - e.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 0 {
				speed := 9.0 // Fast projectile
				velX := (dx / dist) * speed
				velY := (dy / dist) * speed
				// Reset lock after shooting
				e.SniperLocked = false
				e.SniperLockTimer = 0
				return NewProjectile(e.X, e.Y+e.Radius, velX, velY, false, 15)
			}
		}
	case EnemyShieldBearer:
		// Shoots straight down, medium speed
		return NewProjectile(e.X, e.Y+e.Radius, 0, 5, false, 15)
	}
	return nil
}

// UpdateBurning updates the burning DoT effect
func (e *Enemy) UpdateBurning() {
	if !e.Burning {
		return
	}

	e.BurnDuration -= 1.0 / 60.0
	e.BurnTickTimer -= 1.0 / 60.0

	if e.BurnDuration <= 0 {
		e.Burning = false
		return
	}

	// Apply burn damage every 0.5 seconds
	if e.BurnTickTimer <= 0 {
		e.Health -= e.BurnDamage
		e.BurnTickTimer = 0.5
	}
}

// ApplyBurn applies a burning DoT effect to the enemy
func (e *Enemy) ApplyBurn(duration float64, damagePerTick int) {
	e.Burning = true
	e.BurnDuration = duration
	e.BurnDamage = damagePerTick
	if e.BurnTickTimer <= 0 {
		e.BurnTickTimer = 0.5 // First tick immediately
	}
}

// GetSplitEnemies returns 2 smaller Scout enemies when a Splitter is destroyed
func (e *Enemy) GetSplitEnemies() []*Enemy {
	if e.Type != EnemySplitter || e.HasSplit {
		return nil
	}

	e.HasSplit = true
	splits := make([]*Enemy, 2)

	// Create 2 scout enemies to the left and right
	for i := 0; i < 2; i++ {
		offsetX := -30.0
		if i == 1 {
			offsetX = 30.0
		}

		scout := NewEnemy(e.X+offsetX, e.Y, EnemyScout)
		// Make them slightly weaker
		scout.Health = scout.Health * 2 / 3
		scout.MaxHealth = scout.MaxHealth * 2 / 3
		scout.Points = scout.Points / 2 // Less points since they're from a split
		splits[i] = scout
	}

	return splits
}

// TakeDamage applies damage to enemy, handling shields for ShieldBearer
func (e *Enemy) TakeDamage(damage int) {
	if e.Type == EnemyShieldBearer && e.ShieldPoints > 0 {
		// Damage shield first
		e.ShieldPoints -= damage
		if e.ShieldPoints < 0 {
			// Overflow damage goes to health
			e.Health += e.ShieldPoints // ShieldPoints is negative here
			e.ShieldPoints = 0
		}
		// Reset shield regen timer when hit
		e.ShieldRegenTimer = 0
	} else {
		// Direct health damage
		e.Health -= damage
	}

	if e.Health < 0 {
		e.Health = 0
	}
}
