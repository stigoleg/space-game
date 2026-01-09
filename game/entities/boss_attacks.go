package entities

import "math"

// executeSpreadShot creates a spread shot pattern with level-based projectile count
func (b *Boss) executeSpreadShot() []*Projectile {
	var projectiles []*Projectile
	spread := 2 // Base spread
	if b.BossLevel >= 2 {
		spread = 3
	}
	for i := -spread; i <= spread; i++ {
		angle := math.Pi/2 + float64(i)*0.2
		projectiles = append(projectiles, b.createProjectile(angle))
	}
	return projectiles
}

// executeAimedShot creates aimed shots at the player with side projectiles
func (b *Boss) executeAimedShot(playerX, playerY float64) []*Projectile {
	var projectiles []*Projectile
	dx := playerX - b.X
	dy := playerY - b.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist > 0 {
		velX := (dx / dist) * 7
		velY := (dy / dist) * 7
		projectiles = append(projectiles, NewProjectile(b.X, b.Y+b.Radius, velX, velY, false, b.Damage))

		// Add side shots for higher levels
		sideCount := 1
		if b.BossLevel >= 3 {
			sideCount = 2
		}
		for s := 1; s <= sideCount; s++ {
			offsetX := float64(s * 30)
			projectiles = append(projectiles, NewProjectile(b.X-offsetX, b.Y+b.Radius, velX*0.8, velY*0.8, false, b.Damage-5))
			projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+b.Radius, velX*0.8, velY*0.8, false, b.Damage-5))
		}
	}

	return projectiles
}

// executeCircularBurst creates a 360-degree burst pattern
func (b *Boss) executeCircularBurst() []*Projectile {
	var projectiles []*Projectile
	count := 8 // Base count
	if b.BossLevel >= 2 {
		count = 12
	}
	if b.BossLevel >= 4 {
		count = 16
	}
	for i := 0; i < count; i++ {
		angle := float64(i) * (2 * math.Pi) / float64(count)
		projectiles = append(projectiles, b.createProjectile(angle))
	}
	return projectiles
}

// executeLaserLines creates vertical laser line pattern
func (b *Boss) executeLaserLines() []*Projectile {
	var projectiles []*Projectile
	count := 3 // Base count
	if b.BossLevel >= 2 {
		count = 5
	}
	for i := 0; i < count; i++ {
		offsetX := float64(i-(count-1)/2) * 40
		projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+b.Radius, 0, 6, false, b.Damage))
	}
	return projectiles
}

// executeSpiralPattern creates a rotating spiral pattern (boss level 2+)
func (b *Boss) executeSpiralPattern() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 2 {
		for i := 0; i < 6; i++ {
			angle := float64(i)*math.Pi/3 + b.AnimTimer*0.1
			velX := math.Cos(angle) * 5
			velY := math.Sin(angle) * 5
			projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage-5))
		}
	}
	return projectiles
}

// executeDoubleArc creates a double arc pattern (boss level 2+)
func (b *Boss) executeDoubleArc() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 2 {
		for i := -3; i <= 3; i++ {
			angle := math.Pi/2 + float64(i)*0.15
			velX := math.Cos(angle) * 6
			velY := math.Sin(angle) * 6
			projectiles = append(projectiles, NewProjectile(b.X-40, b.Y+b.Radius, velX, velY, false, b.Damage-5))
			projectiles = append(projectiles, NewProjectile(b.X+40, b.Y+b.Radius, velX, velY, false, b.Damage-5))
		}
	}
	return projectiles
}

// executeTrackingSpiral creates a faster tracking spiral (boss level 3+)
func (b *Boss) executeTrackingSpiral() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 3 {
		for i := 0; i < 6; i++ {
			angle := float64(i)*math.Pi/3 + b.AnimTimer*0.2
			velX := math.Cos(angle) * 5.5
			velY := math.Sin(angle) * 5.5
			projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage))
		}
	}
	return projectiles
}

// executeWavePattern creates a wave pattern (boss level 3+)
func (b *Boss) executeWavePattern() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 3 {
		waveCount := 7
		for i := 0; i < waveCount; i++ {
			offsetX := float64(i-(waveCount-1)/2) * 30
			waveY := math.Sin(float64(i)*math.Pi/5) * 50
			projectiles = append(projectiles, NewProjectile(b.X+offsetX, b.Y+waveY, 0, 6, false, b.Damage-5))
		}
	}
	return projectiles
}

// executeCrossBurst creates a cross burst pattern (boss level 4+)
func (b *Boss) executeCrossBurst() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 4 {
		for i := 0; i < 4; i++ {
			angle := float64(i)*math.Pi/2 + math.Pi/4
			for j := 0; j < 3; j++ {
				ratio := float64(j) / 2.0
				vel := 4.0 + ratio*3.0
				velX := math.Cos(angle) * vel
				velY := math.Sin(angle) * vel
				projectiles = append(projectiles, NewProjectile(b.X, b.Y+b.Radius, velX, velY, false, b.Damage))
			}
		}
	}
	return projectiles
}

// executeChaosPattern creates a chaotic multi-directional pattern (boss level 4+)
func (b *Boss) executeChaosPattern() []*Projectile {
	var projectiles []*Projectile
	if b.BossLevel >= 4 {
		for i := 0; i < 10; i++ {
			angle := float64(i)*math.Pi/5 + b.AnimTimer*0.3
			speed := 4.0 + math.Sin(b.AnimTimer+float64(i))*2.0
			velX := math.Cos(angle) * speed
			velY := math.Sin(angle) * speed
			projectiles = append(projectiles, NewProjectile(b.X, b.Y, velX, velY, false, b.Damage))
		}
	}
	return projectiles
}
