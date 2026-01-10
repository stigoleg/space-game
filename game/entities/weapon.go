package entities

import "image/color"

// WeaponType represents different weapon types
type WeaponType string

const (
	WeaponTypeSpread          WeaponType = "spread"
	WeaponTypeLaser           WeaponType = "laser"
	WeaponTypeShotgun         WeaponType = "shotgun"
	WeaponTypePlasma          WeaponType = "plasma"
	WeaponTypeHoming          WeaponType = "homing"
	WeaponTypeRailgun         WeaponType = "railgun"
	WeaponTypeEnergyLance     WeaponType = "energy_lance"
	WeaponTypePulse           WeaponType = "pulse"
	WeaponTypeBlaster         WeaponType = "blaster"
	WeaponTypeFollowingRocket WeaponType = "following_rocket"
	WeaponTypeChainLightning  WeaponType = "chain_lightning"
	WeaponTypeFlamethrower    WeaponType = "flamethrower"
	WeaponTypeIonBeam         WeaponType = "ion_beam"
)

// WeaponLevel represents weapon upgrade level (1-5)
type WeaponLevel int

const (
	WeaponLevelMkI   WeaponLevel = 1
	WeaponLevelMkII  WeaponLevel = 2
	WeaponLevelMkIII WeaponLevel = 3
	WeaponLevelMkIV  WeaponLevel = 4
	WeaponLevelMkV   WeaponLevel = 5
)

// Weapon represents a player weapon
type Weapon struct {
	Type            WeaponType
	Level           WeaponLevel
	Name            string
	Description     string
	IconEmoji       string
	Damage          float64
	FireRate        float64 // Shots per second
	FireTimer       float64 // Current cooldown
	ProjectileSpeed float64
	Spread          float64 // Angle spread in radians
	ProjectileCount int     // Number of projectiles per shot
	Unlocked        bool
	Color           color.RGBA // Main projectile color
	GlowColor       color.RGBA // Glow/trail color
}

// WeaponManager manages player weapons
type WeaponManager struct {
	Weapons       map[WeaponType]*Weapon
	CurrentWeapon WeaponType
}

// ShouldUseMixedMode returns true if special weapons should be mixed with side blasters
func (wm *WeaponManager) ShouldUseMixedMode() bool {
	// Mixed mode activates when:
	// 1. Basic gun is at max level (MkV = 8 projectiles)
	// 2. Current weapon is NOT the basic gun
	basicGun := wm.Weapons[WeaponTypeSpread]
	if basicGun == nil {
		return false
	}
	return basicGun.Level == WeaponLevelMkV && wm.CurrentWeapon != WeaponTypeSpread
}

// GetBasicGun returns the basic spread weapon
func (wm *WeaponManager) GetBasicGun() *Weapon {
	return wm.Weapons[WeaponTypeSpread]
}

// NewWeaponManager creates a new weapon manager
func NewWeaponManager() *WeaponManager {
	wm := &WeaponManager{
		Weapons:       make(map[WeaponType]*Weapon),
		CurrentWeapon: WeaponTypeSpread,
	}

	// Initialize base weapons - starts with single shot
	wm.Weapons[WeaponTypeSpread] = &Weapon{
		Type:            WeaponTypeSpread,
		Level:           WeaponLevelMkI,
		Name:            "Basic Gun",
		Description:     "Single forward shot",
		IconEmoji:       "💥",
		Damage:          30,  // Increased from 25 (+20% to compensate for fire rate reduction)
		FireRate:        4.5, // Reduced from 5.5 for performance (-18%)
		ProjectileSpeed: 300,
		Spread:          0.2, // 0.2 radians ≈ 11 degrees
		ProjectileCount: 1,   // Start with single shot
		Unlocked:        true,
		Color:           color.RGBA{0, 255, 255, 255}, // Cyan
		GlowColor:       color.RGBA{0, 200, 255, 180}, // Cyan glow
	}

	return wm
}

// AddWeapon adds a new weapon to the arsenal
func (wm *WeaponManager) AddWeapon(weaponType WeaponType) bool {
	if _, exists := wm.Weapons[weaponType]; exists {
		return false // Already have this weapon
	}

	var weapon *Weapon
	switch weaponType {
	case WeaponTypeLaser:
		weapon = &Weapon{
			Type:            WeaponTypeLaser,
			Level:           WeaponLevelMkI,
			Name:            "Laser Rifle",
			Description:     "Continuous beam, high damage",
			IconEmoji:       "🔴",
			Damage:          25,
			FireRate:        8.0,
			ProjectileSpeed: 400,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{255, 0, 0, 255},   // Red
			GlowColor:       color.RGBA{255, 50, 50, 180}, // Red glow
		}
	case WeaponTypeShotgun:
		weapon = &Weapon{
			Type:            WeaponTypeShotgun,
			Level:           WeaponLevelMkI,
			Name:            "Shotgun",
			Description:     "Wide spread, close range",
			IconEmoji:       "🔥",
			Damage:          45,  // Increased from 35 (+30% to compensate for fire rate)
			FireRate:        2.5, // Reduced from 3.0 for performance (-17%)
			ProjectileSpeed: 250,
			Spread:          0.8,
			ProjectileCount: 5, // Odd number ensures center shot goes forward
			Unlocked:        true,
			Color:           color.RGBA{255, 136, 0, 255},  // Orange
			GlowColor:       color.RGBA{255, 180, 50, 180}, // Orange glow
		}
	case WeaponTypePlasma:
		weapon = &Weapon{
			Type:            WeaponTypePlasma,
			Level:           WeaponLevelMkI,
			Name:            "Plasma Burst",
			Description:     "Explosive projectiles with splash",
			IconEmoji:       "⚡",
			Damage:          50,  // Increased from 40 (+25% to compensate for fire rate)
			FireRate:        3.5, // Reduced from 4.0 for performance (-12%)
			ProjectileSpeed: 280,
			Spread:          0.3,
			ProjectileCount: 3, // Odd number ensures center shot goes forward
			Unlocked:        true,
			Color:           color.RGBA{0, 255, 136, 255},  // Green
			GlowColor:       color.RGBA{50, 255, 150, 180}, // Green glow
		}
	case WeaponTypeHoming:
		weapon = &Weapon{
			Type:            WeaponTypeHoming,
			Level:           WeaponLevelMkI,
			Name:            "Homing Missiles",
			Description:     "Track enemies automatically",
			IconEmoji:       "🚀",
			Damage:          50,
			FireRate:        1.8, // Further reduced from 2.5 to prevent lag
			ProjectileSpeed: 200,
			Spread:          0.2,
			ProjectileCount: 1, // Reduced from 2 to prevent lag
			Unlocked:        true,
			Color:           color.RGBA{255, 255, 0, 255},  // Yellow
			GlowColor:       color.RGBA{255, 220, 50, 180}, // Yellow glow
		}
	case WeaponTypeRailgun:
		weapon = &Weapon{
			Type:            WeaponTypeRailgun,
			Level:           WeaponLevelMkI,
			Name:            "Railgun",
			Description:     "Pierces through enemies",
			IconEmoji:       "🔵",
			Damage:          50,
			FireRate:        2.5,
			ProjectileSpeed: 500,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{170, 0, 255, 255},  // Purple
			GlowColor:       color.RGBA{200, 50, 255, 180}, // Purple glow
		}
	case WeaponTypeEnergyLance:
		weapon = &Weapon{
			Type:            WeaponTypeEnergyLance,
			Level:           WeaponLevelMkI,
			Name:            "Energy Lance",
			Description:     "Charges up for massive damage",
			IconEmoji:       "⚔️",
			Damage:          80,
			FireRate:        1.5,
			ProjectileSpeed: 350,
			Spread:          0.1,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{255, 255, 255, 255}, // White
			GlowColor:       color.RGBA{240, 240, 255, 200}, // White glow
		}
	case WeaponTypePulse:
		weapon = &Weapon{
			Type:            WeaponTypePulse,
			Level:           WeaponLevelMkI,
			Name:            "Pulse Cannon",
			Description:     "Rapid burst fire",
			IconEmoji:       "💫",
			Damage:          28,  // Increased from 21 (+33% to compensate for fire rate)
			FireRate:        5.0, // Reduced from 6.0 for performance (-17%)
			ProjectileSpeed: 320,
			Spread:          0.1,
			ProjectileCount: 1, // Reduced from 2 for performance (-50%)
			Unlocked:        true,
			Color:           color.RGBA{255, 0, 255, 255},   // Pink/Magenta
			GlowColor:       color.RGBA{255, 100, 255, 180}, // Pink glow
		}
	case WeaponTypeBlaster:
		weapon = &Weapon{
			Type:            WeaponTypeBlaster,
			Level:           WeaponLevelMkI,
			Name:            "Blaster",
			Description:     "High-damage single shot",
			IconEmoji:       "🔶",
			Damage:          60,
			FireRate:        2.0,
			ProjectileSpeed: 350,
			Spread:          0.05,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{255, 100, 0, 255},  // Bright Orange
			GlowColor:       color.RGBA{255, 150, 50, 180}, // Orange glow
		}
	case WeaponTypeFollowingRocket:
		weapon = &Weapon{
			Type:            WeaponTypeFollowingRocket,
			Level:           WeaponLevelMkI,
			Name:            "Following Rockets",
			Description:     "Smart missiles that track enemies",
			IconEmoji:       "🚀",
			Damage:          40,  // Increased from 35
			FireRate:        1.8, // Further reduced from 2.5 to prevent lag
			ProjectileSpeed: 180,
			Spread:          0.15,
			ProjectileCount: 1, // Reduced from 2 to prevent lag
			Unlocked:        true,
			Color:           color.RGBA{255, 200, 0, 255},   // Yellow/Orange
			GlowColor:       color.RGBA{255, 180, 100, 180}, // Warm glow
		}
	case WeaponTypeChainLightning:
		weapon = &Weapon{
			Type:            WeaponTypeChainLightning,
			Level:           WeaponLevelMkI,
			Name:            "Chain Lightning",
			Description:     "Electric bolts that chain enemies",
			IconEmoji:       "⚡",
			Damage:          25,
			FireRate:        4.0,
			ProjectileSpeed: 400,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{100, 200, 255, 255}, // Electric Blue
			GlowColor:       color.RGBA{200, 230, 255, 200}, // White-blue glow
		}
	case WeaponTypeFlamethrower:
		weapon = &Weapon{
			Type:            WeaponTypeFlamethrower,
			Level:           WeaponLevelMkI,
			Name:            "Flamethrower",
			Description:     "Short range flame stream",
			IconEmoji:       "🔥",
			Damage:          18,  // Increased from 14 (+30% to compensate for fire rate)
			FireRate:        3.5, // Reduced from 4.5 for performance (-22%)
			ProjectileSpeed: 200,
			Spread:          0.6,
			ProjectileCount: 2, // Reduced from 3 for performance (-33%)
			Unlocked:        true,
			Color:           color.RGBA{255, 100, 0, 255},  // Red/Orange
			GlowColor:       color.RGBA{255, 200, 50, 180}, // Yellow glow
		}
	case WeaponTypeIonBeam:
		weapon = &Weapon{
			Type:            WeaponTypeIonBeam,
			Level:           WeaponLevelMkI,
			Name:            "Ion Beam",
			Description:     "Continuous penetrating beam",
			IconEmoji:       "🌟",
			Damage:          12,  // Increased from 10 (+15% to compensate for fire rate)
			FireRate:        6.0, // Reduced from 8.0 for performance (-25%)
			ProjectileSpeed: 600,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
			Color:           color.RGBA{0, 255, 255, 255},   // Cyan
			GlowColor:       color.RGBA{150, 255, 255, 200}, // Bright cyan glow
		}
	default:
		return false
	}

	wm.Weapons[weaponType] = weapon
	return true
}

// CanFireWeapon checks if the current weapon can fire
func (wm *WeaponManager) CanFireWeapon() bool {
	if weapon, exists := wm.Weapons[wm.CurrentWeapon]; exists {
		return weapon.FireTimer <= 0
	}
	return false
}

// FireWeapon fires the current weapon and resets cooldown
func (wm *WeaponManager) FireWeapon() bool {
	if !wm.CanFireWeapon() {
		return false
	}

	weapon := wm.Weapons[wm.CurrentWeapon]
	weapon.FireTimer = 1.0 / weapon.FireRate
	return true
}

// SwitchWeapon changes to a different weapon
func (wm *WeaponManager) SwitchWeapon(weaponType WeaponType) bool {
	if weapon, exists := wm.Weapons[weaponType]; exists && weapon.Unlocked {
		wm.CurrentWeapon = weaponType
		return true
	}
	return false
}

// UpgradeWeapon upgrades a weapon to next level (up to Mk V)
func (wm *WeaponManager) UpgradeWeapon(weaponType WeaponType) bool {
	if weapon, exists := wm.Weapons[weaponType]; exists {
		if weapon.Level < WeaponLevelMkV {
			weapon.Level++

			// Special handling for basic gun progression
			// Progression: Early levels boost fire rate, later levels add projectiles
			// This ensures progression feels rewarding (more shots over time)
			if weapon.Type == WeaponTypeSpread {
				switch weapon.Level {
				case WeaponLevelMkI:
					// Level 1: 1 blast straight forward (base stats)
					weapon.ProjectileCount = 1
					weapon.Name = "Basic Gun"
					weapon.Description = "Single forward shot"
				case WeaponLevelMkII:
					// Level 2: Same projectile, but faster fire rate
					weapon.ProjectileCount = 1
					weapon.Name = "Rapid Gun"
					weapon.Description = "Faster single shot"
					weapon.FireRate *= 1.25 // +25% fire rate
					weapon.Damage *= 1.1    // +10% damage
				case WeaponLevelMkIII:
					// Level 3: Spread unlocked - 1 center + 2 angled = 3 total
					weapon.ProjectileCount = 3
					weapon.Name = "Spread Shot"
					weapon.Description = "1 forward + 2 angled"
					weapon.Damage *= 1.15 // +15% damage
				case WeaponLevelMkIV:
					// Level 4: Same 3 projectiles, but wider spread and more damage
					weapon.ProjectileCount = 3
					weapon.Name = "Wide Spread"
					weapon.Description = "Wider angle coverage"
					weapon.Spread *= 1.3   // Wider spread angle
					weapon.Damage *= 1.2   // +20% damage
					weapon.FireRate *= 1.1 // +10% fire rate
				case WeaponLevelMkV:
					// Level 5: Full spread - 1 center + 4 angled = 5 total
					weapon.ProjectileCount = 5
					weapon.Name = "Maximum Spread"
					weapon.Description = "1 forward + 4 angled"
					weapon.Damage *= 1.25 // +25% damage
					// Brighter colors for max level
					weapon.Color.R = uint8(minInt(int(weapon.Color.R)+30, 255))
					weapon.Color.G = uint8(minInt(int(weapon.Color.G)+30, 255))
					weapon.Color.B = uint8(minInt(int(weapon.Color.B)+30, 255))
				}
				return true
			}

			// Apply upgrade bonuses for other weapons (more gradual scaling)
			weapon.Damage *= 1.15          // +15% damage per level
			weapon.FireRate *= 1.08        // +8% fire rate per level
			weapon.ProjectileSpeed *= 1.05 // +5% speed per level

			// Special bonuses at level 4 - add extra projectiles (but keep it balanced)
			// Use odd counts to ensure center shot goes forward
			if weapon.Level == WeaponLevelMkIV {
				switch weapon.Type {
				case WeaponTypeShotgun:
					weapon.ProjectileCount = 7 // Odd: 1 center + 3 each side
				case WeaponTypePlasma:
					weapon.ProjectileCount = 3 // Odd: 1 center + 1 each side
				case WeaponTypePulse:
					weapon.ProjectileCount = 1 // Keep single for performance
				case WeaponTypeFollowingRocket:
					weapon.ProjectileCount = 1 // Keep at 1 to prevent lag
				case WeaponTypeFlamethrower:
					weapon.ProjectileCount = 3 // Odd: 1 center + 1 each side
				}
			}

			// Special bonuses at level 5 - even more projectiles (but keep it reasonable)
			// Use odd counts to ensure center shot goes forward
			if weapon.Level == WeaponLevelMkV {
				switch weapon.Type {
				case WeaponTypeShotgun:
					weapon.ProjectileCount = 9 // Odd: 1 center + 4 each side
				case WeaponTypePlasma:
					weapon.ProjectileCount = 5 // Odd: 1 center + 2 each side
				case WeaponTypePulse:
					weapon.ProjectileCount = 3 // Odd: 1 center + 1 each side
				case WeaponTypeHoming:
					weapon.ProjectileCount = 1 // Keep single for performance
				case WeaponTypeFollowingRocket:
					weapon.ProjectileCount = 1 // Keep single for performance
				case WeaponTypeFlamethrower:
					weapon.ProjectileCount = 5 // Odd: 1 center + 2 each side
				}

				// Level 5 weapons get brighter, more intense colors
				weapon.Color.R = uint8(minInt(int(weapon.Color.R)+30, 255))
				weapon.Color.G = uint8(minInt(int(weapon.Color.G)+30, 255))
				weapon.Color.B = uint8(minInt(int(weapon.Color.B)+30, 255))
			}

			return true
		}
	}
	return false
}

// Helper function for color intensification
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Update updates weapon cooldowns
func (wm *WeaponManager) Update() {
	for _, weapon := range wm.Weapons {
		if weapon.FireTimer > 0 {
			weapon.FireTimer -= 1.0 / 60.0 // 60 FPS
		}
	}
}

// GetCurrentWeapon returns the currently equipped weapon
func (wm *WeaponManager) GetCurrentWeapon() *Weapon {
	if weapon, exists := wm.Weapons[wm.CurrentWeapon]; exists {
		return weapon
	}
	return nil
}

// GetWeapon returns a weapon by type
func (wm *WeaponManager) GetWeapon(weaponType WeaponType) *Weapon {
	return wm.Weapons[weaponType]
}

// GetUnlockedWeapons returns all unlocked weapons
func (wm *WeaponManager) GetUnlockedWeapons() []*Weapon {
	var weapons []*Weapon
	for _, weapon := range wm.Weapons {
		if weapon.Unlocked {
			weapons = append(weapons, weapon)
		}
	}
	return weapons
}

// HasWeapon checks if a weapon is unlocked
func (wm *WeaponManager) HasWeapon(weaponType WeaponType) bool {
	if weapon, exists := wm.Weapons[weaponType]; exists {
		return weapon.Unlocked
	}
	return false
}

// UnlockWeapon unlocks a weapon for use
func (wm *WeaponManager) UnlockWeapon(weaponType WeaponType) bool {
	if weapon, exists := wm.Weapons[weaponType]; exists {
		weapon.Unlocked = true
		return true
	}
	// Try to add weapon if it doesn't exist yet
	return wm.AddWeapon(weaponType)
}

// ApplyFireRateModifier applies a temporary fire rate multiplier
func (wm *WeaponManager) ApplyFireRateModifier(weaponType WeaponType, multiplier float64) {
	if weapon, exists := wm.Weapons[weaponType]; exists {
		weapon.FireRate *= multiplier
	}
}

// GetWeaponFireCooldown returns current fire cooldown (0.0-1.0)
func (wm *WeaponManager) GetWeaponFireCooldown() float64 {
	weapon := wm.GetCurrentWeapon()
	if weapon == nil || weapon.FireRate == 0 {
		return 0
	}
	maxCooldown := 1.0 / weapon.FireRate
	return weapon.FireTimer / maxCooldown
}
