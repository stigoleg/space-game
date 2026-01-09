package entities

import (
	"math"
	"math/rand"
)

type EnemyType int

const (
	EnemyScout EnemyType = iota
	EnemyDrone
	EnemyHunter
	EnemyTank
	EnemyBomber
	EnemySniper       // Long-range precise shooter
	EnemySplitter     // Splits into smaller enemies when destroyed
	EnemyShieldBearer // Heavily armored with regenerating shield
)

// FormationType represents different enemy formation patterns
type FormationType int

const (
	FormationTypeNone FormationType = iota
	FormationTypeVFormation
	FormationTypeCircular
	FormationTypeWave
	FormationTypePincer
	FormationTypeConvoy
)

// NewEnemy creates a new enemy with default stats based on type
func NewEnemy(x, y float64, enemyType EnemyType) *Enemy {
	e := &Enemy{
		X:         x,
		Y:         y,
		Type:      enemyType,
		Active:    true,
		Phase:     rand.Float64() * math.Pi * 2,
		AnimTimer: 0,
	}

	switch enemyType {
	case EnemyScout:
		e.Radius = 15
		e.Speed = 4
		e.Health = 20
		e.MaxHealth = 20
		e.Points = 100
		e.ShootRate = 0 // Doesn't shoot
	case EnemyDrone:
		e.Radius = 18
		e.Speed = 2.5
		e.Health = 30
		e.MaxHealth = 30
		e.Points = 150
		e.ShootRate = 2.0
	case EnemyHunter:
		e.Radius = 20
		e.Speed = 3
		e.Health = 50
		e.MaxHealth = 50
		e.Points = 250
		e.ShootRate = 1.5
	case EnemyTank:
		e.Radius = 30
		e.Speed = 1.5
		e.Health = 100
		e.MaxHealth = 100
		e.Points = 400
		e.ShootRate = 1.0
	case EnemyBomber:
		e.Radius = 22
		e.Speed = 3.5
		e.Health = 40
		e.MaxHealth = 40
		e.Points = 300
		e.ShootRate = 0 // Explodes instead
	case EnemySniper:
		e.Radius = 16
		e.Speed = 1.0 // Very slow, stays at top
		e.Health = 35
		e.MaxHealth = 35
		e.Points = 350
		e.ShootRate = 3.0 // Slow but precise shots
		e.SniperLockTimer = 0
		e.SniperLocked = false
	case EnemySplitter:
		e.Radius = 20
		e.Speed = 2.0
		e.Health = 45
		e.MaxHealth = 45
		e.Points = 200  // Lower points since it splits
		e.ShootRate = 0 // Doesn't shoot
		e.HasSplit = false
	case EnemyShieldBearer:
		e.Radius = 25
		e.Speed = 1.2 // Slow like tank
		e.Health = 80
		e.MaxHealth = 80
		e.ShieldPoints = 50 // Starts with shield
		e.MaxShieldPoints = 50
		e.Points = 500
		e.ShootRate = 2.5
		e.ShieldRegenTimer = 0
	}

	return e
}

// NewEnemyWithDifficulty creates an enemy with difficulty adjustments
func NewEnemyWithDifficulty(x, y float64, enemyType EnemyType, healthMult, speedMult float64) *Enemy {
	e := NewEnemy(x, y, enemyType)

	// Apply difficulty multipliers
	e.Health = int(float64(e.Health) * healthMult)
	e.MaxHealth = e.Health
	e.Speed = e.Speed * speedMult

	return e
}
