package systems

import (
	"math"
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

	if ws.spawnTimer >= ws.spawnDelay && ws.enemiesLeft > 0 {
		ws.spawnTimer = 0
		ws.enemiesLeft--

		newEnemies := ws.spawnEnemyBatch()

		if ws.enemiesLeft <= 0 {
			ws.WaveCompleted = true
		}

		return newEnemies
	}

	return nil
}

// spawnEnemyBatch spawns multiple enemies based on wave difficulty
func (ws *WaveSpawner) spawnEnemyBatch() []*entities.Enemy {
	spawnCount := getSpawnCount(ws.currentWave)
	var newEnemies []*entities.Enemy

	for i := 0; i < spawnCount && ws.enemiesLeft >= 0; i++ {
		enemy := ws.spawnEnemy()
		if enemy != nil {
			newEnemies = append(newEnemies, enemy)
		}
		if ws.enemiesLeft > 0 {
			ws.enemiesLeft--
		}
	}

	return newEnemies
}

// calculateSpawnPosition returns a random spawn position at the top of the screen
func (ws *WaveSpawner) calculateSpawnPosition() (float64, float64) {
	margin := 50.0
	x := margin + rand.Float64()*(float64(ws.width)-margin*2)
	y := -30.0
	return x, y
}

func (ws *WaveSpawner) spawnEnemy() *entities.Enemy {
	// Use configuration-based enemy selection
	enemyType := selectEnemyType(ws.currentWave)

	// Calculate spawn position
	x, y := ws.calculateSpawnPosition()

	return entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
}

// SpawnFormation spawns a group of enemies in a coordinated formation
func (ws *WaveSpawner) SpawnFormation(wave int) []*entities.Enemy {
	// Determine formation type based on wave progression
	formationType := ws.selectFormationType(wave)

	// Create formation based on type
	switch formationType {
	case entities.FormationTypeVFormation:
		return ws.spawnVFormation(wave)
	case entities.FormationTypeCircular:
		return ws.spawnCircularFormation(wave)
	case entities.FormationTypeWave:
		return ws.spawnWaveFormation(wave)
	case entities.FormationTypePincer:
		return ws.spawnPincerFormation(wave)
	case entities.FormationTypeConvoy:
		return ws.spawnConvoyFormation(wave)
	default:
		return nil
	}
}

// selectFormationType chooses appropriate formation based on wave
func (ws *WaveSpawner) selectFormationType(wave int) entities.FormationType {
	r := rand.Float64()

	// Earlier waves: simpler formations
	if wave < 8 {
		if r < 0.4 {
			return entities.FormationTypeVFormation
		} else if r < 0.7 {
			return entities.FormationTypeWave
		}
		return entities.FormationTypeConvoy
	}

	// Later waves: all formations available
	if r < 0.2 {
		return entities.FormationTypeVFormation
	} else if r < 0.4 {
		return entities.FormationTypeCircular
	} else if r < 0.6 {
		return entities.FormationTypeWave
	} else if r < 0.8 {
		return entities.FormationTypePincer
	}
	return entities.FormationTypeConvoy
}

// Helper function to set formation properties on an enemy
func setFormationProperties(enemy *entities.Enemy, formationType entities.FormationType, formationID, index int, isLeader bool) {
	enemy.FormationType = formationType
	enemy.FormationID = formationID
	enemy.FormationIndex = index
	enemy.IsFormationLeader = isLeader
}

// spawnVFormation creates a V-shaped formation
func (ws *WaveSpawner) spawnVFormation(wave int) []*entities.Enemy {
	formationID := rand.Intn(10000)
	count := 3 + rand.Intn(3) // 3-5 enemies

	// Center position
	centerX := float64(ws.width) / 2.0
	centerY := -50.0

	enemies := make([]*entities.Enemy, count)

	// Create leader
	enemyType := ws.selectFormationEnemyType(wave)
	enemies[0] = entities.NewEnemyWithDifficulty(centerX, centerY, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
	setFormationProperties(enemies[0], entities.FormationTypeVFormation, formationID, 0, true)

	// Create V wings
	spacing := 60.0
	for i := 1; i < count; i++ {
		side := float64(1)
		if i%2 == 0 {
			side = -1
		}
		offset := float64((i+1)/2) * spacing

		x := centerX + side*offset
		y := centerY + float64((i+1)/2)*30.0

		enemyType := ws.selectFormationEnemyType(wave)
		enemies[i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[i], entities.FormationTypeVFormation, formationID, i, false)
		enemies[i].FormationTargetX = x
		enemies[i].FormationTargetY = y
	}

	return enemies
}

// spawnCircularFormation creates enemies in a circular pattern
func (ws *WaveSpawner) spawnCircularFormation(wave int) []*entities.Enemy {
	formationID := rand.Intn(10000)
	count := 4 + rand.Intn(3) // 4-6 enemies

	centerX := float64(ws.width) / 2.0
	centerY := 100.0 // Start lower for circular formation
	radius := 80.0

	enemies := make([]*entities.Enemy, count)

	for i := 0; i < count; i++ {
		angle := float64(i) * (2.0 * math.Pi / float64(count))
		x := centerX + math.Cos(angle)*radius
		y := centerY + math.Sin(angle)*radius

		enemyType := ws.selectFormationEnemyType(wave)
		enemies[i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[i], entities.FormationTypeCircular, formationID, i, i == 0)
		enemies[i].FormationTargetX = x
		enemies[i].FormationTargetY = y
	}

	return enemies
}

// spawnWaveFormation creates enemies in a wave pattern
func (ws *WaveSpawner) spawnWaveFormation(wave int) []*entities.Enemy {
	formationID := rand.Intn(10000)
	count := 5 + rand.Intn(3) // 5-7 enemies

	spacing := 70.0
	startX := float64(ws.width)/2.0 - (float64(count-1)*spacing)/2.0
	y := -50.0

	enemies := make([]*entities.Enemy, count)

	for i := 0; i < count; i++ {
		x := startX + float64(i)*spacing

		enemyType := ws.selectFormationEnemyType(wave)
		enemies[i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[i], entities.FormationTypeWave, formationID, i, i == count/2)
		enemies[i].Phase = float64(i) * 0.5 // Offset wave phase
	}

	return enemies
}

// spawnPincerFormation creates two groups attacking from sides
func (ws *WaveSpawner) spawnPincerFormation(wave int) []*entities.Enemy {
	formationID := rand.Intn(10000)
	countPerSide := 2 + rand.Intn(2) // 2-3 per side
	totalCount := countPerSide * 2

	enemies := make([]*entities.Enemy, totalCount)

	// Left flank
	for i := 0; i < countPerSide; i++ {
		x := 50.0
		y := -50.0 - float64(i)*40.0

		enemyType := ws.selectFormationEnemyType(wave)
		enemies[i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[i], entities.FormationTypePincer, formationID, i*2, i == 0)
	}

	// Right flank
	for i := 0; i < countPerSide; i++ {
		x := float64(ws.width) - 50.0
		y := -50.0 - float64(i)*40.0

		enemyType := ws.selectFormationEnemyType(wave)
		enemies[countPerSide+i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[countPerSide+i], entities.FormationTypePincer, formationID, i*2+1, i == 0)
	}

	return enemies
}

// spawnConvoyFormation creates a line of enemies with a leader
func (ws *WaveSpawner) spawnConvoyFormation(wave int) []*entities.Enemy {
	formationID := rand.Intn(10000)
	count := 3 + rand.Intn(3) // 3-5 enemies

	centerX := float64(ws.width) / 2.0
	spacing := 50.0

	enemies := make([]*entities.Enemy, count)

	for i := 0; i < count; i++ {
		x := centerX
		y := -50.0 - float64(i)*spacing

		// Leader is tougher enemy type
		var enemyType entities.EnemyType
		if i == 0 {
			enemyType = ws.selectToughFormationEnemyType(wave)
		} else {
			enemyType = ws.selectFormationEnemyType(wave)
		}

		enemies[i] = entities.NewEnemyWithDifficulty(x, y, enemyType, ws.enemyHealthMult, ws.enemySpeedMult)
		setFormationProperties(enemies[i], entities.FormationTypeConvoy, formationID, i, i == 0)
		enemies[i].FormationTargetX = x
		enemies[i].FormationTargetY = float64(100 + i*60) // Target Y positions
	}

	return enemies
}

// selectFormationEnemyType chooses enemy type suitable for formations
func (ws *WaveSpawner) selectFormationEnemyType(wave int) entities.EnemyType {
	// Formations use more organized enemy types (not splitters/bombers)
	r := rand.Float64()

	if wave <= 5 {
		if r < 0.5 {
			return entities.EnemyScout
		}
		return entities.EnemyDrone
	} else if wave <= 10 {
		if r < 0.3 {
			return entities.EnemyScout
		} else if r < 0.6 {
			return entities.EnemyDrone
		}
		return entities.EnemyHunter
	} else {
		if r < 0.25 {
			return entities.EnemyDrone
		} else if r < 0.5 {
			return entities.EnemyHunter
		} else if r < 0.75 {
			return entities.EnemySniper
		}
		return entities.EnemyTank
	}
}

// selectToughFormationEnemyType chooses tougher enemy for formation leaders
func (ws *WaveSpawner) selectToughFormationEnemyType(wave int) entities.EnemyType {
	r := rand.Float64()

	if wave <= 8 {
		if r < 0.5 {
			return entities.EnemyTank
		}
		return entities.EnemyHunter
	} else {
		if r < 0.4 {
			return entities.EnemyTank
		} else if r < 0.7 {
			return entities.EnemySniper
		}
		return entities.EnemyShieldBearer
	}
}
