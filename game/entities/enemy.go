package entities

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type EnemyType int

const (
	EnemyScout EnemyType = iota
	EnemyDrone
	EnemyHunter
	EnemyTank
	EnemyBomber
)

// FormationType represents different enemy formation patterns
type FormationType int

const (
	FormationTypeNone FormationType = iota
	FormationTypeVFormation
	FormationTypeCircular
	FormationTypeWave
	FormationTypePincer
	FormationTypeConvoy
)

type Enemy struct {
	X, Y       float64
	VelX, VelY float64
	Radius     float64
	Speed      float64
	Health     int
	MaxHealth  int
	Points     int
	Type       EnemyType
	Active     bool
	ShootTimer float64
	ShootRate  float64
	AnimTimer  float64
	Phase      float64 // For wave movement

	// Formation system
	FormationType     FormationType
	FormationID       int     // Groups enemies in same formation
	IsFormationLeader bool    // Is this the formation leader?
	FormationTargetX  float64 // Target position for formation
	FormationTargetY  float64
	FormationIndex    int      // Position in formation (0 = leader)
	NearbyAllies      []*Enemy // References to nearby allies in formation
	LastShootTime     float64
	CoorditatedShoot  bool // Should coordinate fire with formation
}

func NewEnemy(x, y float64, enemyType EnemyType) *Enemy {
	e := &Enemy{
		X:         x,
		Y:         y,
		Type:      enemyType,
		Active:    true,
		Phase:     rand.Float64() * math.Pi * 2,
		AnimTimer: 0,
	}

	switch enemyType {
	case EnemyScout:
		e.Radius = 15
		e.Speed = 4
		e.Health = 20
		e.MaxHealth = 20
		e.Points = 100
		e.ShootRate = 0 // Doesn't shoot
	case EnemyDrone:
		e.Radius = 18
		e.Speed = 2.5
		e.Health = 30
		e.MaxHealth = 30
		e.Points = 150
		e.ShootRate = 2.0
	case EnemyHunter:
		e.Radius = 20
		e.Speed = 3
		e.Health = 50
		e.MaxHealth = 50
		e.Points = 250
		e.ShootRate = 1.5
	case EnemyTank:
		e.Radius = 30
		e.Speed = 1.5
		e.Health = 100
		e.MaxHealth = 100
		e.Points = 400
		e.ShootRate = 1.0
	case EnemyBomber:
		e.Radius = 22
		e.Speed = 3.5
		e.Health = 40
		e.MaxHealth = 40
		e.Points = 300
		e.ShootRate = 0 // Explodes instead
	}

	return e
}

// NewEnemyWithDifficulty creates an enemy with difficulty adjustments
func NewEnemyWithDifficulty(x, y float64, enemyType EnemyType, healthMult, speedMult float64) *Enemy {
	e := NewEnemy(x, y, enemyType)

	// Apply difficulty multipliers
	e.Health = int(float64(e.Health) * healthMult)
	e.MaxHealth = e.Health
	e.Speed = e.Speed * speedMult

	return e
}

func (e *Enemy) Update(playerX, playerY float64, screenWidth, screenHeight int) {
	e.AnimTimer += 0.1
	e.ShootTimer += 1.0 / 60.0

	switch e.Type {
	case EnemyScout:
		// Straight down movement
		e.VelY = e.Speed
		e.Y += e.VelY
	case EnemyDrone:
		// Wave pattern
		e.Phase += 0.05
		e.VelX = math.Sin(e.Phase) * 3
		e.VelY = e.Speed
		e.X += e.VelX
		e.Y += e.VelY
	case EnemyHunter:
		// Track player horizontally
		dx := playerX - e.X
		if math.Abs(dx) > 5 {
			if dx > 0 {
				e.VelX = e.Speed * 0.8
			} else {
				e.VelX = -e.Speed * 0.8
			}
		} else {
			e.VelX = 0
		}
		e.VelY = e.Speed * 0.6
		e.X += e.VelX
		e.Y += e.VelY
	case EnemyTank:
		// Slow descent
		e.VelY = e.Speed
		e.Y += e.VelY
	case EnemyBomber:
		// Dive toward player
		dx := playerX - e.X
		dy := playerY - e.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			e.VelX = (dx / dist) * e.Speed
			e.VelY = (dy / dist) * e.Speed
		}
		e.X += e.VelX
		e.Y += e.VelY
	}

	// Keep in bounds horizontally
	if e.X < e.Radius {
		e.X = e.Radius
	}
	if e.X > float64(screenWidth)-e.Radius {
		e.X = float64(screenWidth) - e.Radius
	}

	// Deactivate if off screen (bottom)
	if e.Y > float64(screenHeight)+50 {
		e.Active = false
	}
}

func (e *Enemy) TryShoot() *Projectile {
	if e.ShootRate <= 0 || e.ShootTimer < e.ShootRate {
		return nil
	}
	e.ShootTimer = 0

	switch e.Type {
	case EnemyDrone, EnemyHunter:
		return NewProjectile(e.X, e.Y+e.Radius, 0, 6, false, 10)
	case EnemyTank:
		return NewProjectile(e.X, e.Y+e.Radius, 0, 5, false, 20)
	}
	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image, shakeX, shakeY float64, sprite *ebiten.Image) {
	// Simple screen coordinates with shake
	x := float32(e.X + shakeX)
	y := float32(e.Y + shakeY)

	// Pulsing effect
	pulse := float32(1.0 + 0.1*math.Sin(e.AnimTimer*2))

	// Damage-based color shift (redder when damaged)
	healthRatio := float32(e.Health) / float32(e.MaxHealth)

	// If sprite is provided, use sprite-based rendering
	if sprite != nil {
		e.drawSpriteBased(screen, x, y, pulse, healthRatio, sprite)
	} else {
		// Fallback to procedural rendering
		e.drawProcedural(screen, x, y, pulse, healthRatio)
	}

	// Health bar for tanks
	if e.Type == EnemyTank && e.Health < e.MaxHealth {
		barWidth := float32(60)
		barHeight := float32(6)
		healthRatioBar := float32(e.Health) / float32(e.MaxHealth)
		vector.DrawFilledRect(screen, x-barWidth/2, y-float32(e.Radius)-15, barWidth, barHeight, color.RGBA{50, 50, 50, 200}, true)
		vector.DrawFilledRect(screen, x-barWidth/2, y-float32(e.Radius)-15, barWidth*healthRatioBar, barHeight, color.RGBA{255, 50, 50, 255}, true)
	}

	// Formation indicator: glow ring if in formation
	if e.FormationType != FormationTypeNone {
		formationGlowColor := color.RGBA{100, 255, 200, 80}
		if e.IsFormationLeader {
			formationGlowColor = color.RGBA{255, 255, 100, 100} // Gold for leader
		}
		vector.StrokeCircle(screen, x, y, float32(e.Radius)+8, 1.5, formationGlowColor, true)
	}

	// Crown indicator for formation leader
	if e.IsFormationLeader {
		crownY := y - float32(e.Radius) - 8
		vector.DrawFilledCircle(screen, x, crownY, 4, color.RGBA{255, 255, 100, 255}, true)
		vector.DrawFilledCircle(screen, x-6, crownY-3, 2, color.RGBA{255, 255, 100, 255}, true)
		vector.DrawFilledCircle(screen, x+6, crownY-3, 2, color.RGBA{255, 255, 100, 255}, true)
	}
}

func (e *Enemy) drawSpriteBased(screen *ebiten.Image, x, y, pulse, healthRatio float32, sprite *ebiten.Image) {
	// Draw shadow (depth indicator)
	radius := float32(e.Radius) * pulse
	shadowColor := color.RGBA{20, 20, 30, 100}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Draw glow effect BEFORE sprite (so sprite appears on top)
	var glowColor color.RGBA
	switch e.Type {
	case EnemyScout:
		glowColor = color.RGBA{255, 100, 50, 80}
	case EnemyDrone:
		glowColor = color.RGBA{200, 100, 255, 80}
	case EnemyHunter:
		glowColor = color.RGBA{100, 255, 100, 80}
	case EnemyTank:
		glowColor = color.RGBA{255, 150, 50, 80}
	case EnemyBomber:
		glowColor = color.RGBA{255, 200, 0, 100}
	}

	// Draw glow as background
	glowSize := radius * 1.4
	vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)

	// Draw sprite with options
	op := &ebiten.DrawImageOptions{}

	// Scale sprite to match enemy size
	spriteBounds := sprite.Bounds()
	spriteWidth := float64(spriteBounds.Dx())
	spriteHeight := float64(spriteBounds.Dy())

	// Calculate scale to match radius (sprite should be about 2x radius)
	targetSize := float64(e.Radius) * 2.0 * float64(pulse)
	scaleX := targetSize / spriteWidth
	scaleY := targetSize / spriteHeight

	op.GeoM.Scale(scaleX, scaleY)

	// Translate to enemy position (center sprite)
	op.GeoM.Translate(float64(x)-targetSize/2, float64(y)-targetSize/2)

	// Apply damage color shift (redder when damaged)
	if healthRatio < 1.0 {
		damageShift := 1.0 - healthRatio
		op.ColorScale.Scale(1.0+damageShift*0.3, 1.0-damageShift*0.5, 1.0-damageShift*0.5, 1.0)
	}

	screen.DrawImage(sprite, op)

	// Draw outline stroke for better visibility (after sprite)
	outlineColor := glowColor
	outlineColor.A = 200
	vector.StrokeCircle(screen, x, y, radius*1.2, 3, outlineColor, true)

	// Damage indicator: flickering outline when heavily damaged
	if healthRatio < 0.3 {
		damageAlpha := uint8(150 + 100*math.Sin(float64(e.AnimTimer)*4))
		damageColor := color.RGBA{255, 50, 50, damageAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.3, 2, damageColor, true)
	}
}

func (e *Enemy) drawProcedural(screen *ebiten.Image, x, y, pulse, healthRatio float32) {
	damageShift := 1.0 - healthRatio

	var mainColor, coreColor, glowColor color.RGBA

	switch e.Type {
	case EnemyScout:
		// Scout: Simple fast wedge shape - red/orange
		mainColor = color.RGBA{uint8(220 + damageShift*30), uint8(80 - damageShift*30), 60, 255}
		coreColor = color.RGBA{255, 150, 100, 255}
		glowColor = color.RGBA{255, 100, 50, 100}
	case EnemyDrone:
		// Drone: Rounded diamond - purple
		mainColor = color.RGBA{uint8(180 + damageShift*20), uint8(80 - damageShift*30), uint8(220 - damageShift*50), 255}
		coreColor = color.RGBA{255, 150, 255, 255}
		glowColor = color.RGBA{200, 100, 255, 100}
	case EnemyHunter:
		// Hunter: Angular with fins - green
		mainColor = color.RGBA{uint8(80 - damageShift*40), uint8(220 - damageShift*50), 120, 255}
		coreColor = color.RGBA{150, 255, 150, 255}
		glowColor = color.RGBA{100, 255, 100, 100}
	case EnemyTank:
		// Tank: Massive hexagon - gray with gold core
		mainColor = color.RGBA{uint8(140 + damageShift*40), uint8(140 - damageShift*30), uint8(140 - damageShift*30), 255}
		coreColor = color.RGBA{255, 200, 50, 255}
		glowColor = color.RGBA{255, 150, 50, 100}
	case EnemyBomber:
		// Bomber: Bulbous - orange with aggression
		mainColor = color.RGBA{255, uint8(140 - damageShift*40), 40, 255}
		coreColor = color.RGBA{255, 255, 100, 255}
		glowColor = color.RGBA{255, 200, 0, 150}
	}

	radius := float32(e.Radius) * pulse

	// Draw shadow (depth indicator)
	shadowColor := color.RGBA{20, 20, 30, 100}
	vector.DrawFilledCircle(screen, x, y+radius+5, radius*0.5, shadowColor, true)

	// Inner core (pulsing glow)
	coreSize := radius * 0.35 * pulse

	// Draw ship type-specific designs
	switch e.Type {
	case EnemyScout:
		// Scout: Small fast wedge pointing down
		drawTriangle(screen, x, y+radius*0.8, x-radius*0.6, y-radius*0.6, x+radius*0.6, y-radius*0.6, mainColor)
		vector.DrawFilledCircle(screen, x, y, coreSize*0.4, coreColor, true)

	case EnemyDrone:
		// Drone: Diamond shape
		drawTriangle(screen, x, y-radius*0.8, x-radius*0.8, y, x, y+radius*0.8, mainColor)
		drawTriangle(screen, x, y-radius*0.8, x+radius*0.8, y, x, y+radius*0.8, mainColor)
		vector.DrawFilledCircle(screen, x, y, radius*0.3*pulse, coreColor, true)

	case EnemyHunter:
		// Hunter: Angular shape with fins
		// Main body (triangle pointing down)
		drawTriangle(screen, x, y+radius*0.9, x-radius*0.5, y-radius*0.4, x+radius*0.5, y-radius*0.4, mainColor)
		// Left fin
		drawTriangle(screen, x-radius*0.5, y-radius*0.4, x-radius*1.0, y, x-radius*0.5, y+radius*0.3, mainColor)
		// Right fin
		drawTriangle(screen, x+radius*0.5, y-radius*0.4, x+radius*1.0, y, x+radius*0.5, y+radius*0.3, mainColor)
		vector.DrawFilledCircle(screen, x, y, radius*0.3, coreColor, true)

	case EnemyTank:
		// Tank: Heavy hexagon
		// Draw as layered circles to simulate hexagon
		for i := 0; i < 6; i++ {
			angle := float64(i) * math.Pi / 3
			px := x + float32(math.Cos(angle))*radius*0.7
			py := y + float32(math.Sin(angle))*radius*0.7
			vector.DrawFilledCircle(screen, px, py, radius*0.5, mainColor, true)
		}
		// Heavy core
		vector.DrawFilledCircle(screen, x, y, radius*0.5*pulse, coreColor, true)

	case EnemyBomber:
		// Bomber: Large bulbous oval shape
		// Top
		vector.DrawFilledCircle(screen, x, y-radius*0.5, radius*0.6, mainColor, true)
		// Middle (largest)
		vector.DrawFilledCircle(screen, x, y, radius*0.9, mainColor, true)
		// Bottom
		vector.DrawFilledCircle(screen, x, y+radius*0.6, radius*0.7, mainColor, true)
		vector.DrawFilledCircle(screen, x, y, radius*0.4*pulse, coreColor, true)
	}

	// Pulsing core
	vector.DrawFilledCircle(screen, x, y, coreSize, coreColor, true)

	// Outer glow
	glowSize := radius + 3
	vector.DrawFilledCircle(screen, x, y, glowSize, glowColor, true)

	// Highlight edge (3D effect)
	highlightSize := radius * 0.25
	highlightColor := color.RGBA{mainColor.R + 40, mainColor.G + 40, mainColor.B + 40, 150}
	vector.DrawFilledCircle(screen, x-radius*0.3, y-radius*0.3, highlightSize, highlightColor, true)

	// Damage indicator: flickering outline when heavily damaged
	if healthRatio < 0.3 {
		damageAlpha := uint8(150 + 100*math.Sin(float64(e.AnimTimer)*4))
		damageColor := color.RGBA{255, 50, 50, damageAlpha}
		vector.StrokeCircle(screen, x, y, radius*1.1, 2, damageColor, true)
	}
}

// UpdateFormation updates enemy position based on formation behavior
func (e *Enemy) UpdateFormation(allEnemies []*Enemy) {
	if e.FormationType == FormationTypeNone {
		return
	}

	// Update nearby allies
	e.NearbyAllies = nil
	const allyDetectionRange = 150.0

	for _, other := range allEnemies {
		if other != e && other.FormationType == e.FormationType && other.FormationID == e.FormationID && other.Active {
			dx := other.X - e.X
			dy := other.Y - e.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < allyDetectionRange {
				e.NearbyAllies = append(e.NearbyAllies, other)
			}
		}
	}

	// Apply formation-specific behavior
	switch e.FormationType {
	case FormationTypeVFormation:
		e.updateVFormation()
	case FormationTypeCircular:
		e.updateCircularFormation()
	case FormationTypeWave:
		e.updateWaveFormation()
	case FormationTypePincer:
		e.updatePincerFormation()
	case FormationTypeConvoy:
		e.updateConvoyFormation()
	}
}

func (e *Enemy) updateVFormation() {
	// V-shape formation around leader
	// Enemies follow leader and maintain V-shape
	if !e.IsFormationLeader {
		// Calculate angle offset from leader based on index
		angleOffset := float64(e.FormationIndex) * 0.3

		// Move toward formation target
		targetDx := e.FormationTargetX - e.X
		targetDy := e.FormationTargetY - e.Y
		dist := math.Sqrt(targetDx*targetDx + targetDy*targetDy)

		if dist > 5 {
			e.VelX = (targetDx / dist) * e.Speed * 0.8
			e.VelY = (targetDy / dist) * e.Speed * 0.8
		}

		// Maintain formation angle
		if len(e.NearbyAllies) > 0 && e.FormationIndex%2 == 0 {
			e.VelX += math.Sin(angleOffset) * 1.5
		}
	}
}

func (e *Enemy) updateCircularFormation() {
	// Enemies orbit around a center point
	if !e.IsFormationLeader && len(e.NearbyAllies) > 0 {
		// Calculate rotation around center
		centerX := (e.X + e.FormationTargetX) / 2.0
		centerY := (e.Y + e.FormationTargetY) / 2.0

		angle := math.Atan2(e.Y-centerY, e.X-centerX)
		angle += 0.02 // Rotation speed

		orbitRadius := 80.0
		e.FormationTargetX = centerX + orbitRadius*math.Cos(angle)
		e.FormationTargetY = centerY + orbitRadius*math.Sin(angle)

		// Move toward target
		targetDx := e.FormationTargetX - e.X
		targetDy := e.FormationTargetY - e.Y
		dist := math.Sqrt(targetDx*targetDx + targetDy*targetDy)

		if dist > 5 {
			e.VelX = (targetDx / dist) * e.Speed
			e.VelY = (targetDy / dist) * e.Speed
		}
	}
}

func (e *Enemy) updateWaveFormation() {
	// All enemies move together in undulating pattern
	e.Phase += 0.03
	e.VelX = math.Sin(e.Phase) * 2.5
	e.VelY = e.Speed
}

func (e *Enemy) updatePincerFormation() {
	// Split formation: left and right flanks
	if e.FormationIndex%2 == 0 {
		e.VelX = -2.0 // Move left
	} else {
		e.VelX = 2.0 // Move right
	}
	e.VelY = e.Speed * 0.8
}

func (e *Enemy) updateConvoyFormation() {
	// Follow the leader in single file
	if e.IsFormationLeader {
		e.VelY = e.Speed
	} else if len(e.NearbyAllies) > 0 {
		// Follow leader
		leader := e.NearbyAllies[0]
		targetDx := leader.X - e.X
		targetDy := leader.Y - e.Y
		dist := math.Sqrt(targetDx*targetDx + targetDy*targetDy)

		// Maintain distance behind leader
		if dist > 70 {
			e.VelX = (targetDx / dist) * e.Speed
			e.VelY = (targetDy / dist) * e.Speed
		} else if dist < 50 {
			e.VelX = 0
			e.VelY = e.Speed * 0.5
		}
	}
}

// JoinFormation makes this enemy join a formation
func (e *Enemy) JoinFormation(formationType FormationType, formationID int, isLeader bool, index int) {
	e.FormationType = formationType
	e.FormationID = formationID
	e.IsFormationLeader = isLeader
	e.FormationIndex = index
}
