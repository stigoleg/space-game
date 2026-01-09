package systems

import (
	"io"
	"math"
	"time"
)

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
