package core

import (
	"image/color"
	"stellar-siege/game/entities"
)

// EntityManager handles entity lifecycle management (creation, cleanup, etc.)
type EntityManager struct{}

// NewEntityManager creates a new entity manager
func NewEntityManager() *EntityManager {
	return &EntityManager{}
}

// CleanupInactive removes inactive entities from all slices
func (em *EntityManager) CleanupInactive(
	projectiles []*entities.Projectile,
	enemies []*entities.Enemy,
	explosions []*entities.Explosion,
	powerups []*entities.PowerUp,
	asteroids []*entities.Asteroid,
	floatingTexts []*entities.FloatingText,
	impactEffects []*entities.ImpactEffect,
) (
	[]*entities.Projectile,
	[]*entities.Enemy,
	[]*entities.Explosion,
	[]*entities.PowerUp,
	[]*entities.Asteroid,
	[]*entities.FloatingText,
	[]*entities.ImpactEffect,
) {
	// Clean projectiles
	activeProjectiles := projectiles[:0]
	for _, p := range projectiles {
		if p.Active {
			activeProjectiles = append(activeProjectiles, p)
		}
	}

	// Clean enemies
	activeEnemies := enemies[:0]
	for _, e := range enemies {
		if e.Active {
			activeEnemies = append(activeEnemies, e)
		}
	}

	// Clean explosions
	activeExplosions := explosions[:0]
	for _, ex := range explosions {
		if ex.Active {
			activeExplosions = append(activeExplosions, ex)
		}
	}

	// Clean powerups
	activePowerups := powerups[:0]
	for _, pu := range powerups {
		if pu.Active {
			activePowerups = append(activePowerups, pu)
		}
	}

	// Clean asteroids
	activeAsteroids := asteroids[:0]
	for _, a := range asteroids {
		if a.Active {
			activeAsteroids = append(activeAsteroids, a)
		}
	}

	// Clean floating text
	activeFloatingText := floatingTexts[:0]
	for _, ft := range floatingTexts {
		if ft.Active {
			activeFloatingText = append(activeFloatingText, ft)
		}
	}

	// Clean impact effects
	activeImpacts := impactEffects[:0]
	for _, ie := range impactEffects {
		if ie.Active {
			activeImpacts = append(activeImpacts, ie)
		}
	}

	return activeProjectiles, activeEnemies, activeExplosions, activePowerups, activeAsteroids, activeFloatingText, activeImpacts
}

// SpawnExplosion creates a new explosion entity
func (em *EntityManager) SpawnExplosion(x, y, size float64, explosions []*entities.Explosion) []*entities.Explosion {
	explosion := entities.NewExplosion(x, y, size)
	return append(explosions, explosion)
}

// SpawnExplosionWithType creates a new explosion entity with a specific type
func (em *EntityManager) SpawnExplosionWithType(x, y, size float64, expType entities.ExplosionType, explosions []*entities.Explosion) []*entities.Explosion {
	explosion := entities.NewExplosionWithType(x, y, size, expType)
	return append(explosions, explosion)
}

// SpawnFloatingScore creates a floating score indicator
func (em *EntityManager) SpawnFloatingScore(x, y float64, score int, floatingTexts []*entities.FloatingText) []*entities.FloatingText {
	ft := entities.NewFloatingScore(x, y, score)
	return append(floatingTexts, ft)
}

// SpawnFloatingDamage creates a floating damage indicator
func (em *EntityManager) SpawnFloatingDamage(x, y float64, damage int, floatingTexts []*entities.FloatingText) []*entities.FloatingText {
	ft := entities.NewFloatingDamage(x, y, damage)
	return append(floatingTexts, ft)
}

// SpawnFloatingUpgrade creates a floating weapon level indicator
func (em *EntityManager) SpawnFloatingUpgrade(x, y float64, level int, floatingTexts []*entities.FloatingText) []*entities.FloatingText {
	ft := entities.NewFloatingUpgrade(x, y, level)
	return append(floatingTexts, ft)
}

// SpawnFloatingText creates a floating text indicator with custom message and color
func (em *EntityManager) SpawnFloatingText(x, y float64, message string, col color.RGBA, floatingTexts []*entities.FloatingText) []*entities.FloatingText {
	ft := entities.NewFloatingText(x, y, message, col)
	return append(floatingTexts, ft)
}

// SpawnPowerUp creates a new powerup entity
func (em *EntityManager) SpawnPowerUp(x, y float64, powerups []*entities.PowerUp) []*entities.PowerUp {
	pu := entities.NewPowerUp(x, y)
	return append(powerups, pu)
}

// SpawnImpactEffect creates a new impact effect
func (em *EntityManager) SpawnImpactEffect(x, y, size float64, col color.RGBA, impactEffects []*entities.ImpactEffect) []*entities.ImpactEffect {
	ie := entities.NewImpactEffect(x, y, size, col)
	return append(impactEffects, ie)
}

// SpawnAsteroid creates a new asteroid entity
func (em *EntityManager) SpawnAsteroid(x, y float64, size entities.AsteroidSize, asteroids []*entities.Asteroid) []*entities.Asteroid {
	asteroid := entities.NewAsteroid(x, y, size)
	return append(asteroids, asteroid)
}
