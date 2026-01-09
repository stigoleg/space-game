package systems

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// ============================================================================
// ENEMY SPRITE GENERATORS
// ============================================================================

// generateScoutSprite creates a sharp triangular scout ship (aggressive, fast)
func (sm *SpriteManager) generateScoutSprite() *ebiten.Image {
	size := 64
	img := ebiten.NewImage(size, size)

	// Create RGBA image for drawing
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Draw sharp triangle pointing down with wings
	// Main body triangle
	drawFilledTriangle(rgba,
		centerX, centerY+24, // Bottom point
		centerX-16, centerY-18, // Top left
		centerX+16, centerY-18, // Top right
		color.RGBA{220, 80, 60, 255})

	// Wings (sharp angular)
	drawFilledTriangle(rgba,
		centerX-16, centerY-18,
		centerX-26, centerY-8,
		centerX-16, centerY+2,
		color.RGBA{200, 60, 40, 255})

	drawFilledTriangle(rgba,
		centerX+16, centerY-18,
		centerX+26, centerY-8,
		centerX+16, centerY+2,
		color.RGBA{200, 60, 40, 255})

	// Engine glow at back
	drawFilledCircle(rgba, centerX, centerY-18, 4, color.RGBA{255, 150, 100, 255})

	// Cockpit (orange core)
	drawFilledCircle(rgba, centerX, centerY, 6, color.RGBA{255, 150, 100, 255})

	// Bright outline for visibility
	drawCircleOutline(rgba, centerX, centerY, 28, 2, color.RGBA{255, 100, 50, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateDroneSprite creates a diamond-shaped drone with rotating ring effect
func (sm *SpriteManager) generateDroneSprite() *ebiten.Image {
	size := 72
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Draw diamond shape (4 triangles)
	drawFilledTriangle(rgba,
		centerX, centerY-28, // Top
		centerX-28, centerY, // Left
		centerX, centerY, // Center
		color.RGBA{180, 80, 220, 255})

	drawFilledTriangle(rgba,
		centerX, centerY-28, // Top
		centerX+28, centerY, // Right
		centerX, centerY, // Center
		color.RGBA{180, 80, 220, 255})

	drawFilledTriangle(rgba,
		centerX-28, centerY, // Left
		centerX, centerY+28, // Bottom
		centerX, centerY, // Center
		color.RGBA{180, 80, 220, 255})

	drawFilledTriangle(rgba,
		centerX+28, centerY, // Right
		centerX, centerY+28, // Bottom
		centerX, centerY, // Center
		color.RGBA{180, 80, 220, 255})

	// Rotating ring (will be animated in-game)
	drawCircleOutline(rgba, centerX, centerY, 32, 2, color.RGBA{200, 100, 255, 200})

	// Core
	drawFilledCircle(rgba, centerX, centerY, 10, color.RGBA{255, 150, 255, 255})

	// Bright outline
	drawCircleOutline(rgba, centerX, centerY, 30, 3, color.RGBA{200, 100, 255, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateHunterSprite creates an elongated predator-like ship with scanning beam
func (sm *SpriteManager) generateHunterSprite() *ebiten.Image {
	size := 80
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Elongated body (2:1 aspect ratio)
	drawFilledEllipse(rgba, centerX, centerY, 16, 32, color.RGBA{80, 220, 120, 255})

	// Predator fins (swept back)
	drawFilledTriangle(rgba,
		centerX-12, centerY-10,
		centerX-28, centerY-5,
		centerX-12, centerY+10,
		color.RGBA{60, 200, 100, 255})

	drawFilledTriangle(rgba,
		centerX+12, centerY-10,
		centerX+28, centerY-5,
		centerX+12, centerY+10,
		color.RGBA{60, 200, 100, 255})

	// Cockpit (front)
	drawFilledCircle(rgba, centerX, centerY-24, 6, color.RGBA{150, 255, 150, 255})

	// Scanning beam indicator (green line)
	drawLine(rgba, centerX, centerY-32, centerX, centerY+32, 2, color.RGBA{100, 255, 100, 150})

	// Bright green outline
	drawCircleOutline(rgba, centerX, centerY, 36, 3, color.RGBA{100, 255, 100, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateTankSprite creates a heavy octagonal fortress ship
func (sm *SpriteManager) generateTankSprite() *ebiten.Image {
	size := 96
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Draw octagon using 8 triangles
	for i := 0; i < 8; i++ {
		angle1 := float64(i) * math.Pi / 4
		angle2 := float64(i+1) * math.Pi / 4

		x1 := centerX + int(math.Cos(angle1)*32)
		y1 := centerY + int(math.Sin(angle1)*32)
		x2 := centerX + int(math.Cos(angle2)*32)
		y2 := centerY + int(math.Sin(angle2)*32)

		drawFilledTriangle(rgba, centerX, centerY, x1, y1, x2, y2, color.RGBA{140, 140, 140, 255})
	}

	// Armor plates (darker segments)
	for i := 0; i < 8; i += 2 {
		angle := float64(i) * math.Pi / 4
		x := centerX + int(math.Cos(angle)*28)
		y := centerY + int(math.Sin(angle)*28)
		drawFilledRect(rgba, x-6, y-6, 12, 12, color.RGBA{100, 100, 100, 255})
	}

	// Gold core
	drawFilledCircle(rgba, centerX, centerY, 16, color.RGBA{255, 200, 50, 255})

	// Metallic outline
	drawCircleOutline(rgba, centerX, centerY, 40, 4, color.RGBA{255, 150, 50, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateBomberSprite creates a bulky rectangular ship with visible payload
func (sm *SpriteManager) generateBomberSprite() *ebiten.Image {
	size := 80
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Main rectangular body
	drawFilledRect(rgba, centerX-20, centerY-24, 40, 48, color.RGBA{255, 140, 40, 255})

	// Wings
	drawFilledRect(rgba, centerX-32, centerY-8, 12, 24, color.RGBA{230, 120, 30, 255})
	drawFilledRect(rgba, centerX+20, centerY-8, 12, 24, color.RGBA{230, 120, 30, 255})

	// Bomb payload underneath (visible spheres)
	drawFilledCircle(rgba, centerX-8, centerY+20, 6, color.RGBA{60, 60, 60, 255})
	drawFilledCircle(rgba, centerX+8, centerY+20, 6, color.RGBA{60, 60, 60, 255})

	// Pulsing warning light
	drawFilledCircle(rgba, centerX, centerY-16, 5, color.RGBA{255, 255, 100, 255})

	// Orange outline
	drawRectOutline(rgba, centerX-20, centerY-24, 40, 48, 3, color.RGBA{255, 200, 0, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateSniperSprite creates a sleek sniper ship with targeting scope
func (sm *SpriteManager) generateSniperSprite() *ebiten.Image {
	size := 72
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Slim elongated body (sniper rifle profile)
	drawFilledEllipse(rgba, centerX, centerY, 10, 30, color.RGBA{40, 60, 120, 255})

	// Barrel/scope (extended front)
	drawFilledRect(rgba, centerX-4, centerY-35, 8, 20, color.RGBA{60, 80, 140, 255})

	// Side stabilizers (small)
	drawFilledTriangle(rgba,
		centerX-8, centerY,
		centerX-16, centerY-8,
		centerX-8, centerY-16,
		color.RGBA{30, 50, 100, 255})

	drawFilledTriangle(rgba,
		centerX+8, centerY,
		centerX+16, centerY-8,
		centerX+8, centerY-16,
		color.RGBA{30, 50, 100, 255})

	// Targeting scope (cyan laser indicator)
	drawFilledCircle(rgba, centerX, centerY-32, 5, color.RGBA{0, 255, 255, 255})

	// Lock-on indicator ring
	drawCircleOutline(rgba, centerX, centerY, 32, 2, color.RGBA{255, 100, 100, 180})

	// Cyan core (targeting computer)
	drawFilledCircle(rgba, centerX, centerY, 8, color.RGBA{100, 200, 255, 255})

	// Dark blue outline
	drawCircleOutline(rgba, centerX, centerY, 34, 3, color.RGBA{60, 100, 200, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateSplitterSprite creates a ship that visually shows split potential
func (sm *SpriteManager) generateSplitterSprite() *ebiten.Image {
	size := 84
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Main body (two halves with separation line)
	// Left half
	drawFilledTriangle(rgba,
		centerX-2, centerY-28,
		centerX-2, centerY+28,
		centerX-24, centerY,
		color.RGBA{255, 180, 40, 255})

	// Right half
	drawFilledTriangle(rgba,
		centerX+2, centerY-28,
		centerX+2, centerY+28,
		centerX+24, centerY,
		color.RGBA{255, 160, 20, 255})

	// Separation crack (visual indicator of split)
	drawLine(rgba, centerX, centerY-28, centerX, centerY+28, 3, color.RGBA{50, 50, 50, 255})

	// Split indicators (small circles at split points)
	drawFilledCircle(rgba, centerX, centerY-20, 4, color.RGBA{255, 255, 100, 255})
	drawFilledCircle(rgba, centerX, centerY, 5, color.RGBA{255, 255, 100, 255})
	drawFilledCircle(rgba, centerX, centerY+20, 4, color.RGBA{255, 255, 100, 255})

	// Outer wings
	drawFilledTriangle(rgba,
		centerX-24, centerY,
		centerX-32, centerY-10,
		centerX-24, centerY-20,
		color.RGBA{220, 140, 20, 255})

	drawFilledTriangle(rgba,
		centerX+24, centerY,
		centerX+32, centerY-10,
		centerX+24, centerY-20,
		color.RGBA{220, 140, 20, 255})

	// Orange/yellow glow outline
	drawCircleOutline(rgba, centerX, centerY, 36, 3, color.RGBA{255, 200, 50, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// generateShieldBearerSprite creates a heavily armored ship with shield visuals
func (sm *SpriteManager) generateShieldBearerSprite() *ebiten.Image {
	size := 100
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Large hexagonal shield shape
	for i := 0; i < 6; i++ {
		angle1 := float64(i) * math.Pi / 3
		angle2 := float64(i+1) * math.Pi / 3

		x1 := centerX + int(math.Cos(angle1)*35)
		y1 := centerY + int(math.Sin(angle1)*35)
		x2 := centerX + int(math.Cos(angle2)*35)
		y2 := centerY + int(math.Sin(angle2)*35)

		drawFilledTriangle(rgba, centerX, centerY, x1, y1, x2, y2, color.RGBA{180, 180, 180, 255})
	}

	// Shield armor plates (darker segments)
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3
		x := centerX + int(math.Cos(angle)*30)
		y := centerY + int(math.Sin(angle)*30)
		drawFilledRect(rgba, x-8, y-8, 16, 16, color.RGBA{140, 140, 140, 255})
	}

	// Shield barrier (blue energy field visual)
	drawCircleOutline(rgba, centerX, centerY, 42, 4, color.RGBA{100, 150, 255, 200})
	drawCircleOutline(rgba, centerX, centerY, 45, 2, color.RGBA{150, 200, 255, 150})

	// Armored core
	drawFilledCircle(rgba, centerX, centerY, 18, color.RGBA{160, 160, 160, 255})

	// Central reactor (glowing)
	drawFilledCircle(rgba, centerX, centerY, 12, color.RGBA{100, 150, 255, 255})

	// Heavy metallic outline
	drawCircleOutline(rgba, centerX, centerY, 38, 4, color.RGBA{200, 200, 200, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// ============================================================================
// ASTEROID SPRITE GENERATORS
// ============================================================================

func (sm *SpriteManager) generateAsteroidSmallSprite() *ebiten.Image {
	return sm.generateAsteroidSprite(32, 10, color.RGBA{150, 120, 100, 255})
}

func (sm *SpriteManager) generateAsteroidMediumSprite() *ebiten.Image {
	return sm.generateAsteroidSprite(56, 20, color.RGBA{130, 100, 80, 255})
}

func (sm *SpriteManager) generateAsteroidLargeSprite() *ebiten.Image {
	return sm.generateAsteroidSprite(88, 35, color.RGBA{110, 80, 60, 255})
}

// generateAsteroidSprite creates an irregular, lumpy asteroid
func (sm *SpriteManager) generateAsteroidSprite(size int, radius int, baseColor color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Create irregular shape with multiple overlapping circles
	numCircles := 6 + radius/10
	for i := 0; i < numCircles; i++ {
		angle := float64(i) * math.Pi * 2 / float64(numCircles)
		offsetX := int(math.Cos(angle) * float64(radius) * 0.3)
		offsetY := int(math.Sin(angle) * float64(radius) * 0.3)
		r := radius*7/10 + (i%3)*radius/10

		drawFilledCircle(rgba, centerX+offsetX, centerY+offsetY, r, baseColor)
	}

	// Add crater texture
	craterCount := radius / 5
	for i := 0; i < craterCount; i++ {
		angle := float64(i) * math.Pi * 2 / float64(craterCount)
		craterX := centerX + int(math.Cos(angle)*float64(radius)*0.5)
		craterY := centerY + int(math.Sin(angle)*float64(radius)*0.5)
		craterSize := radius / 4

		craterColor := color.RGBA{baseColor.R / 2, baseColor.G / 2, baseColor.B / 2, 255}
		drawFilledCircle(rgba, craterX, craterY, craterSize, craterColor)
	}

	// Highlight for lighting
	drawFilledCircle(rgba, centerX-radius/3, centerY-radius/3, radius/4, color.RGBA{180, 150, 130, 200})

	// Dotted outline (non-hostile indicator)
	drawDottedCircle(rgba, centerX, centerY, radius+2, 2, color.RGBA{100, 100, 100, 200})

	img.WritePixels(rgba.Pix)
	return img
}

// ============================================================================
// POWER-UP SPRITE GENERATORS
// ============================================================================

func (sm *SpriteManager) generateHealthPowerUpSprite() *ebiten.Image {
	return sm.generatePowerUpSprite(60, color.RGBA{50, 255, 50, 255}, "health")
}

func (sm *SpriteManager) generateShieldPowerUpSprite() *ebiten.Image {
	return sm.generatePowerUpSprite(60, color.RGBA{50, 150, 255, 255}, "shield")
}

func (sm *SpriteManager) generateWeaponPowerUpSprite() *ebiten.Image {
	return sm.generatePowerUpSprite(60, color.RGBA{255, 200, 50, 255}, "weapon")
}

func (sm *SpriteManager) generateSpeedPowerUpSprite() *ebiten.Image {
	return sm.generatePowerUpSprite(60, color.RGBA{255, 100, 255, 255}, "speed")
}

func (sm *SpriteManager) generateMysteryPowerUpSprite() *ebiten.Image {
	return sm.generatePowerUpSprite(60, color.RGBA{180, 50, 255, 255}, "mystery")
}

// generatePowerUpSprite creates a power-up with enhanced visibility
func (sm *SpriteManager) generatePowerUpSprite(size int, mainColor color.RGBA, powerType string) *ebiten.Image {
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2
	radius := size / 3

	// Outer glow
	glowColor := color.RGBA{mainColor.R, mainColor.G, mainColor.B, 100}
	drawFilledCircle(rgba, centerX, centerY, radius+8, glowColor)

	// Main circle
	drawFilledCircle(rgba, centerX, centerY, radius, mainColor)

	// Animated rotating border (dashed circle)
	drawDottedCircle(rgba, centerX, centerY, radius+6, 3, mainColor)

	// Inner highlight
	drawFilledCircle(rgba, centerX-4, centerY-4, radius/3, color.RGBA{255, 255, 255, 220})

	// Type-specific icon (larger and more prominent)
	iconColor := color.RGBA{255, 255, 255, 255}
	switch powerType {
	case "health":
		// Large plus sign
		drawFilledRect(rgba, centerX-10, centerY-3, 20, 6, iconColor)
		drawFilledRect(rgba, centerX-3, centerY-10, 6, 20, iconColor)
	case "shield":
		// Hexagonal border
		for i := 0; i < 6; i++ {
			angle1 := float64(i) * math.Pi / 3
			angle2 := float64(i+1) * math.Pi / 3
			x1 := centerX + int(math.Cos(angle1)*float64(radius/2))
			y1 := centerY + int(math.Sin(angle1)*float64(radius/2))
			x2 := centerX + int(math.Cos(angle2)*float64(radius/2))
			y2 := centerY + int(math.Sin(angle2)*float64(radius/2))
			drawLine(rgba, x1, y1, x2, y2, 3, iconColor)
		}
	case "weapon":
		// Large up arrow with glow
		drawFilledTriangle(rgba, centerX, centerY-10, centerX-8, centerY+4, centerX+8, centerY+4, iconColor)
		drawFilledRect(rgba, centerX-3, centerY+4, 6, 6, iconColor)
	case "speed":
		// Lightning bolt
		points := []struct{ x, y int }{
			{centerX - 4, centerY - 10},
			{centerX + 4, centerY - 2},
			{centerX, centerY - 2},
			{centerX + 4, centerY + 10},
			{centerX - 4, centerY + 2},
			{centerX, centerY + 2},
		}
		for i := 0; i < len(points)-1; i++ {
			drawLine(rgba, points[i].x, points[i].y, points[i+1].x, points[i+1].y, 3, iconColor)
		}
	case "mystery":
		// Question mark
		// Top arc
		for i := 0; i < 8; i++ {
			angle := math.Pi - float64(i)*math.Pi/7
			x1 := centerX + int(math.Cos(angle)*6)
			y1 := centerY - 8 + int(math.Sin(angle)*6)
			x2 := centerX + int(math.Cos(angle-math.Pi/7)*6)
			y2 := centerY - 8 + int(math.Sin(angle-math.Pi/7)*6)
			drawLine(rgba, x1, y1, x2, y2, 3, iconColor)
		}
		// Stem
		drawFilledRect(rgba, centerX-2, centerY-2, 4, 6, iconColor)
		// Dot
		drawFilledCircle(rgba, centerX, centerY+6, 2, iconColor)
	}

	// Strong outer outline for maximum visibility
	drawCircleOutline(rgba, centerX, centerY, radius+2, 3, mainColor)

	img.WritePixels(rgba.Pix)
	return img
}

// ============================================================================
// PROJECTILE SPRITE GENERATORS
// ============================================================================

func (sm *SpriteManager) generatePlayerProjectileSprite() *ebiten.Image {
	size := 24
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Large outer glow
	drawFilledCircle(rgba, centerX, centerY, 10, color.RGBA{50, 150, 255, 100})

	// Mid glow
	drawFilledCircle(rgba, centerX, centerY, 7, color.RGBA{100, 200, 255, 200})

	// Main projectile (larger than before)
	drawFilledCircle(rgba, centerX, centerY, 5, color.RGBA{100, 200, 255, 255})

	// Bright core
	drawFilledCircle(rgba, centerX, centerY, 3, color.RGBA{255, 255, 255, 255})

	img.WritePixels(rgba.Pix)
	return img
}

func (sm *SpriteManager) generateEnemyProjectileSprite() *ebiten.Image {
	size := 28 // Larger than player projectiles
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2

	// Large outer glow (danger color)
	drawFilledCircle(rgba, centerX, centerY, 12, color.RGBA{255, 50, 50, 100})

	// Mid glow
	drawFilledCircle(rgba, centerX, centerY, 9, color.RGBA{255, 100, 100, 200})

	// Main projectile (larger for visibility)
	drawFilledCircle(rgba, centerX, centerY, 6, color.RGBA{255, 100, 100, 255})

	// Yellow/orange core (danger indicator)
	drawFilledCircle(rgba, centerX, centerY, 4, color.RGBA{255, 200, 50, 255})

	// Bright center
	drawFilledCircle(rgba, centerX, centerY, 2, color.RGBA{255, 255, 255, 255})

	img.WritePixels(rgba.Pix)
	return img
}

// ============================================================================
// EFFECT SPRITE GENERATORS
// ============================================================================

func (sm *SpriteManager) generateExplosionFrame(frame int) *ebiten.Image {
	size := 128
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2
	progress := float64(frame) / 8.0

	// Explosion grows then fades
	radius := int(progress * 50)
	alpha := uint8(255 * (1 - progress))

	// Fire colors (red -> orange -> yellow)
	numParticles := 20 - frame*2
	for i := 0; i < numParticles; i++ {
		angle := float64(i) * math.Pi * 2 / float64(numParticles)
		dist := float64(radius) * (0.5 + progress*0.5)
		x := centerX + int(math.Cos(angle)*dist)
		y := centerY + int(math.Sin(angle)*dist)

		particleSize := 8 - frame

		// Color gradient
		var c color.RGBA
		if progress < 0.3 {
			c = color.RGBA{255, 255, 255, alpha} // White hot
		} else if progress < 0.6 {
			c = color.RGBA{255, 200, 50, alpha} // Orange
		} else {
			c = color.RGBA{255, 100, 50, alpha} // Red
		}

		drawFilledCircle(rgba, x, y, particleSize, c)
	}

	img.WritePixels(rgba.Pix)
	return img
}

func (sm *SpriteManager) generateSparkleFrame(frame int) *ebiten.Image {
	size := 32
	img := ebiten.NewImage(size, size)
	rgba := image.NewRGBA(image.Rect(0, 0, size, size))

	centerX, centerY := size/2, size/2
	progress := float64(frame) / 6.0

	// Rotating sparkle
	angle := progress * math.Pi * 2
	for i := 0; i < 4; i++ {
		a := angle + float64(i)*math.Pi/2
		x := centerX + int(math.Cos(a)*8)
		y := centerY + int(math.Sin(a)*8)

		drawFilledCircle(rgba, x, y, 2, color.RGBA{255, 255, 255, 255})
	}

	img.WritePixels(rgba.Pix)
	return img
}
