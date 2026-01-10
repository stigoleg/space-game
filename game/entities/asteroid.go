package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type AsteroidSize int

const (
	AsteroidSmall AsteroidSize = iota
	AsteroidMedium
	AsteroidLarge
)

type Asteroid struct {
	X, Y       float64
	VelX, VelY float64
	Radius     float64
	Size       AsteroidSize
	Health     int
	MaxHealth  int
	Rotation   float64
	RotSpeed   float64
	Active     bool
}

func NewAsteroid(x, y float64, size AsteroidSize) *Asteroid {
	var radius, health int
	switch size {
	case AsteroidSmall:
		radius = 10
		health = 1
	case AsteroidMedium:
		radius = 20
		health = 2
	case AsteroidLarge:
		radius = 35
		health = 4
	}

	return &Asteroid{
		X:         x,
		Y:         y,
		VelX:      (rand.Float64() - 0.5) * 2,
		VelY:      rand.Float64()*1.5 + 0.5, // Always moving down
		Radius:    float64(radius),
		Size:      size,
		Health:    health,
		MaxHealth: health,
		Rotation:  rand.Float64() * math.Pi * 2,
		RotSpeed:  (rand.Float64() - 0.5) * 0.1,
		Active:    true,
	}
}

func (a *Asteroid) Update() {
	a.X += a.VelX
	a.Y += a.VelY
	a.Rotation += a.RotSpeed

	// Wrap around screen
	if a.X < -a.Radius {
		a.X = 1280 + a.Radius
	}
	if a.X > 1280+a.Radius {
		a.X = -a.Radius
	}

	// Destroy if off bottom
	if a.Y > 800 {
		a.Active = false
	}
}

func (a *Asteroid) TakeDamage(damage int) {
	a.Health -= damage
	if a.Health <= 0 {
		a.Active = false
	}
}

// Reset resets the asteroid to default state (for object pooling)
func (a *Asteroid) Reset() {
	a.X = 0
	a.Y = 0
	a.VelX = 0
	a.VelY = 1.0
	a.Radius = 20
	a.Size = AsteroidMedium
	a.Health = 2
	a.MaxHealth = 2
	a.Rotation = 0
	a.RotSpeed = 0.05
	a.Active = false
}

// IsActive returns whether the asteroid is active
func (a *Asteroid) IsActive() bool {
	return a.Active
}

// SetActive sets the active state
func (a *Asteroid) SetActive(active bool) {
	a.Active = active
}

func (a *Asteroid) Draw(screen *ebiten.Image, shakeX, shakeY float64, perspectiveScale float64, sprite *ebiten.Image) {
	if !a.Active {
		return
	}

	x := float32(a.X + shakeX)
	y := float32(a.Y + shakeY)
	radius := float32(a.Radius * perspectiveScale)

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		a.drawSpriteBased(screen, x, y, radius, sprite)
	} else {
		// Fallback to procedural rendering
		a.drawProcedural(screen, x, y, radius)
	}
}

func (a *Asteroid) drawSpriteBased(screen *ebiten.Image, x, y, radius float32, sprite *ebiten.Image) {
	// Draw shadow
	shadowColor := color.RGBA{20, 20, 30, 80}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.6, shadowColor, true)

	// Draw sprite with rotation
	op := &ebiten.DrawImageOptions{}

	// Scale sprite to match asteroid size
	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	targetSize := float64(radius) * 2.0
	scaleX := targetSize / spriteWidth
	scaleY := targetSize / spriteHeight

	// Apply rotation
	op.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
	op.GeoM.Rotate(a.Rotation)
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(sprite, op)

	// Health indicator (glow) when damaged
	if a.Health < a.MaxHealth {
		healthRatio := float32(a.Health) / float32(a.MaxHealth)
		glowColor := color.RGBA{255, uint8(100 * healthRatio), 50, uint8(100 * (1 - healthRatio))}
		glowSize := radius + float32(3*(1-healthRatio))
		vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)
	}
}

func (a *Asteroid) drawProcedural(screen *ebiten.Image, x, y, radius float32) {
	// Draw shadow
	shadowColor := color.RGBA{20, 20, 30, 80}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.6, shadowColor, true)

	// Draw asteroid with rocky appearance
	// Base color varies by size
	var baseColor color.RGBA
	switch a.Size {
	case AsteroidSmall:
		baseColor = color.RGBA{150, 120, 100, 255}
	case AsteroidMedium:
		baseColor = color.RGBA{130, 100, 80, 255}
	case AsteroidLarge:
		baseColor = color.RGBA{110, 80, 60, 255}
	}

	// Main body with crater effect
	vector.DrawFilledCircle(screen, x, y, radius, baseColor, true)

	// Add rocky texture with offset circles
	craterCount := int(a.Radius) / 5
	for i := 0; i < craterCount; i++ {
		angle := float64(i) * (math.Pi * 2 / float64(craterCount))
		angle += a.Rotation
		craterX := x + float32(math.Cos(angle)*float64(radius)*0.6)
		craterY := y + float32(math.Sin(angle)*float64(radius)*0.6)
		craterSize := radius * 0.25

		craterColor := color.RGBA{baseColor.R / 2, baseColor.G / 2, baseColor.B / 2, 200}
		vector.DrawFilledCircle(screen, craterX, craterY, craterSize, craterColor, true)
	}

	// Highlight edge
	highlightColor := color.RGBA{180, 150, 130, 150}
	highlightRadius := radius * 0.15
	vector.DrawFilledCircle(screen, x-radius*0.3, y-radius*0.3, highlightRadius, highlightColor, true)

	// Health indicator (glow) when damaged
	if a.Health < a.MaxHealth {
		healthRatio := float32(a.Health) / float32(a.MaxHealth)
		glowColor := color.RGBA{255, uint8(100 * healthRatio), 50, uint8(100 * (1 - healthRatio))}
		glowSize := radius + float32(3*(1-healthRatio))
		vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)
	}
}
