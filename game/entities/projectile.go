package entities

import (
	"image/color"

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
}

func NewProjectile(x, y, velX, velY float64, friendly bool, damage int) *Projectile {
	return &Projectile{
		X:        x,
		Y:        y,
		VelX:     velX,
		VelY:     velY,
		Radius:   7, // Increased from 5 to 7
		Damage:   damage,
		Friendly: friendly,
		Active:   true,
		Trail:    make([]struct{ X, Y float64 }, 0, 8), // Increased from 5 to 8
	}
}

func (p *Projectile) Update() {
	// Store trail position
	p.Trail = append(p.Trail, struct{ X, Y float64 }{p.X, p.Y})
	if len(p.Trail) > 8 { // Increased trail length
		p.Trail = p.Trail[1:]
	}

	p.X += p.VelX
	p.Y += p.VelY
}

func (p *Projectile) Draw(screen *ebiten.Image, shakeX, shakeY float64, sprite *ebiten.Image) {
	// Simple screen coordinates with shake
	x := float32(p.X + shakeX)
	y := float32(p.Y + shakeY)

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		p.drawSpriteBased(screen, x, y, sprite, shakeX, shakeY)
	} else {
		// Fallback to procedural rendering
		p.drawProcedural(screen, x, y, shakeX, shakeY)
	}
}

func (p *Projectile) drawSpriteBased(screen *ebiten.Image, x, y float32, sprite *ebiten.Image, shakeX, shakeY float64) {
	var mainColor, glowColor color.RGBA

	if p.Friendly {
		mainColor = color.RGBA{100, 200, 255, 255}
		glowColor = color.RGBA{50, 150, 255, 180}
	} else {
		mainColor = color.RGBA{255, 100, 100, 255}
		glowColor = color.RGBA{255, 50, 50, 180}
	}

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
	var mainColor, glowColor color.RGBA

	if p.Friendly {
		mainColor = color.RGBA{100, 200, 255, 255}
		glowColor = color.RGBA{50, 150, 255, 180}
	} else {
		mainColor = color.RGBA{255, 100, 100, 255}
		glowColor = color.RGBA{255, 50, 50, 180}
	}

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
