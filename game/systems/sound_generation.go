package systems

import (
	"math"
	"time"
)

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

// sineWave generates a sine wave
func sineWave(t, freq, vol float64) float32 {
	return float32(math.Sin(2*math.Pi*freq*t) * vol)
}

// squareWave generates a square wave
func squareWave(t, freq, vol float64) float32 {
	cycle := math.Mod(t*freq, 1.0)
	if cycle < 0.5 {
		return float32(vol)
	}
	return float32(-vol)
}

// noisyWave generates a noisy wave (pseudo-random)
func noisyWave(t, freq, vol float64) float32 {
	// Pseudo-random based on time
	seed := int64(t*1000000) % 1000000
	noise := float64((seed*9173+3517)%1000000) / 500000.0
	return float32((noise - 1.0) * vol)
}
