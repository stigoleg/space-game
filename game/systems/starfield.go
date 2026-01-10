package systems

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Star struct {
	X, Y    float64
	Size    float64
	Speed   float64
	Bright  float64
	Twinkle float64 // Pre-computed twinkle value to avoid per-frame math.Sin
}

type StarField struct {
	layers [][]Star
	width  int
	height int
}

func NewStarField(width, height int) *StarField {
	sf := &StarField{
		width:  width,
		height: height,
		layers: make([][]Star, 3),
	}

	// Layer 0: Far stars (slow, small, dim)
	sf.layers[0] = make([]Star, 100)
	for i := range sf.layers[0] {
		x := rand.Float64() * float64(width)
		y := rand.Float64() * float64(height)
		sf.layers[0][i] = Star{
			X:       x,
			Y:       y,
			Size:    rand.Float64()*1 + 0.5,
			Speed:   0.3,
			Bright:  rand.Float64()*0.3 + 0.2,
			Twinkle: 0.8 + 0.2*math.Sin(x+y+100), // Pre-compute twinkle for layer 0
		}
	}

	// Layer 1: Mid stars (medium speed, medium size)
	sf.layers[1] = make([]Star, 70)
	for i := range sf.layers[1] {
		x := rand.Float64() * float64(width)
		y := rand.Float64() * float64(height)
		sf.layers[1][i] = Star{
			X:       x,
			Y:       y,
			Size:    rand.Float64()*1.5 + 1,
			Speed:   0.8,
			Bright:  rand.Float64()*0.4 + 0.4,
			Twinkle: 0.8 + 0.2*math.Sin(x+y+200), // Pre-compute twinkle for layer 1
		}
	}

	// Layer 2: Close stars (fast, large, bright)
	sf.layers[2] = make([]Star, 40)
	for i := range sf.layers[2] {
		x := rand.Float64() * float64(width)
		y := rand.Float64() * float64(height)
		sf.layers[2][i] = Star{
			X:       x,
			Y:       y,
			Size:    rand.Float64()*2 + 1.5,
			Speed:   1.5,
			Bright:  rand.Float64()*0.3 + 0.7,
			Twinkle: 0.8 + 0.2*math.Sin(x+y+300), // Pre-compute twinkle for layer 2
		}
	}

	return sf
}

func (sf *StarField) Update() {
	for l := range sf.layers {
		layerOffset := float64((l + 1) * 100) // 100, 200, 300 for layers 0, 1, 2
		for i := range sf.layers[l] {
			sf.layers[l][i].Y += sf.layers[l][i].Speed
			if sf.layers[l][i].Y > float64(sf.height) {
				sf.layers[l][i].Y = 0
				sf.layers[l][i].X = rand.Float64() * float64(sf.width)
				// Recalculate twinkle when star wraps
				sf.layers[l][i].Twinkle = 0.8 + 0.2*math.Sin(sf.layers[l][i].X+sf.layers[l][i].Y+layerOffset)
			}
		}
	}
}

func (sf *StarField) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	for l := range sf.layers {
		for _, star := range sf.layers[l] {
			x := float32(star.X + shakeX*float64(l+1)*0.3)
			y := float32(star.Y + shakeY*float64(l+1)*0.3)

			// Use pre-computed twinkle value (avoids math.Sin per star per frame)
			alpha := uint8(star.Bright * star.Twinkle * 255)

			// Color tint based on layer
			var c color.RGBA
			switch l {
			case 0: // Blue tint (distant)
				c = color.RGBA{180, 180, 255, alpha}
			case 1: // White
				c = color.RGBA{255, 255, 255, alpha}
			case 2: // Slight yellow (close)
				c = color.RGBA{255, 255, 220, alpha}
			}

			// Draw star with glow
			size := float32(star.Size)
			if size > 2 {
				// Glow for larger stars
				glowColor := color.RGBA{c.R, c.G, c.B, alpha / 4}
				vector.DrawFilledCircle(screen, x, y, size*2, glowColor, true)
			}
			vector.DrawFilledCircle(screen, x, y, size, c, true)
		}
	}
}
