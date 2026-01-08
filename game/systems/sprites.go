package systems

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteManager handles loading and managing all game sprites
type SpriteManager struct {
	// Enemy sprites
	ScoutSprite  *ebiten.Image
	DroneSprite  *ebiten.Image
	HunterSprite *ebiten.Image
	TankSprite   *ebiten.Image
	BomberSprite *ebiten.Image

	// Asteroid sprites (3 sizes)
	AsteroidSmallSprite  *ebiten.Image
	AsteroidMediumSprite *ebiten.Image
	AsteroidLargeSprite  *ebiten.Image

	// Power-up sprites
	PowerUpHealthSprite *ebiten.Image
	PowerUpShieldSprite *ebiten.Image
	PowerUpWeaponSprite *ebiten.Image
	PowerUpSpeedSprite  *ebiten.Image

	// Projectile sprites
	PlayerProjectileSprite *ebiten.Image
	EnemyProjectileSprite  *ebiten.Image

	// Effect sprites
	ExplosionFrames []*ebiten.Image
	SparkleFrames   []*ebiten.Image
}

// NewSpriteManager creates and initializes the sprite manager
func NewSpriteManager() *SpriteManager {
	sm := &SpriteManager{}

	// Try to load sprites from files, fallback to procedural generation
	sm.loadOrGenerateSprites()

	return sm
}

// loadOrGenerateSprites attempts to load sprites from disk, generates them if files don't exist
func (sm *SpriteManager) loadOrGenerateSprites() {
	assetsDir := "assets/sprites"

	// Enemy sprites
	sm.ScoutSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_scout.png"), sm.generateScoutSprite)
	sm.DroneSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_drone.png"), sm.generateDroneSprite)
	sm.HunterSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_hunter.png"), sm.generateHunterSprite)
	sm.TankSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_tank.png"), sm.generateTankSprite)
	sm.BomberSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_bomber.png"), sm.generateBomberSprite)

	// Asteroid sprites
	sm.AsteroidSmallSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_small.png"), sm.generateAsteroidSmallSprite)
	sm.AsteroidMediumSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_medium.png"), sm.generateAsteroidMediumSprite)
	sm.AsteroidLargeSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_large.png"), sm.generateAsteroidLargeSprite)

	// Power-up sprites
	sm.PowerUpHealthSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_health.png"), sm.generateHealthPowerUpSprite)
	sm.PowerUpShieldSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_shield.png"), sm.generateShieldPowerUpSprite)
	sm.PowerUpWeaponSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_weapon.png"), sm.generateWeaponPowerUpSprite)
	sm.PowerUpSpeedSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_speed.png"), sm.generateSpeedPowerUpSprite)

	// Projectile sprites
	sm.PlayerProjectileSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "projectile_player.png"), sm.generatePlayerProjectileSprite)
	sm.EnemyProjectileSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "projectile_enemy.png"), sm.generateEnemyProjectileSprite)

	// Generate explosion animation frames
	sm.ExplosionFrames = make([]*ebiten.Image, 8)
	for i := 0; i < 8; i++ {
		sm.ExplosionFrames[i] = sm.generateExplosionFrame(i)
	}

	// Generate sparkle animation frames for power-ups
	sm.SparkleFrames = make([]*ebiten.Image, 6)
	for i := 0; i < 6; i++ {
		sm.SparkleFrames[i] = sm.generateSparkleFrame(i)
	}
}

// loadOrGenerate attempts to load a sprite from file, generates it if not found
func (sm *SpriteManager) loadOrGenerate(path string, generator func() *ebiten.Image) *ebiten.Image {
	// Try to load from file
	if img := sm.loadSprite(path); img != nil {
		return img
	}

	// Generate sprite
	img := generator()

	// Note: We can't save sprites during initialization because Ebiten's
	// ReadPixels can't be called before the game starts. Sprites will be
	// generated fresh each time. To save sprites, run the game once and
	// use a separate tool to export them.

	return img
}

// loadSprite loads a sprite from a PNG file
func (sm *SpriteManager) loadSprite(path string) *ebiten.Image {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil
	}

	return ebiten.NewImageFromImage(img)
}

// saveSprite saves a sprite to a PNG file
// Note: This cannot be called during initialization because Ebiten's ReadPixels
// can't be called before the game starts. This is kept for future use if we
// want to export sprites at runtime.
func (sm *SpriteManager) saveSprite(path string, sprite *ebiten.Image) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Convert ebiten.Image to image.Image
	bounds := sprite.Bounds()
	img := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, sprite.At(x, y))
		}
	}

	// Encode to PNG
	png.Encode(file, img)
}

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

// ============================================================================
// DRAWING HELPER FUNCTIONS
// ============================================================================

func drawFilledCircle(img *image.RGBA, cx, cy, radius int, c color.RGBA) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= radius*radius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawFilledEllipse(img *image.RGBA, cx, cy, rx, ry int, c color.RGBA) {
	for y := cy - ry; y <= cy+ry; y++ {
		for x := cx - rx; x <= cx+rx; x++ {
			dx := float64(x - cx)
			dy := float64(y - cy)
			if (dx*dx)/(float64(rx*rx))+(dy*dy)/(float64(ry*ry)) <= 1.0 {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawCircleOutline(img *image.RGBA, cx, cy, radius, thickness int, c color.RGBA) {
	for y := cy - radius - thickness; y <= cy+radius+thickness; y++ {
		for x := cx - radius - thickness; x <= cx+radius+thickness; x++ {
			dx := x - cx
			dy := y - cy
			dist := dx*dx + dy*dy
			innerRadius := (radius - thickness) * (radius - thickness)
			outerRadius := (radius + thickness) * (radius + thickness)
			if dist >= innerRadius && dist <= outerRadius {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func drawDottedCircle(img *image.RGBA, cx, cy, radius, dotSize int, c color.RGBA) {
	numDots := 24
	for i := 0; i < numDots; i++ {
		angle := float64(i) * math.Pi * 2 / float64(numDots)
		x := cx + int(math.Cos(angle)*float64(radius))
		y := cy + int(math.Sin(angle)*float64(radius))
		drawFilledCircle(img, x, y, dotSize, c)
	}
}

func drawFilledRect(img *image.RGBA, x, y, w, h int, c color.RGBA) {
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
				img.Set(px, py, c)
			}
		}
	}
}

func drawRectOutline(img *image.RGBA, x, y, w, h, thickness int, c color.RGBA) {
	// Top
	drawFilledRect(img, x, y, w, thickness, c)
	// Bottom
	drawFilledRect(img, x, y+h-thickness, w, thickness, c)
	// Left
	drawFilledRect(img, x, y, thickness, h, c)
	// Right
	drawFilledRect(img, x+w-thickness, y, thickness, h, c)
}

func drawLine(img *image.RGBA, x1, y1, x2, y2, thickness int, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Sqrt(float64(dx*dx + dy*dy)))

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + int(float64(dx)*t)
		y := y1 + int(float64(dy)*t)

		for tx := -thickness / 2; tx <= thickness/2; tx++ {
			for ty := -thickness / 2; ty <= thickness/2; ty++ {
				px := x + tx
				py := y + ty
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, c)
				}
			}
		}
	}
}

func drawFilledTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, c color.RGBA) {
	// Find bounding box
	minX := min(x1, min(x2, x3))
	maxX := max(x1, max(x2, x3))
	minY := min(y1, min(y2, y3))
	maxY := max(y1, max(y2, y3))

	// Check each pixel in bounding box
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(x, y, x1, y1, x2, y2, x3, y3) {
				if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
					img.Set(x, y, c)
				}
			}
		}
	}
}

func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 int) bool {
	// Barycentric coordinate method
	denominator := float64((y2-y3)*(x1-x3) + (x3-x2)*(y1-y3))
	if denominator == 0 {
		return false
	}

	a := float64((y2-y3)*(px-x3)+(x3-x2)*(py-y3)) / denominator
	b := float64((y3-y1)*(px-x3)+(x1-x3)*(py-y3)) / denominator
	c := 1 - a - b

	return a >= 0 && a <= 1 && b >= 0 && b <= 1 && c >= 0 && c <= 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetSpriteForEnemy returns the appropriate sprite for an enemy type
func (sm *SpriteManager) GetSpriteForEnemy(enemyType int) *ebiten.Image {
	switch enemyType {
	case 0: // Scout
		return sm.ScoutSprite
	case 1: // Drone
		return sm.DroneSprite
	case 2: // Hunter
		return sm.HunterSprite
	case 3: // Tank
		return sm.TankSprite
	case 4: // Bomber
		return sm.BomberSprite
	default:
		return sm.ScoutSprite
	}
}

// GetSpriteForAsteroid returns the appropriate sprite for an asteroid size
func (sm *SpriteManager) GetSpriteForAsteroid(size int) *ebiten.Image {
	switch size {
	case 0: // Small
		return sm.AsteroidSmallSprite
	case 1: // Medium
		return sm.AsteroidMediumSprite
	case 2: // Large
		return sm.AsteroidLargeSprite
	default:
		return sm.AsteroidSmallSprite
	}
}

// GetSpriteForPowerUp returns the appropriate sprite for a power-up type
func (sm *SpriteManager) GetSpriteForPowerUp(powerUpType int) *ebiten.Image {
	switch powerUpType {
	case 0: // Health
		return sm.PowerUpHealthSprite
	case 1: // Shield
		return sm.PowerUpShieldSprite
	case 2: // Weapon
		return sm.PowerUpWeaponSprite
	case 3: // Speed
		return sm.PowerUpSpeedSprite
	default:
		return sm.PowerUpHealthSprite
	}
}

// Debug function to print sprite generation status
func (sm *SpriteManager) PrintStatus() {
	fmt.Println("Sprite Manager initialized successfully:")
	fmt.Println("  - Enemy sprites: 5")
	fmt.Println("  - Asteroid sprites: 3")
	fmt.Println("  - Power-up sprites: 4")
	fmt.Println("  - Projectile sprites: 2")
	fmt.Println("  - Explosion frames: 8")
	fmt.Println("  - Sparkle frames: 6")
}
