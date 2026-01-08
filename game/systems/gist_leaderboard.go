package systems

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// OnlineScore represents a score entry in the online leaderboard
type OnlineScore struct {
	PlayerName string    `json:"player_name"`
	Score      int64     `json:"score"`
	Difficulty string    `json:"difficulty"`
	Date       time.Time `json:"date"`
	Wave       int       `json:"wave"`
}

// GistLeaderboard manages the online leaderboard via GitHub Gist
type GistLeaderboard struct {
	GistID      string
	GitHubToken string
	client      *http.Client
	cachedData  []OnlineScore
	lastFetch   time.Time
	cacheExpiry time.Duration
}

// NewGistLeaderboard creates a new Gist-based leaderboard manager
func NewGistLeaderboard(gistID, githubToken string) *GistLeaderboard {
	return &GistLeaderboard{
		GistID:      gistID,
		GitHubToken: githubToken,
		client:      &http.Client{Timeout: 10 * time.Second},
		cachedData:  make([]OnlineScore, 0),
		cacheExpiry: 30 * time.Second, // Cache for 30 seconds
	}
}

// GetTopScores fetches the top scores from the online leaderboard (with local caching)
func (gl *GistLeaderboard) GetTopScores(limit int) ([]OnlineScore, error) {
	// Check if cache is still valid
	if time.Since(gl.lastFetch) < gl.cacheExpiry && len(gl.cachedData) > 0 {
		// Return from cache
		if limit > len(gl.cachedData) {
			return gl.cachedData, nil
		}
		return gl.cachedData[:limit], nil
	}

	// Fetch from Gist
	scores, err := gl.fetchFromGist()
	if err != nil {
		// Return cached data even if fetch fails
		if len(gl.cachedData) > 0 {
			if limit > len(gl.cachedData) {
				return gl.cachedData, nil
			}
			return gl.cachedData[:limit], nil
		}
		return nil, err
	}

	// Sort by score (highest first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Update cache
	gl.cachedData = scores
	gl.lastFetch = time.Now()

	// Return limited results
	if limit > len(scores) {
		limit = len(scores)
	}
	return scores[:limit], nil
}

// SubmitScore adds a new score to the online leaderboard
func (gl *GistLeaderboard) SubmitScore(playerName string, score int64, difficulty string, wave int) error {
	if gl.GitHubToken == "" {
		return fmt.Errorf("GitHub token not configured")
	}

	// Fetch current scores
	scores, err := gl.fetchFromGist()
	if err != nil {
		return fmt.Errorf("failed to fetch current scores: %w", err)
	}

	// Add new score
	newScore := OnlineScore{
		PlayerName: playerName,
		Score:      score,
		Difficulty: difficulty,
		Date:       time.Now(),
		Wave:       wave,
	}
	scores = append(scores, newScore)

	// Keep only top 100 scores to prevent gist from getting too large
	if len(scores) > 100 {
		sort.Slice(scores, func(i, j int) bool {
			return scores[i].Score > scores[j].Score
		})
		scores = scores[:100]
	}

	// Upload to gist
	return gl.uploadToGist(scores)
}

// fetchFromGist retrieves the score data from GitHub Gist
func (gl *GistLeaderboard) fetchFromGist() ([]OnlineScore, error) {
	// Construct the raw content URL
	url := fmt.Sprintf("https://gist.githubusercontent.com/raw/%s/leaderboard.json", gl.GistID)

	resp, err := gl.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gist: %w", err)
	}
	defer resp.Body.Close()

	// Handle 404 (gist doesn't exist yet)
	if resp.StatusCode == 404 {
		return []OnlineScore{}, nil
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gist fetch failed with status %d: %s", resp.StatusCode, string(body))
	}

	var scores []OnlineScore
	if err := json.NewDecoder(resp.Body).Decode(&scores); err != nil {
		// If decode fails, return empty slice (gist might be empty)
		return []OnlineScore{}, nil
	}

	return scores, nil
}

// uploadToGist updates the score data in GitHub Gist
func (gl *GistLeaderboard) uploadToGist(scores []OnlineScore) error {
	// Prepare the JSON payload
	jsonData, err := json.MarshalIndent(scores, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal scores: %w", err)
	}

	// Prepare the Gist update request
	updatePayload := map[string]interface{}{
		"files": map[string]interface{}{
			"leaderboard.json": map[string]interface{}{
				"content": string(jsonData),
			},
		},
	}

	payloadBytes, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	url := fmt.Sprintf("https://api.github.com/gists/%s", gl.GistID)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication header
	req.Header.Set("Authorization", fmt.Sprintf("token %s", gl.GitHubToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := gl.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update gist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gist update failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Invalidate cache to force refresh on next read
	gl.lastFetch = time.Time{}

	return nil
}

// GetPlayerRank returns the rank of a player if they're on the leaderboard
func (gl *GistLeaderboard) GetPlayerRank(playerName string, minScore int64) (int, int64, bool) {
	scores, err := gl.GetTopScores(100)
	if err != nil {
		return 0, 0, false
	}

	for rank, score := range scores {
		if score.PlayerName == playerName && score.Score == minScore {
			return rank + 1, score.Score, true
		}
	}

	return 0, 0, false
}

// ClearCache forces a refresh on next fetch
func (gl *GistLeaderboard) ClearCache() {
	gl.lastFetch = time.Time{}
	gl.cachedData = make([]OnlineScore, 0)
}
