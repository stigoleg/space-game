package systems

import (
	"math/rand"

	"stellar-siege/game/entities"
)

type WaveDefinition struct {
	EnemyTypes []entities.EnemyType
	Count      int
	Delay      float64 // Delay between spawns
}

type WaveSpawner struct {
	width            int
	height           int
	currentWave      int
	spawnTimer       float64
	spawnDelay       float64
	enemiesLeft      int
	WaveCompleted    bool
	waveStarted      bool
	spawnMultiplier  float64
	enemyHealthMult  float64
	enemySpeedMult   float64
	damageMultiplier float64
}

func NewWaveSpawner(width, height int) *WaveSpawner {
	return &WaveSpawner{
		width:            width,
		height:           height,
		currentWave:      0,
		spawnTimer:       0,
		spawnDelay:       1.0,
		enemiesLeft:      0,
		WaveCompleted:    true,
		waveStarted:      false,
		spawnMultiplier:  1.0,
		enemyHealthMult:  1.0,
		enemySpeedMult:   1.0,
		damageMultiplier: 1.0,
	}
}

func (ws *WaveSpawner) SetDifficultyMultipliers(spawnMult, healthMult, speedMult, damageMult float64) {
	ws.spawnMultiplier = spawnMult
	ws.enemyHealthMult = healthMult
	ws.enemySpeedMult = speedMult
	ws.damageMultiplier = damageMult
}

func (ws *WaveSpawner) StartWave(wave int) {
	ws.currentWave = wave
	ws.WaveCompleted = false
	ws.waveStarted = true

	// Calculate enemies based on wave with difficulty multiplier
	baseCount := int(float64(5+wave*2) * ws.spawnMultiplier)
	if baseCount > 30 {
		baseCount = 30
	}
	ws.enemiesLeft = baseCount

	// Decrease spawn delay as waves progress
	ws.spawnDelay = 1.5 - float64(wave)*0.1
	if ws.spawnDelay < 0.3 {
		ws.spawnDelay = 0.3
	}
	ws.spawnTimer = 0
}

func (ws *WaveSpawner) Update(gameTime float64, currentWave int) []*entities.Enemy {
	// Start first wave
	if !ws.waveStarted && currentWave == 0 {
		ws.StartWave(1)
	}

	if ws.WaveCompleted {
		return nil
	}

	ws.spawnTimer += 1.0 / 60.0

	var newEnemies []*entities.Enemy

	if ws.spawnTimer >= ws.spawnDelay && ws.enemiesLeft > 0 {
		ws.spawnTimer = 0
		ws.enemiesLeft--

		// Spawn 1-3 enemies at once
		spawnCount := 1
		if ws.currentWave > 3 {
			spawnCount = rand.Intn(2) + 1
		}
		if ws.currentWave > 7 {
			spawnCount = rand.Intn(3) + 1
		}

		for i := 0; i < spawnCount && ws.enemiesLeft >= 0; i++ {
			enemy := ws.spawnEnemy()
			if enemy != nil {
				newEnemies = append(newEnemies, enemy)
			}
			if ws.enemiesLeft > 0 {
				ws.enemiesLeft--
			}
		}

		if ws.enemiesLeft <= 0 {
			ws.WaveCompleted = true
		}
	}

	return newEnemies
}

func (ws *WaveSpawner) spawnEnemy() *entities.Enemy {
	// Choose enemy type based on wave
	var enemyType entities.EnemyType

	r := rand.Float64()
	wave := ws.currentWave

	if wave <= 2 {
		// Early waves: mostly scouts
		if r < 0.7 {
			enemyType = entities.EnemyScout
		} else {
			enemyType = entities.EnemyDrone
		}
	} else if wave <= 5 {
		// Mid-early: add hunters
		if r < 0.4 {
			enemyType = entities.EnemyScout
		} else if r < 0.7 {
			enemyType = entities.EnemyDrone
		} else {
			enemyType = entities.EnemyHunter
		}
	} else if wave <= 8 {
		// Mid: add tanks
		if r < 0.25 {
			enemyType = entities.EnemyScout
		} else if r < 0.45 {
			enemyType = entities.EnemyDrone
		} else if r < 0.7 {
			enemyType = entities.EnemyHunter
		} else {
			enemyType = entities.EnemyTank
		}
	} else {
		// Late: all types including bombers
		if r < 0.15 {
			enemyType = entities.EnemyScout
		} else if r < 0.35 {
			enemyType = entities.EnemyDrone
		} else if r < 0.55 {
			enemyType = entities.EnemyHunter
		} else if r < 0.75 {
			enemyType = entities.EnemyTank
		} else {
			enemyType = entities.EnemyBomber
		}
	}

	// Random X position
	margin := 50.0
	x := margin + rand.Float64()*(float64(ws.width)-margin*2)
	y := -30.0

	return entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
}
