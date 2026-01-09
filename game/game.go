package game

import (
	"image/color"
	"math"
	"math/rand"

	"stellar-siege/game/entities"
	"stellar-siege/game/systems"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
)

type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
)

type Game struct {
	state       GameState
	player      *entities.Player
	enemies     []*entities.Enemy
	boss        *entities.Boss
	projectiles []*entities.Projectile
	explosions  []*entities.Explosion
	powerups    []*entities.PowerUp
	asteroids   []*entities.Asteroid
	stars       *systems.StarField
	spawner     *systems.WaveSpawner
	hud         *systems.HUD
	leaderboard *systems.Leaderboard
	menu        *systems.Menu
	sound       *systems.SoundManager
	sprites     *systems.SpriteManager

	score                int64
	wave                 int
	multiplier           float64
	comboTimer           float64
	screenShake          float64
	gameTime             float64
	bossWave             bool // Is this a boss wave?
	asteroidSpawn        float64
	damageFlash          float64 // For red screen flash on damage
	lastLowHealthWarning float64 // Track when we last played low health warning

	playerName    string
	nameInputMode bool

	// Difficulty system
	selectedDifficulty DifficultyMode
	difficultyConfig   DifficultyConfig

	// 3D Camera view
	cameraDistance float64 // Distance from top of play area
	cameraHeight   float64 // Height above play area for perspective

	// Advanced camera system
	cameraZoom           float64 // Current zoom level (1.0 = normal)
	cameraTargetZoom     float64 // Target zoom for smooth transitions
	cameraShakeAmount    float64 // Screen shake intensity
	cameraCinematicMode  bool    // Cinematic mode active (for boss, etc)
	cameraCinematicTimer float64 // Time in cinematic mode

	// Floating text for damage and score
	floatingTexts []*entities.FloatingText

	// Impact effects for hit feedback
	impactEffects []*entities.ImpactEffect

	// Mini-boss spawning during boss battles
	miniBossSpawnTimer float64 // Timer for spawning mini-bosses during boss fight
	miniBossesSpawned  int     // Number of mini-bosses spawned in current boss wave

	// Online leaderboard (GitHub Gist)
	onlineLeaderboard *systems.GistLeaderboard
	gistConfig        *systems.GistConfig
	submitScorePrompt bool                  // Whether to prompt user to submit score
	scoreSubmitted    bool                  // Whether score was submitted this session
	onlineScores      []systems.OnlineScore // Cached online scores
}

func NewGame() *Game {
	g := &Game{
		state:              StateMenu,
		stars:              systems.NewStarField(ScreenWidth, ScreenHeight),
		selectedDifficulty: DifficultyNormal,                  // Default to Normal difficulty
		cameraDistance:     100.0,                             // How far back to view from
		cameraHeight:       60.0,                              // How high to view from (for angle)
		cameraZoom:         1.0,                               // Normal zoom
		cameraTargetZoom:   1.0,                               // Target zoom
		cameraShakeAmount:  0.0,                               // No screen shake initially
		floatingTexts:      make([]*entities.FloatingText, 0), // Initialize floating text list
		impactEffects:      make([]*entities.ImpactEffect, 0), // Initialize impact effects list
	}
	g.leaderboard = systems.NewLeaderboard("data/leaderboard.json")
	g.sound, _ = systems.NewSoundManager()
	g.sprites = systems.NewSpriteManager()
	g.menu = systems.NewMenu(g.sprites)

	// Load Gist configuration for online leaderboard from environment variables
	gistConfig, _ := systems.LoadGistConfig("")
	g.gistConfig = gistConfig
	if gistConfig.Enabled && gistConfig.GistID != "" && gistConfig.GitHubToken != "" {
		g.onlineLeaderboard = systems.NewGistLeaderboard(gistConfig.GistID, gistConfig.GitHubToken)
		// Pre-fetch online scores in background
		go func() {
			if scores, err := g.onlineLeaderboard.GetTopScores(100); err == nil {
				g.onlineScores = scores
			}
		}()
	}

	return g
}

func (g *Game) startGame() {
	// Get difficulty config
	g.difficultyConfig = GetDifficultyConfig(g.selectedDifficulty)

	g.state = StatePlaying
	g.player = entities.NewPlayer(ScreenWidth/2, ScreenHeight-100)

	// Apply difficulty settings to player
	g.player.Health = g.difficultyConfig.PlayerHealth
	g.player.MaxHealth = g.difficultyConfig.PlayerHealth
	g.player.Shield = g.difficultyConfig.PlayerMaxShield
	g.player.MaxShield = g.difficultyConfig.PlayerMaxShield
	g.player.ShieldRegenRate = g.difficultyConfig.ShieldRegenRate
	g.player.InvincibilityTime = g.difficultyConfig.InvincibilityTime
	g.player.ShieldRegenDelay = g.difficultyConfig.ShieldRegenDelay
	g.player.LastDamageTime = -999 // Start with regen available

	g.enemies = nil
	g.boss = nil
	g.projectiles = nil
	g.explosions = nil
	g.powerups = nil
	g.asteroids = nil
	g.score = 0
	g.wave = 0
	g.multiplier = 1.0
	g.comboTimer = 0
	g.screenShake = 0
	g.gameTime = 0
	g.bossWave = false
	g.asteroidSpawn = 0
	g.miniBossSpawnTimer = 0
	g.miniBossesSpawned = 0
	g.lastLowHealthWarning = 0
	g.spawner = systems.NewWaveSpawner(ScreenWidth, ScreenHeight)
	g.spawner.SetDifficultyMultipliers(g.difficultyConfig.SpawnMultiplier, g.difficultyConfig.EnemyHealthMult, g.difficultyConfig.EnemySpeedMult, g.difficultyConfig.DamageMultiplier)
	g.hud = systems.NewHUD()
	g.nameInputMode = false
	g.playerName = ""
	g.submitScorePrompt = false
	g.scoreSubmitted = false
}

func (g *Game) Update() error {
	// Update starfield always (visual effect)
	g.stars.Update()

	switch g.state {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateGameOver:
		g.updateGameOver()
	}

	return nil
}

func (g *Game) updateMenu() {
	// Update menu input handling
	g.menu.Update()

	// If showing difficulty select, allow selection
	if g.menu.ShowDifficultySelect {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			// Set the selected difficulty and start game
			g.selectedDifficulty = DifficultyMode(g.menu.SelectedDifficulty)
			g.sound.PlaySound(systems.SoundUIClick)
			g.startGame()
		}
	} else {
		// Main menu
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// Show difficulty selection screen
			g.menu.ShowDifficultySelectMenu()
			g.sound.PlaySound(systems.SoundUIClick)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyL) {
			// Toggle leaderboard view in menu
			g.menu.ToggleLeaderboard()
			g.sound.PlaySound(systems.SoundUIClick)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyI) {
			// Show information menu
			g.menu.ShowInfo()
			g.sound.PlaySound(systems.SoundUIClick)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			// Toggle sound
			g.menu.SoundEnabled = !g.menu.SoundEnabled
			g.sound.SetEnabled(g.menu.SoundEnabled)
			g.sound.PlaySound(systems.SoundUIClick)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			// Could exit, but for now just ignore
		}
	}
}

func (g *Game) updatePlaying() {
	g.gameTime += 1.0 / 60.0

	// Update camera system
	g.updateCamera()

	// Pause
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = StatePaused
		return
	}

	// Update game systems in order
	g.updatePlayerState()
	g.updateBossWave()
	g.updateRegularWaveSpawning()
	g.updateEnemies()
	g.updateProjectiles()
	g.updateExplosions()
	g.updatePowerups()
	g.updateAsteroids()
	g.checkCollisions()
	g.updateComboSystem()
	g.updateVisualEffects()
	g.updateLowHealthWarning()
	g.cleanupEntities()
	g.checkGameOver()
}

// updatePlayerState handles player update, shield recharge, and shooting
func (g *Game) updatePlayerState() {
	if g.player != nil && g.player.Active {
		// Store previous shield value
		prevShield := g.player.Shield

		g.player.Update(ScreenWidth, ScreenHeight, g.gameTime)

		// Check if shield reached max from a lower value (fully recharged)
		if g.player.Shield >= g.player.MaxShield && prevShield < g.player.MaxShield && prevShield > 0 {
			g.sound.PlaySound(systems.SoundShieldRecharge)
		}

		// Player shooting
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			newProjectiles := g.player.Shoot()
			if len(newProjectiles) > 0 {
				g.projectiles = append(g.projectiles, newProjectiles...)
				g.sound.PlaySound(systems.SoundPlayerShoot)
			}
		}
	}
}

// updateBossWave handles boss wave logic, including boss updates and defeat
func (g *Game) updateBossWave() {
	if !g.bossWave {
		return
	}

	if g.boss != nil && g.boss.Active && !g.boss.IsDead() {
		// Track previous phase to detect transitions
		prevPhase := g.boss.Phase

		bossProjectiles := g.boss.Update(g.player.X, g.player.Y, ScreenWidth, ScreenHeight)
		g.projectiles = append(g.projectiles, bossProjectiles...)

		// Play phase transition sounds
		if prevPhase != g.boss.Phase {
			if g.boss.Phase == entities.BossPhaseRage {
				g.sound.PlaySound(systems.SoundBossRage)
			} else if g.boss.Phase == entities.BossPhaseSpecialAttack {
				g.sound.PlaySound(systems.SoundBossSpecial)
			}
		}

		// Play attack sound every few attacks
		if len(bossProjectiles) > 0 && g.boss.AttackPattern == 0 {
			// Play boss attack sound periodically
			if g.boss.AttackTimer < 0.1 {
				if g.boss.Phase == entities.BossPhaseSpecialAttack {
					g.sound.PlaySound(systems.SoundBossSpecial)
				} else {
					g.sound.PlaySound(systems.SoundBossAttack)
				}
			}
		}

		// Spawn mini-bosses during boss fight for higher difficulties
		g.updateMiniBossSpawning()
	} else if g.boss != nil && g.boss.IsDead() {
		// Boss defeated!
		g.spawnExplosion(g.boss.X, g.boss.Y, 100)
		g.spawnExplosion(g.boss.X-40, g.boss.Y-20, 60)
		g.spawnExplosion(g.boss.X+40, g.boss.Y+20, 60)
		g.addScore(int64(g.boss.Points))
		g.screenShake = 30
		g.sound.PlaySound(systems.SoundExplosionBoss) // Boss explosion
		g.sound.PlaySound(systems.SoundBossDefeat)    // Victory fanfare
		g.boss = nil
		g.bossWave = false
		g.miniBossSpawnTimer = 0
		g.miniBossesSpawned = 0
		// Spawn health powerup after boss
		g.powerups = append(g.powerups, entities.NewPowerUp(ScreenWidth/2, 200))
	}
}

// updateRegularWaveSpawning handles regular wave spawning and progression
func (g *Game) updateRegularWaveSpawning() {
	if g.bossWave {
		return
	}

	// Regular wave spawning
	newEnemies := g.spawner.Update(g.gameTime, g.wave)
	g.enemies = append(g.enemies, newEnemies...)
	if g.spawner.WaveCompleted && len(g.enemies) == 0 {
		g.wave++
		g.score += int64(g.wave * 1000) // Wave bonus

		// Every 5 waves, spawn a boss
		if g.wave%5 == 0 {
			g.bossWave = true
			g.boss = entities.NewBoss(ScreenWidth, g.wave/5)
			g.sound.PlaySound(systems.SoundBossAppear)
		} else {
			g.spawner.StartWave(g.wave)
			g.sound.PlaySound(systems.SoundWaveStart)
		}
	}
}

// updateEnemies handles enemy updates and shooting
func (g *Game) updateEnemies() {
	for _, e := range g.enemies {
		if e.Active {
			e.Update(g.player.X, g.player.Y, ScreenWidth, ScreenHeight)
			// Enemy shooting
			if proj := e.TryShoot(); proj != nil {
				// Scale projectile damage by difficulty multiplier
				proj.Damage = int(float64(proj.Damage) * g.difficultyConfig.DamageMultiplier)
				g.projectiles = append(g.projectiles, proj)
				g.sound.PlaySound(systems.SoundEnemyShoot)
			}
		}
	}
}

// updateProjectiles handles projectile updates and off-screen cleanup
func (g *Game) updateProjectiles() {
	for _, p := range g.projectiles {
		if p.Active {
			p.Update()
			// Off-screen check
			if p.Y < -20 || p.Y > ScreenHeight+20 || p.X < -20 || p.X > ScreenWidth+20 {
				p.Active = false
			}
		}
	}
}

// updateExplosions handles explosion animation updates
func (g *Game) updateExplosions() {
	for _, ex := range g.explosions {
		if ex.Active {
			ex.Update()
		}
	}
}

// updatePowerups handles powerup updates and off-screen cleanup
func (g *Game) updatePowerups() {
	for _, pu := range g.powerups {
		if pu.Active {
			pu.Update()
			if pu.Y > ScreenHeight+20 {
				pu.Active = false
			}
		}
	}
}

// updateAsteroids handles asteroid spawning and updates
func (g *Game) updateAsteroids() {
	// Spawn asteroids
	g.asteroidSpawn += 1.0 / 60.0
	if g.asteroidSpawn > 2.0 { // Spawn every 2 seconds
		g.asteroidSpawn = 0
		// Spawn 1-2 asteroids per spawn
		spawnCount := 1
		if rand.Float64() < 0.3 {
			spawnCount = 2
		}
		for i := 0; i < spawnCount; i++ {
			x := rand.Float64() * ScreenWidth
			size := entities.AsteroidSize(rand.Intn(3))
			g.asteroids = append(g.asteroids, entities.NewAsteroid(x, -40, size))
		}
	}

	// Update asteroids
	for _, a := range g.asteroids {
		if a.Active {
			a.Update()
		}
	}
}

// updateComboSystem handles combo timer and multiplier decay
func (g *Game) updateComboSystem() {
	if g.comboTimer > 0 {
		g.comboTimer -= 1.0 / 60.0
		if g.comboTimer <= 0 {
			g.multiplier = 1.0
		}
	}
}

// updateVisualEffects handles screen shake, damage flash, floating text, and impact effects
func (g *Game) updateVisualEffects() {
	// Update screen shake
	if g.screenShake > 0 {
		g.screenShake -= 0.5
		if g.screenShake < 0 {
			g.screenShake = 0
		}
	}

	// Update damage flash
	if g.damageFlash > 0 {
		g.damageFlash -= 1.0 / 60.0
		if g.damageFlash < 0 {
			g.damageFlash = 0
		}
	}

	// Update floating text (damage/score indicators)
	for _, ft := range g.floatingTexts {
		if ft.Active {
			ft.Update()
		}
	}

	// Update impact effects
	for _, ie := range g.impactEffects {
		if ie.Active {
			ie.Update()
		}
	}
}

// updateLowHealthWarning plays low health warning sound when appropriate
func (g *Game) updateLowHealthWarning() {
	if g.player != nil && g.player.Active {
		healthPercent := float64(g.player.Health) / float64(g.player.MaxHealth)
		if healthPercent < 0.3 && g.gameTime-g.lastLowHealthWarning > 3.0 {
			g.sound.PlaySound(systems.SoundLowHealthWarn)
			g.lastLowHealthWarning = g.gameTime
		}
	}
}

// checkGameOver handles game over condition and cleanup
func (g *Game) checkGameOver() {
	if g.player == nil || !g.player.Active {
		g.state = StateGameOver
		g.nameInputMode = true
		g.sound.PlaySound(systems.SoundGameOver)

		// Refresh online leaderboard scores for qualification check
		if g.onlineLeaderboard != nil {
			go func() {
				if scores, err := g.onlineLeaderboard.GetTopScores(100); err == nil {
					g.onlineScores = scores
				}
			}()
		}
	}
}

func (g *Game) updatePaused() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = StatePlaying
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.state = StateMenu
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		// Toggle sound in pause menu
		g.menu.SoundEnabled = !g.menu.SoundEnabled
		g.sound.SetEnabled(g.menu.SoundEnabled)
		g.sound.PlaySound(systems.SoundUIClick)
	}
}

func (g *Game) updateGameOver() {
	if g.nameInputMode {
		// Handle name input
		for _, r := range ebiten.AppendInputChars(nil) {
			if len(g.playerName) < 10 {
				g.playerName += string(r)
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.playerName) > 0 {
			g.playerName = g.playerName[:len(g.playerName)-1]
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(g.playerName) > 0 {
			g.leaderboard.AddEntry(g.playerName, g.score, g.wave)
			g.nameInputMode = false

			// Automatically submit to online leaderboard if score qualifies
			g.autoSubmitScoreIfQualified()
			g.scoreSubmitted = true // Mark as submitted (or attempted)
		}
	} else {
		// Game over controls
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.startGame()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.state = StateMenu
		}
	}
}

// submitScoreOnline submits the player's score to the online leaderboard
func (g *Game) submitScoreOnline() {
	if g.onlineLeaderboard == nil || g.playerName == "" {
		return
	}

	difficulty := ""
	switch g.selectedDifficulty {
	case DifficultyEasy:
		difficulty = "Easy"
	case DifficultyNormal:
		difficulty = "Normal"
	case DifficultyHard:
		difficulty = "Hard"
	}

	// Submit asynchronously to not block game
	go func() {
		if err := g.onlineLeaderboard.SubmitScore(g.playerName, g.score, difficulty, g.wave); err != nil {
			// Silently fail - no UI feedback for now
			return
		}
		// Clear cache to force refresh on next leaderboard view
		g.onlineLeaderboard.ClearCache()
	}()
}

// autoSubmitScoreIfQualified automatically submits score if it's high enough
func (g *Game) autoSubmitScoreIfQualified() {
	if g.onlineLeaderboard == nil || g.playerName == "" {
		return
	}

	// Define minimum thresholds for auto-submission
	const minScore = 1000 // Minimum score to consider
	const minWave = 3     // Minimum wave reached

	// Check if score meets minimum threshold
	if g.score < minScore || g.wave < minWave {
		return // Score too low to submit
	}

	// Check if score would make it into top 100
	// If we have cached scores, check against them
	if len(g.onlineScores) > 0 {
		// If we have less than 100 scores, always submit
		if len(g.onlineScores) < 100 {
			g.submitScoreOnline()
			return
		}

		// Check if our score beats the 100th place
		if len(g.onlineScores) >= 100 {
			lowestScore := g.onlineScores[99].Score
			if g.score > lowestScore {
				g.submitScoreOnline()
				return
			}
		}
	} else {
		// No cached scores - submit if meets minimum threshold
		g.submitScoreOnline()
	}
}

// drawOnlineLeaderboard renders the online leaderboard on screen
func (g *Game) drawOnlineLeaderboard(screen *ebiten.Image) {
	startY := 150
	lineHeight := 20

	// Title
	systems.DrawTextCentered(screen, "GLOBAL LEADERBOARD", ScreenWidth/2, startY, 2, color.RGBA{100, 200, 255, 255})

	// Column headers
	systems.DrawText(screen, "Rank", 100, startY+30, 1, color.RGBA{150, 150, 150, 255})
	systems.DrawText(screen, "Player", 160, startY+30, 1, color.RGBA{150, 150, 150, 255})
	systems.DrawText(screen, "Score", 400, startY+30, 1, color.RGBA{150, 150, 150, 255})
	systems.DrawText(screen, "Difficulty", 550, startY+30, 1, color.RGBA{150, 150, 150, 255})
	systems.DrawText(screen, "Wave", 750, startY+30, 1, color.RGBA{150, 150, 150, 255})

	// Display top 10
	limit := 10
	if len(g.onlineScores) < 10 {
		limit = len(g.onlineScores)
	}

	for i := 0; i < limit; i++ {
		score := g.onlineScores[i]
		y := startY + 60 + (i * lineHeight)

		// Highlight if this is the current player's score
		textColor := color.RGBA{200, 200, 200, 255}
		if score.PlayerName == g.playerName && score.Score == g.score {
			textColor = color.RGBA{255, 255, 100, 255}
		}

		rankStr := systems.FormatNumber(int64(i + 1))
		scoreStr := systems.FormatNumber(score.Score)
		waveStr := systems.FormatNumber(int64(score.Wave))

		systems.DrawText(screen, rankStr, 100, y, 0.8, textColor)
		systems.DrawText(screen, score.PlayerName, 160, y, 0.8, textColor)
		systems.DrawText(screen, scoreStr, 400, y, 0.8, textColor)
		systems.DrawText(screen, score.Difficulty, 550, y, 0.8, textColor)
		systems.DrawText(screen, waveStr, 750, y, 0.8, textColor)
	}
}

func (g *Game) checkCollisions() {
	// Player projectiles vs enemies
	for _, p := range g.projectiles {
		if !p.Active || !p.Friendly {
			continue
		}
		for _, e := range g.enemies {
			if !e.Active {
				continue
			}
			if g.checkCircleCollision(p.X, p.Y, p.Radius, e.X, e.Y, e.Radius) {
				p.Active = false
				e.Health -= p.Damage

				// Add impact effect
				g.impactEffects = append(g.impactEffects, entities.NewImpactEffect(e.X, e.Y, 30, color.RGBA{100, 200, 255, 255}))

				if e.Health <= 0 {
					e.Active = false
					g.spawnExplosion(e.X, e.Y, e.Radius)

					// Play appropriate explosion sound based on enemy type
					switch e.Type {
					case entities.EnemyScout:
						g.sound.PlaySound(systems.SoundExplosionSmall)
					case entities.EnemyDrone:
						g.sound.PlaySound(systems.SoundExplosionSmall)
					case entities.EnemyHunter:
						g.sound.PlaySound(systems.SoundExplosionMedium)
					case entities.EnemyTank:
						g.sound.PlaySound(systems.SoundExplosionLarge)
					case entities.EnemyBomber:
						g.sound.PlaySound(systems.SoundExplosionMedium)
					default:
						g.sound.PlaySound(systems.SoundExplosionSmall)
					}

					points := int64(e.Points)
					g.addScore(points)
					g.spawnFloatingScore(e.X, e.Y, int(points)) // Show score popup
					g.screenShake = 5

					// Chance to spawn powerup
					if rand.Float64() < 0.15 {
						g.powerups = append(g.powerups, entities.NewPowerUp(e.X, e.Y))
					}
				}
			}
		}

		// Player projectiles vs boss
		if g.boss != nil && g.boss.Active && !g.boss.IsDead() {
			if g.checkCircleCollision(p.X, p.Y, p.Radius, g.boss.X, g.boss.Y, g.boss.Radius) {
				p.Active = false

				// Add impact effect for boss
				g.impactEffects = append(g.impactEffects, entities.NewImpactEffect(g.boss.X, g.boss.Y, 40, color.RGBA{255, 150, 100, 255}))

				if g.boss.TakeDamage(p.Damage) {
					// Boss defeated
					g.screenShake = 20
				} else {
					g.screenShake = 3
				}
			}
		}
	}

	// Enemy projectiles vs player
	if g.player != nil && g.player.Active {
		for _, p := range g.projectiles {
			if !p.Active || p.Friendly {
				continue
			}
			if g.checkCircleCollision(p.X, p.Y, p.Radius, g.player.X, g.player.Y, g.player.Radius) {
				p.Active = false
				g.player.TakeDamage(p.Damage, g.gameTime)
				g.spawnFloatingDamage(g.player.X, g.player.Y-20, p.Damage) // Show damage popup
				g.screenShake = 10
				g.damageFlash = 0.2 // Red flash for 0.2 seconds
				g.sound.PlaySound(systems.SoundHitPlayer)
				if g.player.Health <= 0 {
					g.spawnExplosionWithType(g.player.X, g.player.Y, 40, entities.ExplosionBlast)
					g.sound.PlaySound(systems.SoundExplosionLarge)
					g.player.Active = false
				}
			}
		}

		// Enemies vs player (collision)
		for _, e := range g.enemies {
			if !e.Active {
				continue
			}
			if g.checkCircleCollision(e.X, e.Y, e.Radius, g.player.X, g.player.Y, g.player.Radius) {
				e.Active = false
				g.spawnExplosion(e.X, e.Y, e.Radius)
				// Play appropriate explosion sound based on enemy type
				switch e.Type {
				case entities.EnemyScout:
					g.sound.PlaySound(systems.SoundExplosionSmall)
				case entities.EnemyDrone:
					g.sound.PlaySound(systems.SoundExplosionSmall)
				case entities.EnemyHunter:
					g.sound.PlaySound(systems.SoundExplosionMedium)
				case entities.EnemyTank:
					g.sound.PlaySound(systems.SoundExplosionLarge)
				case entities.EnemyBomber:
					g.sound.PlaySound(systems.SoundExplosionMedium)
				}
				g.sound.PlaySound(systems.SoundHitPlayer)
				collisionDamage := int(float64(30) * g.difficultyConfig.DamageMultiplier)
				g.player.TakeDamage(collisionDamage, g.gameTime)
				g.spawnFloatingDamage(g.player.X, g.player.Y-20, collisionDamage) // Show damage popup
				g.screenShake = 15
				if g.player.Health <= 0 {
					g.spawnExplosion(g.player.X, g.player.Y, 40)
					g.sound.PlaySound(systems.SoundExplosionLarge)
					g.player.Active = false
				}
			}
		}

		// Boss vs player (collision)
		if g.boss != nil && g.boss.Active {
			if g.checkCircleCollision(g.boss.X, g.boss.Y, g.boss.Radius*0.5, g.player.X, g.player.Y, g.player.Radius) {
				bossDamage := int(float64(50) * g.difficultyConfig.DamageMultiplier)
				g.player.TakeDamage(bossDamage, g.gameTime)
				g.screenShake = 20
				if g.player.Health <= 0 {
					g.spawnExplosion(g.player.X, g.player.Y, 40)
					g.player.Active = false
				}
			}
		}

		// Powerups vs player
		for _, pu := range g.powerups {
			if !pu.Active {
				continue
			}
			if g.checkCircleCollision(pu.X, pu.Y, pu.Radius, g.player.X, g.player.Y, g.player.Radius) {
				pu.Active = false

				// Play appropriate sound based on power-up type
				switch pu.Type {
				case entities.PowerUpHealth:
					g.sound.PlaySound(systems.SoundPowerUpCollect)
				case entities.PowerUpShield:
					g.sound.PlaySound(systems.SoundShieldRecharge)
				case entities.PowerUpWeapon:
					// Store old weapon level to check if it actually leveled up
					oldLevel := g.player.WeaponLevel
					g.player.ApplyPowerUp(pu.Type)
					if g.player.WeaponLevel > oldLevel {
						g.sound.PlaySound(systems.SoundWeaponLevelUp)
						g.spawnFloatingUpgrade(g.player.X, g.player.Y-30, g.player.WeaponLevel)
					}
					return // Skip ApplyPowerUp call below since we already called it
				case entities.PowerUpSpeed:
					g.sound.PlaySound(systems.SoundPowerUpCollect)
				}

				// Apply power-up effect (skip for weapon since we already handled it)
				if pu.Type != entities.PowerUpWeapon {
					g.player.ApplyPowerUp(pu.Type)
				}
			}
		}

		// Asteroids vs player
		for _, a := range g.asteroids {
			if !a.Active {
				continue
			}
			if g.checkCircleCollision(a.X, a.Y, a.Radius, g.player.X, g.player.Y, g.player.Radius) {
				asteroidDamage := int(float64(15) * g.difficultyConfig.DamageMultiplier)
				g.player.TakeDamage(asteroidDamage, g.gameTime)
				g.spawnExplosion(a.X, a.Y, a.Radius)
				g.sound.PlaySound(systems.SoundHitAsteroid)
				// Play appropriate explosion sound based on asteroid size
				switch a.Size {
				case entities.AsteroidSmall:
					g.sound.PlaySound(systems.SoundExplosionSmall)
				case entities.AsteroidMedium:
					g.sound.PlaySound(systems.SoundExplosionMedium)
				case entities.AsteroidLarge:
					g.sound.PlaySound(systems.SoundExplosionLarge)
				}
				a.Active = false
				g.screenShake = 8
				if g.player.Health <= 0 {
					g.spawnExplosion(g.player.X, g.player.Y, 40)
					g.sound.PlaySound(systems.SoundExplosionLarge)
					g.player.Active = false
				}
			}
		}
	}

	// Player projectiles vs asteroids
	for _, p := range g.projectiles {
		if !p.Active || !p.Friendly {
			continue
		}
		for _, a := range g.asteroids {
			if !a.Active {
				continue
			}
			if g.checkCircleCollision(p.X, p.Y, p.Radius, a.X, a.Y, a.Radius) {
				p.Active = false
				a.TakeDamage(p.Damage)

				// Add impact effect for asteroid
				g.impactEffects = append(g.impactEffects, entities.NewImpactEffect(a.X, a.Y, 25, color.RGBA{200, 100, 50, 255}))

				if !a.Active {
					g.spawnExplosion(a.X, a.Y, a.Radius)
					// Play appropriate explosion sound based on asteroid size
					switch a.Size {
					case entities.AsteroidSmall:
						g.sound.PlaySound(systems.SoundExplosionSmall)
					case entities.AsteroidMedium:
						g.sound.PlaySound(systems.SoundExplosionMedium)
					case entities.AsteroidLarge:
						g.sound.PlaySound(systems.SoundExplosionLarge)
					}
					points := int64(10 + int(a.Radius))
					g.addScore(points)
					g.spawnFloatingScore(a.X, a.Y, int(points)) // Show score popup
				}
			}
		}
	}
}

func (g *Game) checkCircleCollision(x1, y1, r1, x2, y2, r2 float64) bool {
	dx := x2 - x1
	dy := y2 - y1
	dist := math.Sqrt(dx*dx + dy*dy)
	return dist < r1+r2
}

func (g *Game) addScore(points int64) {
	g.score += int64(float64(points) * g.multiplier)
	g.comboTimer = 2.0
	g.multiplier = math.Min(g.multiplier+0.1, 5.0)
}

func (g *Game) spawnExplosion(x, y, size float64) {
	g.explosions = append(g.explosions, entities.NewExplosion(x, y, size))
	// Sound is played by the caller based on enemy/asteroid type
}

func (g *Game) spawnExplosionWithType(x, y, size float64, expType entities.ExplosionType) {
	g.explosions = append(g.explosions, entities.NewExplosionWithType(x, y, size, expType))
	// Sound is played by the caller based on enemy/asteroid type
}

func (g *Game) cleanupEntities() {
	// Clean projectiles
	activeProjectiles := g.projectiles[:0]
	for _, p := range g.projectiles {
		if p.Active {
			activeProjectiles = append(activeProjectiles, p)
		}
	}
	g.projectiles = activeProjectiles

	// Clean enemies
	activeEnemies := g.enemies[:0]
	for _, e := range g.enemies {
		if e.Active {
			activeEnemies = append(activeEnemies, e)
		}
	}
	g.enemies = activeEnemies

	// Clean explosions
	activeExplosions := g.explosions[:0]
	for _, ex := range g.explosions {
		if ex.Active {
			activeExplosions = append(activeExplosions, ex)
		}
	}
	g.explosions = activeExplosions

	// Clean powerups
	activePowerups := g.powerups[:0]
	for _, pu := range g.powerups {
		if pu.Active {
			activePowerups = append(activePowerups, pu)
		}
	}
	g.powerups = activePowerups

	// Clean asteroids
	activeAsteroids := g.asteroids[:0]
	for _, a := range g.asteroids {
		if a.Active {
			activeAsteroids = append(activeAsteroids, a)
		}
	}
	g.asteroids = activeAsteroids

	// Clean floating text
	activeFloatingText := g.floatingTexts[:0]
	for _, ft := range g.floatingTexts {
		if ft.Active {
			activeFloatingText = append(activeFloatingText, ft)
		}
	}
	g.floatingTexts = activeFloatingText

	// Clean impact effects
	activeImpacts := g.impactEffects[:0]
	for _, ie := range g.impactEffects {
		if ie.Active {
			activeImpacts = append(activeImpacts, ie)
		}
	}
	g.impactEffects = activeImpacts
}

// getPerspectiveScale returns a scale factor based on Y position (depth)
// Objects further back (lower Y) are smaller, objects closer (higher Y) are larger
// This creates a 3D sense of depth
func (g *Game) getPerspectiveScale(worldY float64) float64 {
	// Normalize Y position: 0 = far back, 720 = bottom of screen
	normalizedY := worldY / ScreenHeight

	// Scale from 0.5 (far away, 50% size) to 1.2 (very close, 120% size)
	// Objects further away appear smaller
	scale := 0.5 + normalizedY*0.7

	return math.Min(1.2, math.Max(0.5, scale))
}

// getDepthY returns the Y position adjusted for perspective depth rendering
func (g *Game) getDepthY(worldY float64) float64 {
	// Objects further away should appear higher on screen
	// This creates a sense of perspective
	scale := g.getPerspectiveScale(worldY)
	depthOffset := (1.0 - scale) * 30 // Offset based on scale
	return worldY - depthOffset
}

func (g *Game) updateCamera() {
	// Smooth zoom transitions
	zoomDifference := g.cameraTargetZoom - g.cameraZoom
	if math.Abs(zoomDifference) > 0.01 {
		// Smooth interpolation towards target zoom
		g.cameraZoom += zoomDifference * 0.1 // Smooth lerp
	} else {
		g.cameraZoom = g.cameraTargetZoom
	}

	// Handle boss cinematic mode
	if g.bossWave && g.boss != nil && g.boss.Active && !g.cameraCinematicMode {
		// Trigger cinematic zoom on boss appearance
		g.cameraCinematicMode = true
		g.cameraCinematicTimer = 0
		g.cameraTargetZoom = 0.75 // Zoom in more for dramatic effect
		g.cameraShakeAmount = 8.0 // Stronger initial shake
	}

	// Update cinematic mode timer
	if g.cameraCinematicMode {
		g.cameraCinematicTimer += 1.0 / 60.0
		if g.cameraCinematicTimer > 2.5 {
			// Exit cinematic mode after 2.5 seconds
			g.cameraCinematicMode = false
			g.cameraTargetZoom = 1.0 // Return to normal zoom
			g.cameraShakeAmount = 0.0
		}
	}

	// Dynamic zoom based on player danger level
	// More danger = zoom out to see more
	if !g.cameraCinematicMode && !g.bossWave && g.player != nil {
		enemyCount := len(g.enemies)

		// Zoom out when many enemies present
		if enemyCount > 15 {
			g.cameraTargetZoom = 1.15 // More zoom out
		} else if enemyCount > 10 {
			g.cameraTargetZoom = 1.1
		} else if enemyCount > 5 {
			g.cameraTargetZoom = 1.05
		} else {
			g.cameraTargetZoom = 1.0 // Normal
		}

		// Extra zoom out if player health low
		if g.player.Health < g.player.MaxHealth/4 {
			g.cameraTargetZoom += 0.05
		}
	}

	// Camera zoom effects on wave completion
	if g.spawner != nil && g.spawner.WaveCompleted && len(g.enemies) == 0 {
		// Slight zoom out on wave completion for celebration
		g.cameraTargetZoom = 1.1
	} else if !g.cameraCinematicMode && !g.bossWave {
		// Dynamic zoom already handled above
	}

	// Decay screen shake more gradually for better feel
	if g.cameraShakeAmount > 0.1 {
		g.cameraShakeAmount *= 0.92 // Slightly slower decay for better feel
	} else {
		g.cameraShakeAmount = 0
	}

	// Add environmental shake (impacts, explosions)
	if g.screenShake > 0 {
		g.cameraShakeAmount += g.screenShake * 0.7 // More impact shake
	}
}

// applyZoom scales coordinates around screen center based on camera zoom
func (g *Game) applyZoom(x, y float64) (float64, float64) {
	if g.cameraZoom == 1.0 {
		return x, y
	}

	// Center point of screen
	centerX := float64(ScreenWidth) / 2
	centerY := float64(ScreenHeight) / 2

	// Translate to center, apply zoom, translate back
	scaledX := centerX + (x-centerX)*g.cameraZoom
	scaledY := centerY + (y-centerY)*g.cameraZoom

	return scaledX, scaledY
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear with deep space color
	screen.Fill(color.RGBA{5, 5, 15, 255})

	// Calculate screen shake offset (includes camera shake)
	totalShake := g.screenShake + g.cameraShakeAmount
	shakeX, shakeY := 0.0, 0.0
	if totalShake > 0 {
		shakeX = (rand.Float64() - 0.5) * totalShake * 2
		shakeY = (rand.Float64() - 0.5) * totalShake * 2
	}

	// Draw starfield (always)
	g.stars.Draw(screen, shakeX, shakeY)

	switch g.state {
	case StateMenu:
		g.menu.Draw(screen, g.leaderboard, ScreenWidth, ScreenHeight)
		g.menu.InfoMenu.Draw(screen, ScreenWidth, ScreenHeight)
	case StatePlaying, StatePaused:
		g.drawGameplay(screen, shakeX, shakeY)
		if g.state == StatePaused {
			g.drawPauseOverlay(screen)
		}
	case StateGameOver:
		g.drawGameplay(screen, shakeX, shakeY)
		g.drawGameOverOverlay(screen)
	}
}

func (g *Game) drawGameplay(screen *ebiten.Image, shakeX, shakeY float64) {
	// Implement depth-sorted drawing using Y-coordinate (painter's algorithm)
	// Lower Y values (further back in isometric) drawn first

	type drawableEntity struct {
		y    float64
		draw func()
	}

	var entities []drawableEntity

	// Add all drawable entities with their Y positions for sorting

	// Powerups
	for _, pu := range g.powerups {
		if pu.Active {
			puCopy := pu // Capture in closure
			sprite := g.sprites.GetSpriteForPowerUp(int(puCopy.Type))
			entities = append(entities, drawableEntity{
				y: pu.Y,
				draw: func() {
					puCopy.Draw(screen, shakeX, shakeY, sprite, g.sprites.SparkleFrames)
				},
			})
		}
	}

	// Enemies
	for _, e := range g.enemies {
		if e.Active {
			eCopy := e
			spriteCopy := g.sprites.GetSpriteForEnemy(int(e.Type))
			entities = append(entities, drawableEntity{
				y: e.Y,
				draw: func() {
					eCopy.Draw(screen, shakeX, shakeY, spriteCopy)
				},
			})
		}
	}

	// Boss
	if g.boss != nil && g.boss.Active {
		bossCopy := g.boss
		entities = append(entities, drawableEntity{
			y: g.boss.Y,
			draw: func() {
				bossCopy.Draw(screen, shakeX, shakeY)
			},
		})
	}

	// Projectiles
	for _, p := range g.projectiles {
		if p.Active {
			pCopy := p
			var sprite *ebiten.Image
			if pCopy.Friendly {
				sprite = g.sprites.PlayerProjectileSprite
			} else {
				sprite = g.sprites.EnemyProjectileSprite
			}
			entities = append(entities, drawableEntity{
				y: p.Y,
				draw: func() {
					pCopy.Draw(screen, shakeX, shakeY, sprite)
				},
			})
		}
	}

	// Explosions
	for _, ex := range g.explosions {
		if ex.Active {
			exCopy := ex
			entities = append(entities, drawableEntity{
				y: ex.Y,
				draw: func() {
					exCopy.Draw(screen, shakeX, shakeY)
				},
			})
		}
	}

	// Asteroids
	for _, a := range g.asteroids {
		if a.Active {
			aCopy := a
			sprite := g.sprites.GetSpriteForAsteroid(int(aCopy.Size))
			entities = append(entities, drawableEntity{
				y: a.Y,
				draw: func() {
					perspScale := g.getPerspectiveScale(aCopy.Y)
					aCopy.Draw(screen, shakeX, shakeY, perspScale, sprite)
				},
			})
		}
	}

	// Player
	if g.player != nil && g.player.Active {
		playerCopy := g.player
		entities = append(entities, drawableEntity{
			y: g.player.Y,
			draw: func() {
				playerCopy.Draw(screen, shakeX, shakeY)
			},
		})
	}

	// Sort by Y position (lower Y = further back = drawn first)
	// Simple bubble sort is fine for game entities
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			if entities[j].y < entities[i].y {
				entities[i], entities[j] = entities[j], entities[i]
			}
		}
	}

	// Draw in sorted order (back to front)
	for _, entity := range entities {
		entity.draw()
	}

	// Draw HUD (always on top, screen space)
	if g.hud != nil {
		health := 0
		maxHealth := 100
		shield := 0
		weaponLevel := 1
		if g.player != nil {
			health = g.player.Health
			maxHealth = g.player.MaxHealth
			shield = g.player.Shield
			weaponLevel = g.player.WeaponLevel
		}
		g.hud.Draw(screen, g.score, g.wave, g.multiplier, health, maxHealth, shield, weaponLevel, ScreenWidth)

		// Boss indicator
		if g.bossWave && g.boss != nil {
			systems.DrawTextCentered(screen, "!! BOSS BATTLE !!", ScreenWidth/2, 60, 2, color.RGBA{255, 50, 50, 255})
		}
	}

	// Draw floating text (damage/score indicators)
	for _, ft := range g.floatingTexts {
		if ft.Active {
			ft.Draw(screen, shakeX, shakeY)
		}
	}

	// Draw impact effects (hit rings)
	for _, ie := range g.impactEffects {
		if ie.Active {
			ie.Draw(screen, shakeX, shakeY)
		}
	}

	// Draw damage flash overlay
	if g.damageFlash > 0 {
		alpha := uint8(255 * (g.damageFlash / 0.2)) // Fade over 0.2 seconds
		overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
		overlay.Fill(color.RGBA{255, 50, 50, alpha})
		screen.DrawImage(overlay, nil)
	}
}

func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 150})
	screen.DrawImage(overlay, nil)

	systems.DrawTextCentered(screen, "PAUSED", ScreenWidth/2, ScreenHeight/2-40, 4, color.RGBA{255, 255, 255, 255})
	systems.DrawTextCentered(screen, "Press P to Resume", ScreenWidth/2, ScreenHeight/2+20, 2, color.RGBA{200, 200, 200, 255})
	systems.DrawTextCentered(screen, "Press Q to Quit", ScreenWidth/2, ScreenHeight/2+50, 2, color.RGBA{200, 200, 200, 255})

	// Sound toggle in pause menu
	soundStatus := "ON"
	soundColor := color.RGBA{100, 255, 100, 255}
	if !g.menu.SoundEnabled {
		soundStatus = "OFF"
		soundColor = color.RGBA{255, 100, 100, 255}
	}
	systems.DrawTextCentered(screen, "Press S to Toggle Sound: "+soundStatus, ScreenWidth/2, ScreenHeight/2+90, 1.5, soundColor)
}

func (g *Game) drawGameOverOverlay(screen *ebiten.Image) {
	// Semi-transparent overlay
	overlay := ebiten.NewImage(ScreenWidth, ScreenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 180})
	screen.DrawImage(overlay, nil)

	systems.DrawTextCentered(screen, "GAME OVER", ScreenWidth/2, 150, 4, color.RGBA{255, 50, 50, 255})

	scoreText := systems.FormatNumber(g.score)
	systems.DrawTextCentered(screen, "Final Score: "+scoreText, ScreenWidth/2, 220, 2, color.RGBA{255, 255, 100, 255})
	systems.DrawTextCentered(screen, "Wave Reached: "+systems.FormatNumber(int64(g.wave)), ScreenWidth/2, 260, 2, color.RGBA{200, 200, 200, 255})

	if g.nameInputMode {
		systems.DrawTextCentered(screen, "Enter Your Name:", ScreenWidth/2, 320, 2, color.RGBA{255, 255, 255, 255})
		nameDisplay := g.playerName + "_"
		systems.DrawTextCentered(screen, nameDisplay, ScreenWidth/2, 360, 3, color.RGBA{100, 255, 100, 255})
		systems.DrawTextCentered(screen, "Press ENTER to confirm", ScreenWidth/2, 420, 1.5, color.RGBA{150, 150, 150, 255})
	} else {
		// Show leaderboard (local)
		g.leaderboard.Draw(screen, ScreenWidth/2, 320, g.score)

		// Show online leaderboard if available
		if len(g.onlineScores) > 0 {
			g.drawOnlineLeaderboard(screen)
		}

		systems.DrawTextCentered(screen, "Press ENTER to Play Again", ScreenWidth/2, ScreenHeight-100, 2, color.RGBA{100, 255, 100, 255})
		systems.DrawTextCentered(screen, "Press Q for Menu", ScreenWidth/2, ScreenHeight-60, 1.5, color.RGBA{150, 150, 150, 255})
	}
}

// spawnFloatingScore creates a floating score text at the given position
func (g *Game) spawnFloatingScore(x, y float64, score int) {
	ft := entities.NewFloatingScore(x, y, score)
	g.floatingTexts = append(g.floatingTexts, ft)
}

// spawnFloatingDamage creates a floating damage indicator at the given position
func (g *Game) spawnFloatingDamage(x, y float64, damage int) {
	ft := entities.NewFloatingDamage(x, y, damage)
	g.floatingTexts = append(g.floatingTexts, ft)
}

// spawnFloatingUpgrade creates a floating weapon level indicator at the given position
func (g *Game) spawnFloatingUpgrade(x, y float64, level int) {
	ft := entities.NewFloatingUpgrade(x, y, level)
	g.floatingTexts = append(g.floatingTexts, ft)
}

// updateMiniBossSpawning handles spawning mini-bosses during boss battles
func (g *Game) updateMiniBossSpawning() {
	if g.boss == nil || g.boss.BossLevel < 2 {
		return // No mini-bosses for boss level 1
	}

	// Mini-boss spawn configuration based on boss level
	var spawnInterval float64
	var maxMiniBosses int

	switch g.boss.BossLevel {
	case 2: // Wave 10
		spawnInterval = 8.0 // Spawn every 8 seconds
		maxMiniBosses = 3   // Max 3 mini-bosses at a time
	case 3: // Wave 15
		spawnInterval = 6.0 // Spawn every 6 seconds
		maxMiniBosses = 4   // Max 4 mini-bosses at a time
	default: // Wave 20+
		spawnInterval = 4.0 // Spawn every 4 seconds
		maxMiniBosses = 5   // Max 5 mini-bosses at a time
	}

	// Update spawn timer
	g.miniBossSpawnTimer += 1.0 / 60.0

	// Count active mini-bosses (enemies that are hunter or tank types - we'll mark them as mini-bosses)
	activeMiniBosses := 0
	for _, e := range g.enemies {
		if e.Active && (e.Type == entities.EnemyHunter || e.Type == entities.EnemyTank) {
			// Estimate if it's a mini-boss (will be higher health)
			if e.Health > 50 {
				activeMiniBosses++
			}
		}
	}

	// Spawn new mini-boss if conditions are met
	if g.miniBossSpawnTimer >= spawnInterval && activeMiniBosses < maxMiniBosses {
		g.miniBossSpawnTimer = 0

		// Alternate between two enemy types for mini-bosses
		var miniBossType entities.EnemyType
		if g.miniBossesSpawned%2 == 0 {
			miniBossType = entities.EnemyHunter
		} else {
			miniBossType = entities.EnemyTank
		}

		// Spawn mini-boss at random side of screen
		spawnX := float64(ScreenWidth / 2)
		if g.miniBossesSpawned%2 == 0 {
			spawnX = float64(ScreenWidth - 100)
		}

		miniBoss := entities.NewEnemy(spawnX, 50, miniBossType)
		// Enhance it to be a mini-boss (increased health and points)
		miniBoss.MaxHealth = 75 + g.boss.BossLevel*25
		miniBoss.Health = miniBoss.MaxHealth
		miniBoss.Points = 500 * g.boss.BossLevel
		miniBoss.Speed *= 1.2     // Slightly faster
		miniBoss.ShootRate *= 0.8 // Shoots more often

		g.enemies = append(g.enemies, miniBoss)
		g.miniBossesSpawned++
		g.sound.PlaySound(systems.SoundWaveStart) // Alert sound for mini-boss spawn
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// Isometric coordinate transformation functions
// ToIsometric converts 2D world coordinates to isometric projection
func ToIsometric(x, y float64) (float64, float64) {
	// Classic isometric projection (30-60-90 angle)
	isoX := (x - y) * 0.866 // cos(30°) ≈ 0.866
	isoY := (x + y) * 0.5   // sin(30°) / 2
	return isoX, isoY
}

// FromIsometric converts isometric back to world coordinates
func FromIsometric(isoX, isoY float64) (float64, float64) {
	// Inverse transformation
	x := (isoX/0.866 + isoY*2) / 2
	y := (-isoX/0.866 + isoY*2) / 2
	return x, y
}
