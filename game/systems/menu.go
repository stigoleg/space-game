package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Menu struct {
	ShowDifficultySelect bool // Exported so Game can access it
	SelectedDifficulty   int  // 0=Easy, 1=Normal, 2=Hard
	showLeaderboard      bool
	InfoMenu             *InfoMenu // Pointer to info menu - exported
	animTimer            float64
	SoundEnabled         bool           // Track sound toggle state
	spriteManager        *SpriteManager // For info menu sprites

	// Update banner fields
	updateStatus     UpdateStatus
	updateVersion    string
	downloadProgress float64
	showUpdateBanner bool
	bannerPulse      float64
	updateManager    *UpdateManager // Reference to call InstallUpdate()
}

func NewMenu(spriteManager *SpriteManager) *Menu {
	infoMenu := NewInfoMenu(spriteManager)
	return &Menu{
		showLeaderboard:      false,
		ShowDifficultySelect: false,
		SelectedDifficulty:   1, // Default to normal
		InfoMenu:             infoMenu,
		animTimer:            0,
		SoundEnabled:         true, // Sound enabled by default
		spriteManager:        spriteManager,
	}
}

func (m *Menu) ToggleLeaderboard() {
	m.showLeaderboard = !m.showLeaderboard
}

func (m *Menu) ShowDifficultySelectMenu() {
	m.ShowDifficultySelect = true
	m.showLeaderboard = false
	m.InfoMenu.Hide()
}

func (m *Menu) ShowInfo() {
	m.InfoMenu.Show()
	m.showLeaderboard = false
}

// SetUpdateManager sets the update manager reference for the menu
func (m *Menu) SetUpdateManager(updateManager *UpdateManager) {
	m.updateManager = updateManager
}

// SetUpdateStatus updates the menu with the current update status
func (m *Menu) SetUpdateStatus(status UpdateStatus, version string, progress float64) {
	m.updateStatus = status
	m.updateVersion = version
	m.downloadProgress = progress

	// Show banner when update is ready or downloading
	if status == UpdateStatusReady || status == UpdateStatusDownloading {
		m.showUpdateBanner = true
	}
}

func (m *Menu) Update() {
	// Update banner pulse animation
	m.bannerPulse += 0.05

	// Handle update banner click
	if m.showUpdateBanner && m.updateStatus == UpdateStatusReady {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			_, my := ebiten.CursorPosition()
			// Banner is at top of screen, 60 pixels tall
			if my < 60 && m.updateManager != nil {
				m.updateManager.InstallUpdate()
				return
			}
		}
	}

	// Update info menu if active
	if m.InfoMenu.IsActive() {
		m.InfoMenu.Update()
		return // Don't process menu input while info menu is active
	}

	// Handle difficulty selection input
	if m.ShowDifficultySelect {
		// Arrow keys or A/D to move selection
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			if m.SelectedDifficulty > 0 {
				m.SelectedDifficulty--
			} else {
				m.SelectedDifficulty = 2 // Hard
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if m.SelectedDifficulty < 2 {
				m.SelectedDifficulty++
			} else {
				m.SelectedDifficulty = 0 // Easy
			}
		}
	}
}

func (m *Menu) Draw(screen *ebiten.Image, leaderboard *Leaderboard, screenWidth, screenHeight int) {
	m.animTimer += 0.02

	// Title with pulsing effect
	titleScale := 4.0 + 0.3*math.Sin(m.animTimer*2)
	titleY := 150 + int(math.Sin(m.animTimer)*5)

	// Draw title glow
	for i := 3; i > 0; i-- {
		alpha := uint8(50 / i)
		DrawTextCentered(screen, "STELLAR SIEGE", screenWidth/2, titleY-i, titleScale+float64(i)*0.2, color.RGBA{100, 150, 255, alpha})
	}
	DrawTextCentered(screen, "STELLAR SIEGE", screenWidth/2, titleY, titleScale, color.RGBA{100, 200, 255, 255})

	// Subtitle
	DrawTextCentered(screen, "Defend the Frontier", screenWidth/2, titleY+60, 2, color.RGBA{200, 200, 200, 255})

	if m.ShowDifficultySelect {
		// Draw difficulty selection screen
		m.drawDifficultySelection(screen, screenWidth, screenHeight)
	} else if m.showLeaderboard {
		// Draw leaderboard
		leaderboard.Draw(screen, screenWidth/2, 280, 0)
		DrawTextCentered(screen, "Press L to return to menu", screenWidth/2, screenHeight-60, 1.5, color.RGBA{150, 150, 150, 255})
	} else {
		// Menu options
		y := 350

		// Pulsing "Press ENTER to Start"
		pulse := 0.8 + 0.2*math.Sin(m.animTimer*4)
		startColor := color.RGBA{uint8(100 * pulse), uint8(255 * pulse), uint8(100 * pulse), 255}
		DrawTextCentered(screen, ">> Press ENTER to Start <<", screenWidth/2, y, 2.5, startColor)

		y += 80
		DrawTextCentered(screen, "Press L for Leaderboard", screenWidth/2, y, 1.5, color.RGBA{150, 150, 200, 255})

		y += 60
		DrawTextCentered(screen, "Press I for Information", screenWidth/2, y, 1.5, color.RGBA{150, 200, 150, 255})

		y += 60
		// Sound toggle display
		soundStatus := "ON"
		soundColor := color.RGBA{100, 255, 100, 255}
		if !m.SoundEnabled {
			soundStatus = "OFF"
			soundColor = color.RGBA{255, 100, 100, 255}
		}
		DrawTextCentered(screen, "Press S to Toggle Sound: "+soundStatus, screenWidth/2, y, 1.5, soundColor)

		// Controls info
		y = screenHeight - 150
		DrawTextCentered(screen, "=== CONTROLS ===", screenWidth/2, y, 1.5, color.RGBA{255, 200, 100, 255})
		y += 30
		DrawTextCentered(screen, "WASD / Arrow Keys - Move", screenWidth/2, y, 1.2, color.RGBA{180, 180, 180, 255})
		y += 25
		DrawTextCentered(screen, "SPACE / Left Click - Fire", screenWidth/2, y, 1.2, color.RGBA{180, 180, 180, 255})
		y += 25
		DrawTextCentered(screen, "P / ESC - Pause", screenWidth/2, y, 1.2, color.RGBA{180, 180, 180, 255})
	}

	// Decorative elements
	m.drawDecorations(screen, screenWidth, screenHeight)

	// Draw update banner last (on top of everything)
	m.DrawUpdateBanner(screen, screenWidth, screenHeight)
}

// DrawUpdateBanner draws the update notification banner at the top of the screen
func (m *Menu) DrawUpdateBanner(screen *ebiten.Image, width, height int) {
	if !m.showUpdateBanner {
		return
	}

	bannerHeight := float32(60)

	// Determine banner color and text based on status
	var bannerColor color.RGBA
	var text string

	switch m.updateStatus {
	case UpdateStatusChecking:
		bannerColor = color.RGBA{100, 150, 200, 200}
		text = "Checking for updates..."
	case UpdateStatusDownloading:
		// Pulsing blue during download
		alpha := uint8(180 + 50*math.Sin(m.bannerPulse))
		bannerColor = color.RGBA{80, 120, 255, alpha}
		text = fmt.Sprintf("Downloading update %.0f%%...", m.downloadProgress*100)
	case UpdateStatusReady:
		// Pulsing green when ready
		alpha := uint8(180 + 50*math.Sin(m.bannerPulse))
		bannerColor = color.RGBA{80, 255, 80, alpha}
		text = fmt.Sprintf("Update Available: %s - Click here to restart and install", m.updateVersion)
	case UpdateStatusError:
		bannerColor = color.RGBA{255, 100, 100, 150}
		text = "Update check failed"
	default:
		return
	}

	// Draw banner background
	vector.DrawFilledRect(screen, 0, 0, float32(width), bannerHeight, bannerColor, true)

	// Draw border at bottom of banner
	borderColor := color.RGBA{bannerColor.R / 2, bannerColor.G / 2, bannerColor.B / 2, 255}
	vector.StrokeLine(screen, 0, bannerHeight, float32(width), bannerHeight, 2, borderColor, true)

	// Draw text centered in banner
	textY := int(bannerHeight / 2)
	textColor := color.RGBA{255, 255, 255, 255}
	if m.updateStatus == UpdateStatusReady {
		// Add extra emphasis for ready state
		textColor = color.RGBA{0, 0, 0, 255}
	}
	DrawTextCentered(screen, text, width/2, textY, 1.5, textColor)
}

func (m *Menu) drawDecorations(screen *ebiten.Image, width, height int) {
	// Animated corner decorations
	cornerSize := float32(30)
	cornerColor := color.RGBA{100, 150, 255, 100}

	// Top-left
	offset := float32(math.Sin(m.animTimer) * 5)
	vector.StrokeLine(screen, 20+offset, 20, 20+cornerSize+offset, 20, 2, cornerColor, true)
	vector.StrokeLine(screen, 20+offset, 20, 20+offset, 20+cornerSize, 2, cornerColor, true)

	// Top-right
	vector.StrokeLine(screen, float32(width)-20-offset, 20, float32(width)-20-cornerSize-offset, 20, 2, cornerColor, true)
	vector.StrokeLine(screen, float32(width)-20-offset, 20, float32(width)-20-offset, 20+cornerSize, 2, cornerColor, true)

	// Bottom-left
	vector.StrokeLine(screen, 20+offset, float32(height)-20, 20+cornerSize+offset, float32(height)-20, 2, cornerColor, true)
	vector.StrokeLine(screen, 20+offset, float32(height)-20, 20+offset, float32(height)-20-cornerSize, 2, cornerColor, true)

	// Bottom-right
	vector.StrokeLine(screen, float32(width)-20-offset, float32(height)-20, float32(width)-20-cornerSize-offset, float32(height)-20, 2, cornerColor, true)
	vector.StrokeLine(screen, float32(width)-20-offset, float32(height)-20, float32(width)-20-offset, float32(height)-20-cornerSize, 2, cornerColor, true)
}

func (m *Menu) drawDifficultySelection(screen *ebiten.Image, screenWidth, screenHeight int) {
	DrawTextCentered(screen, "SELECT DIFFICULTY", screenWidth/2, 250, 3, color.RGBA{255, 150, 100, 255})

	y := 350
	difficulties := []string{"EASY", "NORMAL", "HARD"}
	diffColors := []color.RGBA{
		color.RGBA{100, 255, 100, 255}, // Easy - green
		color.RGBA{255, 255, 100, 255}, // Normal - yellow
		color.RGBA{255, 100, 100, 255}, // Hard - red
	}

	boxWidth := float32(150)
	boxHeight := float32(60)
	spacing := float32(200)
	startX := float32(screenWidth/2) - spacing

	for i := 0; i < 3; i++ {
		x := startX + float32(i)*spacing
		isSelected := m.SelectedDifficulty == i

		// Draw box with selection indicator
		boxColor := diffColors[i]
		if isSelected {
			// Highlight selected difficulty
			boxColor.A = 255
			vector.DrawFilledRect(screen, x-boxWidth/2, float32(y-30), boxWidth, boxHeight, boxColor, true)
		} else {
			// Darker unselected boxes
			vector.StrokeRect(screen, x-boxWidth/2, float32(y-30), boxWidth, boxHeight, 2, color.RGBA{boxColor.R / 2, boxColor.G / 2, boxColor.B / 2, 150}, true)
		}

		// Draw difficulty text
		textColor := color.RGBA{0, 0, 0, 255}
		if !isSelected {
			textColor = boxColor
		}
		DrawTextCentered(screen, difficulties[i], int(x), y, 2, textColor)
	}

	DrawTextCentered(screen, "Use LEFT/RIGHT or A/D to select", screenWidth/2, screenHeight-150, 1.5, color.RGBA{200, 200, 200, 255})
	DrawTextCentered(screen, "Press ENTER to confirm", screenWidth/2, screenHeight-100, 1.5, color.RGBA{100, 255, 100, 255})
}
