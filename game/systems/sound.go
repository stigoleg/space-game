package systems

import (
	"io"
	"math"
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
	enabled      bool
	volume       float64
	audioContext *audio.Context
	players      []*audio.Player
	mutex        sync.Mutex
}

// NewSoundManager creates a new sound manager
func NewSoundManager() (*SoundManager, error) {
	ctx := audio.NewContext(44100) // 44.1kHz sample rate
	return &SoundManager{
		enabled:      true,
		volume:       0.5,
		audioContext: ctx,
		players:      make([]*audio.Player, 0),
	}, nil
}

// PlaySound plays a procedural sound effect
func (sm *SoundManager) PlaySound(soundType SoundType) {
	if !sm.enabled {
		return
	}

	sound := sm.createSound(soundType)
	if sound != nil {
		go sm.playSoundAsync(sound)
	}
}

// playSoundAsync plays a sound asynchronously
func (sm *SoundManager) playSoundAsync(sound *Sound) {
	// Create audio stream for this sound
	sampleRate := 44100
	numSamples := int(float64(sampleRate) * sound.duration.Seconds())

	// Create a reader that generates the audio samples
	reader := &SoundReader{
		sound:       sound,
		sampleRate:  sampleRate,
		numSamples:  numSamples,
		sampleIndex: 0,
		sm:          sm,
	}

	player, err := sm.audioContext.NewPlayer(reader)
	if err != nil {
		return
	}

	sm.mutex.Lock()
	sm.players = append(sm.players, player)
	sm.mutex.Unlock()

	player.Play()

	// Clean up finished player
	go func() {
		time.Sleep(sound.duration)
		sm.mutex.Lock()
		for i, p := range sm.players {
			if p == player {
				sm.players = append(sm.players[:i], sm.players[i+1:]...)
				break
			}
		}
		sm.mutex.Unlock()
	}()
}

// SoundReader implements io.Reader for audio streaming
type SoundReader struct {
	sound       *Sound
	sampleRate  int
	numSamples  int
	sampleIndex int
	sm          *SoundManager
}

// Read generates audio samples
func (sr *SoundReader) Read(p []byte) (int, error) {
	if sr.sampleIndex >= sr.numSamples {
		return 0, io.EOF
	}

	// We need 2 bytes per sample (16-bit mono)
	// But Ebiten audio typically expects 4 bytes per sample (2 channels, 16-bit each)
	bytesPerSample := 4
	numSamples := len(p) / bytesPerSample
	if numSamples > sr.numSamples-sr.sampleIndex {
		numSamples = sr.numSamples - sr.sampleIndex
	}

	for i := 0; i < numSamples; i++ {
		// Calculate time in seconds
		t := float64(sr.sampleIndex) / float64(sr.sampleRate)

		// Get base waveform
		sample := sr.sound.waveFunc(t, sr.sound.frequency, sr.sound.volume*sr.sm.volume)

		// Apply modulators
		modulationFactor := 1.0
		if sr.sound.duration.Seconds() > 0 {
			timeRatio := t / sr.sound.duration.Seconds()
			if timeRatio <= 1.0 {
				for _, mod := range sr.sound.modulators {
					modulationFactor *= mod.applyFunc(timeRatio)
				}
			}
		}

		sample = float32(float64(sample) * modulationFactor)

		// Clamp sample to [-1, 1]
		if sample > 1.0 {
			sample = 1.0
		}
		if sample < -1.0 {
			sample = -1.0
		}

		// Convert float32 sample to 16-bit PCM (little-endian)
		pcmSample := int16(sample * 32767)

		// Write stereo (both channels get same sample)
		idx := i * bytesPerSample
		p[idx] = byte(pcmSample)
		p[idx+1] = byte(pcmSample >> 8)
		p[idx+2] = byte(pcmSample)
		p[idx+3] = byte(pcmSample >> 8)

		sr.sampleIndex++
	}

	return numSamples * bytesPerSample, nil
}

// createSound creates a new procedural sound based on type
func (sm *SoundManager) createSound(soundType SoundType) *Sound {
	now := time.Now()

	switch soundType {
	case SoundPlayerShoot:
		// High-energy laser shot - crisp and satisfying
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  150 * time.Millisecond,
			frequency: 900,
			volume:    0.8 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Frequency sweep up for "laser charging" feel
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.5 * t)
					},
				},
				{
					// Sharp attack, medium sustain, quick decay
					applyFunc: func(t float64) float64 {
						if t < 0.05 {
							return t * 20 // Sharp attack
						}
						return math.Exp(-12 * (t - 0.05)) // Quick decay
					},
				},
			},
		}
		return sound

	case SoundEnemyShoot:
		// Lower, more ominous enemy shot
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  140 * time.Millisecond,
			frequency: 350,
			volume:    0.7 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Slightly lower pitch sweep
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.3 * t)
					},
				},
				{
					// Sharp attack, medium decay
					applyFunc: func(t float64) float64 {
						if t < 0.08 {
							return t * 12.5
						}
						return math.Exp(-10 * (t - 0.08))
					},
				},
			},
		}
		return sound

	case SoundExplosion:
		// Much more dramatic, satisfying explosion with multiple layers
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  600 * time.Millisecond,
			frequency: 250,
			volume:    0.85 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Dramatic pitch sweep from low to mid
					applyFunc: func(t float64) float64 {
						if t < 0.15 {
							return 1.0
						} else if t < 0.4 {
							return 1.0 + (0.4 * (t - 0.15) / 0.25)
						}
						return 1.4 - (0.4 * (t - 0.4) / 0.6)
					},
				},
				{
					// Powerful attack with sustained rumble
					applyFunc: func(t float64) float64 {
						if t < 0.08 {
							return t * 12.5 // Immediate impact
						} else if t < 0.25 {
							return 1.0 // Sustain rumble
						}
						return math.Exp(-4 * (t - 0.25)) // Slow decay for impact feel
					},
				},
				{
					// Harmonic resonance
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.3 * math.Sin(t*20*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundExplosionSmall:
		// Quick, crisp explosion for small enemies/asteroids
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  200 * time.Millisecond,
			frequency: 600,
			volume:    0.65 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Sharp quick burst
					applyFunc: func(t float64) float64 {
						if t < 0.05 {
							return t * 20
						}
						return math.Exp(-15 * (t - 0.05))
					},
				},
			},
		}
		return sound

	case SoundExplosionMedium:
		// Medium explosion with more impact
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  400 * time.Millisecond,
			frequency: 300,
			volume:    0.75 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Bass rumble
					applyFunc: func(t float64) float64 {
						if t < 0.1 {
							return t * 10
						}
						return math.Exp(-8 * (t - 0.1))
					},
				},
				{
					// Add crackle
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.2 * math.Sin(t*40*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundExplosionLarge:
		// Deep, powerful explosion for large enemies
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  800 * time.Millisecond,
			frequency: 150,
			volume:    0.9 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Deep rumble with sustain
					applyFunc: func(t float64) float64 {
						if t < 0.08 {
							return t * 12.5
						} else if t < 0.3 {
							return 1.0
						}
						return math.Exp(-3 * (t - 0.3))
					},
				},
				{
					// Low frequency wobble
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.4 * math.Sin(t*10*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundExplosionBoss:
		// Epic multi-layered boss explosion
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  1500 * time.Millisecond,
			frequency: 100,
			volume:    1.0 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Massive sustained explosion
					applyFunc: func(t float64) float64 {
						if t < 0.1 {
							return t * 10
						} else if t < 0.5 {
							return 1.0
						}
						return math.Exp(-2 * (t - 0.5))
					},
				},
				{
					// Deep harmonic resonance
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.5 * math.Sin(t*8*math.Pi))
					},
				},
				{
					// Rising pitch sweep
					applyFunc: func(t float64) float64 {
						if t < 0.7 {
							return 1.0 + (t * 0.5)
						}
						return 1.35
					},
				},
			},
		}
		return sound

	case SoundHitPlayer:
		// Painful impact sound when player is hit
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  180 * time.Millisecond,
			frequency: 400,
			volume:    0.8 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Harsh descending pitch
					applyFunc: func(t float64) float64 {
						return 1.0 - (0.6 * t)
					},
				},
				{
					// Sharp attack with medium decay
					applyFunc: func(t float64) float64 {
						if t < 0.05 {
							return t * 20
						}
						return math.Exp(-8 * (t - 0.05))
					},
				},
			},
		}
		return sound

	case SoundHitAsteroid:
		// Rock collision sound
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  100 * time.Millisecond,
			frequency: 200,
			volume:    0.7 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Quick sharp impact
					applyFunc: func(t float64) float64 {
						if t < 0.02 {
							return t * 50
						}
						return math.Exp(-30 * (t - 0.02))
					},
				},
			},
		}
		return sound

	case SoundPowerUpCollect:
		// Satisfying collection sound
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  250 * time.Millisecond,
			frequency: 800,
			volume:    0.75 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Quick rising arpeggio
					applyFunc: func(t float64) float64 {
						return 1.0 + (t * 1.5)
					},
				},
				{
					// Bell-like envelope
					applyFunc: func(t float64) float64 {
						return math.Exp(-6 * t)
					},
				},
			},
		}
		return sound

	case SoundShieldRecharge:
		// Pulsing shield regeneration
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  300 * time.Millisecond,
			frequency: 500,
			volume:    0.6 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Pulsing charge-up
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.5 * math.Sin(t*20*math.Pi))
					},
				},
				{
					// Fade in and out
					applyFunc: func(t float64) float64 {
						if t < 0.5 {
							return t * 2
						}
						return 2 - (t * 2)
					},
				},
			},
		}
		return sound

	case SoundWeaponLevelUp:
		// Triumphant level up sound
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  600 * time.Millisecond,
			frequency: 600,
			volume:    0.85 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Major chord arpeggio
					applyFunc: func(t float64) float64 {
						if t < 0.2 {
							return 1.0
						} else if t < 0.4 {
							return 1.26 // Major third
						} else if t < 0.6 {
							return 1.5 // Fifth
						}
						return 2.0 // Octave
					},
				},
				{
					// Bright envelope
					applyFunc: func(t float64) float64 {
						if t < 0.1 {
							return t * 10
						}
						return math.Exp(-3 * (t - 0.1))
					},
				},
			},
		}
		return sound

	case SoundLowHealthWarn:
		// Urgent pulsing warning beep
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  150 * time.Millisecond,
			frequency: 1000,
			volume:    0.7 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Sharp on-off pulse
					applyFunc: func(t float64) float64 {
						if math.Mod(t*8, 1.0) < 0.5 {
							return 1.0
						}
						return 0.2
					},
				},
			},
		}
		return sound

	case SoundHit:
		// Crisp impact with satisfying punch - more addictive
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  120 * time.Millisecond,
			frequency: 700,
			volume:    0.75 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Frequency drop for impact feel
					applyFunc: func(t float64) float64 {
						return 1.0 - (0.4 * t)
					},
				},
				{
					// Sharp attack with quick decay - very snappy
					applyFunc: func(t float64) float64 {
						if t < 0.03 {
							return 1.0 // Immediate attack
						}
						return math.Exp(-20 * (t - 0.03)) // Very quick decay
					},
				},
			},
		}
		return sound

	case SoundPowerUp:
		// Much more satisfying ascending tone with reverb-like effect
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  450 * time.Millisecond,
			frequency: 400,
			volume:    0.8 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Aggressive frequency sweep upward
					applyFunc: func(t float64) float64 {
						return 1.0 + (2.0 * t)
					},
				},
				{
					// Pulsing volume for excitement
					applyFunc: func(t float64) float64 {
						envelope := 0.0
						if t < 0.5 {
							envelope = t * 2.0
						} else {
							envelope = 1.0 - ((t - 0.5) * 2.0)
						}
						// Add pulsing tremolo
						pulse := 1.0 + (0.3 * math.Sin(t*16*math.Pi))
						return envelope * pulse
					},
				},
			},
		}
		return sound

	case SoundWaveStart:
		// Fanfare-like sound - multiple notes
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  500 * time.Millisecond,
			frequency: 600,
			volume:    0.7 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Frequency modulation for "fanfare" effect
					applyFunc: func(t float64) float64 {
						if t < 0.3 {
							return 1.0
						} else if t < 0.6 {
							return 1.3
						}
						return 1.6
					},
				},
				{
					// Rhythmic volume
					applyFunc: func(t float64) float64 {
						return math.Sin(t*math.Pi) * math.Exp(-2*t)
					},
				},
			},
		}
		return sound

	case SoundBossAppear:
		// Deep dramatic tone
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  600 * time.Millisecond,
			frequency: 150,
			volume:    0.8 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Slow rise in frequency
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.8 * t)
					},
				},
				{
					// Gradual decay
					applyFunc: func(t float64) float64 {
						return math.Exp(-1.5 * t)
					},
				},
			},
		}
		return sound

	case SoundUIClick:
		// Crisp UI click - modern and satisfying
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  70 * time.Millisecond,
			frequency: 1200,
			volume:    0.5 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Bright click with quick fade
					applyFunc: func(t float64) float64 {
						return math.Exp(-35 * t)
					},
				},
			},
		}
		return sound

	case SoundGameOver:
		// Game over - sad, descending tone with finality
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  1200 * time.Millisecond,
			frequency: 600,
			volume:    0.75 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Dramatic descending pitch
					applyFunc: func(t float64) float64 {
						return 1.0 - (0.8 * t)
					},
				},
				{
					// Slow, sorrowful decay
					applyFunc: func(t float64) float64 {
						if t < 0.1 {
							return 1.0
						}
						return math.Exp(-1.2 * (t - 0.1))
					},
				},
				{
					// Sad harmonics
					applyFunc: func(t float64) float64 {
						return 1.0 - (0.2 * math.Sin(t*4*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundBossDefeat:
		// Boss defeated - triumphant, victorious fanfare
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  1000 * time.Millisecond,
			frequency: 700,
			volume:    0.85 * sm.volume,
			waveFunc:  sineWave,
			modulators: []SoundModulator{
				{
					// Rising pitch for triumph
					applyFunc: func(t float64) float64 {
						if t < 0.4 {
							return 1.0 + (0.5 * t)
						}
						return 1.2 + (0.3 * (t - 0.4))
					},
				},
				{
					// Build and sustain
					applyFunc: func(t float64) float64 {
						if t < 0.3 {
							return t * 3.33
						} else if t < 0.7 {
							return 1.0
						}
						return math.Exp(-3 * (t - 0.7))
					},
				},
				{
					// Harmonic bells
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.3 * math.Sin(t*10*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundBossAttack:
		// Boss attack pattern - powerful synth sound with heavy bass
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  250 * time.Millisecond,
			frequency: 220,
			volume:    0.85 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Deep bass sweep up with menacing tone
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.8 * t)
					},
				},
				{
					// Powerful attack with impact
					applyFunc: func(t float64) float64 {
						if t < 0.08 {
							return t * 12.5 // Immediate heavy impact
						}
						return math.Exp(-6 * (t - 0.08)) // Medium decay for resonance
					},
				},
				{
					// Sub-bass harmonic rumble
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.4 * math.Sin(t*12*math.Pi))
					},
				},
			},
		}
		return sound

	case SoundBossSpecial:
		// Boss special attack - intense sci-fi burst with urgency
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  400 * time.Millisecond,
			frequency: 500,
			volume:    0.9 * sm.volume,
			waveFunc:  noisyWave,
			modulators: []SoundModulator{
				{
					// Aggressive frequency sweep with urgency
					applyFunc: func(t float64) float64 {
						if t < 0.15 {
							return 1.0 + (1.5 * t)
						}
						return 1.0 + (1.5 * 0.15) - (0.5 * (t - 0.15))
					},
				},
				{
					// Intense attack with rapid pulses
					applyFunc: func(t float64) float64 {
						if t < 0.03 {
							return 1.0
						}
						pulse := 1.0 + (0.6 * math.Sin(t*40*math.Pi))
						return pulse * math.Exp(-3*(t-0.03))
					},
				},
			},
		}
		return sound

	case SoundBossRage:
		// Boss rage mode - menacing, deep, terrifying sound
		sound := &Sound{
			soundType: soundType,
			startTime: now,
			duration:  600 * time.Millisecond,
			frequency: 150,
			volume:    0.9 * sm.volume,
			waveFunc:  squareWave,
			modulators: []SoundModulator{
				{
					// Deep menacing tone with slight pitch rise
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.4 * t)
					},
				},
				{
					// Pulsing intensity - like a monster roar
					applyFunc: func(t float64) float64 {
						pulse := 1.0 + (0.7 * math.Sin(t*8*math.Pi))
						return pulse * math.Exp(-1.5*t)
					},
				},
				{
					// Harmonic rumble for threat feel
					applyFunc: func(t float64) float64 {
						return 1.0 + (0.5 * math.Sin(t*3*math.Pi))
					},
				},
			},
		}
		return sound

	default:
		return nil
	}
}

// Wave generation functions
func sineWave(t, freq, vol float64) float32 {
	return float32(math.Sin(2*math.Pi*freq*t) * vol)
}

func squareWave(t, freq, vol float64) float32 {
	cycle := math.Mod(t*freq, 1.0)
	if cycle < 0.5 {
		return float32(vol)
	}
	return float32(-vol)
}

func noisyWave(t, freq, vol float64) float32 {
	// Pseudo-random based on time
	seed := int64(t*1000000) % 1000000
	noise := float64((seed*9173+3517)%1000000) / 500000.0
	return float32((noise - 1.0) * vol)
}

// PlaySoundWithDistance plays a sound with volume based on distance
// maxDistance is the distance at which sound becomes silent
func (sm *SoundManager) PlaySoundWithDistance(soundType SoundType, sourceX, sourceY, listenerX, listenerY, maxDistance float64) {
	if !sm.enabled {
		return
	}

	// Calculate distance from sound source to listener
	dx := sourceX - listenerX
	dy := sourceY - listenerY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Calculate volume attenuation based on distance
	var volumeMultiplier float64
	if distance <= 0 {
		volumeMultiplier = 1.0 // Full volume at source
	} else if distance >= maxDistance {
		return // Sound is too far away, don't play
	} else {
		// Linear falloff: volume decreases linearly with distance
		volumeMultiplier = 1.0 - (distance / maxDistance)
	}

	// Create a modified sound with adjusted volume
	sound := sm.createSound(soundType)
	if sound != nil {
		// Apply distance attenuation
		sound.volume *= volumeMultiplier
		go sm.playSoundAsync(sound)
	}
}

// PlaySoundVariation plays a sound with varied parameters
func (sm *SoundManager) PlaySoundVariation(soundType SoundType, frequencyVariation float64) {
	if !sm.enabled {
		return
	}

	sound := sm.createSound(soundType)
	if sound != nil {
		// Vary the frequency
		sound.frequency *= (1.0 + (frequencyVariation * 0.1))
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
