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

	// Telegraph system
	TelegraphTimer  float64 // Timer for pre-attack warning
	TelegraphActive bool    // Is telegraph warning active
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
		attackRateMult = 1.1 // Reduced from 1.3 to 1.1 for more balanced difficulty
	case 3: // Wave 15 - Advanced boss
		health = 1500
		baseDamage = 25
		patternCount = 8
		speedMult = 1.4
		attackRateMult = 1.3 // Reduced from 1.6 to 1.3 for more balanced difficulty
	default: // Wave 20+ - Extreme boss
		health = 2000 + (bossLevel-4)*500
		baseDamage = 30
		patternCount = 10
		speedMult = 1.6 + float64(bossLevel-4)*0.1      // Reduced speed scaling
		attackRateMult = 1.5 + float64(bossLevel-4)*0.1 // Reduced from 2.0 base and 0.2 scaling
	}

	return &Boss{
		X:               float64(screenWidth) / 2,
		Y:               -100,
		Radius:          60,
		Health:          health,
		MaxHealth:       health,
		Points:          10000 * bossLevel,
		Active:          true,
		Phase:           BossPhaseEntering,
		EntryY:          150, // Boss enters further into screen so HP bar is visible
		BossLevel:       bossLevel,
		Speed:           speedMult,
		AttackRate:      attackRateMult,
		Damage:          baseDamage,
		PatternCount:    patternCount,
		SpecialTimer:    0,
		MinionsActive:   0,
		TelegraphTimer:  0,
		TelegraphActive: false,
	}
}

func (b *Boss) Update(playerX, playerY float64, screenWidth, screenHeight int) []*Projectile {
	b.AnimTimer += 0.05
	var projectiles []*Projectile

	switch b.Phase {
	case BossPhaseEntering:
		b.updateEntryPhase(playerX)

	case BossPhaseAttacking, BossPhaseRage, BossPhaseSpecialAttack:
		projectiles = b.updateAttackingPhase(playerX, playerY, screenWidth)

	case BossPhaseDying:
		// Explosion sequence handled elsewhere
	}

	return projectiles
}

func (b *Boss) executeAttack(playerX, playerY float64) []*Projectile {
	pattern := b.AttackPattern % b.PatternCount

	switch pattern {
	case 0:
		return b.executeSpreadShot()
	case 1:
		return b.executeAimedShot(playerX, playerY)
	case 2:
		return b.executeCircularBurst()
	case 3:
		return b.executeLaserLines()
	case 4:
		return b.executeSpiralPattern()
	case 5:
		return b.executeDoubleArc()
	case 6:
		return b.executeTrackingSpiral()
	case 7:
		return b.executeWavePattern()
	case 8:
		return b.executeCrossBurst()
	case 9:
		return b.executeChaosPattern()
	}

	return nil
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

	// Get colors based on phase
	mainColor, coreColor, glowColor := b.getPhaseColors()

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

	// Telegraph warning effect
	if b.TelegraphActive {
		b.drawTelegraphWarning(screen, x, y, radius)
	}

	// Shield effect
	if b.ShieldUp {
		b.drawShieldEffect(screen, x, y, radius)
	}

	// Health bar
	b.drawHealthBar(screen, x, y, radius)

	// Boss level indicator
	b.drawLevelIndicator(screen, x, y, radius)
}

// drawTelegraphWarning draws the telegraph warning rings
func (b *Boss) drawTelegraphWarning(screen *ebiten.Image, x, y, radius float32) {
	telegraphIntensity := float32(b.TelegraphTimer / 0.5)
	if telegraphIntensity > 1.0 {
		telegraphIntensity = 1.0
	}
	telegraphAlpha := uint8(180 * telegraphIntensity)
	telegraphPulse := float32(1.0 + 0.3*math.Sin(b.AnimTimer*12))

	telegraphColor := color.RGBA{255, 220, 0, telegraphAlpha}
	vector.StrokeCircle(screen, x, y, radius+20*telegraphPulse, 3, telegraphColor, true)
	vector.StrokeCircle(screen, x, y, radius+30*telegraphPulse, 2, color.RGBA{255, 180, 0, telegraphAlpha / 2}, true)
}

// drawShieldEffect draws the shield effect
func (b *Boss) drawShieldEffect(screen *ebiten.Image, x, y, radius float32) {
	shieldPulse := float32(0.8 + 0.2*math.Sin(b.AnimTimer*8))
	shieldColor := color.RGBA{100, 200, 255, uint8(150 * shieldPulse)}
	vector.StrokeCircle(screen, x, y, radius+30, 4, shieldColor, true)
	vector.StrokeCircle(screen, x, y, radius+35, 2, color.RGBA{150, 220, 255, 100}, true)
}

// drawHealthBar draws the health bar above the boss
func (b *Boss) drawHealthBar(screen *ebiten.Image, x, y, radius float32) {
	barWidth := float32(140)
	barHeight := float32(12)
	healthRatio := float32(b.Health) / float32(b.MaxHealth)

	barX := x - barWidth/2
	barY := y - radius - 40

	// Background
	vector.DrawFilledRect(screen, barX, barY, barWidth, barHeight, color.RGBA{30, 30, 30, 200}, true)
	// Health fill
	healthColor := b.getHealthBarColor()
	vector.DrawFilledRect(screen, barX, barY, barWidth*healthRatio, barHeight, healthColor, true)
	// Border
	vector.StrokeRect(screen, barX, barY, barWidth, barHeight, 2, color.RGBA{255, 255, 255, 150}, true)
}

// drawLevelIndicator draws the boss level stars
func (b *Boss) drawLevelIndicator(screen *ebiten.Image, x, y, radius float32) {
	levelIndicatorY := y - radius - 65
	for i := 0; i < b.BossLevel; i++ {
		starX := x - float32((b.BossLevel-1)*8) + float32(i*16)
		starSize := float32(6)
		vector.DrawFilledCircle(screen, starX, levelIndicatorY, starSize, color.RGBA{255, 215, 0, 255}, true)
	}
}

// Helper to check if boss is still active for the game loop
func (b *Boss) IsDead() bool {
	return b.Health <= 0 || b.Phase == BossPhaseDying
}

// Interface implementation methods

// IsActive returns whether the boss is active
func (b *Boss) IsActive() bool {
	return b.Active
}

// GetPosition returns the boss's position
func (b *Boss) GetPosition() (x, y float64) {
	return b.X, b.Y
}

// GetCollisionBounds returns the boss's collision bounds
func (b *Boss) GetCollisionBounds() (x, y, radius float64) {
	return b.X, b.Y, b.Radius
}

// GetHealth returns the boss's current health
func (b *Boss) GetHealth() int {
	return b.Health
}

// GetMaxHealth returns the boss's maximum health
func (b *Boss) GetMaxHealth() int {
	return b.MaxHealth
}

// GetVelocity returns the boss's velocity
func (b *Boss) GetVelocity() (vx, vy float64) {
	return b.VelX, b.VelY
}

// SetVelocity sets the boss's velocity
func (b *Boss) SetVelocity(vx, vy float64) {
	b.VelX = vx
	b.VelY = vy
}

// GetSpeed returns the boss's speed
func (b *Boss) GetSpeed() float64 {
	return b.Speed
}

// SetSpeed sets the boss's speed
func (b *Boss) SetSpeed(speed float64) {
	b.Speed = speed
}
