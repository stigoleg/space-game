package core

import (
	"image/color"
	"math/rand"

	"stellar-siege/game/entities"
	"stellar-siege/game/systems"
)

// CollisionManager handles all collision detection and resolution in the game
type CollisionManager struct {
	// Callback functions to notify game of collision events
	OnEnemyKilled       func(enemy *entities.Enemy, points int64)
	OnPlayerDamaged     func(damage int)
	OnBossDamaged       func(damage int) bool
	OnScoreAdded        func(points int64)
	OnExplosionSpawned  func(x, y, size float64)
	OnImpactSpawned     func(x, y, size float64, color color.RGBA)
	OnFloatingTextAdded func(x, y float64, text string, color color.RGBA)
	OnSoundPlayed       func(soundType systems.SoundType)
	OnScreenShake       func(amount float64)
	OnPowerUpSpawned    func(x, y float64)
	OnChainLightning    func(proj *entities.Projectile, target *entities.Enemy)

	// Spatial grid for optimization
	spatialGrid    *SpatialGrid
	useSpatialGrid bool
}

// NewCollisionManager creates a new collision manager
func NewCollisionManager(screenWidth, screenHeight float64) *CollisionManager {
	// Cell size of 100 pixels provides good balance between grid overhead and collision reduction
	// Smaller cells = more overhead but fewer checks per cell
	// Larger cells = less overhead but more checks per cell
	return &CollisionManager{
		spatialGrid:    NewSpatialGrid(screenWidth, screenHeight, 100.0),
		useSpatialGrid: true, // Enable spatial grid by default
	}
}

// CheckAllCollisions performs all collision detection for the game
func (cm *CollisionManager) CheckAllCollisions(
	player *entities.Player,
	enemies []*entities.Enemy,
	boss *entities.Boss,
	projectiles []*entities.Projectile,
	powerups []*entities.PowerUp,
	asteroids []*entities.Asteroid,
	gameTime float64,
	dashInvincibility float64,
	damageMultiplier float64,
	powerupSpawnRate float64,
) ([]*entities.Enemy, []*entities.PowerUp) {

	// Populate spatial grid for optimized collision detection
	if cm.useSpatialGrid {
		cm.spatialGrid.PopulateGrid(enemies, projectiles, powerups, asteroids)
	}

	// Check player projectiles vs enemies
	enemies = cm.handleProjectileEnemyCollisions(enemies, projectiles, powerupSpawnRate)

	// Check player projectiles vs boss
	if boss != nil && boss.Active && !boss.IsDead() {
		cm.handleProjectileBossCollisions(boss, projectiles)
	}

	// Check enemy projectiles vs player
	if player != nil && player.Active {
		cm.handleEnemyProjectilePlayerCollisions(player, projectiles, gameTime, dashInvincibility)

		// Check enemy collision with player
		cm.handleEnemyPlayerCollisions(player, enemies, gameTime, dashInvincibility, damageMultiplier)

		// Check boss collision with player
		if boss != nil && boss.Active {
			cm.handleBossPlayerCollision(player, boss, gameTime, damageMultiplier)
		}

		// Check powerup collection
		powerups = cm.handlePowerUpCollisions(player, powerups)

		// Check asteroid collisions with player
		cm.handleAsteroidPlayerCollisions(player, asteroids, gameTime, damageMultiplier)
	}

	// Check player projectiles vs asteroids
	cm.handleProjectileAsteroidCollisions(projectiles, asteroids)

	return enemies, powerups
}

// handleProjectileEnemyCollisions checks player projectiles against enemies
func (cm *CollisionManager) handleProjectileEnemyCollisions(
	enemies []*entities.Enemy,
	projectiles []*entities.Projectile,
	powerupSpawnRate float64,
) []*entities.Enemy {
	if cm.useSpatialGrid {
		// Spatial grid optimization: only check enemies near each projectile
		for _, p := range projectiles {
			if !p.Active || !p.Friendly {
				continue
			}

			// Get only nearby enemies instead of checking all enemies
			nearbyEnemies := cm.spatialGrid.GetNearbyEnemies(p.X, p.Y, p.Radius+50) // +50 for safety margin
			for _, e := range nearbyEnemies {
				if !e.Active {
					continue
				}

				if checkCircleCollision(p.X, p.Y, p.Radius, e.X, e.Y, e.Radius) {
					cm.handleProjectileHitEnemy(p, e, powerupSpawnRate)
				}
			}
		}
	} else {
		// Fallback: check all projectiles against all enemies (O(n²))
		for _, p := range projectiles {
			if !p.Active || !p.Friendly {
				continue
			}

			for _, e := range enemies {
				if !e.Active {
					continue
				}

				if checkCircleCollision(p.X, p.Y, p.Radius, e.X, e.Y, e.Radius) {
					cm.handleProjectileHitEnemy(p, e, powerupSpawnRate)
				}
			}
		}
	}

	return enemies
}

// handleProjectileHitEnemy processes what happens when a projectile hits an enemy
func (cm *CollisionManager) handleProjectileHitEnemy(
	proj *entities.Projectile,
	enemy *entities.Enemy,
	powerupSpawnRate float64,
) {
	// Don't deactivate projectile if it's piercing (beam)
	if !proj.Piercing {
		proj.Active = false
	}

	// Use TakeDamage method to handle shields properly
	enemy.TakeDamage(proj.Damage)

	// Apply burning DoT if projectile has burning flag
	if proj.Burning {
		enemy.ApplyBurn(proj.BurnDuration, proj.BurnDamage)
	}

	// Add impact effect
	if cm.OnImpactSpawned != nil {
		cm.OnImpactSpawned(enemy.X, enemy.Y, 30, color.RGBA{100, 200, 255, 255})
	}

	// Handle chain lightning
	if proj.Chaining && proj.ChainCount > 0 && cm.OnChainLightning != nil {
		cm.OnChainLightning(proj, enemy)
	}

	// Check if enemy was killed
	if enemy.Health <= 0 {
		cm.handleEnemyKilled(enemy, powerupSpawnRate)
	}
}

// handleEnemyKilled processes enemy death effects
func (cm *CollisionManager) handleEnemyKilled(enemy *entities.Enemy, powerupSpawnRate float64) {
	enemy.Active = false

	// Spawn explosion
	if cm.OnExplosionSpawned != nil {
		cm.OnExplosionSpawned(enemy.X, enemy.Y, enemy.Radius)
	}

	// Handle Splitter splitting into 2 scouts
	if enemy.Type == entities.EnemySplitter {
		splitEnemies := enemy.GetSplitEnemies()
		if splitEnemies != nil && cm.OnFloatingTextAdded != nil {
			cm.OnFloatingTextAdded(enemy.X, enemy.Y, "SPLIT!", color.RGBA{255, 230, 100, 255})
		}
	}

	// Play appropriate explosion sound based on enemy type
	if cm.OnSoundPlayed != nil {
		soundType := getEnemyExplosionSound(enemy.Type)
		cm.OnSoundPlayed(soundType)
	}

	// Add score
	points := int64(enemy.Points)
	if cm.OnScoreAdded != nil {
		cm.OnScoreAdded(points)
	}

	// Screen shake
	if cm.OnScreenShake != nil {
		cm.OnScreenShake(5)
	}

	// Notify game of kill (for achievements, scrap, etc.)
	if cm.OnEnemyKilled != nil {
		cm.OnEnemyKilled(enemy, points)
	}

	// Chance to spawn powerup (modified by challenge config)
	powerupChance := 0.15 * powerupSpawnRate
	if rand.Float64() < powerupChance && cm.OnPowerUpSpawned != nil {
		cm.OnPowerUpSpawned(enemy.X, enemy.Y)
	}
}

// handleProjectileBossCollisions checks player projectiles against boss
func (cm *CollisionManager) handleProjectileBossCollisions(
	boss *entities.Boss,
	projectiles []*entities.Projectile,
) {
	for _, p := range projectiles {
		if !p.Active || !p.Friendly {
			continue
		}

		if checkCircleCollision(p.X, p.Y, p.Radius, boss.X, boss.Y, boss.Radius) {
			// Don't deactivate projectile if it's piercing (beam)
			if !p.Piercing {
				p.Active = false
			}

			// Add impact effect for boss
			if cm.OnImpactSpawned != nil {
				cm.OnImpactSpawned(boss.X, boss.Y, 40, color.RGBA{255, 150, 100, 255})
			}

			bossDefeated := false
			if cm.OnBossDamaged != nil {
				bossDefeated = cm.OnBossDamaged(p.Damage)
			} else {
				bossDefeated = boss.TakeDamage(p.Damage)
			}

			// Screen shake
			if cm.OnScreenShake != nil {
				if bossDefeated {
					cm.OnScreenShake(20)
				} else {
					cm.OnScreenShake(3)
				}
			}
		}
	}
}

// handleEnemyProjectilePlayerCollisions checks enemy projectiles against player
func (cm *CollisionManager) handleEnemyProjectilePlayerCollisions(
	player *entities.Player,
	projectiles []*entities.Projectile,
	gameTime float64,
	dashInvincibility float64,
) {
	for _, p := range projectiles {
		if !p.Active || p.Friendly {
			continue
		}

		if checkCircleCollision(p.X, p.Y, p.Radius, player.X, player.Y, player.Radius) {
			// Check dash invincibility
			if dashInvincibility > 0 {
				p.Active = false
				if cm.OnExplosionSpawned != nil {
					cm.OnExplosionSpawned(p.X, p.Y, 10)
				}
				continue
			}

			p.Active = false
			player.TakeDamage(p.Damage, gameTime)

			// Notify game of damage
			if cm.OnPlayerDamaged != nil {
				cm.OnPlayerDamaged(p.Damage)
			}

			// Screen shake and sound
			if cm.OnScreenShake != nil {
				cm.OnScreenShake(10)
			}
			if cm.OnSoundPlayed != nil {
				cm.OnSoundPlayed(systems.SoundHitPlayer)
			}

			// Check if player died
			if player.Health <= 0 {
				if cm.OnExplosionSpawned != nil {
					cm.OnExplosionSpawned(player.X, player.Y, 40)
				}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundExplosionLarge)
				}
				player.Active = false
			}
		}
	}
}

// handleEnemyPlayerCollisions checks direct enemy-player collisions
func (cm *CollisionManager) handleEnemyPlayerCollisions(
	player *entities.Player,
	enemies []*entities.Enemy,
	gameTime float64,
	dashInvincibility float64,
	damageMultiplier float64,
) {
	for _, e := range enemies {
		if !e.Active {
			continue
		}

		if checkCircleCollision(e.X, e.Y, e.Radius, player.X, player.Y, player.Radius) {
			// Check dash invincibility
			if dashInvincibility > 0 {
				e.Active = false
				if cm.OnExplosionSpawned != nil {
					cm.OnExplosionSpawned(e.X, e.Y, e.Radius)
				}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundExplosionSmall)
				}
				continue
			}

			e.Active = false
			if cm.OnExplosionSpawned != nil {
				cm.OnExplosionSpawned(e.X, e.Y, e.Radius)
			}

			// Play appropriate explosion sound
			if cm.OnSoundPlayed != nil {
				soundType := getEnemyExplosionSound(e.Type)
				cm.OnSoundPlayed(soundType)
				cm.OnSoundPlayed(systems.SoundHitPlayer)
			}

			collisionDamage := int(float64(30) * damageMultiplier)
			player.TakeDamage(collisionDamage, gameTime)

			// Notify game of damage
			if cm.OnPlayerDamaged != nil {
				cm.OnPlayerDamaged(collisionDamage)
			}

			// Screen shake
			if cm.OnScreenShake != nil {
				cm.OnScreenShake(15)
			}

			// Check if player died
			if player.Health <= 0 {
				if cm.OnExplosionSpawned != nil {
					cm.OnExplosionSpawned(player.X, player.Y, 40)
				}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundExplosionLarge)
				}
				player.Active = false
			}
		}
	}
}

// handleBossPlayerCollision checks boss-player collision
func (cm *CollisionManager) handleBossPlayerCollision(
	player *entities.Player,
	boss *entities.Boss,
	gameTime float64,
	damageMultiplier float64,
) {
	if checkCircleCollision(boss.X, boss.Y, boss.Radius*0.5, player.X, player.Y, player.Radius) {
		bossDamage := int(float64(50) * damageMultiplier)
		player.TakeDamage(bossDamage, gameTime)

		// Screen shake
		if cm.OnScreenShake != nil {
			cm.OnScreenShake(20)
		}

		// Check if player died
		if player.Health <= 0 {
			if cm.OnExplosionSpawned != nil {
				cm.OnExplosionSpawned(player.X, player.Y, 40)
			}
			player.Active = false
		}
	}
}

// handlePowerUpCollisions checks powerup collection
func (cm *CollisionManager) handlePowerUpCollisions(
	player *entities.Player,
	powerups []*entities.PowerUp,
) []*entities.PowerUp {
	for _, pu := range powerups {
		if !pu.Active {
			continue
		}

		if checkCircleCollision(pu.X, pu.Y, pu.Radius, player.X, player.Y, player.Radius) {
			pu.Active = false

			// Play appropriate sound and handle special cases
			cm.handlePowerUpEffect(player, pu)
		}
	}

	return powerups
}

// handlePowerUpEffect processes powerup effects
func (cm *CollisionManager) handlePowerUpEffect(player *entities.Player, pu *entities.PowerUp) {
	switch pu.Type {
	case entities.PowerUpHealth:
		if cm.OnSoundPlayed != nil {
			cm.OnSoundPlayed(systems.SoundPowerUpCollect)
		}
		player.ApplyPowerUp(pu.Type)

	case entities.PowerUpShield:
		if cm.OnSoundPlayed != nil {
			cm.OnSoundPlayed(systems.SoundShieldRecharge)
		}
		player.ApplyPowerUp(pu.Type)

	case entities.PowerUpWeapon:
		message, upgraded := player.ApplyPowerUp(pu.Type)
		if upgraded && message != "" {
			if cm.OnSoundPlayed != nil {
				cm.OnSoundPlayed(systems.SoundWeaponLevelUp)
			}
			if cm.OnFloatingTextAdded != nil {
				cm.OnFloatingTextAdded(player.X, player.Y-30, message, color.RGBA{255, 200, 50, 255})
			}
		}

	case entities.PowerUpSpeed:
		if cm.OnSoundPlayed != nil {
			cm.OnSoundPlayed(systems.SoundPowerUpCollect)
		}
		player.ApplyPowerUp(pu.Type)

	case entities.PowerUpMystery:
		message, isPositive := player.ApplyPowerUp(pu.Type)
		if message != "" {
			var textColor color.RGBA
			if isPositive {
				textColor = color.RGBA{50, 255, 50, 255}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundPowerUpCollect)
				}
			} else {
				textColor = color.RGBA{255, 50, 50, 255}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundHitPlayer)
				}
			}
			if cm.OnFloatingTextAdded != nil {
				cm.OnFloatingTextAdded(player.X, player.Y-50, message, textColor)
			}
		}
	}
}

// handleAsteroidPlayerCollisions checks asteroid collisions with player
func (cm *CollisionManager) handleAsteroidPlayerCollisions(
	player *entities.Player,
	asteroids []*entities.Asteroid,
	gameTime float64,
	damageMultiplier float64,
) {
	for _, a := range asteroids {
		if !a.Active {
			continue
		}

		if checkCircleCollision(a.X, a.Y, a.Radius, player.X, player.Y, player.Radius) {
			asteroidDamage := int(float64(15) * damageMultiplier)
			player.TakeDamage(asteroidDamage, gameTime)

			if cm.OnExplosionSpawned != nil {
				cm.OnExplosionSpawned(a.X, a.Y, a.Radius)
			}

			// Play sounds
			if cm.OnSoundPlayed != nil {
				cm.OnSoundPlayed(systems.SoundHitAsteroid)
				soundType := getAsteroidExplosionSound(a.Size)
				cm.OnSoundPlayed(soundType)
			}

			a.Active = false

			// Screen shake
			if cm.OnScreenShake != nil {
				cm.OnScreenShake(8)
			}

			// Check if player died
			if player.Health <= 0 {
				if cm.OnExplosionSpawned != nil {
					cm.OnExplosionSpawned(player.X, player.Y, 40)
				}
				if cm.OnSoundPlayed != nil {
					cm.OnSoundPlayed(systems.SoundExplosionLarge)
				}
				player.Active = false
			}
		}
	}
}

// handleProjectileAsteroidCollisions checks player projectiles against asteroids
func (cm *CollisionManager) handleProjectileAsteroidCollisions(
	projectiles []*entities.Projectile,
	asteroids []*entities.Asteroid,
) {
	for _, p := range projectiles {
		if !p.Active || !p.Friendly {
			continue
		}

		for _, a := range asteroids {
			if !a.Active {
				continue
			}

			if checkCircleCollision(p.X, p.Y, p.Radius, a.X, a.Y, a.Radius) {
				p.Active = false
				a.TakeDamage(p.Damage)

				// Add impact effect
				if cm.OnImpactSpawned != nil {
					cm.OnImpactSpawned(a.X, a.Y, 25, color.RGBA{200, 100, 50, 255})
				}

				if !a.Active {
					if cm.OnExplosionSpawned != nil {
						cm.OnExplosionSpawned(a.X, a.Y, a.Radius)
					}

					// Play sound
					if cm.OnSoundPlayed != nil {
						soundType := getAsteroidExplosionSound(a.Size)
						cm.OnSoundPlayed(soundType)
					}

					// Add score
					points := int64(10 + int(a.Radius))
					if cm.OnScoreAdded != nil {
						cm.OnScoreAdded(points)
					}
				}
			}
		}
	}
}

// Helper functions

// checkCircleCollision performs circle-circle collision detection
func checkCircleCollision(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	distSq := dx*dx + dy*dy
	radiusSum := r1 + r2
	return distSq < radiusSum*radiusSum
}

// getEnemyExplosionSound returns the appropriate sound for enemy explosion
func getEnemyExplosionSound(enemyType entities.EnemyType) systems.SoundType {
	switch enemyType {
	case entities.EnemyScout, entities.EnemyDrone:
		return systems.SoundExplosionSmall
	case entities.EnemyHunter, entities.EnemyBomber, entities.EnemySniper, entities.EnemySplitter:
		return systems.SoundExplosionMedium
	case entities.EnemyTank, entities.EnemyShieldBearer:
		return systems.SoundExplosionLarge
	default:
		return systems.SoundExplosionSmall
	}
}

// getAsteroidExplosionSound returns the appropriate sound for asteroid explosion
func getAsteroidExplosionSound(size entities.AsteroidSize) systems.SoundType {
	switch size {
	case entities.AsteroidSmall:
		return systems.SoundExplosionSmall
	case entities.AsteroidMedium:
		return systems.SoundExplosionMedium
	case entities.AsteroidLarge:
		return systems.SoundExplosionLarge
	default:
		return systems.SoundExplosionSmall
	}
}
