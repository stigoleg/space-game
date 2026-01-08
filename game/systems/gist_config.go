package systems

import (
	"os"
	"strconv"
)

// GistConfig represents the GitHub Gist configuration
type GistConfig struct {
	GistID      string
	GitHubToken string
	Enabled     bool
}

// LoadGistConfig loads the Gist configuration from environment variables
// It looks for: GIST_ID, GITHUB_TOKEN, GIST_ENABLED
// The filePath parameter is kept for backwards compatibility but is ignored
func LoadGistConfig(filePath string) (*GistConfig, error) {
	config := &GistConfig{
		GistID:      os.Getenv("GIST_ID"),
		GitHubToken: os.Getenv("GITHUB_TOKEN"),
		Enabled:     parseEnvBool("GIST_ENABLED", false),
	}

	return config, nil
}

// parseEnvBool parses a boolean environment variable
// Returns the defaultValue if the environment variable is not set
func parseEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolValue
}
