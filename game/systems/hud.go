package systems

import (
	"fmt"
	"image/color"
	"math"
	"stellar-siege/game/entities"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type HUD struct{}

func NewHUD() *HUD {
	return &HUD{}
}

func (h *HUD) Draw(screen *ebiten.Image, score int64, wave int, multiplier float64, health, maxHealth, shield, weaponLevel int, screenWidth int) {
	// Score display (top right)
	scoreText := FormatNumber(score)
	DrawText(screen, "SCORE: "+scoreText, screenWidth-200, 30, 1.5, color.RGBA{255, 255, 255, 255})

	// Wave display (top center)
	waveText := fmt.Sprintf("WAVE %d", wave)
	DrawTextCentered(screen, waveText, screenWidth/2, 30, 2, color.RGBA{255, 200, 50, 255})

	// Multiplier (below score)
	if multiplier > 1.0 {
		multText := fmt.Sprintf("x%.1f", multiplier)
		multColor := color.RGBA{255, 100 + uint8(multiplier*30), 50, 255}
		DrawText(screen, multText, screenWidth-100, 55, 1.2, multColor)
	}

	// Health bar (top left)
	drawBar(screen, 20, 20, 200, 20, float64(health)/float64(maxHealth),
		color.RGBA{50, 200, 50, 255}, color.RGBA{30, 30, 30, 200}, "HP")

	// Shield bar (below health)
	if shield > 0 {
		drawBar(screen, 20, 45, 150, 15, float64(shield)/50.0,
			color.RGBA{50, 150, 255, 255}, color.RGBA{30, 30, 30, 200}, "SHIELD")
	}

	// Weapon level indicator
	weaponText := fmt.Sprintf("WEAPON LV.%d", weaponLevel)
	DrawText(screen, weaponText, 20, 80, 1.0, color.RGBA{255, 200, 100, 255})
}

func drawBar(screen *ebiten.Image, x, y, width, height float32, ratio float64, fillColor, bgColor color.RGBA, label string) {
	// Background
	vector.DrawFilledRect(screen, x, y, width, height, bgColor, true)

	// Fill
	fillWidth := float32(ratio * float64(width))
	if fillWidth > 0 {
		vector.DrawFilledRect(screen, x, y, fillWidth, height, fillColor, true)
	}

	// Border
	vector.StrokeRect(screen, x, y, width, height, 1, color.RGBA{255, 255, 255, 150}, true)

	// Label
	DrawText(screen, label, int(x+5), int(y+height-5), 0.8, color.RGBA{255, 255, 255, 255})
}

// DrawWeaponInfo draws current weapon information with icon and cooldown
func (h *HUD) DrawWeaponInfo(screen *ebiten.Image, weaponName, weaponEmoji string, weaponLevel int, fireTimer, fireRate, gameTime float64) {
	x := float32(20)
	y := float32(80)

	// Draw background panel
	panelWidth := float32(220)
	panelHeight := float32(60)
	bgColor := color.RGBA{20, 20, 30, 200}
	vector.DrawFilledRect(screen, x, y, panelWidth, panelHeight, bgColor, true)
	vector.StrokeRect(screen, x, y, panelWidth, panelHeight, 2, color.RGBA{100, 150, 200, 255}, true)

	// Draw weapon emoji/icon (larger)
	emojiX := int(x + 15)
	emojiY := int(y + 30)
	DrawText(screen, weaponEmoji, emojiX, emojiY, 2.5, color.RGBA{255, 255, 255, 255})

	// Draw weapon name
	nameX := int(x + 55)
	nameY := int(y + 20)
	DrawText(screen, weaponName, nameX, nameY, 1.0, color.RGBA{255, 220, 100, 255})

	// Draw level indicator (Mk I - Mk V)
	levelText := fmt.Sprintf("Mk %s", getLevelRoman(weaponLevel))
	levelY := int(y + 38)
	DrawText(screen, levelText, nameX, levelY, 0.9, color.RGBA{150, 200, 255, 255})

	// Draw cooldown bar
	cooldownBarX := x + 55
	cooldownBarY := y + 45
	cooldownBarWidth := float32(155)
	cooldownBarHeight := float32(8)

	// Calculate cooldown ratio (how much time until next shot)
	cooldownInterval := 1.0 / fireRate // Time between shots
	cooldownRatio := 1.0 - (fireTimer / cooldownInterval)
	if cooldownRatio < 0 {
		cooldownRatio = 0
	}
	if cooldownRatio > 1 {
		cooldownRatio = 1
	}

	// Background
	vector.DrawFilledRect(screen, cooldownBarX, cooldownBarY, cooldownBarWidth, cooldownBarHeight, color.RGBA{40, 40, 50, 200}, true)

	// Fill bar (green when ready, yellow when cooling down)
	var fillColor color.RGBA
	if cooldownRatio >= 0.95 {
		// Ready to fire - pulsing green
		pulse := uint8(200 + 55*math.Sin(gameTime*8))
		fillColor = color.RGBA{50, pulse, 50, 255}
	} else {
		// Cooling down - gradient from red to yellow
		redAmount := uint8(255)
		greenAmount := uint8(cooldownRatio * 200)
		fillColor = color.RGBA{redAmount, greenAmount, 0, 255}
	}

	fillWidth := float32(cooldownRatio * float64(cooldownBarWidth))
	if fillWidth > 0 {
		vector.DrawFilledRect(screen, cooldownBarX, cooldownBarY, fillWidth, cooldownBarHeight, fillColor, true)
	}

	// Border
	vector.StrokeRect(screen, cooldownBarX, cooldownBarY, cooldownBarWidth, cooldownBarHeight, 1, color.RGBA{150, 150, 180, 255}, true)
}

// getLevelRoman converts weapon level (1-5) to Roman numerals
func getLevelRoman(level int) string {
	switch level {
	case 1:
		return "I"
	case 2:
		return "II"
	case 3:
		return "III"
	case 4:
		return "IV"
	case 5:
		return "V"
	default:
		return "I"
	}
}

// FormatNumber formats a number with commas
func FormatNumber(n int64) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	result := ""
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result += ","
		}
		result += string(c)
	}
	return result
}

// DrawText draws text at the specified position
func DrawText(screen *ebiten.Image, str string, x, y int, scale float64, clr color.RGBA) {
	face := text.NewGoXFace(basicfont.Face7x13)
	op := &text.DrawOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, op)
}

// DrawTextCentered draws text centered at the specified position
func DrawTextCentered(screen *ebiten.Image, str string, x, y int, scale float64, clr color.RGBA) {
	face := text.NewGoXFace(basicfont.Face7x13)

	// Calculate text width for centering
	width := float64(len(str)*7) * scale // Approximate width

	op := &text.DrawOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x)-width/2, float64(y))
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, face, op)
}

// DrawAbilities draws ability cooldown indicators
func (h *HUD) DrawAbilities(screen *ebiten.Image, abilities []*entities.Ability, screenWidth, screenHeight int) {
	if len(abilities) == 0 {
		return
	}

	// Position abilities in bottom-right corner
	startX := float32(screenWidth - 320)
	startY := float32(screenHeight - 80)
	iconSize := float32(50)
	spacing := float32(10)

	for i, ability := range abilities {
		x := startX + float32(i)*(iconSize+spacing)
		y := startY

		// Draw background panel
		bgColor := color.RGBA{20, 20, 30, 180}
		vector.DrawFilledRect(screen, x, y, iconSize, iconSize, bgColor, true)

		// Draw border (color changes based on cooldown status)
		var borderColor color.RGBA
		if ability.CooldownTimer > 0 {
			// On cooldown - red
			borderColor = color.RGBA{200, 50, 50, 255}
		} else {
			// Ready - green
			borderColor = color.RGBA{50, 200, 50, 255}
		}
		vector.StrokeRect(screen, x, y, iconSize, iconSize, 2, borderColor, true)

		// Draw ability icon/emoji
		emojiX := int(x + 15)
		emojiY := int(y + 28)
		emojiColor := color.RGBA{255, 255, 255, 255}
		if ability.CooldownTimer > 0 {
			// Dim icon when on cooldown
			emojiColor = color.RGBA{150, 150, 150, 255}
		}
		DrawText(screen, ability.IconEmoji, emojiX, emojiY, 1.8, emojiColor)

		// Draw cooldown overlay
		if ability.CooldownTimer > 0 {
			cooldownRatio := ability.CooldownTimer / ability.Cooldown
			overlayHeight := float32(cooldownRatio) * iconSize
			overlayColor := color.RGBA{0, 0, 0, 150}
			vector.DrawFilledRect(screen, x, y, iconSize, overlayHeight, overlayColor, true)

			// Draw cooldown time text
			timeLeft := fmt.Sprintf("%.1f", ability.CooldownTimer)
			timeX := int(x + iconSize/2 - 10)
			timeY := int(y + iconSize/2 + 5)
			DrawText(screen, timeLeft, timeX, timeY, 0.9, color.RGBA{255, 255, 255, 255})
		}

		// Draw key binding at bottom
		keyX := int(x + iconSize/2 - 5)
		keyY := int(y + iconSize - 8)
		DrawText(screen, ability.KeyBinding, keyX, keyY, 0.8, color.RGBA{200, 200, 200, 255})
	}
}
