package config

import (
	"encoding/json"
	"os"
)

// GameConfig holds all game configuration
type GameConfig struct {
	Game       GameSettings     `json:"game"`
	Player     PlayerConfig     `json:"player"`
	Enemy      EnemyConfig      `json:"enemy"`
	Boss       BossConfig       `json:"boss"`
	Projectile ProjectileConfig `json:"projectile"`
	Powerup    PowerupConfig    `json:"powerup"`
	Wave       WaveConfig       `json:"wave"`
	Audio      AudioConfig      `json:"audio"`
	Graphics   GraphicsConfig   `json:"graphics"`
	Pool       PoolConfig       `json:"pool"`
}

// GameSettings holds general game settings
type GameSettings struct {
	ScreenWidth     int     `json:"screen_width"`
	ScreenHeight    int     `json:"screen_height"`
	TargetFPS       int     `json:"target_fps"`
	ComboTimeout    float64 `json:"combo_timeout"`
	MaxMultiplier   float64 `json:"max_multiplier"`
	MultiplierDecay float64 `json:"multiplier_decay"`
}

// PlayerConfig holds player-related configuration
type PlayerConfig struct {
	StartHealth       int     `json:"start_health"`
	StartShield       int     `json:"start_shield"`
	Speed             float64 `json:"speed"`
	Radius            float64 `json:"radius"`
	ShieldRegenRate   float64 `json:"shield_regen_rate"`
	ShieldRegenDelay  float64 `json:"shield_regen_delay"`
	InvincibilityTime float64 `json:"invincibility_time"`
	FireRate          float64 `json:"fire_rate"`
}

// EnemyConfig holds enemy-related configuration
type EnemyConfig struct {
	ScoutHealth      int     `json:"scout_health"`
	ScoutSpeed       float64 `json:"scout_speed"`
	ScoutPoints      int     `json:"scout_points"`
	DroneHealth      int     `json:"drone_health"`
	DroneSpeed       float64 `json:"drone_speed"`
	DronePoints      int     `json:"drone_points"`
	HunterHealth     int     `json:"hunter_health"`
	HunterSpeed      float64 `json:"hunter_speed"`
	HunterPoints     int     `json:"hunter_points"`
	TankHealth       int     `json:"tank_health"`
	TankSpeed        float64 `json:"tank_speed"`
	TankPoints       int     `json:"tank_points"`
	FormationBonus   float64 `json:"formation_bonus"`
	BurnTickInterval float64 `json:"burn_tick_interval"`
}

// BossConfig holds boss-related configuration
type BossConfig struct {
	BaseHealth     int     `json:"base_health"`
	BaseSpeed      float64 `json:"base_speed"`
	BaseAttackRate float64 `json:"base_attack_rate"`
	BaseDamage     int     `json:"base_damage"`
	BasePoints     int     `json:"base_points"`
	HealthScaling  float64 `json:"health_scaling"`
	DamageScaling  float64 `json:"damage_scaling"`
	ShieldDuration float64 `json:"shield_duration"`
	ShieldCooldown float64 `json:"shield_cooldown"`
	TelegraphTime  float64 `json:"telegraph_time"`
}

// ProjectileConfig holds projectile-related configuration
type ProjectileConfig struct {
	PlayerSpeed     float64 `json:"player_speed"`
	PlayerDamage    int     `json:"player_damage"`
	EnemySpeed      float64 `json:"enemy_speed"`
	EnemyDamage     int     `json:"enemy_damage"`
	DefaultLifetime float64 `json:"default_lifetime"`
	HomingTurnRate  float64 `json:"homing_turn_rate"`
	ChainRange      float64 `json:"chain_range"`
}

// PowerupConfig holds powerup-related configuration
type PowerupConfig struct {
	HealthRestoreAmount  int     `json:"health_restore_amount"`
	ShieldRestoreAmount  int     `json:"shield_restore_amount"`
	SpeedBoostMultiplier float64 `json:"speed_boost_multiplier"`
	SpeedBoostDuration   float64 `json:"speed_boost_duration"`
	DropChance           float64 `json:"drop_chance"`
	Lifetime             float64 `json:"lifetime"`
}

// WaveConfig holds wave-related configuration
type WaveConfig struct {
	StartingEnemies   int     `json:"starting_enemies"`
	EnemiesPerWave    int     `json:"enemies_per_wave"`
	SpawnInterval     float64 `json:"spawn_interval"`
	BossInterval      int     `json:"boss_interval"`
	AsteroidSpawnRate float64 `json:"asteroid_spawn_rate"`
	MiniBossInterval  float64 `json:"miniboss_interval"`
}

// AudioConfig holds audio-related configuration
type AudioConfig struct {
	MasterVolume float64 `json:"master_volume"`
	SFXVolume    float64 `json:"sfx_volume"`
	MusicVolume  float64 `json:"music_volume"`
	SoundEnabled bool    `json:"sound_enabled"`
}

// GraphicsConfig holds graphics-related configuration
type GraphicsConfig struct {
	EnableParticles   bool    `json:"enable_particles"`
	EnableGlow        bool    `json:"enable_glow"`
	EnableScreenShake bool    `json:"enable_screen_shake"`
	MaxParticles      int     `json:"max_particles"`
	ParticleLifetime  float64 `json:"particle_lifetime"`
	ShakeIntensity    float64 `json:"shake_intensity"`
}

// PoolConfig holds object pool configuration
type PoolConfig struct {
	InitialProjectilePoolSize int `json:"initial_projectile_pool_size"`
	InitialExplosionPoolSize  int `json:"initial_explosion_pool_size"`
	InitialParticlePoolSize   int `json:"initial_particle_pool_size"`
	MaxPoolGrowth             int `json:"max_pool_growth"`
}

// DefaultConfig returns the default game configuration
func DefaultConfig() *GameConfig {
	return &GameConfig{
		Game: GameSettings{
			ScreenWidth:     1280,
			ScreenHeight:    720,
			TargetFPS:       60,
			ComboTimeout:    3.0,
			MaxMultiplier:   5.0,
			MultiplierDecay: 0.5,
		},
		Player: PlayerConfig{
			StartHealth:       100,
			StartShield:       50,
			Speed:             6.0,
			Radius:            20.0,
			ShieldRegenRate:   0.5,
			ShieldRegenDelay:  3.5,
			InvincibilityTime: 0.25,
			FireRate:          0.12,
		},
		Enemy: EnemyConfig{
			ScoutHealth:      10,
			ScoutSpeed:       2.0,
			ScoutPoints:      10,
			DroneHealth:      20,
			DroneSpeed:       1.5,
			DronePoints:      20,
			HunterHealth:     30,
			HunterSpeed:      2.5,
			HunterPoints:     30,
			TankHealth:       50,
			TankSpeed:        1.0,
			TankPoints:       50,
			FormationBonus:   1.2,
			BurnTickInterval: 0.5,
		},
		Boss: BossConfig{
			BaseHealth:     500,
			BaseSpeed:      1.0,
			BaseAttackRate: 1.0,
			BaseDamage:     15,
			BasePoints:     1000,
			HealthScaling:  1.5,
			DamageScaling:  1.2,
			ShieldDuration: 2.0,
			ShieldCooldown: 5.0,
			TelegraphTime:  0.5,
		},
		Projectile: ProjectileConfig{
			PlayerSpeed:     8.0,
			PlayerDamage:    10,
			EnemySpeed:      5.0,
			EnemyDamage:     10,
			DefaultLifetime: 3.0,
			HomingTurnRate:  0.1,
			ChainRange:      150.0,
		},
		Powerup: PowerupConfig{
			HealthRestoreAmount:  25,
			ShieldRestoreAmount:  25,
			SpeedBoostMultiplier: 1.5,
			SpeedBoostDuration:   5.0,
			DropChance:           0.15,
			Lifetime:             10.0,
		},
		Wave: WaveConfig{
			StartingEnemies:   5,
			EnemiesPerWave:    2,
			SpawnInterval:     2.0,
			BossInterval:      5,
			AsteroidSpawnRate: 3.0,
			MiniBossInterval:  15.0,
		},
		Audio: AudioConfig{
			MasterVolume: 0.5,
			SFXVolume:    0.7,
			MusicVolume:  0.5,
			SoundEnabled: true,
		},
		Graphics: GraphicsConfig{
			EnableParticles:   true,
			EnableGlow:        true,
			EnableScreenShake: true,
			MaxParticles:      1000,
			ParticleLifetime:  2.0,
			ShakeIntensity:    1.0,
		},
		Pool: PoolConfig{
			InitialProjectilePoolSize: 100,
			InitialExplosionPoolSize:  50,
			InitialParticlePoolSize:   500,
			MaxPoolGrowth:             1000,
		},
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*GameConfig, error) {
	// Start with defaults
	config := DefaultConfig()

	// If no filename provided, return defaults
	if filename == "" {
		return config, nil
	}

	// Try to load from file
	data, err := os.ReadFile(filename)
	if err != nil {
		// File doesn't exist, return defaults
		return config, nil
	}

	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(config *GameConfig, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
