package systems

import (
	"math/rand"

	"stellar-siege/game/entities"
)

// EnemyProbability defines the probability of spawning a specific enemy type
type EnemyProbability struct {
	Type        entities.EnemyType
	Probability float64 // Cumulative probability threshold (0.0 to 1.0)
}

// WaveEnemyConfig defines enemy spawn probabilities for a wave range
type WaveEnemyConfig struct {
	MinWave       int
	MaxWave       int
	Probabilities []EnemyProbability
}

// SpawnCountConfig defines how many enemies spawn at once based on wave
type SpawnCountConfig struct {
	MinWave  int
	MaxWave  int
	MinCount int
	MaxCount int
}

// waveConfigs defines enemy type probabilities for different wave ranges
var waveConfigs = []WaveEnemyConfig{
	{
		MinWave: 1,
		MaxWave: 2,
		Probabilities: []EnemyProbability{
			{Type: entities.EnemyScout, Probability: 0.7},
			{Type: entities.EnemyDrone, Probability: 1.0},
		},
	},
	{
		MinWave: 3,
		MaxWave: 5,
		Probabilities: []EnemyProbability{
			{Type: entities.EnemyScout, Probability: 0.4},
			{Type: entities.EnemyDrone, Probability: 0.7},
			{Type: entities.EnemyHunter, Probability: 1.0},
		},
	},
	{
		MinWave: 6,
		MaxWave: 8,
		Probabilities: []EnemyProbability{
			{Type: entities.EnemyScout, Probability: 0.25},
			{Type: entities.EnemyDrone, Probability: 0.45},
			{Type: entities.EnemyHunter, Probability: 0.7},
			{Type: entities.EnemyTank, Probability: 0.85},
			{Type: entities.EnemySniper, Probability: 1.0},
		},
	},
	{
		MinWave: 9,
		MaxWave: 12,
		Probabilities: []EnemyProbability{
			{Type: entities.EnemyScout, Probability: 0.15},
			{Type: entities.EnemyDrone, Probability: 0.3},
			{Type: entities.EnemyHunter, Probability: 0.5},
			{Type: entities.EnemyTank, Probability: 0.65},
			{Type: entities.EnemySniper, Probability: 0.8},
			{Type: entities.EnemySplitter, Probability: 0.9},
			{Type: entities.EnemyBomber, Probability: 1.0},
		},
	},
	{
		MinWave: 13,
		MaxWave: 9999, // Effectively infinite
		Probabilities: []EnemyProbability{
			{Type: entities.EnemyScout, Probability: 0.1},
			{Type: entities.EnemyDrone, Probability: 0.25},
			{Type: entities.EnemyHunter, Probability: 0.4},
			{Type: entities.EnemyTank, Probability: 0.55},
			{Type: entities.EnemySniper, Probability: 0.68},
			{Type: entities.EnemySplitter, Probability: 0.78},
			{Type: entities.EnemyBomber, Probability: 0.88},
			{Type: entities.EnemyShieldBearer, Probability: 1.0},
		},
	},
}

// spawnCountConfigs defines how many enemies spawn simultaneously
var spawnCountConfigs = []SpawnCountConfig{
	{MinWave: 1, MaxWave: 3, MinCount: 1, MaxCount: 1},
	{MinWave: 4, MaxWave: 7, MinCount: 1, MaxCount: 2},
	{MinWave: 8, MaxWave: 9999, MinCount: 1, MaxCount: 3},
}

// getWaveConfig returns the enemy configuration for a given wave
func getWaveConfig(wave int) *WaveEnemyConfig {
	for i := range waveConfigs {
		if wave >= waveConfigs[i].MinWave && wave <= waveConfigs[i].MaxWave {
			return &waveConfigs[i]
		}
	}
	// Fallback to last config if wave is beyond all ranges
	return &waveConfigs[len(waveConfigs)-1]
}

// selectEnemyType selects an enemy type based on wave configuration
func selectEnemyType(wave int) entities.EnemyType {
	config := getWaveConfig(wave)
	r := rand.Float64()

	for _, prob := range config.Probabilities {
		if r < prob.Probability {
			return prob.Type
		}
	}

	// Fallback to last enemy type in config
	return config.Probabilities[len(config.Probabilities)-1].Type
}

// getSpawnCount returns how many enemies should spawn at once for a given wave
func getSpawnCount(wave int) int {
	for _, config := range spawnCountConfigs {
		if wave >= config.MinWave && wave <= config.MaxWave {
			if config.MinCount == config.MaxCount {
				return config.MinCount
			}
			return config.MinCount + rand.Intn(config.MaxCount-config.MinCount+1)
		}
	}
	// Fallback
	return 1
}
