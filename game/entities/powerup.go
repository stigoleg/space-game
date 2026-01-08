package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PowerUpType int

const (
	PowerUpHealth PowerUpType = iota
	PowerUpShield
	PowerUpWeapon
	PowerUpSpeed
)

type PowerUp struct {
	X, Y      float64
	VelY      float64
	Radius    float64
	Type      PowerUpType
	Active    bool
	AnimTimer float64
}

func NewPowerUp(x, y float64) *PowerUp {
	return &PowerUp{
		X:         x,
		Y:         y,
		VelY:      1.5,
		Radius:    15,
		Type:      PowerUpType(rand.Intn(4)),
		Active:    true,
		AnimTimer: rand.Float64() * math.Pi * 2,
	}
}

func (p *PowerUp) Update() {
	p.Y += p.VelY
	p.AnimTimer += 0.15 // Increased animation speed
}

func (p *PowerUp) Draw(screen *ebiten.Image, shakeX, shakeY float64, sprite *ebiten.Image, sparkleSprites []*ebiten.Image) {
	// Simple screen coordinates with shake
	x := float32(p.X + shakeX)
	y := float32(p.Y + shakeY)

	// Floating effect (increased amplitude for better visibility)
	floatOffset := float32(math.Sin(p.AnimTimer) * 5) // Increased from 3 to 5
	y += floatOffset

	// Larger pulsing effect (40% instead of 20%)
	pulse := float32(1.0 + 0.4*math.Sin(p.AnimTimer*2))

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		p.drawSpriteBased(screen, x, y, pulse, sprite, sparkleSprites)
	} else {
		// Fallback to procedural rendering
		p.drawProcedural(screen, x, y, pulse)
	}
}

func (p *PowerUp) drawSpriteBased(screen *ebiten.Image, x, y, pulse float32, sprite *ebiten.Image, sparkleSprites []*ebiten.Image) {
	// Increased radius by 30%
	baseRadius := float32(p.Radius) * 1.3
	radius := baseRadius * pulse

	// Draw shadow
	shadowColor := color.RGBA{20, 20, 30, 80}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Draw vertical beam of light for easy spotting
	var beamColor color.RGBA
	switch p.Type {
	case PowerUpHealth:
		beamColor = color.RGBA{50, 255, 50, 60}
	case PowerUpShield:
		beamColor = color.RGBA{50, 150, 255, 60}
	case PowerUpWeapon:
		beamColor = color.RGBA{255, 200, 50, 60}
	case PowerUpSpeed:
		beamColor = color.RGBA{255, 100, 255, 60}
	}
	beamWidth := radius * 0.3
	vector.DrawFilledRect(screen, x-beamWidth/2, y-200, beamWidth, 200, beamColor, true)

	// Draw sprite with rotation and scaling
	op := &ebiten.DrawImageOptions{}

	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	// Scale to be 30% larger
	targetSize := float64(radius) * 2.0
	scaleX := targetSize / spriteWidth
	scaleY := targetSize / spriteHeight

	// Slow rotation for visual interest
	op.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)
	op.GeoM.Rotate(p.AnimTimer * 0.5)
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(sprite, op)

	// Draw sparkle particles orbiting the power-up
	if sparkleSprites != nil && len(sparkleSprites) > 0 {
		for i := 0; i < 4; i++ {
			angle := p.AnimTimer + float64(i)*math.Pi/2
			sparkleX := x + float32(math.Cos(angle))*radius*1.5
			sparkleY := y + float32(math.Sin(angle))*radius*1.5

			sparkleOp := &ebiten.DrawImageOptions{}
			sparkleOp.GeoM.Translate(-8, -8) // Center sparkle (16x16 sprite)
			sparkleOp.GeoM.Translate(float64(sparkleX), float64(sparkleY))

			frameIndex := int(p.AnimTimer*10) % len(sparkleSprites)
			screen.DrawImage(sparkleSprites[frameIndex], sparkleOp)
		}
	}

	// Draw animated rotating border (dashed circle)
	numDots := 16
	for i := 0; i < numDots; i++ {
		angle := p.AnimTimer*2 + float64(i)*math.Pi*2/float64(numDots)
		dotX := x + float32(math.Cos(angle))*(radius+10)
		dotY := y + float32(math.Sin(angle))*(radius+10)
		dotColor := beamColor
		dotColor.A = 255
		vector.DrawFilledCircle(screen, dotX, dotY, 3, dotColor, true)
	}
}

func (p *PowerUp) drawProcedural(screen *ebiten.Image, x, y, pulse float32) {
	var mainColor, glowColor color.RGBA

	switch p.Type {
	case PowerUpHealth:
		mainColor = color.RGBA{50, 255, 50, 255}
		glowColor = color.RGBA{50, 255, 50, 100}
	case PowerUpShield:
		mainColor = color.RGBA{50, 150, 255, 255}
		glowColor = color.RGBA{50, 150, 255, 100}
	case PowerUpWeapon:
		mainColor = color.RGBA{255, 200, 50, 255}
		glowColor = color.RGBA{255, 200, 50, 100}
	case PowerUpSpeed:
		mainColor = color.RGBA{255, 100, 255, 255}
		glowColor = color.RGBA{255, 100, 255, 100}
	}

	radius := float32(p.Radius) * pulse

	// Draw shadow
	shadowColor := color.RGBA{20, 20, 30, 80}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Main circle
	vector.DrawFilledCircle(screen, x, y, radius, mainColor, true)

	// Outer glow ring
	vector.StrokeCircle(screen, x, y, radius+6, 2, glowColor, true)

	// Inner highlight
	vector.DrawFilledCircle(screen, x-3, y-3, radius*0.35, color.RGBA{255, 255, 255, 220}, true)

	// Type indicator (simple shapes)
	switch p.Type {
	case PowerUpHealth:
		// Plus sign
		vector.DrawFilledRect(screen, x-6, y-2, 12, 4, color.RGBA{255, 255, 255, 255}, true)
		vector.DrawFilledRect(screen, x-2, y-6, 4, 12, color.RGBA{255, 255, 255, 255}, true)
	case PowerUpShield:
		// Shield shape (circle with stroke)
		vector.StrokeCircle(screen, x, y, 6, 3, color.RGBA{255, 255, 255, 255}, true)
	case PowerUpWeapon:
		// Arrow up
		vector.DrawFilledRect(screen, x-2, y-4, 4, 10, color.RGBA{255, 255, 255, 255}, true)
	case PowerUpSpeed:
		// Lightning bolt (simplified as lines)
		vector.StrokeLine(screen, x-3, y-6, x+3, y, 2, color.RGBA{255, 255, 255, 255}, true)
		vector.StrokeLine(screen, x+3, y, x-3, y+6, 2, color.RGBA{255, 255, 255, 255}, true)
	}
}
