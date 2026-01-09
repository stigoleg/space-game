package systems

import (
	"os"
	"path/filepath"
)

// GetDataPath returns the appropriate path for saving user data
// When running from .app bundle, uses ~/Library/Application Support/StellarSiege/
// When running from terminal, uses ./data/
func GetDataPath(filename string) string {
	// Check if running from .app bundle
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)

		// If running from .app/Contents/MacOS, use Application Support
		if filepath.Base(exeDir) == "MacOS" {
			homeDir, err := os.UserHomeDir()
			if err == nil {
				appSupportDir := filepath.Join(homeDir, "Library", "Application Support", "StellarSiege")
				// Ensure directory exists
				os.MkdirAll(appSupportDir, 0755)
				return filepath.Join(appSupportDir, filename)
			}
		}
	}

	// Default: use local data directory
	// Ensure directory exists
	os.MkdirAll("data", 0755)
	return filepath.Join("data", filename)
}
