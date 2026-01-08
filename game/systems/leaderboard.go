package systems

import (
	"encoding/json"
	"image/color"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type LeaderboardEntry struct {
	Rank    int       `json:"rank"`
	Name    string    `json:"name"`
	Score   int64     `json:"score"`
	Wave    int       `json:"wave"`
	Country string    `json:"country"`
	Date    time.Time `json:"date"`
}

// IP API response structure
type ipApiResponse struct {
	CountryCode string `json:"countryCode"`
	Status      string `json:"status"`
}

type Leaderboard struct {
	Entries  []LeaderboardEntry `json:"entries"`
	FilePath string             `json:"-"`
	ipCache  map[string]string  `json:"-"` // Cache IP -> Country mappings
	cacheMux sync.Mutex         `json:"-"`
}

func NewLeaderboard(filePath string) *Leaderboard {
	lb := &Leaderboard{
		Entries:  make([]LeaderboardEntry, 0),
		FilePath: filePath,
		ipCache:  make(map[string]string),
	}
	lb.Load()
	return lb
}

func (lb *Leaderboard) Load() error {
	data, err := os.ReadFile(lb.FilePath)
	if err != nil {
		// File doesn't exist, that's fine
		return nil
	}

	if err := json.Unmarshal(data, &lb.Entries); err != nil {
		return err
	}

	lb.updateRanks()
	return nil
}

func (lb *Leaderboard) Save() error {
	// Ensure directory exists
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(lb.Entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lb.FilePath, data, 0644)
}

// GetCountryFromIP fetches the country code for an IP address using IP geolocation service
func (lb *Leaderboard) GetCountryFromIP() string {
	lb.cacheMux.Lock()
	defer lb.cacheMux.Unlock()

	// Try to fetch country from IP using ip-api.com (free service)
	// Fallback available at ipapi.io if needed
	resp, err := http.Get("http://ip-api.com/json/?fields=countryCode")
	if err != nil {
		// Fallback to ipapi.io
		resp, err = http.Get("https://ipapi.co/json/")
		if err != nil {
			return "XX" // Unknown country
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &data)

		if country, ok := data["country_code"].(string); ok {
			return country
		}
		return "XX"
	}
	defer resp.Body.Close()

	var result ipApiResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &result)

	if result.Status == "success" && result.CountryCode != "" {
		return result.CountryCode
	}

	return "XX" // Unknown country
}

func (lb *Leaderboard) AddEntry(name string, score int64, wave int) {
	// Get country from IP in background (don't block)
	country := "XX"
	go func() {
		country = lb.GetCountryFromIP()
		// Update entry with country if we got it
		lb.cacheMux.Lock()
		for i := range lb.Entries {
			if lb.Entries[i].Name == name && lb.Entries[i].Score == score {
				lb.Entries[i].Country = country
				break
			}
		}
		lb.cacheMux.Unlock()
		lb.Save()
	}()

	entry := LeaderboardEntry{
		Name:    name,
		Score:   score,
		Wave:    wave,
		Country: country,
		Date:    time.Now(),
	}

	lb.Entries = append(lb.Entries, entry)
	lb.updateRanks()

	// Keep only top 10
	if len(lb.Entries) > 10 {
		lb.Entries = lb.Entries[:10]
	}

	lb.Save()
}

func (lb *Leaderboard) updateRanks() {
	// Sort by score descending
	sort.Slice(lb.Entries, func(i, j int) bool {
		return lb.Entries[i].Score > lb.Entries[j].Score
	})

	// Update ranks
	for i := range lb.Entries {
		lb.Entries[i].Rank = i + 1
	}
}

func (lb *Leaderboard) GetHighScore() int64 {
	if len(lb.Entries) == 0 {
		return 0
	}
	return lb.Entries[0].Score
}

func (lb *Leaderboard) Draw(screen *ebiten.Image, centerX, startY int, currentScore int64) {
	DrawTextCentered(screen, "=== LEADERBOARD ===", centerX, startY, 2, color.RGBA{255, 200, 50, 255})

	if len(lb.Entries) == 0 {
		DrawTextCentered(screen, "No scores yet!", centerX, startY+50, 1.5, color.RGBA{150, 150, 150, 255})
		return
	}

	y := startY + 40
	for _, entry := range lb.Entries {
		// Highlight if this is the current score
		clr := color.RGBA{200, 200, 200, 255}
		if entry.Score == currentScore {
			clr = color.RGBA{100, 255, 100, 255}
		}

		// Build leaderboard entry with country code
		country := entry.Country
		if country == "" || country == "XX" {
			country = "??"
		}

		line := FormatNumber(int64(entry.Rank)) + ". " + entry.Name + " (" + country + ") - " + FormatNumber(entry.Score) + " (Wave " + FormatNumber(int64(entry.Wave)) + ")"
		DrawTextCentered(screen, line, centerX, y, 1.5, clr)
		y += 30
	}
}
