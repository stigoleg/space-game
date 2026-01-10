package systems

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SemanticVersion represents a semantic version (major.minor.patch)
type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string (e.g., "v1.2.3" or "1.2.3") into a SemanticVersion
func ParseVersion(v string) (SemanticVersion, error) {
	// Remove 'v' prefix if present
	v = strings.TrimPrefix(v, "v")

	// Split by dots
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return SemanticVersion{}, fmt.Errorf("invalid version format: %s (expected major.minor.patch)", v)
	}

	// Parse each component
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return SemanticVersion{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}

	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// IsNewerThan returns true if this version is newer than the other version
func (v SemanticVersion) IsNewerThan(other SemanticVersion) bool {
	if v.Major != other.Major {
		return v.Major > other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor > other.Minor
	}
	return v.Patch > other.Patch
}

// String returns the version as a string (e.g., "1.2.3")
func (v SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// StringWithV returns the version with 'v' prefix (e.g., "v1.2.3")
func (v SemanticVersion) StringWithV() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// UpdateStatus represents the current state of the update system
type UpdateStatus int

const (
	UpdateStatusIdle UpdateStatus = iota
	UpdateStatusChecking
	UpdateStatusAvailable
	UpdateStatusDownloading
	UpdateStatusVerifying
	UpdateStatusReady
	UpdateStatusError
	UpdateStatusNoUpdate
)

// String returns a human-readable status name
func (s UpdateStatus) String() string {
	switch s {
	case UpdateStatusIdle:
		return "Idle"
	case UpdateStatusChecking:
		return "Checking"
	case UpdateStatusAvailable:
		return "Available"
	case UpdateStatusDownloading:
		return "Downloading"
	case UpdateStatusVerifying:
		return "Verifying"
	case UpdateStatusReady:
		return "Ready"
	case UpdateStatusError:
		return "Error"
	case UpdateStatusNoUpdate:
		return "NoUpdate"
	default:
		return "Unknown"
	}
}

// UpdateManager orchestrates the update process
type UpdateManager struct {
	currentVersion   SemanticVersion
	latestVersion    SemanticVersion
	status           UpdateStatus
	downloadProgress float64
	errorMessage     string
	updateReady      bool
	downloadPath     string
	checksumPath     string
	releaseInfo      *GitHubRelease

	// GitHub repository info
	gitHubOwner string
	gitHubRepo  string

	// Channels for async communication
	statusChan   chan UpdateStatus
	progressChan chan float64

	// HTTP client
	client *http.Client

	// Platform info
	platform string
	arch     string

	// Mutex for thread safety
	mutex sync.RWMutex

	// Installer
	installer *UpdateInstaller
}

// NewUpdateManager creates a new update manager
func NewUpdateManager(currentVersionStr, gitHubOwner, gitHubRepo string) *UpdateManager {
	currentVer, err := ParseVersion(currentVersionStr)
	if err != nil {
		log.Printf("Warning: Failed to parse current version %q: %v", currentVersionStr, err)
		currentVer = SemanticVersion{Major: 0, Minor: 0, Patch: 0}
	}

	um := &UpdateManager{
		currentVersion: currentVer,
		gitHubOwner:    gitHubOwner,
		gitHubRepo:     gitHubRepo,
		status:         UpdateStatusIdle,
		statusChan:     make(chan UpdateStatus, 10),
		progressChan:   make(chan float64, 10),
		client:         &http.Client{Timeout: 60 * time.Second},
		platform:       runtime.GOOS,
		arch:           runtime.GOARCH,
	}

	um.installer = NewUpdateInstaller()

	return um
}

// CheckForUpdatesAsync starts a background check for updates (non-blocking)
func (um *UpdateManager) CheckForUpdatesAsync() {
	go um.checkForUpdates()
}

// checkForUpdates performs the actual update check (runs in goroutine)
func (um *UpdateManager) checkForUpdates() {
	um.setStatus(UpdateStatusChecking)
	log.Println("Checking for updates...")

	// Fetch latest release from GitHub
	release, err := FetchLatestRelease(um.gitHubOwner, um.gitHubRepo)
	if err != nil {
		um.setError(fmt.Sprintf("Failed to check for updates: %v", err))
		log.Printf("Update check failed: %v", err)
		return
	}

	// Parse latest version
	latestVer, err := ParseVersion(release.TagName)
	if err != nil {
		um.setError(fmt.Sprintf("Invalid version format: %s", release.TagName))
		log.Printf("Failed to parse latest version %q: %v", release.TagName, err)
		return
	}

	um.mutex.Lock()
	um.latestVersion = latestVer
	um.releaseInfo = release
	um.mutex.Unlock()

	log.Printf("Current version: %s, Latest version: %s", um.currentVersion.String(), latestVer.String())

	// Compare versions
	if latestVer.IsNewerThan(um.currentVersion) {
		log.Printf("Update available: %s", latestVer.StringWithV())
		um.setStatus(UpdateStatusAvailable)

		// Start download automatically in background
		go um.downloadUpdate()
	} else {
		log.Println("No update available (already on latest version)")
		um.setStatus(UpdateStatusNoUpdate)
	}
}

// downloadUpdate downloads the update file and verifies checksum
func (um *UpdateManager) downloadUpdate() {
	um.setStatus(UpdateStatusDownloading)
	log.Println("Downloading update...")

	// Find asset for current platform
	asset, err := FindAssetForPlatform(um.releaseInfo)
	if err != nil {
		um.setError(fmt.Sprintf("No update available for your platform: %v", err))
		log.Printf("Failed to find asset: %v", err)
		return
	}

	// Find checksums file
	checksumAsset, err := FindChecksumAsset(um.releaseInfo)
	if err != nil {
		um.setError("Checksum file not found in release")
		log.Printf("Failed to find checksums: %v", err)
		return
	}

	// Create temp directory for download
	tempDir := filepath.Join(os.TempDir(), "stellar-siege-update")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		um.setError(fmt.Sprintf("Failed to create temp directory: %v", err))
		return
	}

	// Download paths
	downloadPath := filepath.Join(tempDir, asset.Name)
	checksumPath := filepath.Join(tempDir, "checksums.txt")

	// Download main asset with progress tracking
	if err := um.downloadFileWithProgress(asset.BrowserDownloadURL, downloadPath, asset.Size); err != nil {
		um.setError(fmt.Sprintf("Download failed: %v", err))
		log.Printf("Failed to download update: %v", err)
		return
	}

	log.Println("Download complete, verifying checksum...")
	um.setStatus(UpdateStatusVerifying)

	// Download checksums file
	if err := um.downloadFile(checksumAsset.BrowserDownloadURL, checksumPath); err != nil {
		um.setError(fmt.Sprintf("Failed to download checksums: %v", err))
		log.Printf("Failed to download checksums: %v", err)
		return
	}

	// Verify checksum
	expectedHash, err := FetchAndParseChecksums(checksumAsset.BrowserDownloadURL, asset.Name)
	if err != nil {
		um.setError(fmt.Sprintf("Failed to parse checksums: %v", err))
		log.Printf("Failed to parse checksums: %v", err)
		return
	}

	actualHash, err := ComputeFileHash(downloadPath)
	if err != nil {
		um.setError(fmt.Sprintf("Failed to compute file hash: %v", err))
		log.Printf("Failed to compute hash: %v", err)
		return
	}

	if actualHash != expectedHash {
		um.setError("Checksum verification failed! Download may be corrupted.")
		log.Printf("Checksum mismatch! Expected: %s, Got: %s", expectedHash, actualHash)
		// Delete corrupted file
		os.Remove(downloadPath)
		return
	}

	log.Println("Checksum verified successfully!")

	// Store paths
	um.mutex.Lock()
	um.downloadPath = downloadPath
	um.checksumPath = checksumPath
	um.updateReady = true
	um.mutex.Unlock()

	um.setStatus(UpdateStatusReady)
	log.Printf("Update ready to install: %s", downloadPath)
}

// downloadFileWithProgress downloads a file and reports progress
func (um *UpdateManager) downloadFileWithProgress(url, destPath string, totalSize int64) error {
	// Create output file
	out, err := os.Create(destPath + ".tmp")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Download
	resp, err := um.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Use totalSize from asset info or response header
	if totalSize == 0 {
		totalSize = resp.ContentLength
	}

	// Create progress reader
	var downloaded int64
	buf := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("failed to write file: %w", writeErr)
			}
			downloaded += int64(n)

			// Update progress
			if totalSize > 0 {
				progress := float64(downloaded) / float64(totalSize)
				um.setProgress(progress)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("download interrupted: %w", err)
		}
	}

	out.Close()

	// Rename from .tmp to final name (atomic operation)
	if err := os.Rename(destPath+".tmp", destPath); err != nil {
		return fmt.Errorf("failed to finalize file: %w", err)
	}

	return nil
}

// downloadFile downloads a file without progress tracking (for small files)
func (um *UpdateManager) downloadFile(url, destPath string) error {
	resp, err := um.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// InstallUpdate triggers the installation process
func (um *UpdateManager) InstallUpdate() error {
	um.mutex.RLock()
	if !um.updateReady {
		um.mutex.RUnlock()
		return fmt.Errorf("no update ready to install")
	}
	downloadPath := um.downloadPath
	um.mutex.RUnlock()

	log.Printf("Installing update from: %s", downloadPath)

	// Delegate to platform-specific installer
	return um.installer.InstallUpdate(downloadPath)
}

// GetStatus returns the current update status (thread-safe)
func (um *UpdateManager) GetStatus() UpdateStatus {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return um.status
}

// GetLatestVersion returns the latest version string (thread-safe)
func (um *UpdateManager) GetLatestVersion() string {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return um.latestVersion.StringWithV()
}

// GetDownloadProgress returns download progress 0.0-1.0 (thread-safe)
func (um *UpdateManager) GetDownloadProgress() float64 {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return um.downloadProgress
}

// GetErrorMessage returns the last error message (thread-safe)
func (um *UpdateManager) GetErrorMessage() string {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return um.errorMessage
}

// GetStatusChannel returns a channel for receiving status updates
func (um *UpdateManager) GetStatusChannel() <-chan UpdateStatus {
	return um.statusChan
}

// GetProgressChannel returns a channel for receiving progress updates
func (um *UpdateManager) GetProgressChannel() <-chan float64 {
	return um.progressChan
}

// IsUpdateReady returns true if an update is downloaded and ready to install
func (um *UpdateManager) IsUpdateReady() bool {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return um.updateReady
}

// setStatus updates the status and notifies listeners (internal use)
func (um *UpdateManager) setStatus(status UpdateStatus) {
	um.mutex.Lock()
	um.status = status
	um.mutex.Unlock()

	// Non-blocking send to channel
	select {
	case um.statusChan <- status:
	default:
	}
}

// setProgress updates download progress (internal use)
func (um *UpdateManager) setProgress(progress float64) {
	um.mutex.Lock()
	um.downloadProgress = progress
	um.mutex.Unlock()

	// Non-blocking send to channel
	select {
	case um.progressChan <- progress:
	default:
	}
}

// setError sets error status and message (internal use)
func (um *UpdateManager) setError(message string) {
	um.mutex.Lock()
	um.status = UpdateStatusError
	um.errorMessage = message
	um.mutex.Unlock()

	// Non-blocking send to channel
	select {
	case um.statusChan <- UpdateStatusError:
	default:
	}
}
