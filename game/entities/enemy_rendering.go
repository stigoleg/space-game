package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw renders the enemy on screen
func (e *Enemy) Draw(screen *ebiten.Image, shakeX, shakeY float64, sprite *ebiten.Image) {
	// Simple screen coordinates with shake
	x := float32(e.X + shakeX)
	y := float32(e.Y + shakeY)

	// Pulsing effect
	pulse := float32(1.0 + 0.1*math.Sin(e.AnimTimer*2))

	// Damage-based color shift (redder when damaged)
	healthRatio := float32(e.Health) / float32(e.MaxHealth)

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		e.drawSpriteBased(screen, x, y, pulse, healthRatio, sprite)
	} else {
		// Fallback to procedural rendering
		e.drawProcedural(screen, x, y, pulse, healthRatio)
	}

	// Health bar for tanks
	if e.Type == EnemyTank && e.Health < e.MaxHealth {
		barWidth := float32(60)
		barHeight := float32(6)
		healthRatioBar := float32(e.Health) / float32(e.MaxHealth)
		vector.DrawFilledRect(screen, x-barWidth/2, y-float32(e.Radius)-15, barWidth, barHeight, color.RGBA{50, 50, 50, 200}, true)
		vector.DrawFilledRect(screen, x-barWidth/2, y-float32(e.Radius)-15, barWidth*healthRatioBar, barHeight, color.RGBA{255, 50, 50, 255}, true)
	}

	// Formation indicator: glow ring if in formation
	if e.FormationType != FormationTypeNone {
		formationGlowColor := color.RGBA{100, 255, 200, 80}
		if e.IsFormationLeader {
			formationGlowColor = color.RGBA{255, 255, 100, 100} // Gold for leader
		}
		vector.StrokeCircle(screen, x, y, float32(e.Radius)+8, 1.5, formationGlowColor, true)
	}

	// Crown indicator for formation leader
	if e.IsFormationLeader {
		crownY := y - float32(e.Radius) - 8
		vector.DrawFilledCircle(screen, x, crownY, 4, color.RGBA{255, 255, 100, 255}, true)
		vector.DrawFilledCircle(screen, x-6, crownY-3, 2, color.RGBA{255, 255, 100, 255}, true)
		vector.DrawFilledCircle(screen, x+6, crownY-3, 2, color.RGBA{255, 255, 100, 255}, true)
	}

	// Burning effect: fire particles
	if e.Burning {
		e.drawBurnEffect(screen, x, y)
	}

	// Sniper lock-on indicator
	if e.Type == EnemySniper && e.SniperLockTimer > 0 && !e.SniperLocked {
		// Show charging lock-on with pulsing circles
		lockProgress := e.SniperLockTimer / 1.5 // 0.0 to 1.0
		lockAlpha := uint8(150 * lockProgress)
		lockColor := color.RGBA{255, 100, 100, lockAlpha}

		// Draw expanding circle to show lock-on progress
		lockRadius := float32(e.Radius) * (1.5 + 0.5*float32(lockProgress))
		vector.StrokeCircle(screen, x, y, lockRadius, 2, lockColor, true)

		// Draw targeting lines pointing to player
		if lockProgress > 0.5 {
			lineLength := float32(15 + 10*lockProgress)
			// Top line
			vector.StrokeLine(screen, x, y-float32(e.Radius)-5, x, y-float32(e.Radius)-5-lineLength, 2, lockColor, true)
			// Bottom line
			vector.StrokeLine(screen, x, y+float32(e.Radius)+5, x, y+float32(e.Radius)+5+lineLength, 2, lockColor, true)
		}
	}
}

func (e *Enemy) drawSpriteBased(screen *ebiten.Image, x, y, pulse, healthRatio float32, sprite *ebiten.Image) {
	// Draw shadow (depth indicator)
	radius := float32(e.Radius) * pulse
	shadowColor := color.RGBA{20, 20, 30, 100}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Draw glow effect BEFORE sprite (so sprite appears on top)
	var glowColor color.RGBA
	switch e.Type {
	case EnemyScout:
		glowColor = color.RGBA{255, 100, 50, 80}
	case EnemyDrone:
		glowColor = color.RGBA{200, 100, 255, 80}
	case EnemyHunter:
		glowColor = color.RGBA{100, 255, 100, 80}
	case EnemyTank:
		glowColor = color.RGBA{255, 150, 50, 80}
	case EnemyBomber:
		glowColor = color.RGBA{255, 200, 0, 100}
	case EnemySniper:
		glowColor = color.RGBA{50, 200, 255, 90}
	case EnemySplitter:
		glowColor = color.RGBA{255, 230, 100, 90}
	case EnemyShieldBearer:
		glowColor = color.RGBA{100, 150, 255, 90}
	}

	// Draw glow as background
	glowSize := radius * 1.4
	vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)

	// Draw sprite with options
	op := &ebiten.DrawImageOptions{}

	// Scale sprite to match enemy size
	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	// Calculate scale to match radius (sprite should be about 2x radius)
	targetSize := float64(e.Radius) * 2.0 * float64(pulse)
	scaleX := targetSize / spriteWidth
	scaleY := targetSize / spriteHeight

	op.GeoM.Scale(scaleX, scaleY)

	// Translate to enemy position (center sprite)
	op.GeoM.Translate(float64(x)-targetSize/2, float64(y)-targetSize/2)

	// Apply damage color shift (redder when damaged)
	if healthRatio < 1.0 {
		damageShift := 1.0 - healthRatio
		op.ColorScale.Scale(1.0+damageShift*0.3, 1.0-damageShift*0.5, 1.0-damageShift*0.5, 1.0)
	}

	screen.DrawImage(sprite, op)

	// Draw outline stroke for better visibility (after sprite)
	outlineColor := glowColor
	outlineColor.A = 200
	vector.StrokeCircle(screen, x, y, radius*1.2, 3, outlineColor, true)

	// Damage indicator: flickering outline when heavily damaged
	if healthRatio < 0.3 {
		damageAlpha := uint8(150 + 100*math.Sin(float64(e.AnimTimer)*4))
		damageColor := color.RGBA{255, 50, 50, damageAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.3, 2, damageColor, true)
	}
}

func (e *Enemy) drawProcedural(screen *ebiten.Image, x, y, pulse, healthRatio float32) {
	damageShift := 1.0 - healthRatio

	var mainColor, coreColor, glowColor color.RGBA

	switch e.Type {
	case EnemyScout:
		// Scout: Simple fast wedge shape - red/orange
		mainColor = color.RGBA{uint8(220 + damageShift*30), uint8(80 - damageShift*30), 60, 255}
		coreColor = color.RGBA{255, 150, 100, 255}
		glowColor = color.RGBA{255, 100, 50, 100}
	case EnemyDrone:
		// Drone: Rounded diamond - purple
		mainColor = color.RGBA{uint8(180 + damageShift*20), uint8(80 - damageShift*30), uint8(220 - damageShift*50), 255}
		coreColor = color.RGBA{255, 150, 255, 255}
		glowColor = color.RGBA{200, 100, 255, 100}
	case EnemyHunter:
		// Hunter: Angular with fins - green
		mainColor = color.RGBA{uint8(80 - damageShift*40), uint8(220 - damageShift*50), 120, 255}
		coreColor = color.RGBA{150, 255, 150, 255}
		glowColor = color.RGBA{100, 255, 100, 100}
	case EnemyTank:
		// Tank: Massive hexagon - gray with gold core
		mainColor = color.RGBA{uint8(140 + damageShift*40), uint8(140 - damageShift*30), uint8(140 - damageShift*30), 255}
		coreColor = color.RGBA{255, 200, 50, 255}
		glowColor = color.RGBA{255, 150, 50, 100}
	case EnemyBomber:
		// Bomber: Bulbous - orange with aggression
		mainColor = color.RGBA{255, uint8(140 - damageShift*40), 40, 255}
		coreColor = color.RGBA{255, 255, 100, 255}
		glowColor = color.RGBA{255, 200, 0, 150}
	case EnemySniper:
		// Sniper: Dark blue with bright cyan core
		mainColor = color.RGBA{uint8(60 + damageShift*40), uint8(100 - damageShift*30), uint8(180 - damageShift*40), 255}
		coreColor = color.RGBA{100, 255, 255, 255}
		glowColor = color.RGBA{50, 200, 255, 120}
	case EnemySplitter:
		// Splitter: Yellow/orange with split indicator
		mainColor = color.RGBA{uint8(255 - damageShift*30), uint8(220 - damageShift*40), uint8(80 + damageShift*40), 255}
		coreColor = color.RGBA{255, 255, 150, 255}
		glowColor = color.RGBA{255, 230, 100, 120}
	case EnemyShieldBearer:
		// ShieldBearer: Silver/gray with blue shield glow
		mainColor = color.RGBA{uint8(160 + damageShift*30), uint8(160 - damageShift*30), uint8(180 - damageShift*30), 255}
		coreColor = color.RGBA{200, 220, 255, 255}
		glowColor = color.RGBA{100, 150, 255, 100}
	}

	radius := float32(e.Radius) * pulse

	// Draw shadow (depth indicator)
	shadowColor := color.RGBA{20, 20, 30, 100}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Inner core (pulsing glow)
	coreSize := radius * 0.35 * pulse

	// Draw ship type-specific designs
	switch e.Type {
	case EnemyScout:
		// Scout: Small fast wedge pointing down
		drawTriangleEnemy(screen, x, y+radius*0.8, x-radius*0.6, y-radius*0.6, x+radius*0.6, y-radius*0.6, mainColor)
		vector.DrawFilledCircle(screen, x, y, coreSize*0.4, coreColor, true)

	case EnemyDrone:
		// Drone: Diamond shape
		drawTriangleEnemy(screen, x, y-radius*0.8, x-radius*0.8, y, x, y+radius*0.8, mainColor)
		drawTriangleEnemy(screen, x, y-radius*0.8, x+radius*0.8, y, x, y+radius*0.8, mainColor)
		vector.DrawFilledCircle(screen, x, y, radius*0.3*pulse, coreColor, true)

	case EnemyHunter:
		// Hunter: Angular shape with fins
		// Main body (triangle pointing down)
		drawTriangleEnemy(screen, x, y+radius*0.9, x-radius*0.5, y-radius*0.4, x+radius*0.5, y-radius*0.4, mainColor)
		// Left fin
		drawTriangleEnemy(screen, x-radius*0.5, y-radius*0.4, x-radius*1.0, y, x-radius*0.5, y+radius*0.3, mainColor)
		// Right fin
		drawTriangleEnemy(screen, x+radius*0.5, y-radius*0.4, x+radius*1.0, y, x+radius*0.5, y+radius*0.3, mainColor)
		vector.DrawFilledCircle(screen, x, y, radius*0.3, coreColor, true)

	case EnemyTank:
		// Tank: Heavy hexagon
		// Draw as layered circles to simulate hexagon
		for i := 0; i < 6; i++ {
			angle := float64(i) * math.Pi / 3
			px := x + float32(math.Cos(angle))*radius*0.7
			py := y + float32(math.Sin(angle))*radius*0.7
			vector.DrawFilledCircle(screen, px, py, radius*0.5, mainColor, true)
		}
		// Heavy core
		vector.DrawFilledCircle(screen, x, y, radius*0.5*pulse, coreColor, true)

	case EnemyBomber:
		// Bomber: Large bulbous oval shape
		// Top
		vector.DrawFilledCircle(screen, x, y-radius*0.5, radius*0.6, mainColor, true)
		// Middle (largest)
		vector.DrawFilledCircle(screen, x, y, radius*0.9, mainColor, true)
		// Bottom
		vector.DrawFilledCircle(screen, x, y+radius*0.6, radius*0.7, mainColor, true)
		vector.DrawFilledCircle(screen, x, y, radius*0.4*pulse, coreColor, true)
	case EnemySniper:
		// Sniper: Long rifle-like shape pointing down
		// Body (thin vertical rectangle)
		drawTriangleEnemy(screen, x, y+radius*0.9, x-radius*0.3, y-radius*0.7, x+radius*0.3, y-radius*0.7, mainColor)
		// Scope at top
		vector.DrawFilledCircle(screen, x, y-radius*0.6, radius*0.4, mainColor, true)
		// Barrel tip (glowing)
		if e.SniperLocked {
			// Red lock-on indicator
			vector.DrawFilledCircle(screen, x, y+radius*0.9, radius*0.3, color.RGBA{255, 50, 50, 255}, true)
		} else {
			vector.DrawFilledCircle(screen, x, y+radius*0.9, radius*0.2, coreColor, true)
		}
		vector.DrawFilledCircle(screen, x, y, radius*0.25, coreColor, true)
	case EnemySplitter:
		// Splitter: Two-part sphere that looks like it can split
		// Left half
		vector.DrawFilledCircle(screen, x-radius*0.3, y, radius*0.7, mainColor, true)
		// Right half
		vector.DrawFilledCircle(screen, x+radius*0.3, y, radius*0.7, mainColor, true)
		// Split line indicator (darker vertical line)
		drawTriangleEnemy(screen, x, y-radius*0.8, x-radius*0.1, y-radius*0.8, x-radius*0.1, y+radius*0.8, color.RGBA{100, 100, 50, 200})
		drawTriangleEnemy(screen, x, y-radius*0.8, x+radius*0.1, y-radius*0.8, x+radius*0.1, y+radius*0.8, color.RGBA{100, 100, 50, 200})
		vector.DrawFilledCircle(screen, x, y, radius*0.3*pulse, coreColor, true)
	case EnemyShieldBearer:
		// ShieldBearer: Heavy octagon shape
		// Draw 8-sided shape with circles
		for i := 0; i < 8; i++ {
			angle := float64(i) * math.Pi / 4
			px := x + float32(math.Cos(angle))*radius*0.65
			py := y + float32(math.Sin(angle))*radius*0.65
			vector.DrawFilledCircle(screen, px, py, radius*0.45, mainColor, true)
		}
		// Heavy core
		vector.DrawFilledCircle(screen, x, y, radius*0.45*pulse, coreColor, true)
		// Shield bar
		if e.ShieldPoints > 0 {
			shieldBarWidth := float32(50)
			shieldBarHeight := float32(5)
			shieldRatio := float32(e.ShieldPoints) / float32(e.MaxShieldPoints)
			barX := x - shieldBarWidth/2
			barY := y + radius + 8
			vector.DrawFilledRect(screen, barX, barY, shieldBarWidth, shieldBarHeight, color.RGBA{40, 40, 60, 200}, true)
			vector.DrawFilledRect(screen, barX, barY, shieldBarWidth*shieldRatio, shieldBarHeight, color.RGBA{100, 150, 255, 255}, true)
		}
	}

	// Pulsing core
	vector.DrawFilledCircle(screen, x, y, coreSize, coreColor, true)

	// Outer glow
	glowSize := radius + 3
	vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)

	// Highlight edge (3D effect)
	highlightSize := radius * 0.25
	highlightColor := color.RGBA{mainColor.R + 40, mainColor.G + 40, mainColor.B + 40, 150}
	vector.DrawFilledCircle(screen, x-radius*0.3, y-radius*0.3, highlightSize, highlightColor, true)

	// Damage indicator: flickering outline when heavily damaged
	if healthRatio < 0.3 {
		damageAlpha := uint8(150 + 100*math.Sin(float64(e.AnimTimer)*4))
		damageColor := color.RGBA{255, 50, 50, damageAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.1, 2, damageColor, true)
	}
}

// drawBurnEffect draws fire particles around a burning enemy
func (e *Enemy) drawBurnEffect(screen *ebiten.Image, x, y float32) {
	// Draw 8-12 fire particles orbiting the enemy
	numFlames := 10
	for i := 0; i < numFlames; i++ {
		// Orbit angle with rotation based on AnimTimer
		angle := e.AnimTimer*4 + float64(i)*math.Pi*2/float64(numFlames)

		// Vary orbit distance slightly per particle
		orbitDist := float64(e.Radius) * (1.2 + 0.2*math.Sin(e.AnimTimer*3+float64(i)))

		offsetX := math.Cos(angle) * orbitDist
		offsetY := math.Sin(angle)*orbitDist - 5 // Rise upward slightly

		flameX := float32(e.X + offsetX)
		flameY := float32(e.Y + offsetY)

		// Flickering orange-red-yellow flames
		// Cycle through fire colors: red (0) -> orange (20) -> yellow (40)
		hue := int(10 + 20*math.Sin(e.AnimTimer*5+float64(i)*0.5))
		if hue < 0 {
			hue = 0
		}
		if hue > 40 {
			hue = 40
		}

		// Convert hue to RGB (simple fire color mapping)
		var flameColor color.RGBA
		if hue < 20 {
			// Red to orange
			flameColor = color.RGBA{255, uint8(100 + hue*7), 0, 255}
		} else {
			// Orange to yellow
			flameColor = color.RGBA{255, uint8(240 + (hue - 20)), uint8((hue - 20) * 5), 255}
		}

		// Flickering alpha
		flameColor.A = uint8(150 + 80*math.Sin(e.AnimTimer*6+float64(i)*0.7))

		// Vary particle size
		flameSize := float32(2.5 + 1.5*math.Sin(e.AnimTimer*4+float64(i)*0.3))

		vector.DrawFilledCircle(screen, flameX, flameY, flameSize, flameColor, true)

		// Add smaller bright core to some particles
		if i%2 == 0 {
			coreColor := color.RGBA{255, 255, 200, 200}
			vector.DrawFilledCircle(screen, flameX, flameY, flameSize*0.5, coreColor, true)
		}
	}

	// Add orange glow overlay on enemy body
	glowPulse := float32(0.8 + 0.2*math.Sin(e.AnimTimer*3))
	glowColor := color.RGBA{255, 100, 0, uint8(60 * glowPulse)}
	vector.DrawFilledCircle(screen, x, y, float32(e.Radius)*1.3*glowPulse, glowColor, true)
}

// drawTriangleEnemy is a helper function to draw filled triangles for enemies
func drawTriangleEnemy(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, col color.RGBA) {
	vector.DrawFilledCircle(screen, x1, y1, 2, col, true)
	vector.DrawFilledCircle(screen, x2, y2, 2, col, true)
	vector.DrawFilledCircle(screen, x3, y3, 2, col, true)
	vector.StrokeLine(screen, x1, y1, x2, y2, 3, col, true)
	vector.StrokeLine(screen, x2, y2, x3, y3, 3, col, true)
	vector.StrokeLine(screen, x3, y3, x1, y1, 3, col, true)
}
