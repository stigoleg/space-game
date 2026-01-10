package entities

import (
	"math"
)

// Update updates enemy AI behavior and movement
func (e *Enemy) Update(playerX, playerY float64, screenWidth, screenHeight int) {
	e.AnimTimer += 0.1
	e.ShootTimer += 1.0 / 60.0

	// Update burning DoT
	e.UpdateBurning()

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
	case EnemySniper:
		// Stay near top of screen, slight horizontal drift
		targetY := 80.0 // Stay near top
		if e.Y < targetY {
			e.VelY = e.Speed
		} else if e.Y > targetY+20 {
			e.VelY = -e.Speed * 0.5
		} else {
			e.VelY = 0
		}
		// Slow drift side to side
		e.Phase += 0.02
		e.VelX = math.Sin(e.Phase) * 0.8
		e.X += e.VelX
		e.Y += e.VelY

		// Lock-on timer
		e.SniperLockTimer += 1.0 / 60.0
		if e.SniperLockTimer >= 1.5 { // 1.5 second lock-on time
			e.SniperLocked = true
			e.SniperTargetX = playerX
			e.SniperTargetY = playerY
		}
	case EnemySplitter:
		// Wave pattern similar to drone
		e.Phase += 0.05
		e.VelX = math.Sin(e.Phase) * 2.5
		e.VelY = e.Speed
		e.X += e.VelX
		e.Y += e.VelY
	case EnemyShieldBearer:
		// Slow advance straight down
		e.VelY = e.Speed
		e.Y += e.VelY

		// Shield regeneration (2 shields per second after 3 seconds)
		e.ShieldRegenTimer += 1.0 / 60.0
		if e.ShieldRegenTimer >= 3.5 && e.ShieldPoints < e.MaxShieldPoints {
			// Regen 1 shield every 0.5 seconds
			e.ShieldPoints++
			e.ShieldRegenTimer = 3.0
			if e.ShieldPoints > e.MaxShieldPoints {
				e.ShieldPoints = e.MaxShieldPoints
			}
		}
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

// UpdateFormation updates enemy position based on formation behavior
// Performance note: This iterates through all enemies to find formation allies.
// Optimized with early termination checks (FormationID, FormationType) so only
// formation members are processed. Typical formations have 3-7 enemies.
func (e *Enemy) UpdateFormation(allEnemies []*Enemy) {
	if e.FormationType == FormationTypeNone {
		return
	}

	// Update nearby allies - reuse slice to avoid allocations
	e.NearbyAllies = e.NearbyAllies[:0] // Clear but keep capacity
	const allyDetectionRange = 150.0
	const allyDetectionRangeSq = allyDetectionRange * allyDetectionRange

	for _, other := range allEnemies {
		if other != e && other.FormationType == e.FormationType && other.FormationID == e.FormationID && other.Active {
			dx := other.X - e.X
			dy := other.Y - e.Y
			distSq := dx*dx + dy*dy // Avoid sqrt for comparison
			if distSq < allyDetectionRangeSq {
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
