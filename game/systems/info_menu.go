package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InfoMenu struct {
	isActive      bool
	scrollY       float64
	maxScroll     float64
	spriteManager *SpriteManager
}

func NewInfoMenu(spriteManager *SpriteManager) *InfoMenu {
	return &InfoMenu{
		isActive:      false,
		scrollY:       0,
		maxScroll:     0,
		spriteManager: spriteManager,
	}
}

func (im *InfoMenu) Show() {
	im.isActive = true
	im.scrollY = 0
}

func (im *InfoMenu) Hide() {
	im.isActive = false
}

func (im *InfoMenu) IsActive() bool {
	return im.isActive
}

func (im *InfoMenu) Update() {
	if !im.isActive {
		return
	}

	// Handle return to menu
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyB) {
		im.Hide()
	}

	// Handle scrolling (arrow keys or mouse wheel)
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		im.scrollY -= 30
		if im.scrollY < 0 {
			im.scrollY = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		im.scrollY += 30
		if im.scrollY > im.maxScroll {
			im.scrollY = im.maxScroll
		}
	}
}

func (im *InfoMenu) Draw(screen *ebiten.Image, screenWidth, screenHeight int) {
	if !im.isActive {
		return
	}

	// Draw dark background with moderate opacity
	backImg := ebiten.NewImage(screenWidth, screenHeight)
	backImg.Fill(color.RGBA{15, 15, 30, 220}) // Dark blue-black with 220 alpha (opacity)
	screen.DrawImage(backImg, nil)

	// Draw title at top (no background bar)
	DrawTextCentered(screen, "=== GAME INFORMATION ===", screenWidth/2, 20, 3, color.RGBA{100, 200, 255, 255})
	DrawTextCentered(screen, "Created by King123", screenWidth/2, 70, 2, color.RGBA{255, 200, 100, 255})

	// Content area starts after title
	y := 140 - int(im.scrollY)
	lineHeight := 54 // Height per item (1.2x higher than 45)
	contentStartY := 140
	contentEndY := screenHeight - 40

	im.maxScroll = 0

	// Helper function to draw section with icon and text aligned to same row
	drawSectionWithIcons := func(title string, titleColor color.RGBA, items []struct {
		name   string
		desc   string
		sprite *ebiten.Image
	}) {
		// Section title
		if y >= contentStartY-50 && y <= contentEndY {
			DrawTextCentered(screen, title, screenWidth/2, y, 2, titleColor)
		}
		y += 45 // Space after section title

		// Section content - each item on its own row with aligned sprite
		for _, item := range items {
			if y >= contentStartY && y <= contentEndY {
				// Draw icon on the left, vertically centered with text baseline
				if item.sprite != nil {
					iconOp := &ebiten.DrawImageOptions{}
					// Reduce scale by 1.5 (0.8 / 1.5 ≈ 0.533)
					spriteScale := 0.533
					iconOp.GeoM.Scale(spriteScale, spriteScale)
					// Position sprite: left side, centered on row
					spriteX := float64(100)
					spriteY := float64(y - 5) // Adjust for smaller sprite
					iconOp.GeoM.Translate(spriteX, spriteY)
					screen.DrawImage(item.sprite, iconOp)
				}

				// Draw name and description (left-aligned text after sprite area)
				nameAndDesc := item.name + " - " + item.desc
				// Draw text at fixed left position (after sprite area)
				x := 160 // Adjusted X position for smaller sprites
				DrawText(screen, nameAndDesc, x, y, 1.3, color.RGBA{220, 220, 200, 255})
			}
			y += lineHeight
		}
		y += 15 // Extra spacing between sections
		im.maxScroll = float64(y - contentEndY)
	}

	// ENEMIES SECTION
	drawSectionWithIcons(">> ENEMIES <<", color.RGBA{255, 150, 100, 255},
		[]struct {
			name   string
			desc   string
			sprite *ebiten.Image
		}{
			{"Scout", "Fast and weak, fires single shots", im.spriteManager.ScoutSprite},
			{"Drone", "Medium speed with coordinated fire", im.spriteManager.DroneSprite},
			{"Hunter", "Fast with tracking capability", im.spriteManager.HunterSprite},
			{"Tank", "Slow, heavily armored, high damage", im.spriteManager.TankSprite},
			{"Bomber", "Drops explosives in patterns", im.spriteManager.BomberSprite},
			{"Sniper", "Locks on and shoots precise fast shots", im.spriteManager.SniperSprite},
			{"Splitter", "Splits into 2 scouts when destroyed", im.spriteManager.SplitterSprite},
			{"Shield Bearer", "Slow tank with regenerating shield", im.spriteManager.ShieldBearerSprite},
		})

	// POWER-UPS SECTION
	drawSectionWithIcons(">> POWER-UPS <<", color.RGBA{255, 200, 100, 255},
		[]struct {
			name   string
			desc   string
			sprite *ebiten.Image
		}{
			{"Health", "Restore player HP", im.spriteManager.PowerUpHealthSprite},
			{"Shield", "Replenish shield capacity", im.spriteManager.PowerUpShieldSprite},
			{"Weapon", "Upgrade weapon to next level", im.spriteManager.PowerUpWeaponSprite},
			{"Speed", "Boost movement speed", im.spriteManager.PowerUpSpeedSprite},
			{"Mystery", "Random effect (60% good, 40% bad)", im.spriteManager.PowerUpMysterySprite},
		})

	// GAME MODES SECTION
	drawSectionWithIcons(">> GAME MODES <<", color.RGBA{150, 255, 200, 255},
		[]struct {
			name   string
			desc   string
			sprite *ebiten.Image
		}{
			{"Easy", "120 HP, 60 Shield, 0.7x Enemy Spawn", nil},
			{"Normal", "85 HP, 40 Shield, 1.0x Enemy Spawn", nil},
			{"Hard", "60 HP, 25 Shield, 1.35x Enemy Spawn", nil},
		})

	// TEXT-ONLY SECTIONS
	y += 10

	// WEAPON PROGRESSION SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> WEAPON PROGRESSION <<", screenWidth/2, y, 2, color.RGBA{150, 200, 255, 255})
	}
	y += 45

	weaponProgressionLines := []string{
		"Level 1 - Basic Gun: 1 shot forward",
		"Level 2 - Dual Guns: 2 shots forward",
		"Level 3 - Spread Shot: 4 shots (2 center + 2 angled)",
		"Level 4 - Wide Spread: 6 shots (2 center + 4 angled)",
		"Level 5 - Maximum Spread: 8 shots (2 center + 6 angled)",
		"",
		"After Level 5: Unlock special weapons!",
	}

	for _, line := range weaponProgressionLines {
		if y >= contentStartY && y <= contentEndY {
			textColor := color.RGBA{180, 200, 220, 255}
			if line == "" {
				y += lineHeight - 30
				continue
			}
			if line == "After Level 5: Unlock special weapons!" {
				textColor = color.RGBA{255, 220, 100, 255}
			}
			DrawTextCentered(screen, line, screenWidth/2, y, 1.2, textColor)
		}
		y += lineHeight - 10
	}
	y += 15
	im.maxScroll = float64(y - contentEndY)

	// SPECIAL WEAPONS SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> SPECIAL WEAPONS <<", screenWidth/2, y, 2, color.RGBA{255, 180, 100, 255})
	}
	y += 45

	specialWeaponLines := []string{
		"Following Rockets - Homing missiles",
		"Chain Lightning - Arcs between enemies",
		"Flamethrower - Burning DoT damage",
		"Ion Beam - Continuous piercing beam",
		"Blaster - High damage single shots",
		"Laser - Focused energy beam",
		"Shotgun - Close-range spread",
		"Plasma - Explosive projectiles",
		"Homing - Target-seeking missiles",
		"Railgun - Armor-piercing rounds",
		"",
		"Mixed Mode: Special weapon + 4 side blasters!",
	}

	for _, line := range specialWeaponLines {
		if y >= contentStartY && y <= contentEndY {
			textColor := color.RGBA{180, 200, 220, 255}
			if line == "" {
				y += lineHeight - 30
				continue
			}
			if line == "Mixed Mode: Special weapon + 6 side blasters!" {
				textColor = color.RGBA{255, 220, 100, 255}
			}
			DrawTextCentered(screen, line, screenWidth/2, y, 1.2, textColor)
		}
		y += lineHeight - 10
	}
	y += 15
	im.maxScroll = float64(y - contentEndY)

	// MYSTERY POWER-UP EFFECTS SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> MYSTERY POWER-UP EFFECTS <<", screenWidth/2, y, 2, color.RGBA{255, 150, 255, 255})
	}
	y += 45

	if y >= contentStartY && y <= contentEndY {
		DrawTextCentered(screen, "POSITIVE (60% chance):", screenWidth/2, y, 1.3, color.RGBA{100, 255, 100, 255})
	}
	y += lineHeight - 5

	positiveEffects := []string{
		"Super Weapon Upgrade - Instant weapon level up",
		"Speed Boost - 2.5x speed for 8 seconds",
		"Shield Overcharge - +30 shield instantly",
		"Rapid Fire - 2.5x fire rate for 10 seconds",
		"Invincibility - Immune to damage for 5 seconds",
		"Score Multiplier - 3x score for 10 seconds",
	}

	for _, line := range positiveEffects {
		if y >= contentStartY && y <= contentEndY {
			DrawTextCentered(screen, line, screenWidth/2, y, 1.1, color.RGBA{150, 220, 150, 255})
		}
		y += lineHeight - 10
	}
	y += 10

	if y >= contentStartY && y <= contentEndY {
		DrawTextCentered(screen, "NEGATIVE (40% chance):", screenWidth/2, y, 1.3, color.RGBA{255, 100, 100, 255})
	}
	y += lineHeight - 5

	negativeEffects := []string{
		"Weapon Downgrade - Lose 1 weapon level",
		"Engine Malfunction - 0.5x speed for 5 seconds",
		"Shield Drain - Lose 20 shield points",
		"Fire Rate Reduction - 0.5x fire rate for 5 sec",
		"Control Reversal - Inverted controls for 4 sec",
	}

	for _, line := range negativeEffects {
		if y >= contentStartY && y <= contentEndY {
			DrawTextCentered(screen, line, screenWidth/2, y, 1.1, color.RGBA{220, 150, 150, 255})
		}
		y += lineHeight - 10
	}
	y += 15
	im.maxScroll = float64(y - contentEndY)

	// ABILITIES SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> ABILITIES <<", screenWidth/2, y, 2, color.RGBA{100, 255, 100, 255})
	}
	y += 45

	abilityLines := []string{
		"Dash - Quickly dodge incoming fire",
		"Bullet Time - Slow time temporarily",
		"Shield Barrier - Temporary protection",
		"Weapon Overcharge - Increase weapon output",
		"EMP Pulse - Disable nearby enemies",
		"Orbital Defense - Automated turret support",
	}

	for _, line := range abilityLines {
		if y >= contentStartY && y <= contentEndY {
			DrawTextCentered(screen, line, screenWidth/2, y, 1.2, color.RGBA{150, 220, 150, 255})
		}
		y += lineHeight - 10
	}
	y += 15
	im.maxScroll = float64(y - contentEndY)

	// HAZARDS SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> HAZARDS <<", screenWidth/2, y, 2, color.RGBA{255, 100, 150, 255})
	}
	y += 45

	hazardLines := []string{
		"Energy Barriers - Block movement/projectiles",
		"Magnetic Fields - Pull enemies toward them",
		"Radiation Zones - Deal damage over time",
		"Black Holes - Gravitational pull hazard",
	}

	for _, line := range hazardLines {
		if y >= contentStartY && y <= contentEndY {
			DrawTextCentered(screen, line, screenWidth/2, y, 1.2, color.RGBA{220, 150, 180, 255})
		}
		y += lineHeight - 10
	}
	y += 15
	im.maxScroll = float64(y - contentEndY)

	// BOSS MECHANICS SECTION
	if y >= contentStartY-50 && y <= contentEndY {
		DrawTextCentered(screen, ">> BOSS MECHANICS <<", screenWidth/2, y, 2, color.RGBA{255, 150, 200, 255})
	}
	y += 45

	bossLines := []string{
		"Progressive difficulty scaling through waves",
		"Wave 5: First boss encounter",
		"Wave 10, 15, 20+: Escalating boss battles",
		"Increasing HP and damage with waves",
	}

	for _, line := range bossLines {
		if y >= contentStartY && y <= contentEndY {
			DrawTextCentered(screen, line, screenWidth/2, y, 1.2, color.RGBA{220, 180, 200, 255})
		}
		y += lineHeight - 10
	}

	im.maxScroll = float64(y - contentEndY)

	// Footer instructions (no background bar)
	footerY := screenHeight - 20
	DrawTextCentered(screen, "Press ESC or B to return | Use UP/DOWN arrows to scroll", screenWidth/2, footerY, 1.0, color.RGBA{150, 200, 200, 255})
}
