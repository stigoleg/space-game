package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Reusable DrawImageOptions to avoid allocation per projectile draw
var projectileDrawOptions = &ebiten.DrawImageOptions{}

// TrailPoint represents a point in the projectile trail
type TrailPoint struct {
	X, Y float64
}

// Maximum trail length constant
const MaxTrailLength = 5

type Projectile struct {
	X, Y       float64
	VelX, VelY float64
	Radius     float64
	Damage     int
	Friendly   bool // true = player's projectile
	Active     bool
	Trail      [MaxTrailLength]TrailPoint // Fixed-size array for ring buffer
	TrailHead  int                        // Current write position in ring buffer
	TrailLen   int                        // Current number of valid trail points
	Color      color.RGBA                 // Custom color for weapon types
	GlowColor  color.RGBA                 // Custom glow color
	Lifetime   float64                    // Time to live in seconds (0 = infinite)
	Age        float64                    // Current age in seconds

	// Special behavior flags
	Homing         bool    // Following rocket behavior
	HomingSpeed    float64 // Turn rate for homing (radians per frame)
	TargetEnemyIdx int     // Index of target enemy (-1 = no target)

	Chaining   bool    // Chain lightning behavior
	ChainCount int     // How many times can chain
	ChainRange float64 // Range to find next target

	Burning      bool    // Causes DoT on hit
	BurnDuration float64 // How long burn lasts (seconds)
	BurnDamage   int     // Damage per tick

	Beam       bool                   // Ion beam behavior
	BeamSource struct{ X, Y float64 } // Start point of beam
	Piercing   bool                   // Penetrates enemies
}

func NewProjectile(x, y, velX, velY float64, friendly bool, damage int) *Projectile {
	// Default colors (blue for friendly, red for enemy)
	var mainColor, glowColor color.RGBA
	if friendly {
		mainColor = color.RGBA{100, 200, 255, 255}
		glowColor = color.RGBA{50, 150, 255, 180}
	} else {
		mainColor = color.RGBA{255, 100, 100, 255}
		glowColor = color.RGBA{255, 50, 50, 180}
	}

	// Calculate speed and set lifetime based on velocity
	speed := math.Sqrt(velX*velX+velY*velY) * 60.0 // Convert to pixels per second
	lifetime := 2.5                                // Default lifetime

	// Fast projectiles get shorter lifetime for better cleanup
	if speed > 350 {
		lifetime = 2.0
	}

	return &Projectile{
		X:              x,
		Y:              y,
		VelX:           velX,
		VelY:           velY,
		Radius:         7, // Increased from 5 to 7
		Damage:         damage,
		Friendly:       friendly,
		Active:         true,
		TrailHead:      0, // Ring buffer starts at 0
		TrailLen:       0, // No trail points yet
		Color:          mainColor,
		GlowColor:      glowColor,
		TargetEnemyIdx: -1,
		Lifetime:       lifetime, // Speed-based lifetime (2.0-2.5s)
		Age:            0.0,      // Start at age 0
	}
}

// NewProjectileWithColor creates a projectile with custom colors (for weapon variety)
func NewProjectileWithColor(x, y, velX, velY float64, friendly bool, damage int, mainColor, glowColor color.RGBA) *Projectile {
	// Calculate speed and set lifetime based on velocity
	speed := math.Sqrt(velX*velX+velY*velY) * 60.0 // Convert to pixels per second
	lifetime := 2.5                                // Default lifetime

	// Fast projectiles get shorter lifetime for better cleanup
	if speed > 350 {
		lifetime = 2.0
	}

	return &Projectile{
		X:              x,
		Y:              y,
		VelX:           velX,
		VelY:           velY,
		Radius:         7,
		Damage:         damage,
		Friendly:       friendly,
		Active:         true,
		TrailHead:      0, // Ring buffer starts at 0
		TrailLen:       0, // No trail points yet
		Color:          mainColor,
		GlowColor:      glowColor,
		TargetEnemyIdx: -1,
		Lifetime:       lifetime, // Speed-based lifetime (2.0-2.5s)
		Age:            0.0,      // Start at age 0
	}
}

func (p *Projectile) Update() {
	// Update age
	p.Age += 1.0 / 60.0 // Assuming 60 FPS

	// Check if projectile has expired
	if p.Lifetime > 0 && p.Age >= p.Lifetime {
		p.Active = false
		return
	}

	// Store trail position using ring buffer (no allocations)
	p.Trail[p.TrailHead] = TrailPoint{X: p.X, Y: p.Y}
	p.TrailHead = (p.TrailHead + 1) % MaxTrailLength
	if p.TrailLen < MaxTrailLength {
		p.TrailLen++
	}

	p.X += p.VelX
	p.Y += p.VelY
}

// IsOffScreen checks if projectile is clearly beyond screen bounds for early culling
// screenWidth and screenHeight should be passed from game
func (p *Projectile) IsOffScreen(screenWidth, screenHeight int) bool {
	const buffer = 100.0 // Pixels beyond screen edge before culling
	return p.X < -buffer || p.X > float64(screenWidth)+buffer ||
		p.Y < -buffer || p.Y > float64(screenHeight)+buffer
}

func (p *Projectile) Draw(screen *ebiten.Image, shakeX, shakeY float64, sprite *ebiten.Image) {
	// Simple screen coordinates with shake
	x := float32(p.X + shakeX)
	y := float32(p.Y + shakeY)

	// Draw beam if this is a beam projectile
	if p.Beam {
		p.DrawBeam(screen, shakeX, shakeY)
		return
	}

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		p.drawSpriteBased(screen, x, y, sprite, shakeX, shakeY)
	} else {
		// Fallback to procedural rendering
		p.drawProcedural(screen, x, y, shakeX, shakeY)
	}
}

// DrawBeam draws a continuous beam from source to current position
func (p *Projectile) DrawBeam(screen *ebiten.Image, shakeX, shakeY float64) {
	if !p.Beam {
		return
	}

	// Draw continuous beam from source to current position
	startX := float32(p.BeamSource.X + shakeX)
	startY := float32(p.BeamSource.Y + shakeY)
	endX := float32(p.X + shakeX)
	endY := float32(p.Y + shakeY)

	// Optimized: Reduced from 4 layers to 3 for better performance
	// Outer glow
	vector.StrokeLine(screen, startX, startY, endX, endY, 8,
		color.RGBA{p.GlowColor.R, p.GlowColor.G, p.GlowColor.B, 50}, true)

	// Core beam
	vector.StrokeLine(screen, startX, startY, endX, endY, 3,
		p.Color, true)

	// Inner bright line
	vector.StrokeLine(screen, startX, startY, endX, endY, 1,
		color.RGBA{255, 255, 255, 255}, true)
}

func (p *Projectile) drawSpriteBased(screen *ebiten.Image, x, y float32, sprite *ebiten.Image, shakeX, shakeY float64) {
	// Use custom colors from the projectile
	mainColor := p.Color
	glowColor := p.GlowColor

	// Optimized: Reduced outer glow size and opacity for better performance
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+5, color.RGBA{glowColor.R, glowColor.G, glowColor.B, 30}, true)

	// Draw enhanced trail with glow using ring buffer
	// Optimized: Only draw every other trail point to reduce draw calls
	for i := 0; i < p.TrailLen; i += 2 {
		// Calculate index in ring buffer (oldest first)
		idx := (p.TrailHead - p.TrailLen + i + MaxTrailLength) % MaxTrailLength
		t := p.Trail[idx]

		alpha := uint8(100 + i*30)
		size := float32(p.Radius) * float32(i+1) / float32(p.TrailLen+2)

		// Trail glow (optimized size)
		glowSize := size * 1.6        // Reduced from 2.2
		glowAlpha := uint8(40 + i*15) // Reduced opacity
		trailGlow := color.RGBA{glowColor.R, glowColor.G, glowColor.B, glowAlpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), glowSize, trailGlow, true)

		// Trail particle
		trailColor := color.RGBA{mainColor.R, mainColor.G, mainColor.B, alpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), size, trailColor, true)
	}

	// Draw sprite on top using reusable DrawImageOptions
	projectileDrawOptions.GeoM.Reset()

	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	// Center sprite on projectile position
	projectileDrawOptions.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
	projectileDrawOptions.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(sprite, projectileDrawOptions)
}

func (p *Projectile) drawProcedural(screen *ebiten.Image, x, y float32, shakeX, shakeY float64) {
	// Use custom colors from the projectile
	mainColor := p.Color
	glowColor := p.GlowColor

	// Draw enhanced trail with glow using ring buffer
	// Optimized: Only draw every other trail point to reduce draw calls
	for i := 0; i < p.TrailLen; i += 2 {
		// Calculate index in ring buffer (oldest first)
		idx := (p.TrailHead - p.TrailLen + i + MaxTrailLength) % MaxTrailLength
		t := p.Trail[idx]

		alpha := uint8(100 + i*30)
		size := float32(p.Radius) * float32(i+1) / float32(p.TrailLen+2)

		// Trail glow
		glowSize := size * 1.8
		glowAlpha := uint8(50 + i*15)
		trailGlow := color.RGBA{glowColor.R, glowColor.G, glowColor.B, glowAlpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), glowSize, trailGlow, true)

		// Trail particle
		trailColor := color.RGBA{mainColor.R, mainColor.G, mainColor.B, alpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), size, trailColor, true)
	}

	// Optimized: Reduced from 5 layers to 3 for better performance

	// Outer glow (reduced size and opacity)
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+4, color.RGBA{glowColor.R, glowColor.G, glowColor.B, 60}, true)

	// Main projectile
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius), mainColor, true)

	// Bright center core
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)*0.5, color.RGBA{255, 255, 255, 200}, true)
}

// Maximum homing range constant - don't search for targets beyond this squared distance
const MaxHomingRangeSq = 500.0 * 500.0

// UpdateHoming updates projectile trajectory to home in on enemies
// enemies parameter should be passed from game loop
func (p *Projectile) UpdateHoming(enemies []*Enemy) {
	if !p.Homing {
		return
	}

	// Find target or validate current target
	var target *Enemy
	if p.TargetEnemyIdx >= 0 && p.TargetEnemyIdx < len(enemies) {
		potentialTarget := enemies[p.TargetEnemyIdx]
		if potentialTarget.Active {
			target = potentialTarget
		}
	}

	// If no valid target, find nearest enemy within range
	if target == nil {
		minDistSq := MaxHomingRangeSq
		for i, e := range enemies {
			if !e.Active {
				continue
			}

			// Use squared distance to avoid sqrt
			dx := e.X - p.X
			dy := e.Y - p.Y
			distSq := dx*dx + dy*dy

			if distSq < minDistSq {
				minDistSq = distSq
				target = e
				p.TargetEnemyIdx = i
			}
		}
	}

	if target == nil {
		return // No target within range, fly straight
	}

	// Calculate angle to target
	dx := target.X - p.X
	dy := target.Y - p.Y
	targetAngle := math.Atan2(dy, dx)

	// Current angle
	currentAngle := math.Atan2(p.VelY, p.VelX)

	// Calculate angle difference
	angleDiff := targetAngle - currentAngle

	// Normalize to [-π, π]
	for angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	for angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}

	// Apply turn rate limit
	turnAmount := math.Min(math.Abs(angleDiff), p.HomingSpeed)
	if angleDiff < 0 {
		turnAmount = -turnAmount
	}

	// Calculate new velocity
	newAngle := currentAngle + turnAmount
	speed := math.Sqrt(p.VelX*p.VelX + p.VelY*p.VelY)
	p.VelX = math.Cos(newAngle) * speed
	p.VelY = math.Sin(newAngle) * speed
}

// Poolable interface implementation

// Reset resets the projectile to default state for reuse
func (p *Projectile) Reset() {
	p.X = 0
	p.Y = 0
	p.VelX = 0
	p.VelY = 0
	p.Radius = 7
	p.Damage = 0
	p.Friendly = false
	p.Active = false
	p.TrailHead = 0 // Reset ring buffer position
	p.TrailLen = 0  // Clear trail
	p.Color = color.RGBA{100, 200, 255, 255}
	p.GlowColor = color.RGBA{50, 150, 255, 180}
	p.Lifetime = 3.0
	p.Age = 0.0

	// Reset special behavior flags
	p.Homing = false
	p.HomingSpeed = 0
	p.TargetEnemyIdx = -1
	p.Chaining = false
	p.ChainCount = 0
	p.ChainRange = 0
	p.Burning = false
	p.BurnDuration = 0
	p.BurnDamage = 0
	p.Beam = false
	p.BeamSource = struct{ X, Y float64 }{0, 0}
	p.Piercing = false
}

// SetActive sets the active state of the projectile
func (p *Projectile) SetActive(active bool) {
	p.Active = active
}

// IsActive returns whether the projectile is active
func (p *Projectile) IsActive() bool {
	return p.Active
}
