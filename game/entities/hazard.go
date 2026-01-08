package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HazardType represents different types of environmental hazards
type HazardType int

const (
	HazardTypeBarrier HazardType = iota
	HazardTypeMagneticField
	HazardTypeRadiationZone
	HazardTypeBlackHole
)

// Hazard represents an environmental hazard
type Hazard struct {
	X, Y        float64
	Radius      float64
	Type        HazardType
	Active      bool
	Health      int
	MaxHealth   int
	AnimTimer   float64
	PullForce   float64 // For magnetic fields and black holes
	DamageRate  float64 // For radiation zones
	LastDamage  float64 // Time since last damage dealt
	Name        string
	Description string
}

// NewHazard creates a new hazard
func NewHazard(x, y float64, hazardType HazardType) *Hazard {
	h := &Hazard{
		X:          x,
		Y:          y,
		Active:     true,
		AnimTimer:  0,
		LastDamage: 0,
	}

	switch hazardType {
	case HazardTypeBarrier:
		h.Type = HazardTypeBarrier
		h.Radius = 40
		h.Health = 100
		h.MaxHealth = 100
		h.Name = "Energy Barrier"
		h.Description = "Blocks projectiles and movement"

	case HazardTypeMagneticField:
		h.Type = HazardTypeMagneticField
		h.Radius = 80
		h.Health = 50
		h.MaxHealth = 50
		h.PullForce = 150.0
		h.Name = "Magnetic Field"
		h.Description = "Pulls in projectiles and enemies"

	case HazardTypeRadiationZone:
		h.Type = HazardTypeRadiationZone
		h.Radius = 100
		h.Health = 30
		h.MaxHealth = 30
		h.DamageRate = 10.0
		h.Name = "Radiation Zone"
		h.Description = "Deals continuous damage"

	case HazardTypeBlackHole:
		h.Type = HazardTypeBlackHole
		h.Radius = 60
		h.Health = 200
		h.MaxHealth = 200
		h.PullForce = 300.0
		h.Name = "Black Hole"
		h.Description = "Extreme gravity - instant death if touched"
	}

	return h
}

// Update updates the hazard state
func (h *Hazard) Update() {
	h.AnimTimer += 0.1
	h.LastDamage += 1.0 / 60.0
}

// TakeDamage applies damage to the hazard
func (h *Hazard) TakeDamage(damage int) {
	h.Health -= damage
	if h.Health <= 0 {
		h.Active = false
	}
}

// Draw renders the hazard
func (h *Hazard) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	x := float32(h.X + shakeX)
	y := float32(h.Y + shakeY)

	healthRatio := float32(h.Health) / float32(h.MaxHealth)
	pulse := float32(1.0 + 0.1*math.Sin(float64(h.AnimTimer)*2))

	switch h.Type {
	case HazardTypeBarrier:
		h.drawBarrier(screen, x, y, pulse)

	case HazardTypeMagneticField:
		h.drawMagneticField(screen, x, y, pulse)

	case HazardTypeRadiationZone:
		h.drawRadiationZone(screen, x, y, healthRatio)

	case HazardTypeBlackHole:
		h.drawBlackHole(screen, x, y, pulse)
	}
}

func (h *Hazard) drawBarrier(screen *ebiten.Image, x, y float32, pulse float32) {
	radius := float32(h.Radius) * pulse

	// Yellow glowing barrier
	barrierColor := color.RGBA{255, 255, 100, 150}

	// Draw concentric circles
	for i := 0; i < 3; i++ {
		r := radius - float32(i*8)
		if r > 0 {
			alpha := uint8(150 - i*40)
			c := color.RGBA{barrierColor.R, barrierColor.G, barrierColor.B, alpha}
			vector.StrokeCircle(screen, x, y, r, 2, c, true)
		}
	}

	// Center core
	vector.DrawFilledCircle(screen, x, y, radius*0.3, barrierColor, true)
}

func (h *Hazard) drawMagneticField(screen *ebiten.Image, x, y float32, pulse float32) {
	radius := float32(h.Radius)

	// Cyan color for magnetic field
	fieldColor := color.RGBA{100, 255, 255, 100}
	coreColor := color.RGBA{150, 255, 255, 200}

	// Draw spiraling lines for magnetic effect
	for i := 0; i < 8; i++ {
		angle := float64(i) * math.Pi / 4
		angle += float64(h.AnimTimer) * 0.05

		startX := x + float32(math.Cos(angle))*radius
		startY := y + float32(math.Sin(angle))*radius
		endX := x + float32(math.Cos(angle))*radius*0.5
		endY := y + float32(math.Sin(angle))*radius*0.5

		vector.StrokeLine(screen, startX, startY, endX, endY, 2, fieldColor, true)
	}

	// Core
	vector.DrawFilledCircle(screen, x, y, radius*0.2*pulse, coreColor, true)

	// Outer glow
	vector.StrokeCircle(screen, x, y, radius, 3, fieldColor, true)
}

func (h *Hazard) drawRadiationZone(screen *ebiten.Image, x, y float32, healthRatio float32) {
	radius := float32(h.Radius)

	// Green/toxic color
	radiationColor := color.RGBA{100, 255, 100, 120}
	warningColor := color.RGBA{255, 200, 0, 150}

	// Draw radiating lines
	lineCount := 12
	for i := 0; i < lineCount; i++ {
		angle := float64(i) * 2 * math.Pi / float64(lineCount)
		angle += float64(h.AnimTimer) * 0.1

		startX := x + float32(math.Cos(angle))*radius
		startY := y + float32(math.Sin(angle))*radius
		endX := x
		endY := y

		vector.StrokeLine(screen, startX, startY, endX, endY, 2, radiationColor, true)
	}

	// Center with pulsing core
	coreRadius := radius * 0.2 * (0.8 + 0.2*float32(math.Sin(float64(h.AnimTimer)*2)))
	vector.DrawFilledCircle(screen, x, y, coreRadius, warningColor, true)

	// Warning ring
	vector.StrokeCircle(screen, x, y, radius, 2, radiationColor, true)
}

func (h *Hazard) drawBlackHole(screen *ebiten.Image, x, y float32, pulse float32) {
	radius := float32(h.Radius) * pulse

	// Black with white/blue accretion disk
	// Draw event horizon
	vector.DrawFilledCircle(screen, x, y, radius*0.7, color.RGBA{20, 20, 30, 255}, true)

	// Accretion disk - rotating rings
	for i := 0; i < 4; i++ {
		ringRadius := radius * (0.4 + float32(i)*0.15)
		angle := float64(h.AnimTimer) * (0.1 - float64(i)*0.02)

		// Draw arc of ring
		steps := 16
		for j := 0; j < steps; j++ {
			a1 := angle + float64(j)*2*math.Pi/float64(steps)
			a2 := angle + float64(j+1)*2*math.Pi/float64(steps)

			p1x := x + float32(math.Cos(a1))*ringRadius
			p1y := y + float32(math.Sin(a1))*ringRadius
			p2x := x + float32(math.Cos(a2))*ringRadius
			p2y := y + float32(math.Sin(a2))*ringRadius

			alpha := uint8(200 - i*40)
			ringColor := color.RGBA{100, 200, 255, alpha}
			vector.StrokeLine(screen, p1x, p1y, p2x, p2y, 2, ringColor, true)
		}
	}

	// Central singularity
	vector.DrawFilledCircle(screen, x, y, radius*0.2, color.RGBA{255, 255, 255, 200}, true)

	// Warning glow
	vector.StrokeCircle(screen, x, y, radius, 2, color.RGBA{255, 100, 100, 150}, true)
}

// GetCollisionRadius returns the effective collision radius
func (h *Hazard) GetCollisionRadius() float64 {
	switch h.Type {
	case HazardTypeBarrier:
		return h.Radius
	case HazardTypeMagneticField:
		return h.Radius * 1.2 // Larger pull radius
	case HazardTypeRadiationZone:
		return h.Radius
	case HazardTypeBlackHole:
		return h.Radius * 0.8 // Tighter collision
	default:
		return h.Radius
	}
}

// IsDangerous returns true if hazard can instantly kill
func (h *Hazard) IsDangerous() bool {
	return h.Type == HazardTypeBlackHole
}

// HazardSpawner manages hazard spawning
type HazardSpawner struct {
	hazards       []*Hazard
	spawnTimer    float64
	spawnInterval float64
}

// NewHazardSpawner creates a new hazard spawner
func NewHazardSpawner() *HazardSpawner {
	return &HazardSpawner{
		hazards:       make([]*Hazard, 0),
		spawnTimer:    0,
		spawnInterval: 15.0, // Spawn every 15 seconds
	}
}

// Update updates hazard spawning
func (hs *HazardSpawner) Update() []*Hazard {
	hs.spawnTimer += 1.0 / 60.0

	newHazards := make([]*Hazard, 0)

	// Periodically spawn new hazards
	if hs.spawnTimer > hs.spawnInterval {
		// Spawn logic here
		hs.spawnTimer = 0
	}

	// Update existing hazards
	var active []*Hazard
	for _, h := range hs.hazards {
		h.Update()
		if h.Active {
			active = append(active, h)
		}
	}
	hs.hazards = active

	return newHazards
}

// AddHazard adds a hazard to be tracked
func (hs *HazardSpawner) AddHazard(h *Hazard) {
	hs.hazards = append(hs.hazards, h)
}

// GetHazards returns all active hazards
func (hs *HazardSpawner) GetHazards() []*Hazard {
	return hs.hazards
}
