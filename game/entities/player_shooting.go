package entities

import (
	"math"
)

// Shoot generates projectiles based on the current weapon
func (p *Player) Shoot() []*Projectile {
	// Check if weapon can fire
	if !p.WeaponMgr.CanFireWeapon() {
		return nil
	}

	weapon := p.WeaponMgr.GetCurrentWeapon()
	if weapon == nil {
		return nil
	}

	// Apply fire rate modifiers from mystery power-ups
	fireRateMultiplier := 1.0
	if p.RapidFireTimer > 0 {
		fireRateMultiplier = p.RapidFireMultiplier
	} else if p.SlowFireTimer > 0 {
		fireRateMultiplier = p.SlowFireMultiplier
	}

	// Temporarily adjust weapon fire rate
	originalFireRate := weapon.FireRate
	weapon.FireRate *= fireRateMultiplier

	// Fire weapon
	p.WeaponMgr.FireWeapon()

	// Restore original fire rate
	weapon.FireRate = originalFireRate

	// Generate projectiles based on weapon type
	return p.createProjectilesForWeapon(weapon)
}

// createProjectilesForWeapon generates projectiles based on weapon type and level
func (p *Player) createProjectilesForWeapon(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	// Check if we should use mixed mode (special weapon + side blasters)
	useMixedMode := p.WeaponMgr.ShouldUseMixedMode()

	// Create the main weapon projectiles
	switch weapon.Type {
	case WeaponTypeFollowingRocket:
		projectiles = p.createFollowingRockets(weapon)
	case WeaponTypeChainLightning:
		projectiles = p.createChainLightning(weapon)
	case WeaponTypeFlamethrower:
		projectiles = p.createFlamethrower(weapon)
	case WeaponTypeIonBeam:
		projectiles = p.createIonBeam(weapon)
	default:
		projectiles = p.createStandardProjectiles(weapon)
	}

	// Add side blasters if in mixed mode
	if useMixedMode {
		sideBlasters := p.createSideBlasters()
		projectiles = append(projectiles, sideBlasters...)
	}

	return projectiles
}

// sideBlasterCounter tracks which side to fire next (alternates each shot)
var sideBlasterCounter int

// createSideBlasters creates 1 angled shot (alternating sides) - optimized for performance
// OPTIMIZED: Reduced from 2 simultaneous shots to 1 alternating shot (-50% projectiles)
func (p *Player) createSideBlasters() []*Projectile {
	var projectiles []*Projectile

	basicGun := p.WeaponMgr.GetBasicGun()
	if basicGun == nil {
		return projectiles
	}

	// Create 1 angled shot, alternating between left and right each shot
	spread := 0.2 // Base spread angle

	// Alternate between left (-1) and right (+1)
	side := float64(-1)
	if sideBlasterCounter%2 == 1 {
		side = 1
	}
	sideBlasterCounter++

	spreadAngle := side * spread * 2.0
	angle := -math.Pi/2 + spreadAngle
	velX := math.Cos(angle) * basicGun.ProjectileSpeed / 60.0
	velY := math.Sin(angle) * basicGun.ProjectileSpeed / 60.0

	proj := NewProjectileWithColor(
		p.X,
		p.Y-p.Radius,
		velX,
		velY,
		true,
		int(basicGun.Damage),
		basicGun.Color,
		basicGun.GlowColor,
	)
	projectiles = append(projectiles, proj)

	return projectiles
}

// createStandardProjectiles creates standard projectiles for most weapons
func (p *Player) createStandardProjectiles(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	// Special patterns for the basic gun progression
	// All patterns ensure at least one projectile goes straight forward
	if weapon.Type == WeaponTypeSpread {
		switch weapon.ProjectileCount {
		case 1:
			// MkI and MkII: Single shot straight forward
			proj := NewProjectileWithColor(
				p.X,
				p.Y-p.Radius,
				0,
				-weapon.ProjectileSpeed/60.0,
				true,
				int(weapon.Damage),
				weapon.Color,
				weapon.GlowColor,
			)
			projectiles = append(projectiles, proj)
			return projectiles

		case 3:
			// MkIII and MkIV: 1 center forward + 2 angled (left and right)
			// Center shot (always straight forward)
			proj := NewProjectileWithColor(
				p.X,
				p.Y-p.Radius,
				0,
				-weapon.ProjectileSpeed/60.0,
				true,
				int(weapon.Damage),
				weapon.Color,
				weapon.GlowColor,
			)
			projectiles = append(projectiles, proj)

			// Two angled shots (left and right)
			angles := []float64{-weapon.Spread * 2.0, weapon.Spread * 2.0}
			for _, spreadAngle := range angles {
				angle := -math.Pi/2 + spreadAngle
				velX := math.Cos(angle) * weapon.ProjectileSpeed / 60.0
				velY := math.Sin(angle) * weapon.ProjectileSpeed / 60.0

				proj := NewProjectileWithColor(
					p.X,
					p.Y-p.Radius,
					velX,
					velY,
					true,
					int(weapon.Damage),
					weapon.Color,
					weapon.GlowColor,
				)
				projectiles = append(projectiles, proj)
			}
			return projectiles

		case 5:
			// MkV: 1 center forward + 4 angled (2 left, 2 right)
			// Center shot (always straight forward)
			proj := NewProjectileWithColor(
				p.X,
				p.Y-p.Radius,
				0,
				-weapon.ProjectileSpeed/60.0,
				true,
				int(weapon.Damage),
				weapon.Color,
				weapon.GlowColor,
			)
			projectiles = append(projectiles, proj)

			// Four angled shots (2 left, 2 right at different angles)
			angles := []float64{
				-weapon.Spread * 1.5, -weapon.Spread * 3.0,
				weapon.Spread * 1.5, weapon.Spread * 3.0,
			}
			for _, spreadAngle := range angles {
				angle := -math.Pi/2 + spreadAngle
				velX := math.Cos(angle) * weapon.ProjectileSpeed / 60.0
				velY := math.Sin(angle) * weapon.ProjectileSpeed / 60.0

				proj := NewProjectileWithColor(
					p.X,
					p.Y-p.Radius,
					velX,
					velY,
					true,
					int(weapon.Damage),
					weapon.Color,
					weapon.GlowColor,
				)
				projectiles = append(projectiles, proj)
			}
			return projectiles
		}
	}

	// Standard spread pattern for other weapons
	// Ensure at least one projectile goes forward by centering the spread
	for i := 0; i < weapon.ProjectileCount; i++ {
		// Calculate spread angle - centered so middle projectile(s) go forward
		spreadAngle := 0.0
		if weapon.ProjectileCount > 1 {
			// Distribute projectiles evenly across spread, centered on forward
			spreadAngle = weapon.Spread * (float64(i) - float64(weapon.ProjectileCount-1)/2.0)
		}

		// Calculate velocity with spread
		angle := -math.Pi/2 + spreadAngle // -90 degrees (up) + spread
		velX := math.Cos(angle) * weapon.ProjectileSpeed / 60.0
		velY := math.Sin(angle) * weapon.ProjectileSpeed / 60.0

		// Create projectile
		proj := NewProjectileWithColor(
			p.X,
			p.Y-p.Radius,
			velX,
			velY,
			true,
			int(weapon.Damage),
			weapon.Color,
			weapon.GlowColor,
		)

		projectiles = append(projectiles, proj)
	}

	return projectiles
}

// createFollowingRockets creates homing missiles
func (p *Player) createFollowingRockets(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	for i := 0; i < weapon.ProjectileCount; i++ {
		// Calculate spread angle
		spreadAngle := 0.0
		if weapon.ProjectileCount > 1 {
			spreadAngle = weapon.Spread * (float64(i) - float64(weapon.ProjectileCount-1)/2.0)
		}

		// Calculate initial velocity with spread
		angle := -math.Pi/2 + spreadAngle
		velX := math.Cos(angle) * weapon.ProjectileSpeed / 60.0
		velY := math.Sin(angle) * weapon.ProjectileSpeed / 60.0

		proj := NewProjectileWithColor(
			p.X,
			p.Y-p.Radius,
			velX,
			velY,
			true,
			int(weapon.Damage),
			weapon.Color,
			weapon.GlowColor,
		)

		// Enable homing behavior
		proj.Homing = true
		proj.HomingSpeed = 0.08 // Turn rate in radians per frame

		projectiles = append(projectiles, proj)
	}

	return projectiles
}

// createChainLightning creates lightning bolts that chain between enemies
func (p *Player) createChainLightning(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	// Single lightning bolt (chains on hit)
	velY := -weapon.ProjectileSpeed / 60.0

	proj := NewProjectileWithColor(
		p.X,
		p.Y-p.Radius,
		0,
		velY,
		true,
		int(weapon.Damage),
		weapon.Color,
		weapon.GlowColor,
	)

	// Enable chaining behavior
	proj.Chaining = true
	proj.ChainCount = 3     // Can chain up to 3 times
	proj.ChainRange = 150.0 // Range to find next target
	proj.Radius = 10        // Larger hitbox for lightning
	// Trail uses fixed-size ring buffer, no need to initialize

	projectiles = append(projectiles, proj)

	return projectiles
}

// createFlamethrower creates short-range flame projectiles with DoT
func (p *Player) createFlamethrower(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	// Always include a center forward flame
	// First create the center flame (straight forward)
	centerProj := NewProjectileWithColor(
		p.X,
		p.Y-p.Radius,
		0,
		-weapon.ProjectileSpeed/60.0,
		true,
		int(weapon.Damage),
		weapon.Color,
		weapon.GlowColor,
	)
	centerProj.Burning = true
	centerProj.BurnDuration = 3.0
	centerProj.BurnDamage = 5
	projectiles = append(projectiles, centerProj)

	// Add side flames if projectile count > 1
	sideCount := weapon.ProjectileCount - 1
	for i := 0; i < sideCount; i++ {
		// Alternate left and right for even distribution
		side := float64(1)
		if i%2 == 0 {
			side = -1
		}
		spreadMultiplier := float64((i/2)+1) * 0.5

		spreadAngle := side * weapon.Spread * spreadMultiplier

		// Calculate velocity with spread
		angle := -math.Pi/2 + spreadAngle
		velX := math.Cos(angle) * weapon.ProjectileSpeed / 60.0
		velY := math.Sin(angle) * weapon.ProjectileSpeed / 60.0

		proj := NewProjectileWithColor(
			p.X,
			p.Y-p.Radius,
			velX,
			velY,
			true,
			int(weapon.Damage),
			weapon.Color,
			weapon.GlowColor,
		)

		// Enable burning DoT
		proj.Burning = true
		proj.BurnDuration = 3.0 // 3 seconds of burn
		proj.BurnDamage = 5     // 5 damage per tick (every 0.5s)

		projectiles = append(projectiles, proj)
	}

	return projectiles
}

// createIonBeam creates continuous beam projectiles
func (p *Player) createIonBeam(weapon *Weapon) []*Projectile {
	var projectiles []*Projectile

	// Single beam projectile
	velY := -weapon.ProjectileSpeed / 60.0

	proj := NewProjectileWithColor(
		p.X,
		p.Y-p.Radius,
		0,
		velY,
		true,
		int(weapon.Damage),
		weapon.Color,
		weapon.GlowColor,
	)

	// Enable beam behavior
	proj.Beam = true
	proj.BeamSource = struct{ X, Y float64 }{p.X, p.Y - p.Radius}
	proj.Piercing = true // Beam penetrates enemies
	proj.Radius = 8      // Thicker beam hitbox

	projectiles = append(projectiles, proj)

	return projectiles
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
// DEPRECATED: This method is from the legacy charged shot system and is no longer used.
// Kept for backward compatibility but should be removed in a future refactor.
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
