package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Entity represents any game entity that can be updated and drawn
type Entity interface {
	Update() error
	Draw(screen *ebiten.Image)
	IsActive() bool
	GetPosition() (x, y float64)
}

// Collidable represents entities that can participate in collision detection
type Collidable interface {
	Entity
	GetCollisionBounds() (x, y, radius float64)
	GetCollisionLayer() CollisionLayer
}

// CollisionLayer represents different layers for collision detection
type CollisionLayer int

const (
	LayerPlayer CollisionLayer = iota
	LayerEnemy
	LayerBoss
	LayerPlayerProjectile
	LayerEnemyProjectile
	LayerPowerUp
	LayerAsteroid
	LayerHazard
)

// Damageable represents entities that can take damage
type Damageable interface {
	Collidable
	TakeDamage(amount int) int // Returns actual damage dealt
	GetHealth() int
	GetMaxHealth() int
	IsDead() bool
}

// Shooter represents entities that can shoot projectiles
type Shooter interface {
	Entity
	CanShoot() bool
	GetShootPosition() (x, y float64)
	GetShootVelocity() (vx, vy float64)
}

// Movable represents entities with physics-based movement
type Movable interface {
	Entity
	GetVelocity() (vx, vy float64)
	SetVelocity(vx, vy float64)
	GetSpeed() float64
	SetSpeed(speed float64)
}

// BoundedEntity represents entities that should stay within screen bounds
type BoundedEntity interface {
	Movable
	ClampToScreen(screenWidth, screenHeight int)
}

// Poolable represents entities that can be pooled for reuse
type Poolable interface {
	Entity
	Reset()
	Activate()
	Deactivate()
}

// Animated represents entities with animation
type Animated interface {
	Entity
	UpdateAnimation(deltaTime float64)
	GetCurrentFrame() int
}

// Temporary represents entities with a limited lifetime
type Temporary interface {
	Entity
	GetLifetime() float64
	GetMaxLifetime() float64
	IsExpired() bool
}
