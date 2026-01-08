package entities

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BossPhase int

const (
	BossPhaseEntering BossPhase = iota
	BossPhaseAttacking
	BossPhaseSpecialAttack // New phase for special attacks
	BossPhaseRage
	BossPhaseDying
)

type Boss struct {
	X, Y          float64
	VelX, VelY    float64
	Radius        float64
	Health        int
	MaxHealth     int
	Points        int
	Active        bool
	Phase         BossPhase
	AttackTimer   float64
	AttackPattern int
	AnimTimer     float64
	ShieldUp      bool
	ShieldTimer   float64
	EntryY        float64

	// Difficulty scaling
	BossLevel    int     // Which boss encounter is this (1, 2, 3, 4+)
	Speed        float64 // Movement speed multiplier
	AttackRate   float64 // Attack speed multiplier
	Damage       int     // Base damage of projectiles
	PatternCount int     // Number of attack patterns available

	// Special attacks
	SpecialTimer  float64 // Timer for special attack cooldown
	SpecialPhase  int     // Current special attack phase
	MinionsActive int     // Number of active minions
}

func NewBoss(screenWidth int, bossLevel int) *Boss {
	// Progressive difficulty scaling
	var health int
	var baseDamage int
	var patternCount int
	var speedMult float64
	var attackRateMult float64

	switch bossLevel {
	case 1: // Wave 5 - First boss
		health = 500
		baseDamage = 15
		patternCount = 4
		speedMult = 1.0
		attackRateMult = 1.0
	case 2: // Wave 10 - Intermediate boss
		health = 1000
		baseDamage = 20
		patternCount = 6
		speedMult = 1.2
		attackRateMult = 1.3
	case 3: // Wave 15 - Advanced boss
		health = 1500
		baseDamage = 25
		patternCount = 8
		speedMult = 1.4
		attackRateMult = 1.6
	default: // Wave 20+ - Extreme boss
		health = 2000 + (bossLevel-4)*500
		baseDamage = 30
		patternCount = 10
		speedMult = 1.6 + float64(bossLevel-4)*0.2
		attackRateMult = 2.0 + float64(bossLevel-4)*0.2
	}

	return &Boss{
		X:             float64(screenWidth) / 2,
		Y:             -100,
		Radius:        60,
		Health:        health,
		MaxHealth:     health,
		Points:        10000 * bossLevel,
		Active:        true,
		Phase:         BossPhaseEntering,
		EntryY:        120,
		BossLevel:     bossLevel,
		Speed:         speedMult,
		AttackRate:    attackRateMult,
		Damage:        baseDamage,
		PatternCount:  patternCount,
		SpecialTimer:  0,
		MinionsActive: 0,
	}
}

func (b *Boss) Update(playerX, playerY float64, screenWidth, screenHeight int) []*Projectile {
	b.AnimTimer += 0.05
	var projectiles []*Projectile

	switch b.Phase {
	case BossPhaseEntering:
		// Move into position
		b.Y += 1
		if b.Y >= b.EntryY {
			b.Y = b.EntryY
			b.Phase = BossPhaseAttacking
		}

	case BossPhaseAttacking, BossPhaseRage, BossPhaseSpecialAttack:
		// Horizontal movement - track player loosely (scaled by boss level)
		targetX := playerX
		dx := targetX - b.X
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

		// Attack patterns with difficulty scaling
		b.AttackTimer += 1.0 / 60.0

		// Base attack interval decreases with boss level
		attackInterval := 1.5 / b.AttackRate
		if b.Phase == BossPhaseRage {
			attackInterval *= 0.6 // Much faster in rage mode
		}
		if b.Phase == BossPhaseSpecialAttack {
			attackInterval *= 0.5 // Even faster in special attack
		}

		if b.AttackTimer >= attackInterval {
			b.AttackTimer = 0
			b.AttackPattern = (b.AttackPattern + 1) % b.PatternCount
			projectiles = b.executeAttack(playerX, playerY)
		}

		// Phase transitions based on health percentage
		healthPercent := float64(b.Health) / float64(b.MaxHealth)

		// Enter special attack at 60% health (for boss level 2+)
		if b.BossLevel >= 2 && b.Phase == BossPhaseAttacking && healthPercent < 0.6 && healthPercent > 0.3 {
			b.Phase = BossPhaseSpecialAttack
			b.SpecialTimer = 0
			b.SpecialPhase = 0
		}

		// Enter rage mode at 30% health
		if (b.Phase == BossPhaseAttacking || b.Phase == BossPhaseSpecialAttack) && healthPercent < 0.3 {
			b.Phase = BossPhaseRage
		}

		// Shield mechanic (varies with boss level)
		b.ShieldTimer += 1.0 / 60.0
		shieldInterval := 5.0 / (float64(b.BossLevel) * 0.5) // More frequent at higher levels
		shieldDuration := 2.0 - float64(b.BossLevel)*0.2     // Shorter duration at higher levels

		if b.ShieldTimer >= shieldInterval && !b.ShieldUp {
			b.ShieldUp = true
			b.ShieldTimer = 0
		}
		if b.ShieldUp && b.ShieldTimer >= shieldDuration {
			b.ShieldUp = false
			b.ShieldTimer = 0
		}

	case BossPhaseDying:
		// Explosion sequence handled elsewhere
	}

	return projectiles
}

func (b *Boss) executeAttack(playerX, playerY float64) []*Projectile {
	var projectiles []*Projectile

	// Use different patterns based on boss level
	pattern := b.AttackPattern % b.PatternCount

	switch pattern {
	case 0:
		// Spread shot - more shots for higher levels
		spread := 3
		if b.BossLevel >= 2 {
			spread = 5
		}
		for i := -spread; i <= spread; i++ {
			angle := math.Pi/2 + float64(i)*0.2
			projectiles = append(projectiles, b.createProjectile(angle))
		}

	case 1:
		// Aimed shot at player with side shots
		dx := playerX - b.X
		dy := playerY - b.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			velX := (dx / dist) * 7
			velY := (dy / dist) * 7
			projectiles = append(projectiles, NewProjectile(b.X, b.Y+b.Radius, velX, velY, false, b.Damage))

			// Add more side shots for higher levels
			sideCount := 1
			if b.BossLevel >= 3 {
				sideCount = 2
			}
			for s := 1; s <= sideCount; s++ {
				offsetX := float64(s * 30)
				projectiles = append(projectiles, NewProjectile(b.X-offsetX, b.Y+b.Radius, velX*0.8, velY*0.8, false, b.Damage-5))
				projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+b.Radius, velX*0.8, velY*0.8, false, b.Damage-5))
			}
		}

	case 2:
		// Circle burst - more projectiles for higher levels
		count := 12
		if b.BossLevel >= 2 {
			count = 16
		}
		if b.BossLevel >= 4 {
			count = 20
		}
		for i := 0; i < count; i++ {
			angle := float64(i) * (2 * math.Pi) / float64(count)
			projectiles = append(projectiles, b.createProjectile(angle))
		}

	case 3:
		// Laser lines - more lines for higher levels
		count := 5
		if b.BossLevel >= 2 {
			count = 7
		}
		for i := 0; i < count; i++ {
			offsetX := float64(i-(count-1)/2) * 40
			projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+b.Radius, 0, 6, false, b.Damage))
		}

	case 4:
		// Spiral pattern (available for boss level 2+)
		if b.BossLevel >= 2 {
			for i := 0; i < 8; i++ {
				angle := float64(i)*math.Pi/4 + b.AnimTimer*0.1
				velX := math.Cos(angle) * 5
				velY := math.Sin(angle) * 5
				projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage-5))
			}
		}

	case 5:
		// Double arc pattern (available for boss level 2+)
		if b.BossLevel >= 2 {
			for i := -4; i <= 4; i++ {
				angle := math.Pi/2 + float64(i)*0.15
				velX := math.Cos(angle) * 6
				velY := math.Sin(angle) * 6
				projectiles = append(projectiles, NewProjectile(b.X-40, b.Y+b.Radius, velX, velY, false, b.Damage-5))
				projectiles = append(projectiles, NewProjectile(b.X+40, b.Y+b.Radius, velX, velY, false, b.Damage-5))
			}
		}

	case 6:
		// Tracking spiral (available for boss level 3+)
		if b.BossLevel >= 3 {
			for i := 0; i < 6; i++ {
				angle := float64(i)*math.Pi/3 + b.AnimTimer*0.2
				velX := math.Cos(angle) * 5.5
				velY := math.Sin(angle) * 5.5
				projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage))
			}
		}

	case 7:
		// Wave pattern (available for boss level 3+)
		if b.BossLevel >= 3 {
			waveCount := 10
			for i := 0; i < waveCount; i++ {
				offsetX := float64(i-(waveCount-1)/2) * 30
				waveY := math.Sin(float64(i)*math.Pi/5) * 50
				projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+waveY, 0, 6, false, b.Damage-5))
			}
		}

	case 8:
		// Cross burst (available for boss level 4+)
		if b.BossLevel >= 4 {
			for i := 0; i < 4; i++ {
				angle := float64(i)*math.Pi/2 + math.Pi/4
				for j := 0; j < 4; j++ {
					ratio := float64(j) / 3.0
					vel := 4.0 + ratio*3.0
					velX := math.Cos(angle) * vel
					velY := math.Sin(angle) * vel
					projectiles = append(projectiles, NewProjectile(b.X, b.Y+b.Radius, velX, velY, false, b.Damage))
				}
			}
		}

	case 9:
		// Chaos pattern (available for boss level 4+)
		if b.BossLevel >= 4 {
			for i := 0; i < 12; i++ {
				angle := float64(i)*math.Pi/6 + b.AnimTimer*0.3
				speed := 4.0 + math.Sin(b.AnimTimer+float64(i))*2.0
				velX := math.Cos(angle) * speed
				velY := math.Sin(angle) * speed
				projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage))
			}
		}
	}

	return projectiles
}

// Helper function to create projectiles with proper damage scaling
func (b *Boss) createProjectile(angle float64) *Projectile {
	speed := 5.0
	if b.Phase == BossPhaseRage {
		speed = 6.5
	}
	if b.Phase == BossPhaseSpecialAttack {
		speed = 7.0
	}
	velX := math.Cos(angle) * speed
	velY := math.Sin(angle) * speed
	return NewProjectile(b.X, b.Y+b.Radius, velX, velY, false, b.Damage)
}

func (b *Boss) TakeDamage(damage int) bool {
	if b.ShieldUp {
		return false // Damage blocked
	}
	b.Health -= damage
	if b.Health <= 0 {
		b.Phase = BossPhaseDying
		return true // Boss defeated
	}
	return false
}

func (b *Boss) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	// Simple screen coordinates with shake
	x := float32(b.X + shakeX)
	y := float32(b.Y + shakeY)

	// Pulsing animation
	pulse := float32(1.0 + 0.05*math.Sin(b.AnimTimer*3))
	radius := float32(b.Radius) * pulse

	// Color based on phase and difficulty level
	var mainColor, coreColor, glowColor color.RGBA

	// Base color changes with boss level
	colorIntensity := uint8(200 - b.BossLevel*20) // Get darker at higher levels

	switch b.Phase {
	case BossPhaseEntering, BossPhaseAttacking:
		mainColor = color.RGBA{colorIntensity, uint8(80 - b.BossLevel*10), 60, 255}
		coreColor = color.RGBA{255, 150, 100, 255}
		glowColor = color.RGBA{255, 80, 60, 100}
	case BossPhaseSpecialAttack:
		// Purple glow for special attack phase
		mainColor = color.RGBA{150, 80, 200, 255}
		coreColor = color.RGBA{200, 150, 255, 255}
		glowColor = color.RGBA{200, 100, 255, 120}
	case BossPhaseRage:
		// Flashing red in rage
		intensity := uint8(220 + 35*math.Sin(b.AnimTimer*10))
		mainColor = color.RGBA{intensity, 60, 40, 255}
		coreColor = color.RGBA{255, 180, 80, 255}
		glowColor = color.RGBA{255, 130, 80, 120}
	case BossPhaseDying:
		mainColor = color.RGBA{120, 120, 120, 200}
		coreColor = color.RGBA{255, 220, 80, 255}
		glowColor = color.RGBA{255, 220, 150, 150}
	}

	// Draw shadow
	shadowColor := color.RGBA{20, 20, 30, 100}
	vector.DrawFilledCircle(screen, x, y+radius+10, radius*0.6, shadowColor, true)

	// Main body
	vector.DrawFilledCircle(screen, x, y, radius, mainColor, true)

	// Core
	coreSize := radius * 0.45
	vector.DrawFilledCircle(screen, x, y, coreSize, coreColor, true)

	// Side pods with 3D effect
	podOffset := float32(60) * pulse
	podRadius := float32(24) * pulse

	// Left pod
	vector.DrawFilledCircle(screen, x-podOffset, y, podRadius, mainColor, true)
	vector.DrawFilledCircle(screen, x-podOffset, y, podRadius*0.55, coreColor, true)

	// Right pod
	vector.DrawFilledCircle(screen, x+podOffset, y, podRadius, mainColor, true)
	vector.DrawFilledCircle(screen, x+podOffset, y, podRadius*0.55, coreColor, true)

	// Outer glow - more intense for higher levels
	glowAlpha := uint8(100 + b.BossLevel*20)
	glowColor.A = glowAlpha
	vector.DrawFilledCircle(screen, x, y, radius+15, glowColor, true)

	// Highlight (3D effect)
	vector.DrawFilledCircle(screen, x-radius*0.3, y-radius*0.3, radius*0.25, color.RGBA{mainColor.R + 50, mainColor.G + 50, mainColor.B + 50, 200}, true)

	// Shield effect
	if b.ShieldUp {
		shieldPulse := float32(0.8 + 0.2*math.Sin(b.AnimTimer*8))
		shieldColor := color.RGBA{100, 200, 255, uint8(150 * shieldPulse)}
		vector.StrokeCircle(screen, x, y, radius+30, 4, shieldColor, true)
		vector.StrokeCircle(screen, x, y, radius+35, 2, color.RGBA{150, 220, 255, 100}, true)
	}

	// Health bar
	barWidth := float32(140)
	barHeight := float32(12)
	healthRatio := float32(b.Health) / float32(b.MaxHealth)

	barX := x - barWidth/2
	barY := y - radius - 40

	// Background
	vector.DrawFilledRect(screen, barX, barY, barWidth, barHeight, color.RGBA{30, 30, 30, 200}, true)
	// Health fill
	healthColor := color.RGBA{255, 80, 60, 255}
	if b.Phase == BossPhaseRage {
		healthColor = color.RGBA{255, 180, 80, 255}
	}
	if b.Phase == BossPhaseSpecialAttack {
		healthColor = color.RGBA{200, 100, 255, 255}
	}
	vector.DrawFilledRect(screen, barX, barY, barWidth*healthRatio, barHeight, healthColor, true)
	// Border
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2, color.RGBA{255, 255, 255, 150}, true)

	// Boss level indicator (stars or rings)
	levelIndicatorY := y - radius - 65
	for i := 0; i < b.BossLevel; i++ {
		starX := x - float32((b.BossLevel-1)*8) + float32(i*16)
		starSize := float32(6)
		// Draw star/diamond shape
		vector.DrawFilledCircle(screen, starX, levelIndicatorY, starSize, color.RGBA{255, 215, 0, 255}, true)
	}
}

// Helper to check if boss is still active for the game loop
func (b *Boss) IsDead() bool {
	return b.Health <= 0 || b.Phase == BossPhaseDying
}
