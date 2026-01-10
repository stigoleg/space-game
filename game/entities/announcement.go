package entities

import (
	"image/color"
)

// ComboAnnouncement represents a combo/achievement announcement on screen
type ComboAnnouncement struct {
	Text      string
	X, Y      float64
	TimeAlive float64
	Duration  float64
	Color     color.RGBA
	Scale     float64
	Type      AnnouncementType
}

// AnnouncementType represents different types of announcements
type AnnouncementType int

const (
	AnnouncementTypeCombo AnnouncementType = iota
	AnnouncementTypeMultiKill
	AnnouncementTypeCritical
	AnnouncementTypeMilestone
	AnnouncementTypePerfectWave
	AnnouncementTypeKillSpree
	AnnouncementTypeMysteryBox
)

// AnnouncementManager manages on-screen announcements
type AnnouncementManager struct {
	Announcements []*ComboAnnouncement
	LastComboTime float64
	LastMultiKill int
}

// NewAnnouncementManager creates a new announcement manager
func NewAnnouncementManager() *AnnouncementManager {
	return &AnnouncementManager{
		Announcements: make([]*ComboAnnouncement, 0),
		LastComboTime: 0,
		LastMultiKill: 0,
	}
}

// AddComboAnnouncement creates a combo announcement
func (am *AnnouncementManager) AddComboAnnouncement(multiplier float64, screenCenterX, screenCenterY float64) {
	text := ""
	announcementColor := color.RGBA{255, 255, 255, 255}

	// Choose text based on multiplier level
	switch {
	case multiplier >= 5.0:
		text = "LEGENDARY!"
		announcementColor = color.RGBA{255, 215, 0, 255} // Gold
	case multiplier >= 4.0:
		text = "UNSTOPPABLE!"
		announcementColor = color.RGBA{255, 100, 100, 255} // Red
	case multiplier >= 3.0:
		text = "ON FIRE!"
		announcementColor = color.RGBA{255, 165, 0, 255} // Orange
	case multiplier >= 2.0:
		text = "DOUBLE COMBO!"
		announcementColor = color.RGBA{255, 200, 100, 255} // Light orange
	default:
		return // Don't announce low multipliers
	}

	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      text,
		X:         screenCenterX,
		Y:         screenCenterY - 100,
		TimeAlive: 0,
		Duration:  2.0,
		Color:     announcementColor,
		Scale:     2.5,
		Type:      AnnouncementTypeCombo,
	})

	am.LastComboTime = 0
}

// AddMultiKillAnnouncement creates a multi-kill announcement
func (am *AnnouncementManager) AddMultiKillAnnouncement(killCount int, screenCenterX, screenCenterY float64) {
	if killCount < 2 {
		return
	}

	text := ""
	announcementColor := color.RGBA{255, 255, 255, 255}

	switch killCount {
	case 2:
		text = "Double Kill!"
		announcementColor = color.RGBA{255, 180, 0, 255}
	case 3:
		text = "Triple Kill!"
		announcementColor = color.RGBA{255, 100, 0, 255}
	case 4:
		text = "Quad Kill!"
		announcementColor = color.RGBA{255, 50, 50, 255}
	case 5:
		text = "Penta Kill!"
		announcementColor = color.RGBA{255, 0, 128, 255}
	default:
		text = "KILLING SPREE!"
		announcementColor = color.RGBA{255, 0, 0, 255}
	}

	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      text,
		X:         screenCenterX,
		Y:         screenCenterY - 50,
		TimeAlive: 0,
		Duration:  2.5,
		Color:     announcementColor,
		Scale:     2.0,
		Type:      AnnouncementTypeMultiKill,
	})

	am.LastMultiKill = killCount
}

// AddCriticalHitAnnouncement creates a critical hit announcement
func (am *AnnouncementManager) AddCriticalHitAnnouncement(screenCenterX, screenCenterY float64) {
	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      "CRITICAL HIT!",
		X:         screenCenterX,
		Y:         screenCenterY,
		TimeAlive: 0,
		Duration:  1.5,
		Color:     color.RGBA{255, 255, 100, 255},
		Scale:     1.8,
		Type:      AnnouncementTypeCritical,
	})
}

// AddMilestoneAnnouncement creates a milestone announcement
func (am *AnnouncementManager) AddMilestoneAnnouncement(milestone string, screenCenterX, screenCenterY float64) {
	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      milestone,
		X:         screenCenterX,
		Y:         screenCenterY - 150,
		TimeAlive: 0,
		Duration:  3.0,
		Color:     color.RGBA{100, 255, 200, 255},
		Scale:     2.0,
		Type:      AnnouncementTypeMilestone,
	})
}

// AddPerfectWaveAnnouncement creates a perfect wave announcement
func (am *AnnouncementManager) AddPerfectWaveAnnouncement(screenCenterX, screenCenterY float64) {
	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      "FLAWLESS VICTORY!",
		X:         screenCenterX,
		Y:         screenCenterY - 200,
		TimeAlive: 0,
		Duration:  3.0,
		Color:     color.RGBA{100, 255, 100, 255},
		Scale:     2.5,
		Type:      AnnouncementTypePerfectWave,
	})
}

// AddMysteryBoxAnnouncement creates a mystery box announcement with appropriate color
func (am *AnnouncementManager) AddMysteryBoxAnnouncement(text string, isPositive bool, screenCenterX, screenCenterY float64) {
	var announcementColor color.RGBA
	var duration float64

	if isPositive {
		// Green for positive effects
		announcementColor = color.RGBA{50, 255, 50, 255}
		duration = 3.0
	} else {
		// Red for negative effects
		announcementColor = color.RGBA{255, 50, 50, 255}
		duration = 3.5 // Show negative effects slightly longer as a warning
	}

	am.Announcements = append(am.Announcements, &ComboAnnouncement{
		Text:      text,
		X:         screenCenterX,
		Y:         screenCenterY - 120,
		TimeAlive: 0,
		Duration:  duration,
		Color:     announcementColor,
		Scale:     2.0,
		Type:      AnnouncementTypeMysteryBox,
	})
}

// Update updates all announcements
func (am *AnnouncementManager) Update() {
	// Only reallocate slice if we actually need to remove announcements
	hasExpired := false
	for _, ann := range am.Announcements {
		ann.TimeAlive += 1.0 / 60.0 // 60 FPS
		if ann.TimeAlive >= ann.Duration {
			hasExpired = true
		}
	}

	// Only filter if there are expired announcements
	if hasExpired {
		// Filter in-place to avoid allocation when possible
		writeIdx := 0
		for _, ann := range am.Announcements {
			if ann.TimeAlive < ann.Duration {
				am.Announcements[writeIdx] = ann
				writeIdx++
			}
		}
		am.Announcements = am.Announcements[:writeIdx]
	}

	am.LastComboTime += 1.0 / 60.0
}

// GetAnnouncements returns all active announcements
func (am *AnnouncementManager) GetAnnouncements() []*ComboAnnouncement {
	return am.Announcements
}

// ResetMultiKill resets the multi-kill counter
func (am *AnnouncementManager) ResetMultiKill() {
	am.LastMultiKill = 0
}

// Clear removes all announcements and resets state
func (am *AnnouncementManager) Clear() {
	am.Announcements = am.Announcements[:0]
	am.LastComboTime = 0
	am.LastMultiKill = 0
}

// GetProgressAlpha calculates fade-out alpha based on announcement progress
func (ann *ComboAnnouncement) GetProgressAlpha() float64 {
	progress := ann.TimeAlive / ann.Duration

	// Fade out in last 30% of duration
	if progress > 0.7 {
		fadePercent := (progress - 0.7) / 0.3
		return 1.0 - fadePercent
	}

	return 1.0
}

// GetDisplayScale returns the current display scale with pop-in effect
func (ann *ComboAnnouncement) GetDisplayScale() float64 {
	progress := ann.TimeAlive / ann.Duration

	// Pop-in effect: scale up quickly then settle
	if progress < 0.2 {
		return ann.Scale * (0.7 + progress*1.5) // 0.7 to 1.0
	}

	return ann.Scale
}

// GetDisplayColor returns the current display color with alpha fading
func (ann *ComboAnnouncement) GetDisplayColor() color.RGBA {
	alpha := uint8(float64(ann.Color.A) * ann.GetProgressAlpha())
	return color.RGBA{
		R: ann.Color.R,
		G: ann.Color.G,
		B: ann.Color.B,
		A: alpha,
	}
}

// GetDisplayY returns the current Y position with floating effect
func (ann *ComboAnnouncement) GetDisplayY() float64 {
	progress := ann.TimeAlive / ann.Duration
	floatAmount := progress * 30.0 // Float up 30 pixels
	return ann.Y - floatAmount
}
