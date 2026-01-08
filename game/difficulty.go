package game

type DifficultyMode int

const (
	DifficultyEasy DifficultyMode = iota
	DifficultyNormal
	DifficultyHard
)

type DifficultyConfig struct {
	Name              string
	PlayerHealth      int
	PlayerMaxShield   int
	ShieldRegenRate   float64 // HP per frame
	InvincibilityTime float64 // seconds
	ShieldRegenDelay  float64 // seconds before regen starts
	SpawnMultiplier   float64 // multiplier on enemy count
	DamageMultiplier  float64 // multiplier on enemy bullet damage
	EnemyHealthMult   float64 // multiplier on enemy health
	EnemySpeedMult    float64 // multiplier on enemy speed
}

func GetDifficultyConfig(mode DifficultyMode) DifficultyConfig {
	switch mode {
	case DifficultyEasy:
		return DifficultyConfig{
			Name:              "EASY",
			PlayerHealth:      120,
			PlayerMaxShield:   60,
			ShieldRegenRate:   0.75,
			InvincibilityTime: 0.4,
			ShieldRegenDelay:  2.0,
			SpawnMultiplier:   0.7,
			DamageMultiplier:  0.8,
			EnemyHealthMult:   0.8,
			EnemySpeedMult:    0.85,
		}
	case DifficultyNormal:
		return DifficultyConfig{
			Name:              "NORMAL",
			PlayerHealth:      85,
			PlayerMaxShield:   40,
			ShieldRegenRate:   0.5,
			InvincibilityTime: 0.25,
			ShieldRegenDelay:  3.5,
			SpawnMultiplier:   1.0,
			DamageMultiplier:  1.0,
			EnemyHealthMult:   1.0,
			EnemySpeedMult:    1.0,
		}
	case DifficultyHard:
		return DifficultyConfig{
			Name:              "HARD",
			PlayerHealth:      60,
			PlayerMaxShield:   25,
			ShieldRegenRate:   0.25,
			InvincibilityTime: 0.15,
			ShieldRegenDelay:  6.0,
			SpawnMultiplier:   1.35,
			DamageMultiplier:  1.3,
			EnemyHealthMult:   1.2,
			EnemySpeedMult:    1.15,
		}
	default:
		return GetDifficultyConfig(DifficultyNormal)
	}
}

func GetDifficultyName(mode DifficultyMode) string {
	return GetDifficultyConfig(mode).Name
}
