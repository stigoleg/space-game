package systems

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	Name        string        `json:"name"`
	Draft       bool          `json:"draft"`
	Prerelease  bool          `json:"prerelease"`
	PublishedAt time.Time     `json:"published_at"`
	Body        string        `json:"body"`
	Assets      []GitHubAsset `json:"assets"`
}

// GitHubAsset represents a release asset
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ContentType        string `json:"content_type"`
}

// FetchLatestRelease fetches the latest stable release from GitHub
func FetchLatestRelease(owner, repo string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Stellar-Siege-Game-Client")

	// Create client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Decode JSON
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Filter out drafts and pre-releases (stable only)
	if release.Draft {
		return nil, fmt.Errorf("latest release is a draft")
	}
	if release.Prerelease {
		return nil, fmt.Errorf("latest release is a pre-release")
	}

	return &release, nil
}

// FindAssetForPlatform finds the appropriate asset for the current platform
func FindAssetForPlatform(release *GitHubRelease) (*GitHubAsset, error) {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	expectedName := getAssetNameForPlatform(platform, arch)
	if expectedName == "" {
		return nil, fmt.Errorf("unsupported platform: %s-%s", platform, arch)
	}

	// Find matching asset
	for i := range release.Assets {
		if release.Assets[i].Name == expectedName {
			return &release.Assets[i], nil
		}
	}

	return nil, fmt.Errorf("no asset found for %s-%s (expected: %s)", platform, arch, expectedName)
}

// FindChecksumAsset finds the checksums.txt asset
func FindChecksumAsset(release *GitHubRelease) (*GitHubAsset, error) {
	for i := range release.Assets {
		if release.Assets[i].Name == "checksums.txt" {
			return &release.Assets[i], nil
		}
	}
	return nil, fmt.Errorf("checksums.txt not found in release")
}

// getAssetNameForPlatform returns the expected asset name for the platform
func getAssetNameForPlatform(platform, arch string) string {
	switch platform {
	case "darwin":
		if arch == "arm64" {
			return "Stellar-Siege-macOS-AppleSilicon.dmg"
		}
		return "Stellar-Siege-macOS-Intel.dmg"

	case "windows":
		return "stellar-siege-windows-amd64.zip"

	case "linux":
		return "stellar-siege-linux-amd64.tar.gz"

	default:
		return ""
	}
}

// GetPlatformInfo returns human-readable platform information
func GetPlatformInfo() string {
	platform := runtime.GOOS
	arch := runtime.GOARCH

	switch platform {
	case "darwin":
		if arch == "arm64" {
			return "macOS Apple Silicon"
		}
		return "macOS Intel"
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	default:
		return fmt.Sprintf("%s-%s", platform, arch)
	}
}
