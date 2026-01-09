package systems

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// loadOrGenerateSprites attempts to load sprites from disk, generates them if files don't exist
func (sm *SpriteManager) loadOrGenerateSprites() {
	assetsDir := "assets/sprites"

	// Enemy sprites
	sm.ScoutSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_scout.png"), sm.generateScoutSprite)
	sm.DroneSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_drone.png"), sm.generateDroneSprite)
	sm.HunterSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_hunter.png"), sm.generateHunterSprite)
	sm.TankSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_tank.png"), sm.generateTankSprite)
	sm.BomberSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_bomber.png"), sm.generateBomberSprite)
	sm.SniperSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_sniper.png"), sm.generateSniperSprite)
	sm.SplitterSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_splitter.png"), sm.generateSplitterSprite)
	sm.ShieldBearerSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "enemy_shield_bearer.png"), sm.generateShieldBearerSprite)

	// Asteroid sprites
	sm.AsteroidSmallSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_small.png"), sm.generateAsteroidSmallSprite)
	sm.AsteroidMediumSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_medium.png"), sm.generateAsteroidMediumSprite)
	sm.AsteroidLargeSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "asteroid_large.png"), sm.generateAsteroidLargeSprite)

	// Power-up sprites
	sm.PowerUpHealthSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_health.png"), sm.generateHealthPowerUpSprite)
	sm.PowerUpShieldSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_shield.png"), sm.generateShieldPowerUpSprite)
	sm.PowerUpWeaponSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_weapon.png"), sm.generateWeaponPowerUpSprite)
	sm.PowerUpSpeedSprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_speed.png"), sm.generateSpeedPowerUpSprite)
	sm.PowerUpMysterySprite = sm.loadOrGenerate(filepath.Join(assetsDir, "powerup_mystery.png"), sm.generateMysteryPowerUpSprite)

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
