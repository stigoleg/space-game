package systems

import (
	"math"
	"math/rand"
)

// CameraSystem manages camera zoom, shake, and cinematic effects
type CameraSystem struct {
	// Zoom
	Zoom       float64 // Current zoom level (1.0 = normal)
	TargetZoom float64 // Target zoom for smooth transitions

	// Screen shake
	ShakeAmount float64 // Screen shake intensity from camera effects
	ScreenShake float64 // External screen shake (from impacts, explosions)

	// Cinematic mode
	CinematicMode  bool    // Cinematic mode active (for boss, etc)
	CinematicTimer float64 // Time in cinematic mode

	// Screen dimensions (for centering calculations)
	ScreenWidth  int
	ScreenHeight int
}

// NewCameraSystem creates a new camera system
func NewCameraSystem(screenWidth, screenHeight int) *CameraSystem {
	return &CameraSystem{
		Zoom:         1.0,
		TargetZoom:   1.0,
		ShakeAmount:  0.0,
		ScreenShake:  0.0,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
}

// Update updates camera state (zoom transitions, shake decay, cinematic mode)
// Parameters:
//   - bossWave: whether a boss wave is active
//   - bossActive: whether the boss entity is active
//   - enemyCount: number of active enemies
//   - playerHealthRatio: player health / max health (0.0 to 1.0)
//   - waveCompleted: whether the current wave is completed
func (c *CameraSystem) Update(bossWave bool, bossActive bool, enemyCount int, playerHealthRatio float64, waveCompleted bool) {
	// Smooth zoom transitions
	zoomDifference := c.TargetZoom - c.Zoom
	if math.Abs(zoomDifference) > 0.01 {
		// Smooth interpolation towards target zoom
		c.Zoom += zoomDifference * 0.1 // Smooth lerp
	} else {
		c.Zoom = c.TargetZoom
	}

	// Handle boss cinematic mode
	if bossWave && bossActive && !c.CinematicMode {
		// Trigger cinematic zoom on boss appearance
		c.CinematicMode = true
		c.CinematicTimer = 0
		c.TargetZoom = 0.75 // Zoom in more for dramatic effect
		c.ShakeAmount = 8.0 // Stronger initial shake
	}

	// Update cinematic mode timer
	if c.CinematicMode {
		c.CinematicTimer += 1.0 / 60.0
		if c.CinematicTimer > 2.5 {
			// Exit cinematic mode after 2.5 seconds
			c.CinematicMode = false
			c.TargetZoom = 1.0 // Return to normal zoom
			c.ShakeAmount = 0.0
		}
	}

	// Dynamic zoom based on player danger level
	// More danger = zoom out to see more
	if !c.CinematicMode && !bossWave {
		// Zoom out when many enemies present
		if enemyCount > 15 {
			c.TargetZoom = 1.15 // More zoom out
		} else if enemyCount > 10 {
			c.TargetZoom = 1.1
		} else if enemyCount > 5 {
			c.TargetZoom = 1.05
		} else {
			c.TargetZoom = 1.0 // Normal
		}

		// Extra zoom out if player health low
		if playerHealthRatio < 0.25 {
			c.TargetZoom += 0.05
		}
	}

	// Camera zoom effects on wave completion
	if waveCompleted && enemyCount == 0 {
		// Slight zoom out on wave completion for celebration
		c.TargetZoom = 1.1
	} else if !c.CinematicMode && !bossWave {
		// Dynamic zoom already handled above
	}

	// Decay screen shake more gradually for better feel
	if c.ShakeAmount > 0.1 {
		c.ShakeAmount *= 0.92 // Slightly slower decay for better feel
	} else {
		c.ShakeAmount = 0
	}

	// Add environmental shake (impacts, explosions)
	if c.ScreenShake > 0 {
		c.ShakeAmount += c.ScreenShake * 0.7 // More impact shake
	}
}

// ApplyZoom scales coordinates around screen center based on camera zoom
func (c *CameraSystem) ApplyZoom(x, y float64) (float64, float64) {
	if c.Zoom == 1.0 {
		return x, y
	}

	// Center point of screen
	centerX := float64(c.ScreenWidth) / 2
	centerY := float64(c.ScreenHeight) / 2

	// Translate to center, apply zoom, translate back
	scaledX := centerX + (x-centerX)*c.Zoom
	scaledY := centerY + (y-centerY)*c.Zoom

	return scaledX, scaledY
}

// GetShakeOffset calculates current screen shake offset (includes both camera shake and external shake)
func (c *CameraSystem) GetShakeOffset() (float64, float64) {
	totalShake := c.ScreenShake + c.ShakeAmount
	if totalShake <= 0 {
		return 0, 0
	}

	shakeX := (rand.Float64() - 0.5) * totalShake * 2
	shakeY := (rand.Float64() - 0.5) * totalShake * 2
	return shakeX, shakeY
}

// AddShake adds screen shake from external sources (impacts, explosions)
func (c *CameraSystem) AddShake(amount float64) {
	c.ScreenShake += amount
}

// SetScreenShake sets the external screen shake amount directly
func (c *CameraSystem) SetScreenShake(amount float64) {
	c.ScreenShake = amount
}

// TriggerCinematicZoom triggers a cinematic zoom effect (e.g., for boss intro)
func (c *CameraSystem) TriggerCinematicZoom(zoomLevel float64, shakeAmount float64) {
	c.CinematicMode = true
	c.CinematicTimer = 0
	c.TargetZoom = zoomLevel
	c.ShakeAmount = shakeAmount
}

// Reset resets camera to default state
func (c *CameraSystem) Reset() {
	c.Zoom = 1.0
	c.TargetZoom = 1.0
	c.ShakeAmount = 0.0
	c.ScreenShake = 0.0
	c.CinematicMode = false
	c.CinematicTimer = 0
}
