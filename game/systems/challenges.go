package systems

import (
	"encoding/json"
	"os"
	"time"
)

// ChallengeMode represents different challenge game modes
type ChallengeMode int

const (
	ChallengeModeEndless ChallengeMode = iota
	ChallengeModeBossRush
	ChallengeModeTimeAttack
	ChallengeModeSurvival
	ChallengeModeDaily
)

// ChallengeConfig represents configuration for a challenge mode
type ChallengeConfig struct {
	Mode              ChallengeMode `json:"mode"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Duration          int           `json:"duration"`   // Seconds (0 = unlimited)
	MaxWaves          int           `json:"max_waves"`  // 0 = unlimited
	MaxBosses         int           `json:"max_bosses"` // For boss rush
	EnemyHealthMult   float64       `json:"enemy_health"`
	EnemySpeedMult    float64       `json:"enemy_speed"`
	PowerUpSpawnRate  float64       `json:"powerup_spawn"`
	AsteroidsEnabled  bool          `json:"asteroids_enabled"`
	HazardsEnabled    bool          `json:"hazards_enabled"`
	ScoringMultiplier float64       `json:"scoring_mult"`
	Icon              string        `json:"icon"`
	Unlocked          bool          `json:"unlocked"`
	Leaderboard       bool          `json:"has_leaderboard"`
}

// ChallengeScore represents a score in a challenge
type ChallengeScore struct {
	PlayerName  string    `json:"player_name"`
	Score       int64     `json:"score"`
	Wave        int       `json:"wave"`
	Bosses      int       `json:"bosses_defeated"`
	TimeSeconds int       `json:"time_seconds"`
	Date        time.Time `json:"date"`
	Difficulty  string    `json:"difficulty"`
}

// ChallengeManager manages challenge modes and scoring
type ChallengeManager struct {
	Config       map[ChallengeMode]ChallengeConfig
	Leaderboards map[ChallengeMode][]*ChallengeScore
	dataPath     string
}

// NewChallengeManager creates a new challenge manager
func NewChallengeManager(dataPath string) *ChallengeManager {
	cm := &ChallengeManager{
		Config:       make(map[ChallengeMode]ChallengeConfig),
		Leaderboards: make(map[ChallengeMode][]*ChallengeScore),
		dataPath:     dataPath,
	}

	cm.initializeChallenges()
	cm.Load()

	return cm
}

// initializeChallenges sets up all challenge modes
func (cm *ChallengeManager) initializeChallenges() {
	challenges := map[ChallengeMode]ChallengeConfig{
		ChallengeModeEndless: {
			Mode:              ChallengeModeEndless,
			Name:              "Endless",
			Description:       "Classic endless wave survival",
			Duration:          0,
			MaxWaves:          0,
			EnemyHealthMult:   1.0,
			EnemySpeedMult:    1.0,
			PowerUpSpawnRate:  1.0,
			AsteroidsEnabled:  true,
			HazardsEnabled:    false,
			ScoringMultiplier: 1.0,
			Icon:              "∞",
			Unlocked:          true,
			Leaderboard:       true,
		},
		ChallengeModeBossRush: {
			Mode:              ChallengeModeBossRush,
			Name:              "Boss Rush",
			Description:       "Face consecutive bosses",
			Duration:          0,
			MaxWaves:          0,
			MaxBosses:         5,
			EnemyHealthMult:   1.3,
			EnemySpeedMult:    1.2,
			PowerUpSpawnRate:  1.2,
			AsteroidsEnabled:  false,
			HazardsEnabled:    false,
			ScoringMultiplier: 2.0,
			Icon:              "🐉",
			Unlocked:          false, // Unlock via achievements
			Leaderboard:       true,
		},
		ChallengeModeTimeAttack: {
			Mode:              ChallengeModeTimeAttack,
			Name:              "Time Attack",
			Description:       "Score as much as possible in 3 minutes",
			Duration:          180, // 3 minutes
			MaxWaves:          0,
			EnemyHealthMult:   0.8,
			EnemySpeedMult:    1.1,
			PowerUpSpawnRate:  0.8,
			AsteroidsEnabled:  true,
			HazardsEnabled:    false,
			ScoringMultiplier: 1.5,
			Icon:              "⏱️",
			Unlocked:          false,
			Leaderboard:       true,
		},
		ChallengeModeSurvival: {
			Mode:              ChallengeModeSurvival,
			Name:              "Survival",
			Description:       "Limited power-ups, manage resources",
			Duration:          0,
			MaxWaves:          0,
			EnemyHealthMult:   1.1,
			EnemySpeedMult:    1.0,
			PowerUpSpawnRate:  0.5, // Fewer power-ups
			AsteroidsEnabled:  true,
			HazardsEnabled:    true, // Hazards make it harder
			ScoringMultiplier: 1.5,
			Icon:              "🛡️",
			Unlocked:          false,
			Leaderboard:       true,
		},
		ChallengeModeDaily: {
			Mode:              ChallengeModeDaily,
			Name:              "Daily Challenge",
			Description:       "Same challenge for all players today",
			Duration:          0,
			MaxWaves:          0,
			EnemyHealthMult:   1.0,
			EnemySpeedMult:    1.0,
			PowerUpSpawnRate:  1.0,
			AsteroidsEnabled:  true,
			HazardsEnabled:    false,
			ScoringMultiplier: 2.0, // Daily has bonus multiplier
			Icon:              "📅",
			Unlocked:          true,
			Leaderboard:       true,
		},
	}

	for mode, config := range challenges {
		cm.Config[mode] = config
		cm.Leaderboards[mode] = make([]*ChallengeScore, 0)
	}
}

// GetChallengeConfig returns config for a challenge mode
func (cm *ChallengeManager) GetChallengeConfig(mode ChallengeMode) ChallengeConfig {
	if config, exists := cm.Config[mode]; exists {
		return config
	}
	return cm.Config[ChallengeModeEndless]
}

// IsChallengeUnlocked checks if a challenge is unlocked
func (cm *ChallengeManager) IsChallengeUnlocked(mode ChallengeMode) bool {
	if config, exists := cm.Config[mode]; exists {
		return config.Unlocked
	}
	return false
}

// UnlockChallenge unlocks a challenge mode
func (cm *ChallengeManager) UnlockChallenge(mode ChallengeMode) {
	if config, exists := cm.Config[mode]; exists {
		config.Unlocked = true
		cm.Config[mode] = config
		cm.Save()
	}
}

// AddScore adds a score to a challenge leaderboard
func (cm *ChallengeManager) AddScore(mode ChallengeMode, score *ChallengeScore) {
	if _, exists := cm.Leaderboards[mode]; !exists {
		cm.Leaderboards[mode] = make([]*ChallengeScore, 0)
	}

	cm.Leaderboards[mode] = append(cm.Leaderboards[mode], score)

	// Keep only top 100 scores
	if len(cm.Leaderboards[mode]) > 100 {
		cm.Leaderboards[mode] = cm.Leaderboards[mode][:100]
	}

	cm.Save()
}

// GetLeaderboard returns top scores for a challenge
func (cm *ChallengeManager) GetLeaderboard(mode ChallengeMode, limit int) []*ChallengeScore {
	if scores, exists := cm.Leaderboards[mode]; exists {
		if limit > len(scores) {
			limit = len(scores)
		}
		return scores[:limit]
	}
	return []*ChallengeScore{}
}

// GetPlayerRank returns a player's rank in a challenge
func (cm *ChallengeManager) GetPlayerRank(mode ChallengeMode, playerName string) int {
	if scores, exists := cm.Leaderboards[mode]; exists {
		for i, score := range scores {
			if score.PlayerName == playerName {
				return i + 1
			}
		}
	}
	return -1
}

// GetPersonalBest returns a player's personal best in a challenge
func (cm *ChallengeManager) GetPersonalBest(mode ChallengeMode, playerName string) *ChallengeScore {
	if scores, exists := cm.Leaderboards[mode]; exists {
		for _, score := range scores {
			if score.PlayerName == playerName {
				return score
			}
		}
	}
	return nil
}

// GetAllUnlockedChallenges returns all unlocked challenges
func (cm *ChallengeManager) GetAllUnlockedChallenges() []ChallengeConfig {
	var unlocked []ChallengeConfig
	for _, config := range cm.Config {
		if config.Unlocked {
			unlocked = append(unlocked, config)
		}
	}
	return unlocked
}

// GetTotalUnlockedChallenges returns count of unlocked challenges
func (cm *ChallengeManager) GetTotalUnlockedChallenges() int {
	count := 0
	for _, config := range cm.Config {
		if config.Unlocked {
			count++
		}
	}
	return count
}

// Save saves challenge data to file
func (cm *ChallengeManager) Save() error {
	// Save leaderboards
	data := make(map[string]interface{})
	data["leaderboards"] = cm.Leaderboards
	data["config"] = cm.Config

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.dataPath, jsonData, 0644)
}

// Load loads challenge data from file
func (cm *ChallengeManager) Load() error {
	jsonData, err := os.ReadFile(cm.dataPath)
	if err != nil {
		// File doesn't exist, that's okay
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return err
	}

	// Load leaderboards
	if lbData, exists := data["leaderboards"]; exists {
		if lbMap, ok := lbData.(map[string]interface{}); ok {
			for modeStr, scores := range lbMap {
				// Would need proper unmarshaling here
				_ = modeStr
				_ = scores
			}
		}
	}

	return nil
}

// GetDailyChallengeHash returns a hash based on current day
func (cm *ChallengeManager) GetDailyChallengeHash() string {
	now := time.Now()
	day := now.Format("2006-01-02")
	return day
}

// GetDailyChallengeVariation returns challenge variation for the day
func (cm *ChallengeManager) GetDailyChallengeVariation() ChallengeConfig {
	config := cm.GetChallengeConfig(ChallengeModeDaily)

	// Use daily hash to seed variation
	hash := cm.GetDailyChallengeHash()
	_ = hash // Would be used to seed RNG for enemy composition, etc

	return config
}
