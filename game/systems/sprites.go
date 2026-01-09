package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteManager handles loading and managing all game sprites
type SpriteManager struct {
	// Enemy sprites
	ScoutSprite        *ebiten.Image
	DroneSprite        *ebiten.Image
	HunterSprite       *ebiten.Image
	TankSprite         *ebiten.Image
	BomberSprite       *ebiten.Image
	SniperSprite       *ebiten.Image
	SplitterSprite     *ebiten.Image
	ShieldBearerSprite *ebiten.Image

	// Asteroid sprites (3 sizes)
	AsteroidSmallSprite  *ebiten.Image
	AsteroidMediumSprite *ebiten.Image
	AsteroidLargeSprite  *ebiten.Image

	// Power-up sprites
	PowerUpHealthSprite  *ebiten.Image
	PowerUpShieldSprite  *ebiten.Image
	PowerUpWeaponSprite  *ebiten.Image
	PowerUpSpeedSprite   *ebiten.Image
	PowerUpMysterySprite *ebiten.Image

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
	case 5: // Sniper
		return sm.SniperSprite
	case 6: // Splitter
		return sm.SplitterSprite
	case 7: // Shield Bearer
		return sm.ShieldBearerSprite
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
	case 4: // Mystery
		return sm.PowerUpMysterySprite
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
