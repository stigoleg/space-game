package interfaces

import (
	"github.com/hajimehoshi/ebiten/v2"
	"stellar-siege/game/entities"
)

// SoundManager interface for audio management
type SoundManager interface {
	PlaySound(soundType interface{})
	SetEnabled(enabled bool)
	IsEnabled() bool
	SetVolume(volume float64)
	GetVolume() float64
	StopAll()
}

// SpriteManager interface for sprite management
type SpriteManager interface {
	GetSpriteForEnemy(enemyType int) *ebiten.Image
	GetSpriteForAsteroid(size int) *ebiten.Image
	GetSpriteForPowerUp(powerUpType int) *ebiten.Image
	GetPlayerProjectileSprite() *ebiten.Image
	GetEnemyProjectileSprite() *ebiten.Image
	GetExplosionFrames() []*ebiten.Image
	GetSparkleFrames() []*ebiten.Image
}

// InputEvent represents a single input event with associated data
type InputEvent struct {
	Action      InputAction
	WeaponType  entities.WeaponType  // For InputActionSwitchWeapon
	AbilityType entities.AbilityType // For InputActionActivateAbility
}

// InputAction represents an action triggered by user input
type InputAction int

const (
	InputActionNone InputAction = iota
	InputActionPause
	InputActionCycleWeapon
	InputActionSwitchWeapon
	InputActionActivateAbility
	InputActionShoot
)

// InputHandler interface for input management
type InputHandler interface {
	PollGameplayInput() []InputEvent
	PollMenuInput() []InputEvent
	GetMovementVector() (x, y float64)
	IsShootPressed() bool
}

// CollisionManager interface for collision detection
type CollisionManager interface {
	CheckCollisions()
	RegisterEntity(entity Collidable)
	UnregisterEntity(entity Collidable)
	Clear()
}

// EntityManager interface for entity lifecycle management
type EntityManager interface {
	SpawnEnemy(enemyType int, x, y float64) Entity
	SpawnPowerUp(powerUpType int, x, y float64) Entity
	SpawnExplosion(x, y float64, radius float64) Entity
	CleanupInactive()
	GetActiveEntities() []Entity
	Clear()
}

// CameraSystem interface for camera management
type CameraSystem interface {
	SetShake(intensity, duration float64)
	Update()
	GetOffset() (x, y float64)
	Reset()
}

// LeaderboardManager interface for leaderboard management
type LeaderboardManager interface {
	SubmitScore(playerName string, score int, wave int) error
	GetTopScores(limit int) ([]LeaderboardEntry, error)
	IsOnline() bool
}

// LeaderboardEntry represents a single leaderboard entry
type LeaderboardEntry struct {
	Rank       int
	PlayerName string
	Score      int
	Wave       int
	Timestamp  string
}

// AchievementManager interface for achievement tracking
type AchievementManager interface {
	CheckAchievement(achievementID string) bool
	UnlockAchievement(achievementID string)
	GetUnlockedAchievements() []Achievement
	GetProgress(achievementID string) (current, total int)
	SaveProgress() error
}

// Achievement represents a game achievement
type Achievement struct {
	ID          string
	Name        string
	Description string
	Unlocked    bool
	Progress    int
	MaxProgress int
}

// ProgressionManager interface for game progression tracking
type ProgressionManager interface {
	GetCurrentLevel() int
	GetExperience() int
	AddExperience(amount int)
	GetUnlockedWeapons() []entities.WeaponType
	GetUnlockedAbilities() []entities.AbilityType
	UnlockWeapon(weaponType entities.WeaponType)
	UnlockAbility(abilityType entities.AbilityType)
	SaveProgress() error
	LoadProgress() error
}
