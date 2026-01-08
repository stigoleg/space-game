package entities

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ImpactEffect struct {
	X, Y      float64
	Radius    float64
	MaxRadius float64
	Life      float64
	MaxLife   float64
	Color     color.RGBA
	Active    bool
	Expanding bool // Is the ring still expanding?
}

// NewImpactEffect creates a new impact effect (hit ring)
func NewImpactEffect(x, y float64, maxRadius float64, color color.RGBA) *ImpactEffect {
	return &ImpactEffect{
		X:         x,
		Y:         y,
		Radius:    2,
		MaxRadius: maxRadius,
		Life:      0.4,
		MaxLife:   0.4,
		Color:     color,
		Active:    true,
		Expanding: true,
	}
}

func (i *ImpactEffect) Update() {
	i.Life -= 1.0 / 60.0

	if i.Expanding {
		i.Radius += i.MaxRadius / (i.MaxLife * 60)
		if i.Radius >= i.MaxRadius {
			i.Radius = i.MaxRadius
			i.Expanding = false
		}
	}

	if i.Life <= 0 {
		i.Active = false
	}
}

func (i *ImpactEffect) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	if !i.Active {
		return
	}

	lifeRatio := i.Life / i.MaxLife
	alpha := uint8(200 * lifeRatio)
	ringColor := color.RGBA{i.Color.R, i.Color.G, i.Color.B, alpha}

	x := float32(i.X + shakeX)
	y := float32(i.Y + shakeY)
	radius := float32(i.Radius)

	// Draw expanding ring
	vector.StrokeCircle(screen, x, y, radius, 2, ringColor, true)

	// Draw fading glow
	glowAlpha := uint8(100 * lifeRatio)
	glowColor := color.RGBA{i.Color.R, i.Color.G, i.Color.B, glowAlpha}
	vector.DrawFilledCircle(screen, x, y, radius, glowColor, true)
}
