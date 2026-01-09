package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Projectile struct {
	X, Y       float64
	VelX, VelY float64
	Radius     float64
	Damage     int
	Friendly   bool // true = player's projectile
	Active     bool
	Trail      []struct{ X, Y float64 }
	Color      color.RGBA // Custom color for weapon types
	GlowColor  color.RGBA // Custom glow color
	Lifetime   float64    // Time to live in seconds (0 = infinite)
	Age        float64    // Current age in seconds

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

	return &Projectile{
		X:              x,
		Y:              y,
		VelX:           velX,
		VelY:           velY,
		Radius:         7, // Increased from 5 to 7
		Damage:         damage,
		Friendly:       friendly,
		Active:         true,
		Trail:          make([]struct{ X, Y float64 }, 0, 8), // Increased from 5 to 8
		Color:          mainColor,
		GlowColor:      glowColor,
		TargetEnemyIdx: -1,  // No target initially
		Lifetime:       3.0, // Reduced from 5.0 to 3.0 for faster cleanup
		Age:            0.0, // Start at age 0
	}
}

// NewProjectileWithColor creates a projectile with custom colors (for weapon variety)
func NewProjectileWithColor(x, y, velX, velY float64, friendly bool, damage int, mainColor, glowColor color.RGBA) *Projectile {
	return &Projectile{
		X:              x,
		Y:              y,
		VelX:           velX,
		VelY:           velY,
		Radius:         7,
		Damage:         damage,
		Friendly:       friendly,
		Active:         true,
		Trail:          make([]struct{ X, Y float64 }, 0, 8),
		Color:          mainColor,
		GlowColor:      glowColor,
		TargetEnemyIdx: -1,  // No target initially
		Lifetime:       3.0, // Reduced from 5.0 to 3.0 for faster cleanup
		Age:            0.0, // Start at age 0
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

	// Store trail position
	p.Trail = append(p.Trail, struct{ X, Y float64 }{p.X, p.Y})
	if len(p.Trail) > 8 { // Increased trail length
		p.Trail = p.Trail[1:]
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

	// Draw multiple layers for glow effect
	// Outer glow
	vector.StrokeLine(screen, startX, startY, endX, endY, 12,
		color.RGBA{p.GlowColor.R, p.GlowColor.G, p.GlowColor.B, 60}, true)

	// Middle glow
	vector.StrokeLine(screen, startX, startY, endX, endY, 6,
		color.RGBA{p.GlowColor.R, p.GlowColor.G, p.GlowColor.B, 150}, true)

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

	// Draw large outer glow FIRST (background)
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+8, color.RGBA{glowColor.R, glowColor.G, glowColor.B, 40}, true)

	// Draw enhanced trail with glow
	for i, t := range p.Trail {
		alpha := uint8(100 + i*30)
		size := float32(p.Radius) * float32(i+1) / float32(len(p.Trail)+2)

		// Trail glow (larger)
		glowSize := size * 2.2 // Increased from 1.8
		glowAlpha := uint8(50 + i*20)
		trailGlow := color.RGBA{glowColor.R, glowColor.G, glowColor.B, glowAlpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), glowSize, trailGlow, true)

		// Trail particle
		trailColor := color.RGBA{mainColor.R, mainColor.G, mainColor.B, alpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), size, trailColor, true)
	}

	// Draw sprite on top
	op := &ebiten.DrawImageOptions{}

	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	// Center sprite on projectile position
	op.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(sprite, op)
}

func (p *Projectile) drawProcedural(screen *ebiten.Image, x, y float32, shakeX, shakeY float64) {
	// Use custom colors from the projectile
	mainColor := p.Color
	glowColor := p.GlowColor

	// Draw enhanced trail with glow
	for i, t := range p.Trail {
		alpha := uint8(100 + i*30)
		size := float32(p.Radius) * float32(i+1) / float32(len(p.Trail)+2)

		// Trail glow
		glowSize := size * 1.8
		glowAlpha := uint8(50 + i*15)
		trailGlow := color.RGBA{glowColor.R, glowColor.G, glowColor.B, glowAlpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), glowSize, trailGlow, true)

		// Trail particle
		trailColor := color.RGBA{mainColor.R, mainColor.G, mainColor.B, alpha}
		vector.DrawFilledCircle(screen, float32(t.X+shakeX), float32(t.Y+shakeY), size, trailColor, true)
	}

	// Large outer glow
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+6, color.RGBA{glowColor.R, glowColor.G, glowColor.B, 80}, true)

	// Main glow
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)+3, glowColor, true)

	// Main projectile
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius), mainColor, true)

	// Bright center core
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)*0.6, color.RGBA{255, 255, 255, 220}, true)

	// Inner bright spot
	vector.DrawFilledCircle(screen, x, y, float32(p.Radius)*0.3, color.RGBA{255, 255, 255, 255}, true)
}

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

	// If no valid target, find nearest enemy
	if target == nil {
		minDist := 999999.0
		for i, e := range enemies {
			if !e.Active {
				continue
			}

			dist := (e.X-p.X)*(e.X-p.X) + (e.Y-p.Y)*(e.Y-p.Y)
			if dist < minDist {
				minDist = dist
				target = e
				p.TargetEnemyIdx = i
			}
		}
	}

	if target == nil {
		return // No target, fly straight
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
	p.Trail = p.Trail[:0] // Keep capacity, clear length
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
