package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw renders the player ship with all visual effects
func (p *Player) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	// Blink when invincible
	if p.InvincTimer > 0 && int(p.InvincTimer*10)%2 == 0 {
		return
	}

	// Simple screen coordinates with shake
	x := float32(p.X + shakeX)
	y := float32(p.Y + shakeY)

	// Draw thruster trail particles (using ring buffer)
	for i := 0; i < p.ThrusterTrailLen; i++ {
		// Calculate index in ring buffer (oldest first)
		idx := (p.ThrusterTrailHead - p.ThrusterTrailLen + i + MaxThrusterTrailLen) % MaxThrusterTrailLen
		trail := p.ThrusterTrail[idx]

		lifeRatio := trail.Life / 0.5
		alpha := uint8(150 * lifeRatio)
		trailColor := color.RGBA{100, 150, 255, alpha}
		trailGlowColor := color.RGBA{150, 200, 255, uint8(100 * lifeRatio)}

		trailX := float32(trail.X + shakeX)
		trailY := float32(trail.Y + shakeY)
		trailSize := float32(3 * (1 - (1-lifeRatio)*(1-lifeRatio)))

		// Draw glow
		vector.DrawFilledCircle(screen, trailX, trailY, trailSize*1.8, trailGlowColor, true)
		// Draw particle
		vector.DrawFilledCircle(screen, trailX, trailY, trailSize, trailColor, true)
	}

	// Draw shadow beneath ship (depth indicator)
	shadowColor := color.RGBA{20, 20, 30, 100}
	shadowSize := float32(p.Radius) * 0.5
	vector.DrawFilledCircle(screen, x, y+float32(p.Radius)+shadowSize, shadowSize, shadowColor, true)

	// Draw polygon-based ship - triangular arrow shape
	shipColor := color.RGBA{100, 160, 220, 255}
	radius := float32(p.Radius)

	// Main hull - triangle pointing up
	// Top vertex (nose)
	topX := x
	topY := y - radius*1.1

	// Bottom-left wing
	leftX := x - radius*0.8
	leftY := y + radius*0.7

	// Bottom-right wing
	rightX := x + radius*0.8
	rightY := y + radius*0.7

	// Draw main hull as filled polygon
	// Using triangles (vector.DrawFilledRect for simpler polygon effects)
	drawTriangle(screen, topX, topY, leftX, leftY, rightX, rightY, shipColor)

	// Draw cockpit window (diamond shape)
	cockpitColor := color.RGBA{220, 240, 255, 255}
	vector.DrawFilledCircle(screen, x, y-radius*0.4, 5, cockpitColor, true)

	// Draw engine thrusters (animated flame effect)
	engineIntensity := float32(0.5 + 0.3*math.Sin(p.EngineGlow))
	engineTrailColor1 := color.RGBA{100, 150, 255, 200}
	engineTrailColor2 := color.RGBA{255, 150, 100, 180}

	// Left engine
	vector.DrawFilledCircle(screen, leftX+radius*0.3, leftY+radius*0.5, float32(6)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, leftX+radius*0.3, leftY+radius*0.8, float32(3)*engineIntensity, engineTrailColor2, true)

	// Right engine
	vector.DrawFilledCircle(screen, rightX-radius*0.3, rightY+radius*0.5, float32(6)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, rightX-radius*0.3, rightY+radius*0.8, float32(3)*engineIntensity, engineTrailColor2, true)

	// Center engine
	vector.DrawFilledCircle(screen, x, y+radius*0.6, float32(5)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, x, y+radius*1.0, float32(2)*engineIntensity, engineTrailColor2, true)

	// Draw wing accents
	wingAccentColor := color.RGBA{80, 140, 200, 255}
	vector.StrokeCircle(screen, leftX, leftY, radius*0.3, 1.5, wingAccentColor, true)
	vector.StrokeCircle(screen, rightX, rightY, radius*0.3, 1.5, wingAccentColor, true)

	// Shield effect
	if p.Shield > 0 {
		shieldAlpha := uint8(50 + float64(p.Shield)/float64(p.MaxShield)*100)
		shieldColor := color.RGBA{100, 200, 255, shieldAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.0, 2, shieldColor, true)
	}

	// Nose highlight for depth
	highlightColor := color.RGBA{200, 230, 255, 200}
	vector.DrawFilledCircle(screen, topX, topY+radius*0.2, 4, highlightColor, true)

	// Draw charge indicator (glow around ship when charging)
	if p.ChargeLevel > 0.1 {
		chargeAlpha := uint8(100 + float64(p.ChargeLevel)*155)
		chargeColor := color.RGBA{255, uint8(100 + float64(p.ChargeLevel)*155), 100, chargeAlpha}
		chargeRadius := radius*1.2 + float32(p.ChargeLevel)*10
		vector.StrokeCircle(screen, x, y, chargeRadius, 2, chargeColor, true)

		// Inner charge aura
		innerChargeColor := color.RGBA{255, 200, 100, uint8(50 + float64(p.ChargeLevel)*100)}
		vector.DrawFilledCircle(screen, x, y, radius*1.1, innerChargeColor, true)
	}

	// Draw ultimate indicator (star glow)
	if p.UltimateCharge > 0.5 {
		ultimateAlpha := uint8(100 + float64(p.UltimateCharge)*155)
		ultimateColor := color.RGBA{255, 50 + uint8(float64(p.UltimateCharge)*200), 255, ultimateAlpha}
		ultimateRadius := radius + float32(p.UltimateCharge)*8
		vector.StrokeCircle(screen, x, y, ultimateRadius, 3, ultimateColor, true)
	}

	// Ultimate active effect - intense glow
	if p.UltimateActive {
		pulseIntensity := float32(0.5 + 0.5*math.Sin(p.UltimateTimer*math.Pi))
		ultimateGlowColor := color.RGBA{200, 100, 255, uint8(200 * pulseIntensity)}
		vector.DrawFilledCircle(screen, x, y, radius*1.3*pulseIntensity, ultimateGlowColor, true)
	}
}

// drawTriangle is a helper function to draw a filled triangle
func drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, col color.Color) {
	// Draw using three filled circles connected - crude but effective
	// For better polygon support, we can use vector paths or draw filled polygons another way

	// Create paths for triangle edges with thick lines to simulate fill
	// This creates a filled triangle effect
	const steps = 20
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)
		// Top to left edge
		edgeX1 := x1*(1-t) + x2*t
		edgeY1 := y1*(1-t) + y2*t

		// Top to right edge
		edgeX2 := x1*(1-t) + x3*t
		edgeY2 := y1*(1-t) + y3*t

		// Draw line between edges
		radius := float32(math.Abs(float64(edgeX2-edgeX1)) / 2)
		if radius > 0.5 {
			midX := (edgeX1 + edgeX2) / 2
			midY := (edgeY1 + edgeY2) / 2
			vector.DrawFilledCircle(screen, midX, midY, radius+1, col, true)
		}
	}

	// Draw bottom edge
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)
		edgeX1 := x2*(1-t) + x3*t
		edgeY1 := y2*(1-t) + y3*t
		vector.DrawFilledCircle(screen, edgeX1, edgeY1, 2, col, true)
	}
}
