package components

import (
	"math"
)

// Component is the base interface for all components
type Component interface {
	Update(deltaTime float64)
}

// PositionComponent handles entity position
type PositionComponent struct {
	X, Y float64
}

// GetPosition returns the current position
func (c *PositionComponent) GetPosition() (x, y float64) {
	return c.X, c.Y
}

// SetPosition sets the position
func (c *PositionComponent) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

// VelocityComponent handles entity velocity and movement
type VelocityComponent struct {
	VelX, VelY float64
	MaxSpeed   float64
}

// Update applies velocity to position
func (c *VelocityComponent) Update(deltaTime float64) {
	// Velocity is already applied in the entity's Update method
	// This component just stores velocity data
}

// ApplyVelocity applies velocity to a position component
func (c *VelocityComponent) ApplyVelocity(pos *PositionComponent, deltaTime float64) {
	pos.X += c.VelX * deltaTime
	pos.Y += c.VelY * deltaTime
}

// ClampSpeed ensures velocity doesn't exceed max speed
func (c *VelocityComponent) ClampSpeed() {
	if c.MaxSpeed <= 0 {
		return
	}

	speed := math.Sqrt(c.VelX*c.VelX + c.VelY*c.VelY)
	if speed > c.MaxSpeed {
		scale := c.MaxSpeed / speed
		c.VelX *= scale
		c.VelY *= scale
	}
}

// HealthComponent handles entity health
type HealthComponent struct {
	Health    int
	MaxHealth int
}

// TakeDamage reduces health by damage amount
func (c *HealthComponent) TakeDamage(damage int) int {
	if damage < 0 {
		damage = 0
	}

	oldHealth := c.Health
	c.Health -= damage
	if c.Health < 0 {
		c.Health = 0
	}

	return oldHealth - c.Health // Return actual damage dealt
}

// Heal increases health up to max
func (c *HealthComponent) Heal(amount int) {
	c.Health += amount
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

// IsDead returns whether health is zero or below
func (c *HealthComponent) IsDead() bool {
	return c.Health <= 0
}

// GetHealthRatio returns health as percentage of max (0.0 to 1.0)
func (c *HealthComponent) GetHealthRatio() float64 {
	if c.MaxHealth <= 0 {
		return 0
	}
	return float64(c.Health) / float64(c.MaxHealth)
}

// TimerComponent handles countdown timers
type TimerComponent struct {
	Time     float64
	Duration float64
	Active   bool
	Loop     bool
}

// NewTimerComponent creates a new timer
func NewTimerComponent(duration float64, loop bool) *TimerComponent {
	return &TimerComponent{
		Time:     0,
		Duration: duration,
		Active:   true,
		Loop:     loop,
	}
}

// Update updates the timer
func (c *TimerComponent) Update(deltaTime float64) {
	if !c.Active {
		return
	}

	c.Time += deltaTime
	if c.Time >= c.Duration {
		if c.Loop {
			c.Time -= c.Duration
		} else {
			c.Active = false
		}
	}
}

// IsExpired returns whether the timer has expired
func (c *TimerComponent) IsExpired() bool {
	return !c.Active && c.Time >= c.Duration
}

// GetProgress returns timer progress (0.0 to 1.0)
func (c *TimerComponent) GetProgress() float64 {
	if c.Duration <= 0 {
		return 1.0
	}
	progress := c.Time / c.Duration
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// Reset resets the timer
func (c *TimerComponent) Reset() {
	c.Time = 0
	c.Active = true
}

// LifetimeComponent handles entity lifetime
type LifetimeComponent struct {
	Age        float64
	MaxAge     float64
	AutoExpire bool // If true, entity deactivates when expired
}

// Update updates the lifetime
func (c *LifetimeComponent) Update(deltaTime float64) {
	c.Age += deltaTime
}

// IsExpired returns whether the entity has exceeded its lifetime
func (c *LifetimeComponent) IsExpired() bool {
	return c.MaxAge > 0 && c.Age >= c.MaxAge
}

// GetLifetimeRatio returns remaining lifetime as percentage (0.0 to 1.0)
func (c *LifetimeComponent) GetLifetimeRatio() float64 {
	if c.MaxAge <= 0 {
		return 1.0
	}
	ratio := 1.0 - (c.Age / c.MaxAge)
	if ratio < 0 {
		return 0
	}
	return ratio
}

// BoundsComponent handles screen bounds clamping
type BoundsComponent struct {
	MinX, MinY float64
	MaxX, MaxY float64
}

// ClampPosition clamps a position to the bounds
func (c *BoundsComponent) ClampPosition(pos *PositionComponent, radius float64) {
	if pos.X-radius < c.MinX {
		pos.X = c.MinX + radius
	}
	if pos.X+radius > c.MaxX {
		pos.X = c.MaxX - radius
	}
	if pos.Y-radius < c.MinY {
		pos.Y = c.MinY + radius
	}
	if pos.Y+radius > c.MaxY {
		pos.Y = c.MaxY - radius
	}
}

// IsOutOfBounds checks if a position is outside bounds
func (c *BoundsComponent) IsOutOfBounds(pos *PositionComponent, radius float64) bool {
	return pos.X+radius < c.MinX ||
		pos.X-radius > c.MaxX ||
		pos.Y+radius < c.MinY ||
		pos.Y-radius > c.MaxY
}
