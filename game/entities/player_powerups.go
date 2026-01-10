package entities

import (
	"math"
	"math/rand"
)

// MysteryEffect represents a mystery power-up effect
type MysteryEffect int

const (
	MysteryEffectSuperWeaponUpgrade MysteryEffect = iota
	MysteryEffectSpeedBoost
	MysteryEffectShieldOvercharge
	MysteryEffectRapidFire
	MysteryEffectInvincibility
	MysteryEffectScoreMultiplier
	MysteryEffectWeaponDowngrade
	MysteryEffectEngineMalfunction
	MysteryEffectShieldDrain
	MysteryEffectFireRateReduction
	MysteryEffectControlReversal
)

// ApplyPowerUp applies a power-up to the player and returns a message and whether it's positive
func (p *Player) ApplyPowerUp(puType PowerUpType) (string, bool) {
	switch puType {
	case PowerUpHealth:
		p.Health = min(p.Health+30, p.MaxHealth)
		return "", false
	case PowerUpShield:
		p.Shield = min(p.Shield+30, p.MaxShield)
		return "", false
	case PowerUpWeapon:
		// First, upgrade basic gun to level 5
		basicGun := p.WeaponMgr.GetWeapon(WeaponTypeSpread)
		if basicGun != nil && basicGun.Level < WeaponLevelMkV {
			// Upgrade basic gun
			p.WeaponMgr.UpgradeWeapon(WeaponTypeSpread)
			// Keep deprecated WeaponLevel in sync
			p.WeaponLevel = int(basicGun.Level)
			return basicGun.IconEmoji + " " + basicGun.Name + " UPGRADED!", true
		}

		// Once basic gun is maxed, unlock special weapons
		weaponTypes := []WeaponType{
			WeaponTypeFollowingRocket,
			WeaponTypeChainLightning,
			WeaponTypeFlamethrower,
			WeaponTypeIonBeam,
			WeaponTypeBlaster,
			WeaponTypeLaser,
			WeaponTypeShotgun,
			WeaponTypePlasma,
			WeaponTypeHoming,
			WeaponTypeRailgun,
		}

		// Try to find a weapon to unlock first
		var unlockedWeapon WeaponType
		var foundUnlocked bool
		for _, wt := range weaponTypes {
			if !p.WeaponMgr.HasWeapon(wt) {
				unlockedWeapon = wt
				foundUnlocked = true
				break
			}
		}

		if foundUnlocked {
			// Unlock new weapon
			p.WeaponMgr.UnlockWeapon(unlockedWeapon)
			p.WeaponMgr.SwitchWeapon(unlockedWeapon)
			weapon := p.WeaponMgr.GetWeapon(unlockedWeapon)
			return weapon.IconEmoji + " " + weapon.Name + " UNLOCKED!", true
		} else {
			// All weapons unlocked, upgrade current weapon
			weapon := p.WeaponMgr.GetCurrentWeapon()
			if weapon != nil && p.WeaponMgr.UpgradeWeapon(weapon.Type) {
				return weapon.IconEmoji + " " + weapon.Name + " UPGRADED!", true
			}
		}
		return "", false
	case PowerUpSpeed:
		p.Speed = math.Min(p.Speed+0.5, 10)
		return "", false
	case PowerUpMystery:
		effect := p.ApplyMysteryEffect()
		effectName := GetMysteryEffectName(effect)
		isPositive := IsPositiveEffect(effect)
		return effectName, isPositive
	}
	return "", false
}

// ApplyMysteryEffect applies a random mystery power-up effect (60% positive, 40% negative)
func (p *Player) ApplyMysteryEffect() MysteryEffect {
	roll := rand.Float64()

	var effect MysteryEffect

	// 60% chance for positive effects
	if roll < 0.60 {
		// Positive effects
		posRoll := rand.Float64()
		switch {
		case posRoll < 0.25: // 15% of total (25% of 60%)
			effect = MysteryEffectSuperWeaponUpgrade
			// Upgrade weapon by 2 levels
			p.WeaponLevel = min(p.WeaponLevel+2, 5)

		case posRoll < 0.42: // 10% of total
			effect = MysteryEffectSpeedBoost
			p.SpeedBoostTimer = 15.0
			p.SpeedBoostMultiplier = 1.5

		case posRoll < 0.58: // 10% of total
			effect = MysteryEffectShieldOvercharge
			p.Shield = min(int(float64(p.MaxShield)*1.5), p.MaxShield*2)

		case posRoll < 0.75: // 10% of total
			effect = MysteryEffectRapidFire
			p.RapidFireTimer = 10.0
			p.RapidFireMultiplier = 2.0

		case posRoll < 0.83: // 5% of total
			effect = MysteryEffectInvincibility
			p.InvincibilityTimer = 5.0

		default: // 5% of total
			effect = MysteryEffectScoreMultiplier
			p.ScoreMultiplierTimer = 20.0
			p.ScoreMultiplier = 2.0
		}
	} else {
		// 40% chance for negative effects
		negRoll := rand.Float64()
		switch {
		case negRoll < 0.25: // 10% of total (25% of 40%)
			effect = MysteryEffectWeaponDowngrade
			// Downgrade weapon by 1 level (minimum 1)
			if p.WeaponLevel > 1 {
				p.WeaponLevel--
			}

		case negRoll < 0.50: // 10% of total
			effect = MysteryEffectEngineMalfunction
			p.SpeedBoostTimer = 10.0
			p.SpeedBoostMultiplier = 0.6 // -40% speed

		case negRoll < 0.70: // 8% of total
			effect = MysteryEffectShieldDrain
			p.Shield = p.Shield / 2 // Lose 50% shield

		case negRoll < 0.88: // 7% of total
			effect = MysteryEffectFireRateReduction
			p.SlowFireTimer = 8.0
			p.SlowFireMultiplier = 0.5 // Half fire rate

		default: // 5% of total
			effect = MysteryEffectControlReversal
			p.ControlReversed = true
			p.ControlReversalTimer = 5.0
		}
	}

	return effect
}

// GetMysteryEffectName returns a display name for a mystery effect
func GetMysteryEffectName(effect MysteryEffect) string {
	switch effect {
	case MysteryEffectSuperWeaponUpgrade:
		return "SUPER WEAPON UPGRADE! +2 Levels"
	case MysteryEffectSpeedBoost:
		return "SPEED BOOST! +50% Speed for 15s"
	case MysteryEffectShieldOvercharge:
		return "SHIELD OVERCHARGE! +50% Max Shield"
	case MysteryEffectRapidFire:
		return "RAPID FIRE! 2x Fire Rate for 10s"
	case MysteryEffectInvincibility:
		return "INVINCIBILITY! 5 seconds"
	case MysteryEffectScoreMultiplier:
		return "SCORE MULTIPLIER! 2x Score for 20s"
	case MysteryEffectWeaponDowngrade:
		return "WEAPON DOWNGRADE! -1 Level"
	case MysteryEffectEngineMalfunction:
		return "ENGINE MALFUNCTION! -40% Speed for 10s"
	case MysteryEffectShieldDrain:
		return "SHIELD DRAIN! -50% Shield"
	case MysteryEffectFireRateReduction:
		return "WEAPON JAM! Half Fire Rate for 8s"
	case MysteryEffectControlReversal:
		return "CONTROL REVERSAL! Reversed for 5s"
	default:
		return "MYSTERY EFFECT!"
	}
}

// IsPositiveEffect returns whether an effect is positive or negative
func IsPositiveEffect(effect MysteryEffect) bool {
	return effect <= MysteryEffectScoreMultiplier
}

// min is a helper function to return the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
