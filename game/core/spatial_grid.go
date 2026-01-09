package core

import (
	"stellar-siege/game/entities"
)

// SpatialGrid implements a uniform grid for spatial partitioning
// This reduces collision detection from O(n²) to O(n) by only checking entities in nearby cells
type SpatialGrid struct {
	cellSize   float64
	gridWidth  int
	gridHeight int
	screenW    float64
	screenH    float64

	// Separate grids for different entity types
	enemyGrid      [][]GridCell
	projectileGrid [][]GridCell
	powerupGrid    [][]GridCell
	asteroidGrid   [][]GridCell
}

// GridCell holds entities in a specific grid cell
type GridCell struct {
	Enemies     []*entities.Enemy
	Projectiles []*entities.Projectile
	Powerups    []*entities.PowerUp
	Asteroids   []*entities.Asteroid
}

// NewSpatialGrid creates a new spatial grid for collision optimization
func NewSpatialGrid(screenWidth, screenHeight, cellSize float64) *SpatialGrid {
	gridW := int(screenWidth/cellSize) + 1
	gridH := int(screenHeight/cellSize) + 1

	sg := &SpatialGrid{
		cellSize:   cellSize,
		gridWidth:  gridW,
		gridHeight: gridH,
		screenW:    screenWidth,
		screenH:    screenHeight,
	}

	sg.initializeGrids()
	return sg
}

// initializeGrids creates the grid cell arrays
func (sg *SpatialGrid) initializeGrids() {
	sg.enemyGrid = make([][]GridCell, sg.gridHeight)
	sg.projectileGrid = make([][]GridCell, sg.gridHeight)
	sg.powerupGrid = make([][]GridCell, sg.gridHeight)
	sg.asteroidGrid = make([][]GridCell, sg.gridHeight)

	for i := 0; i < sg.gridHeight; i++ {
		sg.enemyGrid[i] = make([]GridCell, sg.gridWidth)
		sg.projectileGrid[i] = make([]GridCell, sg.gridWidth)
		sg.powerupGrid[i] = make([]GridCell, sg.gridWidth)
		sg.asteroidGrid[i] = make([]GridCell, sg.gridWidth)
	}
}

// Clear clears all entities from the grid (call before repopulating each frame)
func (sg *SpatialGrid) Clear() {
	for y := 0; y < sg.gridHeight; y++ {
		for x := 0; x < sg.gridWidth; x++ {
			sg.enemyGrid[y][x].Enemies = sg.enemyGrid[y][x].Enemies[:0]
			sg.projectileGrid[y][x].Projectiles = sg.projectileGrid[y][x].Projectiles[:0]
			sg.powerupGrid[y][x].Powerups = sg.powerupGrid[y][x].Powerups[:0]
			sg.asteroidGrid[y][x].Asteroids = sg.asteroidGrid[y][x].Asteroids[:0]
		}
	}
}

// getCellCoords returns grid coordinates for a world position
func (sg *SpatialGrid) getCellCoords(x, y float64) (int, int) {
	cellX := int(x / sg.cellSize)
	cellY := int(y / sg.cellSize)

	// Clamp to grid bounds
	if cellX < 0 {
		cellX = 0
	}
	if cellX >= sg.gridWidth {
		cellX = sg.gridWidth - 1
	}
	if cellY < 0 {
		cellY = 0
	}
	if cellY >= sg.gridHeight {
		cellY = sg.gridHeight - 1
	}

	return cellX, cellY
}

// getCellsInRadius returns all cell coordinates within a radius of a position
// This accounts for entities that span multiple cells
func (sg *SpatialGrid) getCellsInRadius(x, y, radius float64) [][2]int {
	cells := make([][2]int, 0, 9) // Typically 1-9 cells

	// Calculate bounding box
	minX := x - radius
	maxX := x + radius
	minY := y - radius
	maxY := y + radius

	// Get cell range
	minCellX, minCellY := sg.getCellCoords(minX, minY)
	maxCellX, maxCellY := sg.getCellCoords(maxX, maxY)

	// Add all cells in range
	for cy := minCellY; cy <= maxCellY; cy++ {
		for cx := minCellX; cx <= maxCellX; cx++ {
			cells = append(cells, [2]int{cx, cy})
		}
	}

	return cells
}

// AddEnemy adds an enemy to the appropriate grid cells
func (sg *SpatialGrid) AddEnemy(enemy *entities.Enemy) {
	if !enemy.Active {
		return
	}

	// Add to all cells the enemy overlaps
	cells := sg.getCellsInRadius(enemy.X, enemy.Y, enemy.Radius)
	for _, cell := range cells {
		sg.enemyGrid[cell[1]][cell[0]].Enemies = append(
			sg.enemyGrid[cell[1]][cell[0]].Enemies, enemy)
	}
}

// AddProjectile adds a projectile to the grid
func (sg *SpatialGrid) AddProjectile(proj *entities.Projectile) {
	if !proj.Active {
		return
	}

	cells := sg.getCellsInRadius(proj.X, proj.Y, proj.Radius)
	for _, cell := range cells {
		sg.projectileGrid[cell[1]][cell[0]].Projectiles = append(
			sg.projectileGrid[cell[1]][cell[0]].Projectiles, proj)
	}
}

// AddPowerup adds a powerup to the grid
func (sg *SpatialGrid) AddPowerup(powerup *entities.PowerUp) {
	if !powerup.Active {
		return
	}

	cells := sg.getCellsInRadius(powerup.X, powerup.Y, powerup.Radius)
	for _, cell := range cells {
		sg.powerupGrid[cell[1]][cell[0]].Powerups = append(
			sg.powerupGrid[cell[1]][cell[0]].Powerups, powerup)
	}
}

// AddAsteroid adds an asteroid to the grid
func (sg *SpatialGrid) AddAsteroid(asteroid *entities.Asteroid) {
	if !asteroid.Active {
		return
	}

	size := float64(asteroid.Size) * 20 // Radius calculation from asteroid drawing
	cells := sg.getCellsInRadius(asteroid.X, asteroid.Y, size)
	for _, cell := range cells {
		sg.asteroidGrid[cell[1]][cell[0]].Asteroids = append(
			sg.asteroidGrid[cell[1]][cell[0]].Asteroids, asteroid)
	}
}

// GetNearbyEnemies returns all enemies near a position
func (sg *SpatialGrid) GetNearbyEnemies(x, y, radius float64) []*entities.Enemy {
	cells := sg.getCellsInRadius(x, y, radius)
	enemies := make([]*entities.Enemy, 0, 32) // Pre-allocate reasonable size

	// Use map to avoid duplicates (entity can be in multiple cells)
	seen := make(map[*entities.Enemy]bool)

	for _, cell := range cells {
		for _, enemy := range sg.enemyGrid[cell[1]][cell[0]].Enemies {
			if !seen[enemy] {
				enemies = append(enemies, enemy)
				seen[enemy] = true
			}
		}
	}

	return enemies
}

// GetNearbyProjectiles returns all projectiles near a position
func (sg *SpatialGrid) GetNearbyProjectiles(x, y, radius float64) []*entities.Projectile {
	cells := sg.getCellsInRadius(x, y, radius)
	projectiles := make([]*entities.Projectile, 0, 64)

	seen := make(map[*entities.Projectile]bool)

	for _, cell := range cells {
		for _, proj := range sg.projectileGrid[cell[1]][cell[0]].Projectiles {
			if !seen[proj] {
				projectiles = append(projectiles, proj)
				seen[proj] = true
			}
		}
	}

	return projectiles
}

// GetNearbyPowerups returns all powerups near a position
func (sg *SpatialGrid) GetNearbyPowerups(x, y, radius float64) []*entities.PowerUp {
	cells := sg.getCellsInRadius(x, y, radius)
	powerups := make([]*entities.PowerUp, 0, 8)

	seen := make(map[*entities.PowerUp]bool)

	for _, cell := range cells {
		for _, powerup := range sg.powerupGrid[cell[1]][cell[0]].Powerups {
			if !seen[powerup] {
				powerups = append(powerups, powerup)
				seen[powerup] = true
			}
		}
	}

	return powerups
}

// GetNearbyAsteroids returns all asteroids near a position
func (sg *SpatialGrid) GetNearbyAsteroids(x, y, radius float64) []*entities.Asteroid {
	cells := sg.getCellsInRadius(x, y, radius)
	asteroids := make([]*entities.Asteroid, 0, 16)

	seen := make(map[*entities.Asteroid]bool)

	for _, cell := range cells {
		for _, asteroid := range sg.asteroidGrid[cell[1]][cell[0]].Asteroids {
			if !seen[asteroid] {
				asteroids = append(asteroids, asteroid)
				seen[asteroid] = true
			}
		}
	}

	return asteroids
}

// PopulateGrid adds all entities to the grid (call once per frame before collision detection)
func (sg *SpatialGrid) PopulateGrid(
	enemies []*entities.Enemy,
	projectiles []*entities.Projectile,
	powerups []*entities.PowerUp,
	asteroids []*entities.Asteroid,
) {
	sg.Clear()

	for _, enemy := range enemies {
		sg.AddEnemy(enemy)
	}

	for _, proj := range projectiles {
		sg.AddProjectile(proj)
	}

	for _, powerup := range powerups {
		sg.AddPowerup(powerup)
	}

	for _, asteroid := range asteroids {
		sg.AddAsteroid(asteroid)
	}
}
