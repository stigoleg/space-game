package systems

import (
	"fmt"
	"image/color"

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
