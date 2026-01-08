package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	X, Y         float64
	VelX, VelY   float64
	Radius       float64
	Speed        float64
	Health       int
	MaxHealth    int
	Shield       int
	MaxShield    int
	WeaponLevel  int
	FireRate     float64
	FireCooldown float64
	InvincTimer  float64
	Active       bool
	EngineGlow   float64

	// Difficulty-dependent settings
	ShieldRegenRate   float64 // HP per frame
	InvincibilityTime float64 // seconds
	ShieldRegenDelay  float64 // seconds before regen starts
	LastDamageTime    float64 // when damage was last taken
	PrevShield        int     // Track previous shield value for sound effects

	// Special attack mechanics
	ChargeLevel       float64 // 0 to 1.0 (charge for special attack)
	UltimateCharge    float64 // 0 to 1.0 (builds from combat)
	MaxUltimateCharge float64 // 1.0
	UltimateActive    bool    // Ultimate ability activated
	UltimateTimer     float64 // Duration of ultimate effect

	// Thruster trail system
	ThrusterTrail []struct{ X, Y, Life float64 }
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		X:                 x,
		Y:                 y,
		Radius:            20,
		Speed:             6,
		Health:            100,
		MaxHealth:         100,
		Shield:            50,
		MaxShield:         50,
		WeaponLevel:       1,
		FireRate:          0.12,
		Active:            true,
		EngineGlow:        0,
		ShieldRegenRate:   0.5,  // Default (Normal difficulty)
		InvincibilityTime: 0.25, // Default (Normal difficulty)
		ShieldRegenDelay:  3.5,  // Default (Normal difficulty)
		LastDamageTime:    -999, // Initialize to long ago so regen starts immediately
		PrevShield:        50,   // Initialize to starting shield
		ChargeLevel:       0,    // No charge initially
		UltimateCharge:    0,    // No ultimate initially
		MaxUltimateCharge: 1.0,  // Max ultimate charge
		UltimateActive:    false,
		UltimateTimer:     0,
		ThrusterTrail:     make([]struct{ X, Y, Life float64 }, 0, 10),
	}
}

func (p *Player) Update(screenWidth, screenHeight int, gameTime float64) {
	// Handle movement
	p.VelX = 0
	p.VelY = 0

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.VelY = -p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.VelY = p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.VelX = -p.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.VelX = p.Speed
	}

	// Normalize diagonal movement
	if p.VelX != 0 && p.VelY != 0 {
		p.VelX *= 0.707
		p.VelY *= 0.707
	}

	p.X += p.VelX
	p.Y += p.VelY

	// Add thruster trail particles when moving
	if (p.VelX != 0 || p.VelY != 0) && rand.Float64() < 0.6 {
		p.ThrusterTrail = append(p.ThrusterTrail, struct{ X, Y, Life float64 }{
			X:    p.X + (rand.Float64()-0.5)*8,
			Y:    p.Y + p.Radius + (rand.Float64()-0.5)*4,
			Life: 0.5,
		})
		// Keep trail list bounded
		if len(p.ThrusterTrail) > 15 {
			p.ThrusterTrail = p.ThrusterTrail[1:]
		}
	}

	// Update thruster trail
	for i := 0; i < len(p.ThrusterTrail); i++ {
		p.ThrusterTrail[i].Life -= 1.0 / 60.0
		p.ThrusterTrail[i].Y += 1.5 // Trail drifts down slightly
		if p.ThrusterTrail[i].Life <= 0 {
			p.ThrusterTrail = append(p.ThrusterTrail[:i], p.ThrusterTrail[i+1:]...)
			i--
		}
	}

	// Bounds check
	margin := p.Radius
	if p.X < margin {
		p.X = margin
	}
	if p.X > float64(screenWidth)-margin {
		p.X = float64(screenWidth) - margin
	}
	if p.Y < margin {
		p.Y = margin
	}
	if p.Y > float64(screenHeight)-margin {
		p.Y = float64(screenHeight) - margin
	}

	// Update cooldowns
	if p.FireCooldown > 0 {
		p.FireCooldown -= 1.0 / 60.0
	}
	if p.InvincTimer > 0 {
		p.InvincTimer -= 1.0 / 60.0
	}

	// Regenerate shield slowly - only after delay since last damage
	timeSinceDamage := gameTime - p.LastDamageTime
	if p.Shield < p.MaxShield && p.InvincTimer <= 0 && timeSinceDamage >= p.ShieldRegenDelay {
		p.Shield++
		if p.Shield > p.MaxShield {
			p.Shield = p.MaxShield
		}
	}

	// Engine glow animation
	p.EngineGlow += 0.2

	// Charge mechanics
	// Handle charge attack (hold space to charge)
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		// Charging shot (slower than normal shooting)
		if p.ChargeLevel < 1.0 {
			p.ChargeLevel += 0.02 // Charge over ~3 seconds
		}
	} else {
		p.ChargeLevel = 0 // Reset when not charging
	}

	// Ultimate ability mechanics - builds from combat
	if p.UltimateActive {
		// Ultimate is active
		p.UltimateTimer -= 1.0 / 60.0
		if p.UltimateTimer <= 0 {
			p.UltimateActive = false
			p.UltimateTimer = 0
		}
	}

	// Slow ultimate charge regen (1% per second during gameplay)
	if !p.UltimateActive && p.UltimateCharge < p.MaxUltimateCharge {
		p.UltimateCharge += 0.01 / 60.0
		if p.UltimateCharge > p.MaxUltimateCharge {
			p.UltimateCharge = p.MaxUltimateCharge
		}
	}
}

func (p *Player) Shoot() []*Projectile {
	if p.FireCooldown > 0 {
		return nil
	}
	p.FireCooldown = p.FireRate

	var projectiles []*Projectile

	switch p.WeaponLevel {
	case 1:
		projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -12, true, 10))
	case 2:
		projectiles = append(projectiles, NewProjectile(p.X-10, p.Y-p.Radius, 0, -12, true, 10))
		projectiles = append(projectiles, NewProjectile(p.X+10, p.Y-p.Radius, 0, -12, true, 10))
	case 3:
		projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -12, true, 12))
		projectiles = append(projectiles, NewProjectile(p.X-15, p.Y-p.Radius+5, -1, -11, true, 10))
		projectiles = append(projectiles, NewProjectile(p.X+15, p.Y-p.Radius+5, 1, -11, true, 10))
	case 4:
		projectiles = append(projectiles, NewProjectile(p.X-8, p.Y-p.Radius, 0, -13, true, 15))
		projectiles = append(projectiles, NewProjectile(p.X+8, p.Y-p.Radius, 0, -13, true, 15))
		projectiles = append(projectiles, NewProjectile(p.X-20, p.Y-p.Radius+5, -2, -11, true, 12))
		projectiles = append(projectiles, NewProjectile(p.X+20, p.Y-p.Radius+5, 2, -11, true, 12))
	default: // Level 5+
		projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -14, true, 20))
		projectiles = append(projectiles, NewProjectile(p.X-12, p.Y-p.Radius, 0, -13, true, 15))
		projectiles = append(projectiles, NewProjectile(p.X+12, p.Y-p.Radius, 0, -13, true, 15))
		projectiles = append(projectiles, NewProjectile(p.X-25, p.Y-p.Radius+5, -2.5, -11, true, 12))
		projectiles = append(projectiles, NewProjectile(p.X+25, p.Y-p.Radius+5, 2.5, -11, true, 12))
	}

	return projectiles
}

func (p *Player) TakeDamage(damage int, gameTime float64) {
	if p.InvincTimer > 0 {
		return
	}

	// Shield absorbs damage first
	if p.Shield > 0 {
		absorbed := min(p.Shield, damage)
		p.Shield -= absorbed
		damage -= absorbed
	}

	p.Health -= damage
	p.LastDamageTime = gameTime
	p.InvincTimer = p.InvincibilityTime // Use difficulty-dependent invincibility time
}

func (p *Player) ApplyPowerUp(puType PowerUpType) {
	switch puType {
	case PowerUpHealth:
		p.Health = min(p.Health+30, p.MaxHealth)
	case PowerUpShield:
		p.Shield = min(p.Shield+30, p.MaxShield)
	case PowerUpWeapon:
		if p.WeaponLevel < 5 {
			p.WeaponLevel++
		}
	case PowerUpSpeed:
		p.Speed = math.Min(p.Speed+0.5, 10)
	}
}

// ActivateUltimate triggers the ultimate ability
func (p *Player) ActivateUltimate() bool {
	if p.UltimateCharge >= p.MaxUltimateCharge {
		p.UltimateActive = true
		p.UltimateTimer = 3.0 // 3 seconds of ultimate effect
		p.UltimateCharge = 0  // Reset charge
		return true
	}
	return false
}

// GetChargedProjectiles returns projectiles with enhanced power based on charge level
func (p *Player) GetChargedProjectiles() []*Projectile {
	if p.FireCooldown > 0 {
		return nil
	}

	var projectiles []*Projectile
	chargeMultiplier := 1.0 + p.ChargeLevel*2.0 // Damage multiplier from charge
	baseDamage := int(float64(10) * chargeMultiplier)

	if p.ChargeLevel > 0.3 { // Only fire charged if significantly charged
		// Charged shot (slower but more powerful)
		p.FireCooldown = p.FireRate * 1.5 // Longer cooldown for charged shots

		// Single powerful shot at full charge, or multiple weaker shots at lower charge
		if p.ChargeLevel > 0.8 {
			// Full charge - massive central shot
			projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -15, true, baseDamage))
		} else if p.ChargeLevel > 0.5 {
			// Medium charge - 3 shots
			projectiles = append(projectiles, NewProjectile(p.X-10, p.Y-p.Radius, 0, -14, true, baseDamage-5))
			projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -15, true, baseDamage))
			projectiles = append(projectiles, NewProjectile(p.X+10, p.Y-p.Radius, 0, -14, true, baseDamage-5))
		} else {
			// Light charge - standard spread
			projectiles = append(projectiles, NewProjectile(p.X-8, p.Y-p.Radius, 0, -13, true, baseDamage-2))
			projectiles = append(projectiles, NewProjectile(p.X+8, p.Y-p.Radius, 0, -13, true, baseDamage-2))
		}
	} else {
		// Regular shooting if not charged enough
		p.FireCooldown = p.FireRate

		switch p.WeaponLevel {
		case 1:
			projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -12, true, 10))
		case 2:
			projectiles = append(projectiles, NewProjectile(p.X-10, p.Y-p.Radius, 0, -12, true, 10))
			projectiles = append(projectiles, NewProjectile(p.X+10, p.Y-p.Radius, 0, -12, true, 10))
		case 3:
			projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -12, true, 12))
			projectiles = append(projectiles, NewProjectile(p.X-15, p.Y-p.Radius+5, -1, -11, true, 10))
			projectiles = append(projectiles, NewProjectile(p.X+15, p.Y-p.Radius+5, 1, -11, true, 10))
		case 4:
			projectiles = append(projectiles, NewProjectile(p.X-8, p.Y-p.Radius, 0, -13, true, 15))
			projectiles = append(projectiles, NewProjectile(p.X+8, p.Y-p.Radius, 0, -13, true, 15))
			projectiles = append(projectiles, NewProjectile(p.X-20, p.Y-p.Radius+5, -2, -11, true, 12))
			projectiles = append(projectiles, NewProjectile(p.X+20, p.Y-p.Radius+5, 2, -11, true, 12))
		default: // Level 5+
			projectiles = append(projectiles, NewProjectile(p.X, p.Y-p.Radius, 0, -14, true, 20))
			projectiles = append(projectiles, NewProjectile(p.X-12, p.Y-p.Radius, 0, -13, true, 15))
			projectiles = append(projectiles, NewProjectile(p.X+12, p.Y-p.Radius, 0, -13, true, 15))
			projectiles = append(projectiles, NewProjectile(p.X-25, p.Y-p.Radius+5, -2.5, -11, true, 12))
			projectiles = append(projectiles, NewProjectile(p.X+25, p.Y-p.Radius+5, 2.5, -11, true, 12))
		}
	}

	p.ChargeLevel = 0 // Reset charge after firing
	return projectiles
}

func (p *Player) Draw(screen *ebiten.Image, shakeX, shakeY float64) {
	// Blink when invincible
	if p.InvincTimer > 0 && int(p.InvincTimer*10)%2 == 0 {
		return
	}

	// Simple screen coordinates with shake
	x := float32(p.X + shakeX)
	y := float32(p.Y + shakeY)

	// Draw thruster trail particles
	for _, trail := range p.ThrusterTrail {
		lifeRatio := trail.Life / 0.5
		alpha := uint8(150 * lifeRatio)
		trailColor := color.RGBA{100, 150, 255, alpha}
		trailGlowColor := color.RGBA{150, 200, 255, uint8(100 * lifeRatio)}

		trailX := float32(trail.X + shakeX)
		trailY := float32(trail.Y + shakeY)
		trailSize := float32(3 * (1 - (1-lifeRatio)*(1-lifeRatio)))

		// Draw glow
		vector.DrawFilledCircle(screen, trailX, trailY, trailSize*1.8, trailGlowColor, true)
		// Draw particle
		vector.DrawFilledCircle(screen, trailX, trailY, trailSize, trailColor, true)
	}

	// Draw shadow beneath ship (depth indicator)
	shadowColor := color.RGBA{20, 20, 30, 100}
	shadowSize := float32(p.Radius) * 0.5
	vector.DrawFilledCircle(screen, x, y+float32(p.Radius)+shadowSize, shadowSize, shadowColor, true)

	// Draw polygon-based ship - triangular arrow shape
	shipColor := color.RGBA{100, 160, 220, 255}
	radius := float32(p.Radius)

	// Main hull - triangle pointing up
	// Top vertex (nose)
	topX := x
	topY := y - radius*1.1

	// Bottom-left wing
	leftX := x - radius*0.8
	leftY := y + radius*0.7

	// Bottom-right wing
	rightX := x + radius*0.8
	rightY := y + radius*0.7

	// Draw main hull as filled polygon
	// Using triangles (vector.DrawFilledRect for simpler polygon effects)
	drawTriangle(screen, topX, topY, leftX, leftY, rightX, rightY, shipColor)

	// Draw cockpit window (diamond shape)
	cockpitColor := color.RGBA{220, 240, 255, 255}
	vector.DrawFilledCircle(screen, x, y-radius*0.4, 5, cockpitColor, true)

	// Draw engine thrusters (animated flame effect)
	engineIntensity := float32(0.5 + 0.3*math.Sin(p.EngineGlow))
	engineTrailColor1 := color.RGBA{100, 150, 255, 200}
	engineTrailColor2 := color.RGBA{255, 150, 100, 180}

	// Left engine
	vector.DrawFilledCircle(screen, leftX+radius*0.3, leftY+radius*0.5, float32(6)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, leftX+radius*0.3, leftY+radius*0.8, float32(3)*engineIntensity, engineTrailColor2, true)

	// Right engine
	vector.DrawFilledCircle(screen, rightX-radius*0.3, rightY+radius*0.5, float32(6)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, rightX-radius*0.3, rightY+radius*0.8, float32(3)*engineIntensity, engineTrailColor2, true)

	// Center engine
	vector.DrawFilledCircle(screen, x, y+radius*0.6, float32(5)*engineIntensity, engineTrailColor1, true)
	vector.DrawFilledCircle(screen, x, y+radius*1.0, float32(2)*engineIntensity, engineTrailColor2, true)

	// Draw wing accents
	wingAccentColor := color.RGBA{80, 140, 200, 255}
	vector.StrokeCircle(screen, leftX, leftY, radius*0.3, 1.5, wingAccentColor, true)
	vector.StrokeCircle(screen, rightX, rightY, radius*0.3, 1.5, wingAccentColor, true)

	// Shield effect
	if p.Shield > 0 {
		shieldAlpha := uint8(50 + float64(p.Shield)/float64(p.MaxShield)*100)
		shieldColor := color.RGBA{100, 200, 255, shieldAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.0, 2, shieldColor, true)
	}

	// Nose highlight for depth
	highlightColor := color.RGBA{200, 230, 255, 200}
	vector.DrawFilledCircle(screen, topX, topY+radius*0.2, 4, highlightColor, true)

	// Draw charge indicator (glow around ship when charging)
	if p.ChargeLevel > 0.1 {
		chargeAlpha := uint8(100 + float64(p.ChargeLevel)*155)
		chargeColor := color.RGBA{255, uint8(100 + float64(p.ChargeLevel)*155), 100, chargeAlpha}
		chargeRadius := radius*1.2 + float32(p.ChargeLevel)*10
		vector.StrokeCircle(screen, x, y, chargeRadius, 2, chargeColor, true)

		// Inner charge aura
		innerChargeColor := color.RGBA{255, 200, 100, uint8(50 + float64(p.ChargeLevel)*100)}
		vector.DrawFilledCircle(screen, x, y, radius*1.1, innerChargeColor, true)
	}

	// Draw ultimate indicator (star glow)
	if p.UltimateCharge > 0.5 {
		ultimateAlpha := uint8(100 + float64(p.UltimateCharge)*155)
		ultimateColor := color.RGBA{255, 50 + uint8(float64(p.UltimateCharge)*200), 255, ultimateAlpha}
		ultimateRadius := radius + float32(p.UltimateCharge)*8
		vector.StrokeCircle(screen, x, y, ultimateRadius, 3, ultimateColor, true)
	}

	// Ultimate active effect - intense glow
	if p.UltimateActive {
		pulseIntensity := float32(0.5 + 0.5*math.Sin(p.UltimateTimer*math.Pi))
		ultimateGlowColor := color.RGBA{200, 100, 255, uint8(200 * pulseIntensity)}
		vector.DrawFilledCircle(screen, x, y, radius*1.3*pulseIntensity, ultimateGlowColor, true)
	}
}

// Helper function to draw a filled triangle
func drawTriangle(screen *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, col color.Color) {
	// Draw using three filled circles connected - crude but effective
	// For better polygon support, we can use vector paths or draw filled polygons another way

	// Create paths for triangle edges with thick lines to simulate fill
	// This creates a filled triangle effect
	const steps = 20
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)
		// Top to left edge
		edgeX1 := x1*(1-t) + x2*t
		edgeY1 := y1*(1-t) + y2*t

		// Top to right edge
		edgeX2 := x1*(1-t) + x3*t
		edgeY2 := y1*(1-t) + y3*t

		// Draw line between edges
		radius := float32(math.Abs(float64(edgeX2-edgeX1)) / 2)
		if radius > 0.5 {
			midX := (edgeX1 + edgeX2) / 2
			midY := (edgeY1 + edgeY2) / 2
			vector.DrawFilledCircle(screen, midX, midY, radius+1, col, true)
		}
	}

	// Draw bottom edge
	for i := 0; i < steps; i++ {
		t := float32(i) / float32(steps)
		edgeX1 := x2*(1-t) + x3*t
		edgeY1 := y2*(1-t) + y3*t
		vector.DrawFilledCircle(screen, edgeX1, edgeY1, 2, col, true)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
