package systems

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// SoundType represents different sound effects
type SoundType int

const (
	SoundPlayerShoot SoundType = iota
	SoundEnemyShoot
	SoundExplosion
	SoundExplosionSmall  // Scout, small asteroid
	SoundExplosionMedium // Drone, medium asteroid
	SoundExplosionLarge  // Tank, large asteroid
	SoundExplosionBoss   // Boss defeat
	SoundHit
	SoundHitPlayer   // Player taking damage
	SoundHitAsteroid // Asteroid collision
	SoundPowerUp
	SoundPowerUpCollect // Power-up collection
	SoundShieldRecharge // Shield regenerating
	SoundWeaponLevelUp  // Weapon upgrade
	SoundLowHealthWarn  // Health critical
	SoundWaveStart
	SoundBossAppear
	SoundGameOver
	SoundUIClick
	SoundBossAttack  // Boss attack pattern
	SoundBossSpecial // Boss special attack
	SoundBossRage    // Boss entering rage mode
	SoundBossDefeat  // Boss defeated
)

// Sound represents a procedural sound effect
type Sound struct {
	soundType  SoundType
	startTime  time.Time
	duration   time.Duration
	frequency  float64
	volume     float64
	waveFunc   func(float64, float64, float64) float32
	modulators []SoundModulator
}

// SoundModulator applies time-based changes to sound properties
type SoundModulator struct {
	applyFunc func(float64) float64 // time-normalized 0-1, returns amplitude multiplier
}

// SoundManager handles procedural sound effects with real-time generation
type SoundManager struct {
	enabled        bool
	volume         float64
	audioContext   *audio.Context
	players        []*audio.Player
	mutex          sync.Mutex
	lastPlayTime   map[SoundType]time.Time     // Per-sound-type throttling
	soundCooldowns map[SoundType]time.Duration // Cooldown per sound type
}

// MaxConcurrentSounds limits the number of simultaneous sounds to prevent goroutine explosion
const MaxConcurrentSounds = 16

// NewSoundManager creates a new sound manager
func NewSoundManager() (*SoundManager, error) {
	ctx := audio.NewContext(44100) // 44.1kHz sample rate

	// Per-sound-type cooldowns to prevent audio spam during intense gameplay
	cooldowns := map[SoundType]time.Duration{
		SoundPlayerShoot:     50 * time.Millisecond,
		SoundEnemyShoot:      80 * time.Millisecond,
		SoundExplosionSmall:  100 * time.Millisecond,
		SoundExplosionMedium: 100 * time.Millisecond,
		SoundExplosionLarge:  120 * time.Millisecond,
		SoundExplosion:       100 * time.Millisecond,
		SoundHit:             50 * time.Millisecond,
		SoundHitPlayer:       200 * time.Millisecond,
		SoundHitAsteroid:     100 * time.Millisecond,
		SoundPowerUp:         100 * time.Millisecond,
		SoundPowerUpCollect:  100 * time.Millisecond,
		SoundBossAttack:      150 * time.Millisecond,
		SoundBossSpecial:     200 * time.Millisecond,
	}

	return &SoundManager{
		enabled:        true,
		volume:         0.5,
		audioContext:   ctx,
		players:        make([]*audio.Player, 0),
		lastPlayTime:   make(map[SoundType]time.Time),
		soundCooldowns: cooldowns,
	}, nil
}

// PlaySound plays a procedural sound effect
func (sm *SoundManager) PlaySound(soundType SoundType) {
	if !sm.enabled {
		return
	}

	sm.mutex.Lock()
	// Limit concurrent sounds to prevent goroutine explosion during intense gameplay
	if len(sm.players) >= MaxConcurrentSounds {
		sm.mutex.Unlock()
		return
	}

	// Per-sound-type throttling: skip if played too recently
	if cooldown, hasCooldown := sm.soundCooldowns[soundType]; hasCooldown {
		if lastTime, exists := sm.lastPlayTime[soundType]; exists {
			if time.Since(lastTime) < cooldown {
				sm.mutex.Unlock()
				return
			}
		}
	}
	sm.lastPlayTime[soundType] = time.Now()
	sm.mutex.Unlock()

	sound := sm.createSound(soundType)
	if sound != nil {
		go sm.playSoundAsync(sound)
	}
}

// SetEnabled enables or disables sound
func (sm *SoundManager) SetEnabled(enabled bool) {
	sm.enabled = enabled
}

// IsEnabled returns whether sound is enabled
func (sm *SoundManager) IsEnabled() bool {
	return sm.enabled
}

// SetVolume sets the volume level
func (sm *SoundManager) SetVolume(volume float64) {
	if volume < 0 {
		volume = 0
	}
	if volume > 1 {
		volume = 1
	}
	sm.volume = volume
}

// GetVolume returns the current volume level
func (sm *SoundManager) GetVolume() float64 {
	return sm.volume
}

// GetActiveSoundCount returns the number of currently playing sounds
func (sm *SoundManager) GetActiveSoundCount() int {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	return len(sm.players)
}
