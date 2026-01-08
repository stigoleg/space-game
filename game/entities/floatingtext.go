package entities

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// FloatingText represents a temporary text that floats and fades
type FloatingText struct {
	X, Y      float64
	VelY      float64 // Upward velocity
	Text      string
	TextColor color.RGBA
	Life      float64 // Time remaining (seconds)
	MaxLife   float64
	Active    bool
}

// NewFloatingText creates a new floating text
func NewFloatingText(x, y float64, text string, col color.RGBA) *FloatingText {
	return &FloatingText{
		X:         x,
		Y:         y,
		VelY:      -2.0, // Float upward
		Text:      text,
		TextColor: col,
		Life:      2.0, // 2 seconds lifetime
		MaxLife:   2.0,
		Active:    true,
	}
}

// NewFloatingScore creates a floating text for score display
func NewFloatingScore(x, y float64, score int) *FloatingText {
	text := fmt.Sprintf("+%d", score)
	return NewFloatingText(x, y, text, color.RGBA{255, 200, 100, 255})
}

// NewFloatingDamage creates a floating text for damage display
func NewFloatingDamage(x, y float64, damage int) *FloatingText {
	text := fmt.Sprintf("-%d", damage)
	return NewFloatingText(x, y, text, color.RGBA{255, 50, 50, 255})
}

// NewFloatingUpgrade creates a floating text for weapon upgrade
func NewFloatingUpgrade(x, y float64, level int) *FloatingText {
	text := fmt.Sprintf("Lvl %d", level)
	return NewFloatingText(x, y, text, color.RGBA{100, 255, 200, 255})
}

// Update moves the floating text
func (ft *FloatingText) Update() {
	if !ft.Active {
		return
	}

	ft.Life -= 1.0 / 60.0
	if ft.Life <= 0 {
		ft.Active = false
		ft.Life = 0
		return
	}

	ft.Y += ft.VelY
	ft.VelY *= 0.98 // Slow down slightly
}

// Draw renders the floating text
func (ft *FloatingText) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	if !ft.Active {
		return
	}

	// Fade out as life decreases
	lifeRatio := ft.Life / ft.MaxLife
	alpha := uint8(255 * lifeRatio)
	col := color.RGBA{ft.TextColor.R, ft.TextColor.G, ft.TextColor.B, alpha}

	// Scale text slightly based on fade
	x := int(ft.X + shakeX)
	y := int(ft.Y + shakeY)

	// Draw text
	text.Draw(screen, ft.Text, basicfont.Face7x13, x-20, y, col)

	// Optional: Draw a glow effect for important text
	if lifeRatio > 0.5 {
		glowAlpha := uint8(100 * (1 - lifeRatio))
		glowCol := color.RGBA{col.R, col.G, col.B, glowAlpha}
		text.Draw(screen, ft.Text, basicfont.Face7x13, x-22, y-2, glowCol)
		text.Draw(screen, ft.Text, basicfont.Face7x13, x-18, y+2, glowCol)
	}
}

// FloatingParticle represents a visual particle effect that floats
type FloatingParticle struct {
	X, Y       float64
	VelX, VelY float64
	Life       float64
	MaxLife    float64
	Active     bool
	Color      color.RGBA
	Scale      float64
}

// NewFloatingParticle creates a particle at (x, y) moving in a random direction
func NewFloatingParticle(x, y float64, col color.RGBA) *FloatingParticle {
	angle := math.Pi * 2 * rand.Float64()
	speed := 2.0 + rand.Float64()*3.0
	return &FloatingParticle{
		X:       x,
		Y:       y,
		VelX:    math.Cos(angle) * speed,
		VelY:    math.Sin(angle) * speed,
		Life:    1.5,
		MaxLife: 1.5,
		Active:  true,
		Color:   col,
		Scale:   1.0,
	}
}

// Update moves and fades the particle
func (fp *FloatingParticle) Update() {
	if !fp.Active {
		return
	}

	fp.Life -= 1.0 / 60.0
	if fp.Life <= 0 {
		fp.Active = false
		return
	}

	fp.X += fp.VelX
	fp.Y += fp.VelY
	fp.VelX *= 0.95
	fp.VelY *= 0.95

	// Scale down as life decreases
	fp.Scale = (fp.Life / fp.MaxLife) * 1.0
}

// Draw renders the particle
func (fp *FloatingParticle) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	if !fp.Active {
		return
	}

	lifeRatio := fp.Life / fp.MaxLife
	alpha := uint8(255 * lifeRatio)
	col := color.RGBA{fp.Color.R, fp.Color.G, fp.Color.B, alpha}

	// Simple circle particle
	x := float32(fp.X + shakeX)
	y := float32(fp.Y + shakeY)

	// Would use vector.DrawFilledCircle here but keeping it simple
	// This can be enhanced with proper drawing
	_ = x
	_ = y
	_ = col
}
