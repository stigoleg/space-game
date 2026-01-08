package entities

import "image/color"

// WeaponType represents different weapon types
type WeaponType string

const (
	WeaponTypeSpread      WeaponType = "spread"
	WeaponTypeLaser       WeaponType = "laser"
	WeaponTypeShotgun     WeaponType = "shotgun"
	WeaponTypePlasma      WeaponType = "plasma"
	WeaponTypeHoming      WeaponType = "homing"
	WeaponTypeRailgun     WeaponType = "railgun"
	WeaponTypeEnergyLance WeaponType = "energy_lance"
	WeaponTypePulse       WeaponType = "pulse"
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

// NewWeaponManager creates a new weapon manager
func NewWeaponManager() *WeaponManager {
	wm := &WeaponManager{
		Weapons:       make(map[WeaponType]*Weapon),
		CurrentWeapon: WeaponTypeSpread,
	}

	// Initialize base weapons
	wm.Weapons[WeaponTypeSpread] = &Weapon{
		Type:            WeaponTypeSpread,
		Level:           WeaponLevelMkI,
		Name:            "Spread Shot",
		Description:     "Standard multi-directional fire",
		IconEmoji:       "💥",
		Damage:          20,
		FireRate:        6.0,
		ProjectileSpeed: 300,
		Spread:          0.2, // 0.2 radians ≈ 11 degrees
		ProjectileCount: 3,
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
			Damage:          35,
			FireRate:        3.0,
			ProjectileSpeed: 250,
			Spread:          0.8,
			ProjectileCount: 8,
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
			Damage:          40,
			FireRate:        4.0,
			ProjectileSpeed: 280,
			Spread:          0.3,
			ProjectileCount: 3,
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
			Damage:          45,
			FireRate:        3.5,
			ProjectileSpeed: 200,
			Spread:          0.2,
			ProjectileCount: 2,
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
			Damage:          15,
			FireRate:        12.0,
			ProjectileSpeed: 320,
			Spread:          0.1,
			ProjectileCount: 2,
			Unlocked:        true,
			Color:           color.RGBA{255, 0, 255, 255},   // Pink/Magenta
			GlowColor:       color.RGBA{255, 100, 255, 180}, // Pink glow
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

			// Apply upgrade bonuses (more gradual scaling)
			weapon.Damage *= 1.15          // +15% damage per level
			weapon.FireRate *= 1.08        // +8% fire rate per level
			weapon.ProjectileSpeed *= 1.05 // +5% speed per level

			// Special bonuses at level 4 - add extra projectiles
			if weapon.Level == WeaponLevelMkIV {
				switch weapon.Type {
				case WeaponTypeSpread:
					weapon.ProjectileCount = 4 // 3 -> 4
				case WeaponTypeShotgun:
					weapon.ProjectileCount = 10 // 8 -> 10
				case WeaponTypePlasma:
					weapon.ProjectileCount = 4 // 3 -> 4
				case WeaponTypePulse:
					weapon.ProjectileCount = 3 // 2 -> 3
				}
			}

			// Special bonuses at level 5 - even more projectiles
			if weapon.Level == WeaponLevelMkV {
				switch weapon.Type {
				case WeaponTypeSpread:
					weapon.ProjectileCount = 5 // 4 -> 5
				case WeaponTypeShotgun:
					weapon.ProjectileCount = 12 // 10 -> 12
				case WeaponTypePlasma:
					weapon.ProjectileCount = 5 // 4 -> 5
				case WeaponTypePulse:
					weapon.ProjectileCount = 4 // 3 -> 4
				case WeaponTypeHoming:
					weapon.ProjectileCount = 3 // 2 -> 3
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
