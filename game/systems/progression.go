package systems

import (
	"encoding/json"
	"os"
	"time"
)

// ProgressionData represents persistent player progression
type ProgressionData struct {
	TotalScrap        int                     `json:"total_scrap"`
	Prestige          int                     `json:"prestige"`
	PrestigePoints    int                     `json:"prestige_points"`
	Upgrades          map[string]UpgradeLevel `json:"upgrades"`
	UnlockedCosmetics map[string]bool         `json:"unlocked_cosmetics"`
	LastUpdated       time.Time               `json:"last_updated"`
}

// UpgradeLevel represents the level of a specific upgrade
type UpgradeLevel struct {
	Level        int `json:"level"`
	MaxLevel     int `json:"max_level"`
	CostPerLevel int `json:"cost_per_level"`
	CurrentCost  int `json:"current_cost"`
}

// ProgressionManager manages persistent progression
type ProgressionManager struct {
	data      *ProgressionData
	dataPath  string
	scrapGain int // Scrap gained this session
}

// NewProgressionManager creates a new progression manager
func NewProgressionManager(dataPath string) *ProgressionManager {
	pm := &ProgressionManager{
		dataPath: dataPath,
		data: &ProgressionData{
			TotalScrap:        0,
			Prestige:          0,
			PrestigePoints:    0,
			Upgrades:          make(map[string]UpgradeLevel),
			UnlockedCosmetics: make(map[string]bool),
			LastUpdated:       time.Now(),
		},
	}

	// Initialize upgrades
	pm.initializeUpgrades()
	pm.Load()

	return pm
}

// initializeUpgrades sets up all available upgrades
func (pm *ProgressionManager) initializeUpgrades() {
	upgrades := map[string]UpgradeLevel{
		"max_health": {
			Level:        0,
			MaxLevel:     5,
			CostPerLevel: 100,
			CurrentCost:  100,
		},
		"max_shield": {
			Level:        0,
			MaxLevel:     5,
			CostPerLevel: 80,
			CurrentCost:  80,
		},
		"movement_speed": {
			Level:        0,
			MaxLevel:     10,
			CostPerLevel: 75,
			CurrentCost:  75,
		},
		"fire_rate": {
			Level:        0,
			MaxLevel:     10,
			CostPerLevel: 90,
			CurrentCost:  90,
		},
		"damage_multiplier": {
			Level:        0,
			MaxLevel:     5,
			CostPerLevel: 150,
			CurrentCost:  150,
		},
		"scrap_gain_multiplier": {
			Level:        0,
			MaxLevel:     5,
			CostPerLevel: 200,
			CurrentCost:  200,
		},
	}

	for key, upgrade := range upgrades {
		pm.data.Upgrades[key] = upgrade
	}
}

// BuyUpgrade purchases an upgrade level
func (pm *ProgressionManager) BuyUpgrade(upgradeID string) bool {
	upgrade, exists := pm.data.Upgrades[upgradeID]
	if !exists {
		return false
	}

	// Check if already maxed
	if upgrade.Level >= upgrade.MaxLevel {
		return false
	}

	// Check if enough scrap
	if pm.data.TotalScrap < upgrade.CurrentCost {
		return false
	}

	// Purchase
	pm.data.TotalScrap -= upgrade.CurrentCost
	upgrade.Level++

	// Calculate next cost (scales with level)
	upgrade.CurrentCost = upgrade.CostPerLevel + (upgrade.Level * (upgrade.CostPerLevel / 2))

	pm.data.Upgrades[upgradeID] = upgrade
	pm.data.LastUpdated = time.Now()
	pm.Save()

	return true
}

// AddScrap adds scrap metal to the player's total
func (pm *ProgressionManager) AddScrap(amount int) {
	// Apply prestige multiplier
	prestigeMultiplier := 1.0 + (float64(pm.data.Prestige) * 0.1)
	finalAmount := int(float64(amount) * prestigeMultiplier)

	pm.data.TotalScrap += finalAmount
	pm.scrapGain += finalAmount
}

// GetUpgradeBonus returns the bonus value from an upgrade
func (pm *ProgressionManager) GetUpgradeBonus(upgradeID string) float64 {
	upgrade, exists := pm.data.Upgrades[upgradeID]
	if !exists {
		return 0
	}

	switch upgradeID {
	case "max_health":
		return float64(upgrade.Level) * 5.0
	case "max_shield":
		return float64(upgrade.Level) * 3.0
	case "movement_speed":
		return float64(upgrade.Level) * 0.02
	case "fire_rate":
		return float64(upgrade.Level) * 0.05
	case "damage_multiplier":
		return 1.0 + (float64(upgrade.Level) * 0.1)
	case "scrap_gain_multiplier":
		return 1.0 + (float64(upgrade.Level) * 0.2)
	default:
		return 0
	}
}

// GetUpgradeLevel returns the level of an upgrade
func (pm *ProgressionManager) GetUpgradeLevel(upgradeID string) int {
	if upgrade, exists := pm.data.Upgrades[upgradeID]; exists {
		return upgrade.Level
	}
	return 0
}

// GetUpgradeMaxLevel returns the max level of an upgrade
func (pm *ProgressionManager) GetUpgradeMaxLevel(upgradeID string) int {
	if upgrade, exists := pm.data.Upgrades[upgradeID]; exists {
		return upgrade.MaxLevel
	}
	return 0
}

// GetUpgradeCost returns the cost of the next level
func (pm *ProgressionManager) GetUpgradeCost(upgradeID string) int {
	if upgrade, exists := pm.data.Upgrades[upgradeID]; exists {
		return upgrade.CurrentCost
	}
	return 0
}

// Prestige resets progress for prestige bonus
func (pm *ProgressionManager) Prestige() bool {
	// Must have high enough level to prestige
	if pm.data.Prestige == 0 && pm.data.TotalScrap < 5000 {
		return false
	}
	if pm.data.Prestige > 0 && pm.data.TotalScrap < (5000*(pm.data.Prestige+1)) {
		return false
	}

	pm.data.Prestige++
	pm.data.PrestigePoints += 10 * pm.data.Prestige

	// Reset upgrades
	for key, upgrade := range pm.data.Upgrades {
		upgrade.Level = 0
		upgrade.CurrentCost = upgrade.CostPerLevel
		pm.data.Upgrades[key] = upgrade
	}

	// Keep some scrap as prestige reward
	pm.data.TotalScrap = pm.data.TotalScrap / 2

	pm.data.LastUpdated = time.Now()
	pm.Save()

	return true
}

// UnlockCosmetic unlocks a cosmetic item
func (pm *ProgressionManager) UnlockCosmetic(cosmeticID string) {
	pm.data.UnlockedCosmetics[cosmeticID] = true
	pm.Save()
}

// IsCosmeticUnlocked checks if a cosmetic is unlocked
func (pm *ProgressionManager) IsCosmeticUnlocked(cosmeticID string) bool {
	return pm.data.UnlockedCosmetics[cosmeticID]
}

// GetTotalScrap returns total scrap
func (pm *ProgressionManager) GetTotalScrap() int {
	return pm.data.TotalScrap
}

// GetPrestige returns prestige level
func (pm *ProgressionManager) GetPrestige() int {
	return pm.data.Prestige
}

// GetSessionScrapGain returns scrap earned in current session
func (pm *ProgressionManager) GetSessionScrapGain() int {
	return pm.scrapGain
}

// ResetSessionStats resets session-based stats
func (pm *ProgressionManager) ResetSessionStats() {
	pm.scrapGain = 0
}

// GetAllUpgrades returns all upgrades
func (pm *ProgressionManager) GetAllUpgrades() map[string]UpgradeLevel {
	return pm.data.Upgrades
}

// Save saves progression to JSON file
func (pm *ProgressionManager) Save() error {
	pm.data.LastUpdated = time.Now()
	data, err := json.MarshalIndent(pm.data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pm.dataPath, data, 0644)
}

// Load loads progression from JSON file
func (pm *ProgressionManager) Load() error {
	data, err := os.ReadFile(pm.dataPath)
	if err != nil {
		// File doesn't exist yet, that's okay
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, pm.data); err != nil {
		return err
	}

	return nil
}

// GetProgressionSummary returns a string summary of progression
func (pm *ProgressionManager) GetProgressionSummary() string {
	return "Prestige: " + string(rune(pm.data.Prestige)) +
		" | Scrap: " + string(rune(pm.data.TotalScrap)) +
		" | Upgrades Purchased: " + string(rune(pm.getTotalUpgradeLevels()))
}

func (pm *ProgressionManager) getTotalUpgradeLevels() int {
	total := 0
	for _, upgrade := range pm.data.Upgrades {
		total += upgrade.Level
	}
	return total
}
