package systems

import (
	"encoding/json"
	"os"
	"sort"
	"time"
)

// AchievementType represents different types of achievements
type AchievementType string

const (
	AchievementTypeMilestone AchievementType = "milestone"
	AchievementTypeCombat    AchievementType = "combat"
	AchievementTypeBoss      AchievementType = "boss"
	AchievementTypeChallenge AchievementType = "challenge"
	AchievementTypeSecret    AchievementType = "secret"
)

// Achievement represents a single achievement
type Achievement struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        AchievementType   `json:"type"`
	IconEmoji   string            `json:"icon"`
	Unlocked    bool              `json:"unlocked"`
	UnlockedAt  *time.Time        `json:"unlocked_at,omitempty"`
	Progress    int               `json:"progress"`
	ProgressMax int               `json:"progress_max"`
	Reward      AchievementReward `json:"reward"`
}

// AchievementReward represents rewards from unlocking achievements
type AchievementReward struct {
	UnlockWeapon    string `json:"unlock_weapon,omitempty"`
	UnlockAbility   string `json:"unlock_ability,omitempty"`
	UnlockMode      string `json:"unlock_mode,omitempty"`
	ScrapMetalBonus int    `json:"scrap_metal_bonus"`
}

// AchievementManager manages all achievements
type AchievementManager struct {
	Achievements map[string]*Achievement
	dataPath     string
}

// NewAchievementManager creates a new achievement manager
func NewAchievementManager(dataPath string) *AchievementManager {
	am := &AchievementManager{
		Achievements: make(map[string]*Achievement),
		dataPath:     dataPath,
	}
	am.initializeAchievements()
	am.Load()
	return am
}

// initializeAchievements sets up all achievements with their definitions
func (am *AchievementManager) initializeAchievements() {
	achievements := []Achievement{
		// Milestone Achievements
		{
			ID:          "first_victory",
			Name:        "First Victory",
			Description: "Defeat your first enemy",
			Type:        AchievementTypeMilestone,
			IconEmoji:   "🎯",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 50,
			},
		},
		{
			ID:          "wave_5",
			Name:        "Survivor",
			Description: "Reach Wave 5",
			Type:        AchievementTypeMilestone,
			IconEmoji:   "🌊",
			ProgressMax: 5,
			Reward: AchievementReward{
				UnlockWeapon:    "laser",
				ScrapMetalBonus: 100,
			},
		},
		{
			ID:          "wave_10",
			Name:        "Veteran",
			Description: "Reach Wave 10",
			Type:        AchievementTypeMilestone,
			IconEmoji:   "⚔️",
			ProgressMax: 10,
			Reward: AchievementReward{
				UnlockAbility:   "dash",
				ScrapMetalBonus: 200,
			},
		},
		{
			ID:          "wave_20",
			Name:        "Commander",
			Description: "Reach Wave 20",
			Type:        AchievementTypeMilestone,
			IconEmoji:   "👑",
			ProgressMax: 20,
			Reward: AchievementReward{
				UnlockWeapon:    "railgun",
				ScrapMetalBonus: 300,
			},
		},
		{
			ID:          "wave_50",
			Name:        "Legendary Defender",
			Description: "Reach Wave 50",
			Type:        AchievementTypeMilestone,
			IconEmoji:   "🏆",
			ProgressMax: 50,
			Reward: AchievementReward{
				UnlockMode:      "boss_rush",
				ScrapMetalBonus: 500,
			},
		},
		// Combat Achievements
		{
			ID:          "triple_kill",
			Name:        "Triple Kill",
			Description: "Defeat 3 enemies within 2 seconds",
			Type:        AchievementTypeCombat,
			IconEmoji:   "💣",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 75,
			},
		},
		{
			ID:          "max_combo",
			Name:        "Unstoppable Force",
			Description: "Reach 5x combo multiplier",
			Type:        AchievementTypeCombat,
			IconEmoji:   "⚡",
			ProgressMax: 5,
			Reward: AchievementReward{
				ScrapMetalBonus: 150,
			},
		},
		{
			ID:          "perfect_wave",
			Name:        "Flawless",
			Description: "Complete a wave without taking damage",
			Type:        AchievementTypeCombat,
			IconEmoji:   "✨",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 200,
			},
		},
		{
			ID:          "thousand_kills",
			Name:        "Annihilator",
			Description: "Defeat 1000 enemies",
			Type:        AchievementTypeCombat,
			IconEmoji:   "💀",
			ProgressMax: 1000,
			Reward: AchievementReward{
				ScrapMetalBonus: 300,
			},
		},
		{
			ID:          "critical_hits",
			Name:        "Precision Strike",
			Description: "Land 50 critical hits",
			Type:        AchievementTypeCombat,
			IconEmoji:   "🎯",
			ProgressMax: 50,
			Reward: AchievementReward{
				ScrapMetalBonus: 100,
			},
		},
		// Boss Achievements
		{
			ID:          "first_boss",
			Name:        "Boss Slayer",
			Description: "Defeat your first boss",
			Type:        AchievementTypeBoss,
			IconEmoji:   "🐉",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 100,
			},
		},
		{
			ID:          "five_bosses",
			Name:        "Dragon Killer",
			Description: "Defeat 5 bosses",
			Type:        AchievementTypeBoss,
			IconEmoji:   "🗡️",
			ProgressMax: 5,
			Reward: AchievementReward{
				UnlockAbility:   "slow_time",
				ScrapMetalBonus: 250,
			},
		},
		{
			ID:          "boss_no_damage",
			Name:        "Perfect Defense",
			Description: "Defeat a boss without taking damage",
			Type:        AchievementTypeBoss,
			IconEmoji:   "🛡️",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 150,
			},
		},
		// Challenge Achievements
		{
			ID:          "hard_mode_victory",
			Name:        "Elite Warrior",
			Description: "Survive on Hard Mode and reach Wave 5",
			Type:        AchievementTypeChallenge,
			IconEmoji:   "🔥",
			ProgressMax: 5,
			Reward: AchievementReward{
				UnlockMode:      "time_attack",
				ScrapMetalBonus: 200,
			},
		},
		{
			ID:          "score_100k",
			Name:        "Century Club",
			Description: "Score 100,000 points in a single game",
			Type:        AchievementTypeChallenge,
			IconEmoji:   "💯",
			ProgressMax: 100000,
			Reward: AchievementReward{
				ScrapMetalBonus: 250,
			},
		},
		// Secret Achievements
		{
			ID:          "all_weapons",
			Name:        "Arsenal Master",
			Description: "Unlock all weapon types",
			Type:        AchievementTypeSecret,
			IconEmoji:   "🎖️",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 400,
			},
		},
		{
			ID:          "all_abilities",
			Name:        "Mystic Powers",
			Description: "Unlock all abilities",
			Type:        AchievementTypeSecret,
			IconEmoji:   "🌟",
			ProgressMax: 1,
			Reward: AchievementReward{
				ScrapMetalBonus: 400,
			},
		},
	}

	for i := range achievements {
		am.Achievements[achievements[i].ID] = &achievements[i]
	}
}

// Unlock unlocks an achievement
func (am *AchievementManager) Unlock(id string) bool {
	if ach, exists := am.Achievements[id]; exists && !ach.Unlocked {
		ach.Unlocked = true
		now := time.Now()
		ach.UnlockedAt = &now
		am.Save()
		return true
	}
	return false
}

// UpdateProgress updates progress toward an achievement
func (am *AchievementManager) UpdateProgress(id string, progress int) {
	if ach, exists := am.Achievements[id]; exists {
		ach.Progress = progress
		if ach.Progress >= ach.ProgressMax && !ach.Unlocked {
			am.Unlock(id)
		}
		am.Save()
	}
}

// IncrementProgress increments progress for an achievement
func (am *AchievementManager) IncrementProgress(id string, amount int) {
	if ach, exists := am.Achievements[id]; exists {
		ach.Progress += amount
		if ach.Progress >= ach.ProgressMax && !ach.Unlocked {
			am.Unlock(id)
		}
		am.Save()
	}
}

// GetUnlockedAchievements returns all unlocked achievements
func (am *AchievementManager) GetUnlockedAchievements() []*Achievement {
	var unlocked []*Achievement
	for _, ach := range am.Achievements {
		if ach.Unlocked {
			unlocked = append(unlocked, ach)
		}
	}
	sort.Slice(unlocked, func(i, j int) bool {
		return unlocked[i].UnlockedAt.After(*unlocked[j].UnlockedAt)
	})
	return unlocked
}

// GetAchievementByID returns an achievement by ID
func (am *AchievementManager) GetAchievementByID(id string) *Achievement {
	return am.Achievements[id]
}

// GetAllAchievements returns all achievements sorted by unlock status
func (am *AchievementManager) GetAllAchievements() []*Achievement {
	var all []*Achievement
	for _, ach := range am.Achievements {
		all = append(all, ach)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].Unlocked != all[j].Unlocked {
			return all[i].Unlocked
		}
		return all[i].Name < all[j].Name
	})
	return all
}

// GetUnlockCount returns number of unlocked achievements
func (am *AchievementManager) GetUnlockCount() int {
	count := 0
	for _, ach := range am.Achievements {
		if ach.Unlocked {
			count++
		}
	}
	return count
}

// GetTotalAchievements returns total number of achievements
func (am *AchievementManager) GetTotalAchievements() int {
	return len(am.Achievements)
}

// Save saves achievements to JSON file
func (am *AchievementManager) Save() error {
	var achievements []*Achievement
	for _, ach := range am.Achievements {
		achievements = append(achievements, ach)
	}

	data, err := json.MarshalIndent(achievements, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(am.dataPath, data, 0644)
}

// Load loads achievements from JSON file
func (am *AchievementManager) Load() error {
	data, err := os.ReadFile(am.dataPath)
	if err != nil {
		// File doesn't exist yet, that's okay
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var achievements []*Achievement
	if err := json.Unmarshal(data, &achievements); err != nil {
		return err
	}

	// Update loaded data into our achievements map
	for _, ach := range achievements {
		if existing, exists := am.Achievements[ach.ID]; exists {
			existing.Unlocked = ach.Unlocked
			existing.UnlockedAt = ach.UnlockedAt
			existing.Progress = ach.Progress
		}
	}

	return nil
}

// IsWeaponUnlocked checks if a weapon has been unlocked
func (am *AchievementManager) IsWeaponUnlocked(weaponID string) bool {
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockWeapon == weaponID {
			return true
		}
	}
	return false
}

// IsAbilityUnlocked checks if an ability has been unlocked
func (am *AchievementManager) IsAbilityUnlocked(abilityID string) bool {
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockAbility == abilityID {
			return true
		}
	}
	return false
}

// IsModeUnlocked checks if a game mode has been unlocked
func (am *AchievementManager) IsModeUnlocked(modeID string) bool {
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockMode == modeID {
			return true
		}
	}
	return false
}

// GetUnlockedWeapons returns all unlocked weapon IDs
func (am *AchievementManager) GetUnlockedWeapons() []string {
	var weapons []string
	seen := make(map[string]bool)
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockWeapon != "" && !seen[ach.Reward.UnlockWeapon] {
			weapons = append(weapons, ach.Reward.UnlockWeapon)
			seen[ach.Reward.UnlockWeapon] = true
		}
	}
	return weapons
}

// GetUnlockedAbilities returns all unlocked ability IDs
func (am *AchievementManager) GetUnlockedAbilities() []string {
	var abilities []string
	seen := make(map[string]bool)
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockAbility != "" && !seen[ach.Reward.UnlockAbility] {
			abilities = append(abilities, ach.Reward.UnlockAbility)
			seen[ach.Reward.UnlockAbility] = true
		}
	}
	return abilities
}

// GetUnlockedModes returns all unlocked game mode IDs
func (am *AchievementManager) GetUnlockedModes() []string {
	var modes []string
	seen := make(map[string]bool)
	for _, ach := range am.Achievements {
		if ach.Unlocked && ach.Reward.UnlockMode != "" && !seen[ach.Reward.UnlockMode] {
			modes = append(modes, ach.Reward.UnlockMode)
			seen[ach.Reward.UnlockMode] = true
		}
	}
	return modes
}
