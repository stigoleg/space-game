package entities

import (
	"image/color"
	"math"
)

// checkProximityAttack checks if player is close during entry phase and triggers early attack
func (b *Boss) checkProximityAttack(playerX float64) bool {
	const proximityRange = 200.0
	distanceToPlayer := math.Abs(b.X - playerX)
	return b.Y >= 0 && distanceToPlayer < proximityRange
}

// updateEntryPhase handles the entry phase movement and transition logic
func (b *Boss) updateEntryPhase(playerX float64) {
	b.Y += 1

	if b.Y >= b.EntryY {
		// Reached target position
		b.Y = b.EntryY
		b.Phase = BossPhaseAttacking
		b.TelegraphTimer = 0.5
		b.TelegraphActive = true
	} else if b.checkProximityAttack(playerX) {
		// Player is close while entering - start attacking early
		b.Phase = BossPhaseAttacking
		b.TelegraphTimer = 0.3
		b.TelegraphActive = true
	}
}

// updateTelegraphWarning manages the telegraph warning timer
func (b *Boss) updateTelegraphWarning() {
	if b.TelegraphActive {
		b.TelegraphTimer -= 1.0 / 60.0
		if b.TelegraphTimer <= 0 {
			b.TelegraphActive = false
		}
	}
}

// updateMovement handles horizontal movement tracking the player
func (b *Boss) updateMovement(playerX float64, screenWidth int) {
	dx := playerX - b.X
	b.VelX = dx * 0.01 * b.Speed

	if b.Phase == BossPhaseRage {
		b.VelX *= 1.5
	}

	b.X += b.VelX

	// Keep in bounds
	margin := b.Radius + 20
	if b.X < margin {
		b.X = margin
	}
	if b.X > float64(screenWidth)-margin {
		b.X = float64(screenWidth) - margin
	}
}

// getAttackInterval calculates the attack interval based on phase and difficulty
func (b *Boss) getAttackInterval() float64 {
	attackInterval := 1.5 / b.AttackRate

	if b.Phase == BossPhaseRage {
		attackInterval *= 0.6
	}
	if b.Phase == BossPhaseSpecialAttack {
		attackInterval *= 0.5
	}

	return attackInterval
}

// updateAttackTimer manages attack timing and triggers attacks
func (b *Boss) updateAttackTimer(playerX, playerY float64) []*Projectile {
	var projectiles []*Projectile

	if b.TelegraphActive {
		return projectiles
	}

	b.AttackTimer += 1.0 / 60.0
	attackInterval := b.getAttackInterval()

	if b.AttackTimer >= attackInterval {
		b.AttackTimer = 0
		b.AttackPattern = (b.AttackPattern + 1) % b.PatternCount
		projectiles = b.executeAttack(playerX, playerY)
	}

	return projectiles
}

// shouldTransitionToSpecialAttack checks if boss should enter special attack phase
func (b *Boss) shouldTransitionToSpecialAttack() bool {
	healthPercent := float64(b.Health) / float64(b.MaxHealth)
	return b.BossLevel >= 2 && b.Phase == BossPhaseAttacking &&
		healthPercent < 0.6 && healthPercent > 0.3
}

// shouldTransitionToRage checks if boss should enter rage phase
func (b *Boss) shouldTransitionToRage() bool {
	healthPercent := float64(b.Health) / float64(b.MaxHealth)
	return (b.Phase == BossPhaseAttacking || b.Phase == BossPhaseSpecialAttack) &&
		healthPercent < 0.3
}

// checkPhaseTransitions handles all phase transitions based on health
func (b *Boss) checkPhaseTransitions() {
	if b.shouldTransitionToSpecialAttack() {
		b.Phase = BossPhaseSpecialAttack
		b.SpecialTimer = 0
		b.SpecialPhase = 0
	} else if b.shouldTransitionToRage() {
		b.Phase = BossPhaseRage
	}
}

// updateShieldMechanic manages shield activation and deactivation
func (b *Boss) updateShieldMechanic() {
	b.ShieldTimer += 1.0 / 60.0
	shieldInterval := 5.0 / (float64(b.BossLevel) * 0.5)
	shieldDuration := 2.0 - float64(b.BossLevel)*0.2

	if b.ShieldTimer >= shieldInterval && !b.ShieldUp {
		b.ShieldUp = true
		b.ShieldTimer = 0
	}
	if b.ShieldUp && b.ShieldTimer >= shieldDuration {
		b.ShieldUp = false
		b.ShieldTimer = 0
	}
}

// updateAttackingPhase handles all logic for attacking phases
func (b *Boss) updateAttackingPhase(playerX, playerY float64, screenWidth int) []*Projectile {
	b.updateTelegraphWarning()
	b.updateMovement(playerX, screenWidth)
	projectiles := b.updateAttackTimer(playerX, playerY)
	b.checkPhaseTransitions()
	b.updateShieldMechanic()
	return projectiles
}

// getPhaseColors returns the main, core, and glow colors based on current phase
func (b *Boss) getPhaseColors() (mainColor, coreColor, glowColor color.RGBA) {
	colorIntensity := uint8(200 - b.BossLevel*20)

	switch b.Phase {
	case BossPhaseEntering, BossPhaseAttacking:
		mainColor = color.RGBA{colorIntensity, uint8(80 - b.BossLevel*10), 60, 255}
		coreColor = color.RGBA{255, 150, 100, 255}
		glowColor = color.RGBA{255, 80, 60, 100}
	case BossPhaseSpecialAttack:
		mainColor = color.RGBA{150, 80, 200, 255}
		coreColor = color.RGBA{200, 150, 255, 255}
		glowColor = color.RGBA{200, 100, 255, 120}
	case BossPhaseRage:
		intensity := uint8(220 + 35*math.Sin(b.AnimTimer*10))
		mainColor = color.RGBA{intensity, 60, 40, 255}
		coreColor = color.RGBA{255, 180, 80, 255}
		glowColor = color.RGBA{255, 130, 80, 120}
	case BossPhaseDying:
		mainColor = color.RGBA{120, 120, 120, 200}
		coreColor = color.RGBA{255, 220, 80, 255}
		glowColor = color.RGBA{255, 220, 150, 150}
	}

	return mainColor, coreColor, glowColor
}

// getHealthBarColor returns the health bar color based on current phase
func (b *Boss) getHealthBarColor() color.RGBA {
	switch b.Phase {
	case BossPhaseRage:
		return color.RGBA{255, 180, 80, 255}
	case BossPhaseSpecialAttack:
		return color.RGBA{200, 100, 255, 255}
	default:
		return color.RGBA{255, 80, 60, 255}
	}
}
