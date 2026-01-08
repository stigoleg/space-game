package entities

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

// WeaponLevel represents weapon upgrade level (1-3)
type WeaponLevel int

const (
	WeaponLevelBasic    WeaponLevel = 1
	WeaponLevelAdvanced WeaponLevel = 2
	WeaponLevelMaster   WeaponLevel = 3
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
		Level:           WeaponLevelBasic,
		Name:            "Spread Shot",
		Description:     "Standard multi-directional fire",
		IconEmoji:       "💥",
		Damage:          20,
		FireRate:        6.0,
		ProjectileSpeed: 300,
		Spread:          0.2, // 0.2 radians ≈ 11 degrees
		ProjectileCount: 3,
		Unlocked:        true,
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
			Level:           WeaponLevelBasic,
			Name:            "Laser Rifle",
			Description:     "Continuous beam, high damage",
			IconEmoji:       "🔴",
			Damage:          25,
			FireRate:        8.0,
			ProjectileSpeed: 400,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
		}
	case WeaponTypeShotgun:
		weapon = &Weapon{
			Type:            WeaponTypeShotgun,
			Level:           WeaponLevelBasic,
			Name:            "Shotgun",
			Description:     "Wide spread, close range",
			IconEmoji:       "🔥",
			Damage:          35,
			FireRate:        3.0,
			ProjectileSpeed: 250,
			Spread:          0.8,
			ProjectileCount: 8,
			Unlocked:        true,
		}
	case WeaponTypePlasma:
		weapon = &Weapon{
			Type:            WeaponTypePlasma,
			Level:           WeaponLevelBasic,
			Name:            "Plasma Burst",
			Description:     "Explosive projectiles with splash",
			IconEmoji:       "⚡",
			Damage:          40,
			FireRate:        4.0,
			ProjectileSpeed: 280,
			Spread:          0.3,
			ProjectileCount: 3,
			Unlocked:        true,
		}
	case WeaponTypeHoming:
		weapon = &Weapon{
			Type:            WeaponTypeHoming,
			Level:           WeaponLevelBasic,
			Name:            "Homing Missiles",
			Description:     "Track enemies automatically",
			IconEmoji:       "🚀",
			Damage:          45,
			FireRate:        3.5,
			ProjectileSpeed: 200,
			Spread:          0.2,
			ProjectileCount: 2,
			Unlocked:        true,
		}
	case WeaponTypeRailgun:
		weapon = &Weapon{
			Type:            WeaponTypeRailgun,
			Level:           WeaponLevelBasic,
			Name:            "Railgun",
			Description:     "Pierces through enemies",
			IconEmoji:       "🔵",
			Damage:          50,
			FireRate:        2.5,
			ProjectileSpeed: 500,
			Spread:          0.0,
			ProjectileCount: 1,
			Unlocked:        true,
		}
	case WeaponTypeEnergyLance:
		weapon = &Weapon{
			Type:            WeaponTypeEnergyLance,
			Level:           WeaponLevelBasic,
			Name:            "Energy Lance",
			Description:     "Charges up for massive damage",
			IconEmoji:       "⚔️",
			Damage:          80,
			FireRate:        1.5,
			ProjectileSpeed: 350,
			Spread:          0.1,
			ProjectileCount: 1,
			Unlocked:        true,
		}
	case WeaponTypePulse:
		weapon = &Weapon{
			Type:            WeaponTypePulse,
			Level:           WeaponLevelBasic,
			Name:            "Pulse Cannon",
			Description:     "Rapid burst fire",
			IconEmoji:       "💫",
			Damage:          15,
			FireRate:        12.0,
			ProjectileSpeed: 320,
			Spread:          0.1,
			ProjectileCount: 2,
			Unlocked:        true,
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

// UpgradeWeapon upgrades a weapon to next level
func (wm *WeaponManager) UpgradeWeapon(weaponType WeaponType) bool {
	if weapon, exists := wm.Weapons[weaponType]; exists {
		if weapon.Level < WeaponLevelMaster {
			weapon.Level++
			// Apply upgrade bonuses
			weapon.Damage *= 1.2
			weapon.FireRate *= 1.1
			weapon.ProjectileSpeed *= 1.05
			return true
		}
	}
	return false
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
