package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Particle struct {
	X, Y       float64
	VelX, VelY float64
	Size       float64
	Life       float64
	MaxLife    float64
	Color      color.RGBA
}

type ExplosionType int

const (
	ExplosionStandard ExplosionType = iota
	ExplosionBlast                  // Bigger, brighter
	ExplosionSmoke                  // Slower, lingering
	ExplosionEnergy                 // Blue/cyan energy
)

type Explosion struct {
	X, Y       float64
	Particles  []Particle
	Active     bool
	Timer      float64
	ExpType    ExplosionType
	BurstScale float64 // Scale for burst effect
}

func NewExplosion(x, y, size float64) *Explosion {
	return NewExplosionWithType(x, y, size, ExplosionStandard)
}

func NewExplosionWithType(x, y, size float64, expType ExplosionType) *Explosion {
	// Optimized: Reduced particle count for better performance
	numParticles := int(size * 1.2) // Reduced from size * 2
	if numParticles < 6 {           // Reduced from 10
		numParticles = 6
	}
	if numParticles > 40 { // Reduced from 80
		numParticles = 40
	}

	particles := make([]Particle, numParticles)

	for i := range particles {
		angle := rand.Float64() * math.Pi * 2

		var speed, life float64
		var r, g, b uint8

		switch expType {
		case ExplosionBlast:
			// Bigger, faster burst
			speed = rand.Float64()*6 + 3
			life = rand.Float64()*0.6 + 0.2
			r = uint8(255)
			g = uint8(180 + rand.Intn(75))
			b = uint8(rand.Intn(50))

		case ExplosionSmoke:
			// Slower, lingering smoke
			speed = rand.Float64()*2 + 0.5
			life = rand.Float64()*1.2 + 0.8
			gray := uint8(100 + rand.Intn(80))
			r, g, b = gray, gray, gray

		case ExplosionEnergy:
			// Blue/cyan energy burst
			speed = rand.Float64()*5 + 2
			life = rand.Float64()*0.7 + 0.3
			r = uint8(100 + rand.Intn(100))
			g = uint8(150 + rand.Intn(100))
			b = uint8(255)

		default: // ExplosionStandard
			// Standard fire explosion
			speed = rand.Float64()*4 + 2
			life = rand.Float64()*0.5 + 0.3
			r = uint8(200 + rand.Intn(55))
			g = uint8(100 + rand.Intn(100))
			b = uint8(rand.Intn(50))
		}

		particles[i] = Particle{
			X:       x,
			Y:       y,
			VelX:    math.Cos(angle) * speed,
			VelY:    math.Sin(angle) * speed,
			Size:    rand.Float64()*size/5 + 2,
			Life:    life,
			MaxLife: life,
			Color:   color.RGBA{r, g, b, 255},
		}
	}

	return &Explosion{
		X:          x,
		Y:          y,
		Particles:  particles,
		Active:     true,
		Timer:      0,
		ExpType:    expType,
		BurstScale: 1.0 + rand.Float64()*0.5, // Slight variation
	}
}

func (e *Explosion) Update() {
	e.Timer += 1.0 / 60.0
	allDead := true

	for i := range e.Particles {
		p := &e.Particles[i]
		if p.Life > 0 {
			p.X += p.VelX
			p.Y += p.VelY
			p.VelX *= 0.93 // Slow down slightly
			p.VelY *= 0.93
			p.Life -= 1.0 / 60.0
			p.Size *= 0.96 // Shrink slower for better effect
			allDead = false
		}
	}

	if allDead || e.Timer > 3 {
		e.Active = false
	}
}

// Screen bounds for particle culling (matches game.go constants)
const (
	explosionScreenWidth  = 1280
	explosionScreenHeight = 960
	explosionCullBuffer   = 50 // Extra margin for particles with glow
)

func (e *Explosion) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	for _, p := range e.Particles {
		if p.Life <= 0 {
			continue
		}

		// Early culling: skip particles outside screen bounds
		px := p.X + shakeX
		py := p.Y + shakeY
		if px < -explosionCullBuffer || px > explosionScreenWidth+explosionCullBuffer ||
			py < -explosionCullBuffer || py > explosionScreenHeight+explosionCullBuffer {
			continue
		}

		lifeRatio := p.Life / p.MaxLife
		alpha := uint8(255 * lifeRatio)
		c := color.RGBA{p.Color.R, p.Color.G, p.Color.B, alpha}

		x := float32(px)
		y := float32(py)
		size := float32(p.Size)

		// Optimized: Reduced glow complexity - single glow layer only
		var glowSize float32
		var glowAlpha uint8

		switch e.ExpType {
		case ExplosionBlast:
			glowSize = size * 1.8
			glowAlpha = alpha / 2
		case ExplosionSmoke:
			glowSize = size * 1.2
			glowAlpha = alpha / 4
		case ExplosionEnergy:
			glowSize = size * 2.0
			glowAlpha = alpha / 2
		default:
			glowSize = size * 1.4
			glowAlpha = alpha / 3
		}

		// Draw glow (combined effect with particle)
		glowColor := color.RGBA{p.Color.R, p.Color.G, p.Color.B, glowAlpha}
		vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)

		// Draw particle core
		vector.DrawFilledCircle(screen, x, y, size, c, true)
	}
}

// Poolable interface implementation

// Reset resets the explosion to default state for reuse
func (e *Explosion) Reset() {
	e.X = 0
	e.Y = 0
	e.Particles = e.Particles[:0] // Keep capacity, clear length
	e.Active = false
	e.Timer = 0
	e.ExpType = ExplosionStandard
	e.BurstScale = 1.0
}

// IsActive returns whether the explosion is active
func (e *Explosion) IsActive() bool {
	return e.Active
}

// SetActive sets the active state of the explosion
func (e *Explosion) SetActive(active bool) {
	e.Active = active
}
